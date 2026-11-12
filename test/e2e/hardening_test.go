package e2e_test

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
	appconfig "github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/filter"
	"github.com/elsyahtech/gorest/hardening"
	gorestlog "github.com/elsyahtech/gorest/log"
	"github.com/elsyahtech/gorest/monitoring"
	"github.com/elsyahtech/gorest/server"
)

const hardeningSubprocessEnvironment = "GOREST_HARDENING_E2E_SUBPROCESS"

func TestHardeningEndpointE2E(t *testing.T) {
	t.Run("disabled by default", func(t *testing.T) {
		baseURL := startHardeningTestServer(t, nil)
		assertHardeningHTTPStatus(t, baseURL, "/gorest/hardening", nil, http.StatusNotFound)
	})

	t.Run("forwarded headers do not bypass the peer allowlist", func(t *testing.T) {
		baseURL := startHardeningTestServer(t, map[string]string{
			"GOREST_HARDENING_ENABLE":             "true",
			"GOREST_HARDENING_ALLOWED_IPS":        "192.0.2.10",
			"GOREST_DISABLE_MONITORING_ENDPOINTS": "true",
		})
		assertHardeningHTTPStatus(t, baseURL, "/gorest/hardening", map[string]string{
			"X-Forwarded-For": "192.0.2.10",
		}, http.StatusForbidden)
	})

	t.Run("allowlisted peer receives redacted no-store reports", func(t *testing.T) {
		const signingSecret = "hardening-e2e-signing-secret"
		baseURL := startHardeningTestServer(t, map[string]string{
			"GOREST_HARDENING_ENABLE":             "true",
			"GOREST_HARDENING_ALLOWED_IPS":        "127.0.0.1/32",
			"GOREST_DISABLE_MONITORING_ENDPOINTS": "true",
			"GOREST_FILTER_SIGNATURE_KEY":         signingSecret,
		})

		jsonResponse := getHardeningResponse(t, baseURL, "/gorest/hardening", nil)
		defer jsonResponse.Body.Close()
		if jsonResponse.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(jsonResponse.Body)
			t.Fatalf("JSON report status = %d, want %d; body: %s",
				jsonResponse.StatusCode, http.StatusOK, body)
		}
		if contentType := jsonResponse.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "application/json") {
			t.Errorf("JSON Content-Type = %q, want application/json", contentType)
		}
		assertHardeningNoStoreHeaders(t, jsonResponse)

		body, err := io.ReadAll(jsonResponse.Body)
		if err != nil {
			t.Fatalf("read JSON hardening report: %v", err)
		}
		if strings.Contains(string(body), signingSecret) {
			t.Fatal("JSON hardening report exposed the configured signing secret")
		}

		var report struct {
			EvidenceMode string `json:"evidence_mode"`
			Inventory    []struct {
				Path string `json:"field_path"`
			} `json:"inventory"`
		}
		if err := json.Unmarshal(body, &report); err != nil {
			t.Fatalf("decode JSON hardening report: %v", err)
		}
		if report.EvidenceMode != "effective_runtime_config" {
			t.Errorf("evidence_mode = %q, want effective_runtime_config", report.EvidenceMode)
		}
		if len(report.Inventory) == 0 {
			t.Fatal("hardening report contains no configuration inventory")
		}

		htmlResponse := getHardeningResponse(t, baseURL, "/gorest/hardening?format=html", nil)
		defer htmlResponse.Body.Close()
		if htmlResponse.StatusCode != http.StatusOK {
			t.Fatalf("HTML report status = %d, want %d", htmlResponse.StatusCode, http.StatusOK)
		}
		if contentType := htmlResponse.Header.Get("Content-Type"); !strings.HasPrefix(contentType, "text/html") {
			t.Errorf("HTML Content-Type = %q, want text/html", contentType)
		}
		assertHardeningNoStoreHeaders(t, htmlResponse)
		htmlBody, err := io.ReadAll(htmlResponse.Body)
		if err != nil {
			t.Fatalf("read HTML hardening report: %v", err)
		}
		if strings.Contains(string(htmlBody), signingSecret) {
			t.Fatal("HTML hardening report exposed the configured signing secret")
		}

		assertHardeningHTTPStatus(t, baseURL, "/gorest/hardening?format=xml", nil, http.StatusBadRequest)
	})
}

func TestHardeningSubprocess(t *testing.T) {
	if os.Getenv(hardeningSubprocessEnvironment) != "1" {
		return
	}

	hardeningConfig, err := appconfig.LoadHardening()
	if err != nil {
		t.Fatalf("load hardening configuration: %v", err)
	}
	monitoringConfig, err := appconfig.LoadMonitoring()
	if err != nil {
		t.Fatalf("load monitoring configuration: %v", err)
	}
	filterConfig, err := appconfig.LoadFilters()
	if err != nil {
		t.Fatalf("load filter configuration: %v", err)
	}
	serverConfig, err := appconfig.LoadServer()
	if err != nil {
		t.Fatalf("load server configuration: %v", err)
	}

	loggerConfig := gorestlog.ConfigDefault
	loggerConfig.LogDir = filepath.Join(".", "logs")
	app := gorest.New(gorest.Config{
		AppName:         "gorest-api-starter-hardening-e2e",
		Environment:     gorest.Testing,
		Timezone:        "UTC",
		CheckForUpdates: gorest.DisableAutoUpdate,
	})
	app.Use(
		gorestlog.Run(loggerConfig),
		filter.Run(filterConfig),
		server.Run(serverConfig),
		monitoring.Run(monitoringConfig),
		hardening.Run(hardeningConfig),
	)
	app.Start()
}

func startHardeningTestServer(t *testing.T, overrides map[string]string) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("reserve local port for hardening E2E server: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatalf("release local port for hardening E2E server: %v", err)
	}

	values := map[string]string{
		"APP_ENV":                      gorest.Testing,
		"GOREST_SERVER_HOST":           "127.0.0.1",
		"GOREST_SERVER_PORT":           strconv.Itoa(port),
		hardeningSubprocessEnvironment: "1",
	}
	for key, value := range overrides {
		values[key] = value
	}

	cmd := exec.Command(os.Args[0], "-test.run=^TestHardeningSubprocess$")
	cmd.Dir = t.TempDir()
	cmd.Env = sanitizedHardeningTestEnvironment(values)
	var stderr synchronizedBuffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		t.Fatalf("start hardening E2E application: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process == nil {
			return
		}
		_ = cmd.Process.Signal(os.Interrupt)
		done := make(chan error, 1)
		go func() { done <- cmd.Wait() }()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			_ = cmd.Process.Kill()
			<-done
		}
	})

	baseURL := fmt.Sprintf("http://127.0.0.1:%d", port)
	deadline := time.Now().Add(10 * time.Second)
	client := &http.Client{Timeout: time.Second}
	for time.Now().Before(deadline) {
		response, requestErr := client.Get(baseURL + "/gorest/hardening")
		if requestErr == nil {
			response.Body.Close()
			return baseURL
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("hardening E2E application did not start; stderr: %s", stderr.String())
	return ""
}

func sanitizedHardeningTestEnvironment(values map[string]string) []string {
	env := make([]string, 0, len(os.Environ())+len(values))
	for _, entry := range os.Environ() {
		key, _, _ := strings.Cut(entry, "=")
		if key == "APP_ENV" || strings.HasPrefix(key, "GOREST_") {
			continue
		}
		env = append(env, entry)
	}
	for key, value := range values {
		env = append(env, key+"="+value)
	}
	return env
}

type synchronizedBuffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (buffer *synchronizedBuffer) Write(value []byte) (int, error) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	return buffer.buffer.Write(value)
}

func (buffer *synchronizedBuffer) String() string {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	return buffer.buffer.String()
}

func getHardeningResponse(t *testing.T, baseURL, path string, headers map[string]string) *http.Response {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, baseURL+path, nil)
	if err != nil {
		t.Fatalf("create hardening request %s: %v", path, err)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(request)
	if err != nil {
		t.Fatalf("request hardening endpoint %s: %v", path, err)
	}
	return response
}

func assertHardeningHTTPStatus(
	t *testing.T,
	baseURL string,
	path string,
	headers map[string]string,
	wantStatus int,
) {
	t.Helper()

	response := getHardeningResponse(t, baseURL, path, headers)
	defer response.Body.Close()
	if response.StatusCode != wantStatus {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("hardening endpoint %s status = %d, want %d; body: %s",
			path, response.StatusCode, wantStatus, body)
	}
}

func assertHardeningNoStoreHeaders(t *testing.T, response *http.Response) {
	t.Helper()
	if got := response.Header.Get("Cache-Control"); got != "no-store" {
		t.Errorf("Cache-Control = %q, want no-store", got)
	}
	if got := response.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Errorf("X-Content-Type-Options = %q, want nosniff", got)
	}
}
