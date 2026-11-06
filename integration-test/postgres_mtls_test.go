//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/elsyahtech/gorest-api-starter/app/config"
	"github.com/elsyahtech/gorest/database"
)

func TestPostgreSQLMutualTLSRequiresClientCertificate(t *testing.T) {
	testConfig, err := loadPostgreSQLTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	applyPostgreSQLTestConfig(t, testConfig)

	dbConfig, err := config.DatabaseConfig()
	if err != nil {
		t.Fatalf("load PostgreSQL mTLS configuration: %v", err)
	}
	if dbConfig.Disable || dbConfig.Driver != database.POSTGRES ||
		dbConfig.TLS != testConfig.tls || dbConfig.SSLCertPath != testConfig.certPath ||
		dbConfig.SSLKeyPath != testConfig.keyPath {
		t.Fatalf("database configuration is not PostgreSQL mTLS")
	}
	dbConfig.Migration = false
	dbConfig.Seeder = false
	dbConfig.Timeout = 15 * time.Second

	withoutClientCertificate := *dbConfig
	withoutClientCertificate.SSLCertPath = ""
	withoutClientCertificate.SSLKeyPath = ""
	rejectionContext, cancelRejection := context.WithTimeout(context.Background(), dbConfig.Timeout)
	rejectedDB, message, rejectionErr := database.New(rejectionContext, &withoutClientCertificate, "UTC")
	cancelRejection()
	if rejectionErr == nil {
		if rejectedDB != nil {
			if err := rejectedDB.SQL.Close(); err != nil {
				t.Errorf("close PostgreSQL connection accepted without a client certificate: %v", err)
			}
		}
		t.Fatalf("PostgreSQL server accepted a connection without a client certificate (%s); mTLS is not enforced", message)
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbConfig.Timeout)
	defer cancel()
	db, message, err := database.New(ctx, dbConfig, "UTC")
	if err != nil {
		t.Fatalf("connect to configured PostgreSQL server with mTLS (%s): %v", message, err)
	}
	t.Cleanup(func() {
		if err := db.SQL.Close(); err != nil {
			t.Errorf("close PostgreSQL mTLS integration connection: %v", err)
		}
	})

	var encrypted bool
	if err := db.SQL.QueryRowContext(ctx,
		"SELECT ssl FROM pg_stat_ssl WHERE pid = pg_backend_pid()",
	).Scan(&encrypted); err != nil {
		t.Fatalf("query PostgreSQL mTLS session status: %v", err)
	}
	if !encrypted {
		t.Fatal("PostgreSQL mTLS connection did not negotiate TLS")
	}
}
