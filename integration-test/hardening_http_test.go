//go:build integration

package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/elsyahtech/gorest"
	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/filter"
	"github.com/elsyahtech/gorest/hardening"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/monitoring"
	"github.com/elsyahtech/gorest/server"
)

const hardeningTestProcess = "GOREST_HARDENING_TEST_PROCESS"

func TestHardeningEndpointIntegration(t *testing.T) {
	t.Run("disabled by default", func(t *testing.T) {
		baseURL := startHardeningApplication(t, nil)
		assertHardeningStatus(t, baseURL, "/gorest/hardening", nil, http.StatusNotFound)
	})

	t.Run("forwarded headers do not bypass the IP allowlist", func(t *testing.T) {
		baseURL := startHardeningApplication(t, map[string]string{
			"GOREST_HARDENING_ENABLE":      "true",
			"GOREST_HARDENING_ALLOWED_IPS": "192.0.2.10",
			"GOREST_MONITORING_DISABLE":    "true",
		})
		assertHardeningStatus(t, baseURL, "/gorest/hardening", map[string]string{
			"X-Forwarded-For": "192.0.2.10",
		}, http.StatusForbidden)
	})

	t.Run("allowlisted peer receives redacted no-store reports", func(t *testing.T) {
		const signatureKey = "sqlite-starter-hardening-test-signing-key"
		baseURL := startHardeningApplication(t, map[string]string{
			"GOREST_HARDENING_ENABLE":      "true",
			"GOREST_HARDENING_ALLOWED_IPS": "127.0.0.1/32",
			"GOREST_MONITORING_DISABLE":    "true",
			"GOREST_FILTER_SIGNATURE_KEY":  signatureKey,
		})

		jsonResponse := getHardeningResponse(t, baseURL, "/gorest/hardening", nil)
		defer jsonResponse.Body.Close()
		if jsonResponse.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(jsonResponse.Body)
			t.Fatalf("JSON report status = %d, want %d; body: %s",
				jsonResponse.StatusCode, http.StatusOK, body)
		}
		assertHardeningResponseHeaders(t, jsonResponse, "application/json")
		body, err := io.ReadAll(jsonResponse.Body)
		if err != nil {
			t.Fatalf("read JSON report: %v", err)
		}
		if strings.Contains(string(body), signatureKey) {
			t.Fatal("JSON report exposed the configured signing key")
		}
		var report struct {
			EvidenceMode string `json:"evidence_mode"`
			Inventory    []struct {
				Path string `json:"field_path"`
			} `json:"inventory"`
		}
		if err := json.Unmarshal(body, &report); err != nil {
			t.Fatalf("decode hardening report: %v", err)
		}
		if report.EvidenceMode != "effective_runtime_config" || len(report.Inventory) == 0 {
			t.Fatalf("report evidence/inventory = (%q, %d), want runtime config inventory",
				report.EvidenceMode, len(report.Inventory))
		}

		htmlResponse := getHardeningResponse(t, baseURL, "/gorest/hardening?format=html", nil)
		defer htmlResponse.Body.Close()
		if htmlResponse.StatusCode != http.StatusOK {
			t.Fatalf("HTML report status = %d, want %d", htmlResponse.StatusCode, http.StatusOK)
		}
		assertHardeningResponseHeaders(t, htmlResponse, "text/html")
		htmlBody, err := io.ReadAll(htmlResponse.Body)
		if err != nil {
			t.Fatalf("read HTML report: %v", err)
		}
		if strings.Contains(string(htmlBody), signatureKey) {
			t.Fatal("HTML report exposed the configured signing key")
		}

		assertHardeningStatus(t, baseURL, "/gorest/hardening?format=xml", nil, http.StatusBadRequest)
	})
}

func TestHardeningApplicationSubprocess(t *testing.T) {
	if os.Getenv(hardeningTestProcess) != "1" {
		return
	}

	appConfig, err := config.AppConfig()
	if err != nil {
		t.Fatalf("load app config: %v", err)
	}
	loggerConfig, err := config.LoggerConfig()
	if err != nil {
		t.Fatalf("load logger config: %v", err)
	}
	serverConfig, err := config.ServerConfig()
	if err != nil {
		t.Fatalf("load server config: %v", err)
	}
	filterConfig, err := config.FilterConfig()
	if err != nil {
		t.Fatalf("load filter config: %v", err)
	}
	monitoringConfig, err := config.MonitoringConfig()
	if err != nil {
		t.Fatalf("load monitoring config: %v", err)
	}
	hardeningConfig, err := config.HardeningConfig()
	if err != nil {
		t.Fatalf("load hardening config: %v", err)
	}

	app := gorest.New(*appConfig)
	app.Use(
		gorestlog.Run(*loggerConfig),
		server.Run(*serverConfig),
		filter.Run(*filterConfig),
		monitoring.Run(*monitoringConfig),
		hardening.Run(*hardeningConfig),
	)
	app.Start()
}

func startHardeningApplication(t *testing.T, overrides map[string]string) string {
	t.Helper()

	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve hardening test port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("release hardening test port: %v", err)
	}

	environment := map[string]string{
		"GOREST_APP_ENV":               gorest.Testing,
		"GOREST_APP_CHECK_FOR_UPDATES": "false",
		"GOREST_LOG_DIR":               filepath.Join(t.TempDir(), "logs"),
		"GOREST_SERVER_HOST":           "127.0.0.1",
		"GOREST_SERVER_PORT":           strconv.Itoa(port),
		hardeningTestProcess:           "1",
	}
	for key, value := range overrides {
		environment[key] = value
	}

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("locate integration test executable: %v", err)
	}
	command := exec.Command(executable, "-test.run=^TestHardeningApplicationSubprocess$")
	command.Dir = t.TempDir()
	command.Env = sanitizedHardeningEnvironment(environment)
	var output synchronizedTestBuffer
	command.Stdout = &output
	command.Stderr = &output
	if err := command.Start(); err != nil {
		t.Fatalf("start hardening test application: %v", err)
	}
	t.Cleanup(func() {
		_ = command.Process.Signal(os.Interrupt)
		done := make(chan error, 1)
		go func() { done <- command.Wait() }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = command.Process.Kill()
			<-done
		}
	})

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	client := &http.Client{Timeout: time.Second}
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		response, requestErr := client.Get(baseURL + "/gorest/hardening")
		if requestErr == nil {
			response.Body.Close()
			return baseURL
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("hardening test application failed to start: %s", output.String())
	return ""
}

func sanitizedHardeningEnvironment(overrides map[string]string) []string {
	environment := make([]string, 0, len(os.Environ())+len(overrides))
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if key == "GOREST_APP_ENV" || strings.HasPrefix(key, "GOREST_") {
			continue
		}
		environment = append(environment, entry)
	}
	for key, value := range overrides {
		environment = append(environment, key+"="+value)
	}
	return environment
}

func getHardeningResponse(t *testing.T, baseURL, path string, headers map[string]string) *http.Response {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		t.Fatalf("create hardening request: %v", err)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(request)
	if err != nil {
		t.Fatalf("request hardening endpoint: %v", err)
	}
	return response
}

func assertHardeningStatus(t *testing.T, baseURL, path string, headers map[string]string, want int) {
	t.Helper()

	response := getHardeningResponse(t, baseURL, path, headers)
	defer response.Body.Close()
	if response.StatusCode != want {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("hardening endpoint %s status = %d, want %d; body: %s",
			path, response.StatusCode, want, body)
	}
}

func assertHardeningResponseHeaders(t *testing.T, response *http.Response, contentType string) {
	t.Helper()
	if got := response.Header.Get("Content-Type"); !strings.HasPrefix(got, contentType) {
		t.Errorf("Content-Type = %q, want prefix %q", got, contentType)
	}
	if got := response.Header.Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
	if got := response.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q, want nosniff", got)
	}
}

type synchronizedTestBuffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (buffer *synchronizedTestBuffer) Write(value []byte) (int, error) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	return buffer.buffer.Write(value)
}

func (buffer *synchronizedTestBuffer) String() string {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	return buffer.buffer.String()
}
