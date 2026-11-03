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
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest-api-starter/app/entity"
	"github.com/elsyahtech/gorest/database"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/server"
)

const lifecycleHelperEnv = "GOREST_LIFECYCLE_HELPER"

func TestDatabaseConnectionFailureAbortsStartup(t *testing.T) {
	driver := strings.TrimSpace(os.Getenv("GOREST_E2E_DRIVER"))
	if driver == "" {
		driver = database.DefaultDatabaseDriver
	}
	if database.NormalizeDatabaseDriver(driver) == database.SQLITE {
		t.Skip("SQLite is an embedded database; network connection failure is not applicable")
	}

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
	if !strings.Contains(string(output), "failed to run database") {
		t.Fatalf("startup exited for an unexpected reason; output: %s", output)
	}
	if strings.Contains(string(output), "database initialized successfully") {
		t.Fatalf("startup reported database initialization success despite connection failure; output: %s", output)
	}
}

func TestProductionSQLiteAllowsOnlyLocalHosts(t *testing.T) {
	for _, host := range []string{"", "localhost", "127.0.0.1", "::1"} {
		t.Run("allows "+host, func(t *testing.T) {
			logConfig := gorestlog.ConfigDefault
			logConfig.LogDir = filepath.Join(t.TempDir(), "logs")
			databaseConfig := database.ConfigDefault
			databaseConfig.Driver = database.SQLITE
			databaseConfig.Host = host
			databaseConfig.Name = "production-sqlite-host-test"
			databaseConfig.DatabaseDirectory = t.TempDir()
			databaseConfig.Migration = false
			databaseConfig.Seeder = false

			app := gorest.New(gorest.Config{
				AppName:     "production-sqlite-host-test",
				Environment: gorest.Production,
				Timezone:    "UTC",
			})
			app.Use(gorestlog.Run(logConfig), database.Run(databaseConfig))

			t.Cleanup(func() {
				if message, _, err := app.Close(); err != nil {
					t.Errorf("close production SQLite app (%s): %v", message, err)
				}
			})

			if app.Database() == nil {
				t.Fatal("SQLite database session is nil, want local database initialized")
			}
		})
	}
}

func TestSQLiteCreatesDatabaseWithRestrictedPermissions(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX file permission modes are not portable to Windows")
	}

	root := t.TempDir()
	dataDir := filepath.Join(root, "new", "private", "database")
	app := openSQLiteFilesystemTestApp(t, dataDir, "new-file")

	databaseInfo, err := os.Stat(filepath.Join(dataDir, "new-file.db"))
	if err != nil {
		t.Fatalf("stat newly created SQLite database: %v", err)
	}
	probeDirectory := filepath.Join(root, "umask-probe")
	if err := os.Mkdir(probeDirectory, 0o777); err != nil {
		t.Fatalf("create umask probe directory: %v", err)
	}
	probeInfo, err := os.Stat(probeDirectory)
	if err != nil {
		t.Fatalf("stat umask probe directory: %v", err)
	}
	effectiveUmask := ^probeInfo.Mode().Perm() & 0o777
	wantDirectoryMode := os.FileMode(0o700) &^ effectiveUmask
	wantFileMode := os.FileMode(0o600) &^ effectiveUmask

	for _, directory := range []string{
		filepath.Join(root, "new"),
		filepath.Join(root, "new", "private"),
		dataDir,
	} {
		directoryInfo, err := os.Stat(directory)
		if err != nil {
			t.Fatalf("stat newly created SQLite directory %q: %v", directory, err)
		}
		if got := directoryInfo.Mode().Perm(); got != wantDirectoryMode {
			t.Errorf("new SQLite directory %q permissions = %04o, want %04o (0700 subject to umask)",
				directory, got, wantDirectoryMode)
		}
	}
	if got := databaseInfo.Mode().Perm(); got != wantFileMode {
		t.Errorf("new SQLite database permissions = %04o, want %04o (0600 subject to umask)",
			got, wantFileMode)
	}
	if app.Database() == nil {
		t.Fatal("SQLite database session is nil after creating the file")
	}
}

func TestSQLitePreservesExistingPermissionsAndOwnership(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX ownership and permission behavior is not portable to Windows")
	}

	dataDir := filepath.Join(t.TempDir(), "existing")
	if err := os.MkdirAll(dataDir, 0o700); err != nil {
		t.Fatalf("create existing SQLite directory: %v", err)
	}
	const directoryMode = os.FileMode(0o750)
	if err := os.Chmod(dataDir, directoryMode); err != nil {
		t.Fatalf("set existing SQLite directory permissions: %v", err)
	}

	databasePath := filepath.Join(dataDir, "existing-file.db")
	file, err := os.OpenFile(databasePath, os.O_CREATE|os.O_RDWR, 0o600)
	if err != nil {
		t.Fatalf("create existing SQLite database file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("close existing SQLite database file: %v", err)
	}
	const fileMode = os.FileMode(0o640)
	if err := os.Chmod(databasePath, fileMode); err != nil {
		t.Fatalf("set existing SQLite database permissions: %v", err)
	}

	beforeDirectory, err := os.Stat(dataDir)
	if err != nil {
		t.Fatalf("stat existing SQLite directory before startup: %v", err)
	}
	beforeFile, err := os.Stat(databasePath)
	if err != nil {
		t.Fatalf("stat existing SQLite database before startup: %v", err)
	}
	beforeDirectoryOwner := sqliteFileOwner(t, beforeDirectory)
	beforeFileOwner := sqliteFileOwner(t, beforeFile)

	app := openSQLiteFilesystemTestApp(t, dataDir, "existing-file")
	if app.Database() == nil {
		t.Fatal("SQLite database session is nil after opening existing file")
	}

	afterDirectory, err := os.Stat(dataDir)
	if err != nil {
		t.Fatalf("stat existing SQLite directory after startup: %v", err)
	}
	afterFile, err := os.Stat(databasePath)
	if err != nil {
		t.Fatalf("stat existing SQLite database after startup: %v", err)
	}

	if got := afterDirectory.Mode().Perm(); got != directoryMode {
		t.Errorf("existing SQLite directory permissions = %04o, want unchanged %04o", got, directoryMode)
	}
	if got := afterFile.Mode().Perm(); got != fileMode {
		t.Errorf("existing SQLite database permissions = %04o, want unchanged %04o", got, fileMode)
	}
	if got := sqliteFileOwner(t, afterDirectory); got != beforeDirectoryOwner {
		t.Errorf("existing SQLite directory owner changed from %s to %s", beforeDirectoryOwner, got)
	}
	if got := sqliteFileOwner(t, afterFile); got != beforeFileOwner {
		t.Errorf("existing SQLite database owner changed from %s to %s", beforeFileOwner, got)
	}
}

func TestProductionSQLitePersistsLocalDataAcrossRestart(t *testing.T) {
	databaseDirectory := t.TempDir()
	t.Setenv("APP_ENV", gorest.Production)
	t.Setenv("GOREST_DB_DRIVER", database.SQLITE)
	t.Setenv("GOREST_DB_HOST", "")
	t.Setenv("GOREST_DB_PORT", "")
	t.Setenv("GOREST_DB_NAME", "production-persistence")
	t.Setenv("GOREST_DB_DIRECTORY", databaseDirectory)
	t.Setenv("GOREST_DB_MIGRATION", "false")
	t.Setenv("GOREST_DB_SEEDER", "false")

	databaseConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load production SQLite configuration: %v", err)
	}
	if databaseConfig.Migration || databaseConfig.Seeder {
		t.Fatalf("production SQLite startup tasks = migration:%t seeder:%t, want both disabled",
			databaseConfig.Migration, databaseConfig.Seeder)
	}

	openProductionApp := func() (*gorest.App, func()) {
		logConfig := gorestlog.ConfigDefault
		logConfig.LogDir = filepath.Join(databaseDirectory, "logs")
		app := gorest.New(gorest.Config{
			AppName:     "production-sqlite-persistence-test",
			Environment: gorest.Production,
			Timezone:    "UTC",
		})
		app.Use(gorestlog.Run(logConfig), database.Run(databaseConfig))

		closed := false
		closeApp := func() {
			if closed {
				return
			}
			closed = true
			if message, _, err := app.Close(); err != nil {
				t.Errorf("close production SQLite app (%s): %v", message, err)
			}
		}
		t.Cleanup(closeApp)
		return app, closeApp
	}

	firstApp, closeFirstApp := openProductionApp()
	ctx, cancel := firstApp.NewContext(10 * time.Second)
	_, _, _, err = firstApp.Database().ExecSQL(ctx, `CREATE TABLE departments (
		id VARCHAR(36) NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		is_active INTEGER NOT NULL DEFAULT 1,
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME,
		deleted_at DATETIME
	)`)
	cancel()
	if err != nil {
		t.Fatalf("create persistence probe table in production SQLite: %v", err)
	}
	department := entity.Department{
		ID:          "production-sqlite-persistence",
		Name:        "Persistent department",
		Description: "Verifies local database survives application restart",
		IsActive:    true,
		CreatedAt:   time.Now().UTC(),
	}
	query := firstApp.Database().Table("departments").Create(&department)
	if query.Error != nil {
		t.Fatalf("write production SQLite data: %s: %v", query.Message, query.Error)
	}
	closeFirstApp()

	restartedApp, _ := openProductionApp()
	var persisted entity.Department
	query = restartedApp.Database().Table("departments").
		Where("id = ?", department.ID).Find(&persisted)
	if query.Error != nil {
		t.Fatalf("read production SQLite data after restart: %s: %v", query.Message, query.Error)
	}
	if persisted.ID != department.ID || persisted.Name != department.Name {
		t.Fatalf("production SQLite row after restart = %+v, want ID %q and name %q",
			persisted, department.ID, department.Name)
	}
}

func openSQLiteFilesystemTestApp(t *testing.T, dataDir, name string) *gorest.App {
	t.Helper()

	databaseConfig := database.ConfigDefault
	databaseConfig.Driver = database.SQLITE
	databaseConfig.Name = name
	databaseConfig.DatabaseDirectory = dataDir
	databaseConfig.Migration = false
	databaseConfig.Seeder = false
	databaseConfig.MaxOpenConnections = 1
	databaseConfig.MaxIdleConnections = 1
	databaseConfig.ConnectionMaxLifetime = 0

	app := gorest.New(gorest.Config{
		AppName:     "sqlite-filesystem-permissions-test",
		Environment: gorest.Testing,
		Timezone:    "UTC",
	})
	app.Use(database.Run(databaseConfig))
	t.Cleanup(func() {
		if message, _, err := app.Close(); err != nil {
			t.Errorf("close SQLite filesystem test app (%s): %v", message, err)
		}
	})
	if app.Database() == nil {
		t.Fatal("SQLite database session is nil")
	}
	return app
}

func sqliteFileOwner(t *testing.T, info os.FileInfo) string {
	t.Helper()

	systemInfo := reflect.ValueOf(info.Sys())
	if systemInfo.Kind() == reflect.Pointer {
		systemInfo = systemInfo.Elem()
	}
	if !systemInfo.IsValid() || systemInfo.Kind() != reflect.Struct {
		t.Fatalf("filesystem metadata has unsupported owner information: %T", info.Sys())
	}

	uid := systemInfo.FieldByName("Uid")
	gid := systemInfo.FieldByName("Gid")
	if !uid.IsValid() || !gid.IsValid() {
		t.Fatalf("filesystem metadata %T does not expose POSIX UID/GID", info.Sys())
	}
	return fmt.Sprintf("%v:%v", uid.Interface(), gid.Interface())
}

func TestProductionSQLiteRejectsRemoteHost(t *testing.T) {
	logDir := filepath.Join(t.TempDir(), "logs")
	databaseDir := filepath.Join(t.TempDir(), "database")
	command := lifecycleHelperCommand(t, "production-sqlite-remote-host", map[string]string{
		"GOREST_LIFECYCLE_DB_DIRECTORY": databaseDir,
		"GOREST_LIFECYCLE_DB_HOST":      "db.internal.example",
		"GOREST_LIFECYCLE_LOG_DIR":      logDir,
	})
	output, err := command.CombinedOutput()
	if err == nil {
		t.Fatalf("production SQLite startup accepted a remote host; output: %s", output)
	}
	if !strings.Contains(string(output), "SQLite database must use local storage in production") {
		t.Fatalf("production SQLite rejected host for an unexpected reason; output: %s", output)
	}
	if strings.Contains(string(output), "context deadline exceeded") {
		t.Fatalf("production SQLite attempted a network connection; output: %s", output)
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
	app := gorest.New(gorest.Config{AppName: "lifecycle-migration-test", Environment: gorest.Testing, Timezone: "UTC", Debug: true})
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
		app := gorest.New(gorest.Config{AppName: "lifecycle-db-failure", Environment: gorest.Testing, Timezone: "UTC", Debug: true})
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
	case "production-sqlite-remote-host":
		logConfig := gorestlog.ConfigDefault
		logConfig.LogDir = os.Getenv("GOREST_LIFECYCLE_LOG_DIR")
		app := gorest.New(gorest.Config{
			AppName:     "production-sqlite-remote-host",
			Environment: gorest.Production,
			Timezone:    "UTC",
		})
		databaseConfig := database.ConfigDefault
		databaseConfig.Driver = database.SQLITE
		databaseConfig.Host = os.Getenv("GOREST_LIFECYCLE_DB_HOST")
		databaseConfig.Name = "production-sqlite-remote-host"
		databaseConfig.DatabaseDirectory = os.Getenv("GOREST_LIFECYCLE_DB_DIRECTORY")
		databaseConfig.Migration = false
		databaseConfig.Seeder = false
		app.Use(gorestlog.Run(logConfig), database.Run(databaseConfig))
		t.Fatal("production SQLite startup unexpectedly accepted a remote host")
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
		app := gorest.New(gorest.Config{AppName: "lifecycle-shutdown", Environment: gorest.Testing, Timezone: "UTC"})
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
