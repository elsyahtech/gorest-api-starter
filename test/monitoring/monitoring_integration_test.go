package monitoring_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/elsyahtech/gorest/database"
	"github.com/elsyahtech/gorest/monitoring"
	"github.com/elsyahtech/gorest/server"
	"github.com/gofiber/fiber/v3"
)

func TestMonitoringEndpointsRespectAllowlistAndDisableSetting(t *testing.T) {
	dbConfig := database.Config{
		Driver:            database.SQLITE,
		Name:              "monitoring_integration",
		DatabaseDirectory: t.TempDir(),
	}
	db, _, err := database.New(context.Background(), &dbConfig, "UTC")
	if err != nil {
		t.Fatalf("initialize SQLite database for monitoring endpoints: %v", err)
	}
	t.Cleanup(func() {
		if err := db.SQL.Close(); err != nil {
			t.Errorf("close SQLite monitoring database: %v", err)
		}
	})
	if _, err := db.SQL.Exec(`CREATE TABLE monitoring_sample (value TEXT)`); err != nil {
		t.Fatalf("create SQLite monitoring fixture: %v", err)
	}
	if _, err := db.SQL.Exec(`INSERT INTO monitoring_sample VALUES ('sample')`); err != nil {
		t.Fatalf("insert SQLite monitoring fixture: %v", err)
	}

	t.Run("empty allowlist denies endpoints, including forwarded IP", func(t *testing.T) {
		httpServer := newMonitoringHTTPServer(t, dbConfig, db, monitoring.ConfigDefault)

		for _, path := range operationalEndpointPaths {
			request := httptest.NewRequest(http.MethodGet, path, nil)
			request.Header.Set("X-Forwarded-For", "127.0.0.1")
			assertMonitoringResponse(t, httpServer, request, http.StatusForbidden, "forbidden")
		}
	})

	t.Run("allowlisted peer can access monitoring endpoints", func(t *testing.T) {
		config := monitoring.Config{
			HTTPMetricsAllowedIPs:     []string{"0.0.0.0/0", "::/0"},
			EnablePerRouteHTTPMetrics: true,
		}
		httpServer := newMonitoringHTTPServer(t, dbConfig, db, config)

		for _, path := range operationalEndpointPaths {
			assertMonitoringResponse(t, httpServer, httptest.NewRequest(http.MethodGet, path, nil), http.StatusOK, "")
		}
		assertPerRouteHTTPMetricsIncluded(t, httpServer)
	})

	t.Run("disabled configuration does not register endpoints", func(t *testing.T) {
		config := monitoring.Config{DisableMonitoringEndpoints: true}
		httpServer := newMonitoringHTTPServer(t, dbConfig, db, config)

		for _, path := range operationalEndpointPaths {
			assertMonitoringResponse(t, httpServer, httptest.NewRequest(http.MethodGet, path, nil), http.StatusNotFound, "")
		}
	})
}

func TestMonitoringRejectsInvalidAllowlistDuringInitialization(t *testing.T) {
	_, err := monitoring.New(monitoring.Config{
		HTTPMetricsAllowedIPs: []string{"not-an-ip"},
	})
	if err == nil {
		t.Fatal("monitoring.New() succeeded with an invalid IP allowlist entry")
	}
}

var operationalEndpointPaths = []string{
	"/gorest/http",
	"/gorest/storage",
	"/gorest/database",
}

func newMonitoringHTTPServer(
	t *testing.T,
	dbConfig database.Config,
	db *database.Database,
	config monitoring.Config,
) *server.Server {
	t.Helper()

	service, err := monitoring.New(config)
	if err != nil {
		t.Fatalf("initialize monitoring service: %v", err)
	}

	var middleware []fiber.Handler
	if service.EndpointsEnabled() {
		middleware = append(middleware, service.Middleware)
	}
	httpServer, _, err := server.NewWithMiddleware(nil, middleware...)
	if err != nil {
		t.Fatalf("create monitoring HTTP server: %v", err)
	}
	if !service.EndpointsEnabled() {
		return httpServer
	}

	routes := httpServer.Group("/gorest", func(ctx *server.Context) error {
		if !service.AllowsPeer(ctx.RequestCtx().RemoteIP()) {
			return ctx.Status(http.StatusForbidden).JSON(map[string]string{"status": "forbidden"})
		}
		return ctx.Next()
	})
	monitoring.Register(routes, dbConfig, db)
	service.RegisterHTTPMetrics(routes)

	return httpServer
}

func assertMonitoringResponse(
	t *testing.T,
	httpServer *server.Server,
	request *http.Request,
	wantStatus int,
	wantBodyStatus string,
) {
	t.Helper()

	response, err := httpServer.SRV.Test(request)
	if err != nil {
		t.Fatalf("request monitoring endpoint %s: %v", request.URL.Path, err)
	}
	defer response.Body.Close()

	if response.StatusCode != wantStatus {
		body, _ := io.ReadAll(response.Body)
		t.Fatalf("monitoring endpoint %s status = %d, want %d; body: %s",
			request.URL.Path, response.StatusCode, wantStatus, body)
	}
	if wantStatus == http.StatusNotFound {
		return
	}
	var body struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		t.Fatalf("decode monitoring endpoint %s response: %v", request.URL.Path, err)
	}
	if body.Status == "" {
		t.Errorf("monitoring endpoint %s response is missing its status field", request.URL.Path)
	} else if wantBodyStatus != "" && body.Status != wantBodyStatus {
		t.Errorf("monitoring endpoint %s status field = %q, want %q",
			request.URL.Path, body.Status, wantBodyStatus)
	}
}

func assertPerRouteHTTPMetricsIncluded(t *testing.T, httpServer *server.Server) {
	t.Helper()

	response, err := httpServer.SRV.Test(httptest.NewRequest(http.MethodGet, "/gorest/http", nil))
	if err != nil {
		t.Fatalf("request per-route HTTP metrics: %v", err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("per-route HTTP metrics status = %d, want %d", response.StatusCode, http.StatusOK)
	}

	var result struct {
		Metrics struct {
			Routes []struct {
				Method string `json:"method"`
				Route  string `json:"route"`
			} `json:"routes"`
		} `json:"metrics"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatalf("decode per-route HTTP metrics: %v", err)
	}
	if len(result.Metrics.Routes) == 0 {
		t.Fatal("per-route HTTP metrics response has no route breakdown")
	}
}
