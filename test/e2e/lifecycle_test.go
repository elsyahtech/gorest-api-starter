package e2e_test

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest/database"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/server"
)

const lifecycleHelperEnv = "GOREST_LIFECYCLE_HELPER"

func TestDatabaseConnectionFailureAbortsStartup(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve local database port: %v", err)
	}
	databasePort := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("release local database port: %v", err)
	}

	logDir := filepath.Join(t.TempDir(), "logs")
	command := lifecycleHelperCommand(t, "database-failure", map[string]string{
		"GOREST_LIFECYCLE_DB_PORT": strconv.Itoa(databasePort),
		"GOREST_LIFECYCLE_LOG_DIR": logDir,
	})
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("application startup succeeded against an unavailable database; output: %s", output)
	}
	if !strings.Contains(string(output), "failed initializing database") {
		t.Fatalf("startup exited for an unexpected reason; output: %s", output)
	}
	if strings.Contains(string(output), "database initialized successfully") {
		t.Fatalf("startup reported database initialization success despite connection failure; output: %s", output)
	}
}

func TestMigrationFailureAbortsStartupTasks(t *testing.T) {
	databaseDir := t.TempDir()
	migrationDir := filepath.Join(databaseDir, "migrations", database.SQLITE)
	if err := os.MkdirAll(migrationDir, 0o700); err != nil {
		t.Fatalf("create SQLite migration directory: %v", err)
	}
	const migrationName = "999_lifecycle_failure.sql"
	if err := os.WriteFile(filepath.Join(migrationDir, migrationName), []byte("THIS IS NOT VALID SQL;"), 0o600); err != nil {
		t.Fatalf("write invalid SQLite migration: %v", err)
	}

	logConfig := gorestlog.ConfigDefault
	logConfig.LogDir = filepath.Join(t.TempDir(), "logs")
	app := gorest.New(gorest.Config{AppName: "lifecycle-migration-test", Environment: "test", Timezone: "UTC", Debug: true})
	databaseConfig := database.ConfigDefault
	databaseConfig.Driver = database.SQLITE
	databaseConfig.Name = "lifecycle-migration-failure"
	databaseConfig.DatabaseDirectory = databaseDir
	databaseConfig.MigrationDirectory = "migrations"
	databaseConfig.Migration = true
	databaseConfig.Seeder = false
	app.Use(gorestlog.Run(logConfig), database.Run(databaseConfig))
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close app after migration failure (%s): %v", message, err)
		}
	})

	if err := app.RunDatabaseStartupTasks(); err == nil || !strings.Contains(err.Error(), migrationName) {
		t.Fatalf("startup migration error = %v, want failure identifying %q", err, migrationName)
	}
}

func TestShutdownWaitsForInFlightRequest(t *testing.T) {
	portListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve local HTTP port: %v", err)
	}
	port := portListener.Addr().(*net.TCPAddr).Port
	if err := portListener.Close(); err != nil {
		t.Fatalf("release local HTTP port: %v", err)
	}

	workDir := t.TempDir()
	startedPath := filepath.Join(workDir, "request-started")
	releasePath := filepath.Join(workDir, "release-request")
	logDir := filepath.Join(workDir, "logs")
	command := lifecycleHelperCommand(t, "graceful-shutdown", map[string]string{
		"GOREST_LIFECYCLE_HTTP_PORT":    strconv.Itoa(port),
		"GOREST_LIFECYCLE_STARTED_PATH": startedPath,
		"GOREST_LIFECYCLE_RELEASE_PATH": releasePath,
		"GOREST_LIFECYCLE_LOG_DIR":      logDir,
	})
	var output bytes.Buffer
	command.Stdout = io.Discard
	command.Stderr = &output
	if err := command.Start(); err != nil {
		t.Fatalf("start lifecycle subprocess: %v", err)
	}
	commandDone := make(chan error, 1)
	waitFinished := make(chan struct{})
	go func() {
		commandDone <- command.Wait()
		close(waitFinished)
	}()
	t.Cleanup(func() {
		select {
		case <-waitFinished:
		default:
			_ = command.Process.Kill()
			<-waitFinished
		}
	})

	endpoint := fmt.Sprintf("http://127.0.0.1:%d/slow", port)
	responseDone := make(chan struct {
		status int
		body   string
		err    error
	}, 1)
	go func() {
		client := &http.Client{Timeout: 8 * time.Second}
		deadline := time.Now().Add(5 * time.Second)
		for time.Now().Before(deadline) {
			response, requestErr := client.Get(endpoint)
			if requestErr != nil {
				time.Sleep(20 * time.Millisecond)
				continue
			}
			body, readErr := io.ReadAll(response.Body)
			closeErr := response.Body.Close()
			if readErr != nil {
				requestErr = readErr
			} else if closeErr != nil {
				requestErr = closeErr
			}
			responseDone <- struct {
				status int
				body   string
				err    error
			}{status: response.StatusCode, body: string(body), err: requestErr}
			return
		}
		responseDone <- struct {
			status int
			body   string
			err    error
		}{err: fmt.Errorf("server did not accept the HTTP request before timeout")}
	}()

	if err := waitForFile(startedPath, 6*time.Second); err != nil {
		t.Fatalf("wait for in-flight request to enter handler: %v; child output: %s", err, output.String())
	}
	if err := command.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("send shutdown signal: %v", err)
	}
	select {
	case err := <-commandDone:
		t.Fatalf("application exited before in-flight request was released: %v; output: %s", err, output.String())
	case <-time.After(150 * time.Millisecond):
	}

	if err := os.WriteFile(releasePath, []byte("release"), 0o600); err != nil {
		t.Fatalf("release in-flight request: %v", err)
	}
	select {
	case response := <-responseDone:
		if response.err != nil {
			t.Fatalf("in-flight HTTP response: %v", response.err)
		}
		if response.status != http.StatusOK || response.body != "completed" {
			t.Fatalf("in-flight response = %d %q, want %d %q", response.status, response.body, http.StatusOK, "completed")
		}
	case <-time.After(6 * time.Second):
		t.Fatal("in-flight request did not complete after release")
	}
	select {
	case err := <-commandDone:
		if err != nil {
			t.Fatalf("application shutdown failed: %v; output: %s", err, output.String())
		}
	case <-time.After(6 * time.Second):
		t.Fatal("application did not exit after the in-flight request completed")
	}
}

// TestLifecycleSubprocess hosts process-level scenarios because GOREST's
// startup failure path intentionally terminates the process and Start handles
// operating-system shutdown signals.
func TestLifecycleSubprocess(t *testing.T) {
	switch os.Getenv(lifecycleHelperEnv) {
	case "database-failure":
		logConfig := gorestlog.ConfigDefault
		logConfig.LogDir = os.Getenv("GOREST_LIFECYCLE_LOG_DIR")
		app := gorest.New(gorest.Config{AppName: "lifecycle-db-failure", Environment: "test", Timezone: "UTC", Debug: true})
		databaseConfig := database.ConfigDefault
		databaseConfig.Driver = database.POSTGRES
		databaseConfig.Host = "127.0.0.1"
		databaseConfig.Port = os.Getenv("GOREST_LIFECYCLE_DB_PORT")
		databaseConfig.Name = "gorest_lifecycle_test"
		databaseConfig.Username = "test"
		databaseConfig.Password = "test"
		databaseConfig.Timeout = 300 * time.Millisecond
		databaseConfig.Migration = false
		databaseConfig.Seeder = false
		app.Use(gorestlog.Run(logConfig), database.Run(databaseConfig))
		t.Fatal("database initialization unexpectedly succeeded")
	case "graceful-shutdown":
		logConfig := gorestlog.ConfigDefault
		logConfig.LogDir = os.Getenv("GOREST_LIFECYCLE_LOG_DIR")
		serverConfig := server.ConfigDefault
		serverConfig.Host = "127.0.0.1"
		port, err := strconv.Atoi(os.Getenv("GOREST_LIFECYCLE_HTTP_PORT"))
		if err != nil {
			t.Fatalf("parse lifecycle HTTP port: %v", err)
		}
		serverConfig.Port = port
		app := gorest.New(gorest.Config{AppName: "lifecycle-shutdown", Environment: "test", Timezone: "UTC"})
		app.Use(gorestlog.Run(logConfig), server.Run(serverConfig))
		app.RegisterServices([]gorest.RouterRegistrar{
			func(app *gorest.App) {
				app.Get("/slow", func(ctx *server.Context) error {
					if err := os.WriteFile(os.Getenv("GOREST_LIFECYCLE_STARTED_PATH"), []byte("started"), 0o600); err != nil {
						return err
					}
					releasePath := os.Getenv("GOREST_LIFECYCLE_RELEASE_PATH")
					deadline := time.Now().Add(7 * time.Second)
					for time.Now().Before(deadline) {
						if _, err := os.Stat(releasePath); err == nil {
							return ctx.SendString("completed")
						} else if !os.IsNotExist(err) {
							return err
						}
						time.Sleep(10 * time.Millisecond)
					}
					return fmt.Errorf("timed out waiting for request release")
				})
			},
		})
		app.Start()
	default:
		return
	}
}

func lifecycleHelperCommand(t *testing.T, mode string, values map[string]string) *exec.Cmd {
	t.Helper()
	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("locate current test executable: %v", err)
	}
	command := exec.Command(executable, "-test.run=^TestLifecycleSubprocess$", "-test.v")
	environment := append(os.Environ(), lifecycleHelperEnv+"="+mode)
	for key, value := range values {
		environment = append(environment, key+"="+value)
	}
	command.Env = environment
	return command
}

func waitForFile(path string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if _, err := os.Stat(path); err == nil {
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}
		time.Sleep(10 * time.Millisecond)
	}
	return fmt.Errorf("file %s was not created within %s", path, timeout)
}
