package main

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteReadableReport(t *testing.T) {
	directory := t.TempDir()
	jsonPath := filepath.Join(directory, "integration.json")
	textPath := filepath.Join(directory, "integration.txt")
	json := "{\"Action\":\"run\",\"Test\":\"TestExample\"}\n" +
		"{\"Action\":\"output\",\"Test\":\"TestExample\",\"Output\":\"=== RUN   TestExample\\n\"}\n" +
		"{\"Action\":\"output\",\"Test\":\"TestExample\",\"Output\":\"--- PASS: TestExample (0.01s)\\n\"}\n"
	if err := os.WriteFile(jsonPath, []byte(json), 0o600); err != nil {
		t.Fatalf("write JSON report fixture: %v", err)
	}

	if err := writeReadableReport(jsonPath, textPath); err != nil {
		t.Fatalf("convert JSON integration report: %v", err)
	}
	report, err := os.ReadFile(textPath)
	if err != nil {
		t.Fatalf("read readable integration report: %v", err)
	}
	want := "=== RUN   TestExample\n--- PASS: TestExample (0.01s)\n"
	if string(report) != want {
		t.Fatalf("readable report = %q, want %q", report, want)
	}
}

func TestSetEnvironmentReplacesExistingValue(t *testing.T) {
	got := setEnvironment(
		[]string{"PATH=/bin", "GOREST_INTEGRATION_COVERAGE_DIR=old", "HOME=/tmp"},
		"GOREST_INTEGRATION_COVERAGE_DIR",
		"/tmp/coverage",
	)
	if len(got) != 3 || got[1] != "HOME=/tmp" || got[2] != "GOREST_INTEGRATION_COVERAGE_DIR=/tmp/coverage" {
		t.Fatalf("updated environment = %#v", got)
	}
}

func TestInstrumentedApplicationProducesCoverage(t *testing.T) {
	root, err := os.Getwd()
	if err != nil {
		t.Fatalf("resolve scripts directory: %v", err)
	}
	root = filepath.Dir(root)
	temporaryDirectory := t.TempDir()
	coverageDirectory := filepath.Join(temporaryDirectory, "coverage")
	if err := os.Mkdir(coverageDirectory, 0o700); err != nil {
		t.Fatalf("create coverage directory: %v", err)
	}

	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve application port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("release application port: %v", err)
	}

	binaryPath := filepath.Join(temporaryDirectory, "gorest-api-starter")
	build := exec.Command("go", "build", "-cover", "-coverpkg=.,./app/...", "-o", binaryPath, ".")
	build.Dir = root
	var buildOutput bytes.Buffer
	build.Stdout = &buildOutput
	build.Stderr = &buildOutput
	if err := build.Run(); err != nil {
		t.Fatalf("build instrumented application: %v\n%s", err, buildOutput.String())
	}

	logDirectory := filepath.Join(temporaryDirectory, "logs")
	if err := os.Mkdir(logDirectory, 0o700); err != nil {
		t.Fatalf("create application log directory: %v", err)
	}
	environment := os.Environ()
	for key, value := range map[string]string{
		"GOCOVERDIR":                           coverageDirectory,
		"GOREST_APP_ENV":                       "testing",
		"GOREST_APP_NAME":                      "gorest-coverage-smoke",
		"GOREST_APP_DEBUG":                     "false",
		"GOREST_APP_CHECK_FOR_UPDATES":         "false",
		"GOREST_APP_TIMEZONE":                  "UTC",
		"GOREST_SERVER_HOST":                   "127.0.0.1",
		"GOREST_SERVER_PORT":                   fmt.Sprint(port),
		"GOREST_SERVER_LISTENER_NETWORK":       "tcp4",
		"GOREST_SERVER_CERT_FILE":              "",
		"GOREST_SERVER_CERT_KEY_FILE":          "",
		"GOREST_SERVER_ENABLE_PREFORK":         "false",
		"GOREST_DB_DISABLE":                    "true",
		"GOREST_DB_DRIVER":                     "mysql",
		"GOREST_DB_HOST":                       "127.0.0.1",
		"GOREST_DB_PORT":                       "3306",
		"GOREST_DB_NAME":                       "coverage_smoke_test",
		"GOREST_DB_USER":                       "coverage-smoke",
		"GOREST_DB_PASSWORD":                   "coverage-smoke-password",
		"GOREST_DB_TLS":                        "disable",
		"GOREST_REDIS_DISABLE":                 "true",
		"GOREST_REDIS_TLS":                     "false",
		"GOREST_REDIS_CERT_PATH":               "",
		"GOREST_REDIS_KEY_PATH":                "",
		"GOREST_REDIS_CA_PATH":                 "",
		"GOREST_FILTER_SIGNATURE_KEY":          "coverage-smoke-test-signature-key",
		"GOREST_MONITORING_DISABLE":            "true",
		"GOREST_LOG_DIR":                       logDirectory,
		"GOREST_LOG_FILE":                      "coverage-smoke.log",
		"GOREST_LOG_CALLER":                    "false",
		"GOREST_LOG_MAX_FILE_SIZE":             "1",
		"GOREST_LOG_MAX_BACKUPS":               "1",
		"GOREST_LOG_MAX_AGE":                   "1",
		"GOREST_SERVER_CORS_ORIGINS":           "http://127.0.0.1",
		"GOREST_SERVER_CORS_ALLOW_CREDENTIALS": "false",
	} {
		environment = setEnvironment(environment, key, value)
	}

	command := exec.Command(binaryPath)
	command.Dir = temporaryDirectory
	command.Env = environment
	var stdout, stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr
	if err := command.Start(); err != nil {
		t.Fatalf("start instrumented application: %v", err)
	}
	processDone := make(chan struct{})
	var processErr error
	go func() {
		processErr = command.Wait()
		close(processDone)
	}()
	stop := func() {
		select {
		case <-processDone:
			return
		default:
		}
		if err := command.Process.Signal(os.Interrupt); err != nil {
			_ = command.Process.Kill()
		}
		select {
		case <-processDone:
		case <-time.After(10 * time.Second):
			_ = command.Process.Kill()
			<-processDone
			t.Error("instrumented application did not stop gracefully")
		}
	}
	t.Cleanup(stop)

	client := &http.Client{Timeout: time.Second}
	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	deadline := time.Now().Add(20 * time.Second)
	ready := false
	lastStatus := 0
	for time.Now().Before(deadline) {
		select {
		case <-processDone:
			t.Fatalf("instrumented application exited before readiness: %v\nstdout:\n%s\nstderr:\n%s",
				processErr, stdout.String(), stderr.String())
		default:
		}
		response, requestErr := client.Get(baseURL + "/")
		if requestErr == nil {
			lastStatus = response.StatusCode
			_ = response.Body.Close()
			if response.StatusCode == http.StatusOK {
				ready = true
				break
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	if !ready {
		t.Fatalf("instrumented application did not become ready (last HTTP status %d)\nstdout:\n%s\nstderr:\n%s",
			lastStatus, stdout.String(), stderr.String())
	}
	stop()

	entries, err := os.ReadDir(coverageDirectory)
	if err != nil {
		t.Fatalf("read instrumented application coverage data: %v", err)
	}
	if len(entries) == 0 {
		t.Fatalf("instrumented application produced no coverage data (exit=%v)\nstdout:\n%s\nstderr:\n%s",
			processErr, stdout.String(), stderr.String())
	}
	profilePath := filepath.Join(temporaryDirectory, "integration.coverprofile")
	convert := exec.Command("go", "tool", "covdata", "textfmt", "-i="+coverageDirectory, "-o="+profilePath)
	if output, err := convert.CombinedOutput(); err != nil {
		t.Fatalf("convert application coverage data: %v\n%s", err, output)
	}
	cover := exec.Command("go", "tool", "cover", "-func="+profilePath)
	output, err := cover.CombinedOutput()
	if err != nil {
		t.Fatalf("summarize application coverage: %v\n%s", err, output)
	}
	summary := string(output)
	if !strings.Contains(summary, "app/config") || !strings.Contains(summary, "total:") || strings.Contains(summary, "total: (statements) 0.0%") {
		t.Fatalf("application coverage summary does not show exercised app code: %s", summary)
	}
}
