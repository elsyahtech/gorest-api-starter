//go:build integration

package integration_test

import (
	"context"
	"testing"
	"time"

	"github.com/elsyahtech/gorest/database"
)

func TestScyllaMutualTLSRequiresClientCertificate(t *testing.T) {
	testConfig, err := loadScyllaTestConfig()
	if err != nil {
		t.Fatal(err)
	}
	if testConfig.certPath == "" {
		t.Skip("configure ScyllaDB TLS client certificate and key to run the mTLS integration test")
	}
	if testConfig.tls != database.TLSEnabled {
		t.Fatal("ScyllaDB mTLS integration requires GOREST_TEST_SCYLLA_TLS=enable")
	}

	connect := func(withClientCertificate bool) error {
		t.Helper()
		testCfg := database.ConfigDefault
		testCfg.Driver = database.SCYLLA
		testCfg.Host = testConfig.host
		testCfg.Port = testConfig.port
		testCfg.Name = testConfig.keyspace
		testCfg.Username = testConfig.username
		testCfg.Password = testConfig.password
		testCfg.TLS = database.TLSEnabled
		testCfg.SSLCAPath = testConfig.caPath
		if withClientCertificate {
			testCfg.SSLCertPath = testConfig.certPath
			testCfg.SSLKeyPath = testConfig.keyPath
		}
		testCfg.Timeout = 10 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), 12*time.Second)
		defer cancel()
		db, _, err := database.New(ctx, &testCfg, "UTC")
		if err != nil {
			return err
		}
		defer func() {
			closeScyllaDatabase(t, db, &testCfg)
		}()
		_, err = db.Ping(ctx, &testCfg)
		return err
	}

	if err := connect(false); err == nil {
		t.Fatal("ScyllaDB accepted a connection without the client certificate; configure the endpoint to require mTLS")
	}
	if err := connect(true); err != nil {
		t.Fatalf("connect to ScyllaDB with the configured client certificate: %v", err)
	}
}
