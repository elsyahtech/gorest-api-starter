package mysqltls

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"

	"github.com/go-sql-driver/mysql"
)

var profilesMu sync.Mutex
var registeredProfiles = make(map[string]struct{})

func RegisterProfile(name, serverName, caPath, certPath, keyPath string) error {
	if name == "" || serverName == "" || caPath == "" || certPath == "" || keyPath == "" {
		return fmt.Errorf("MySQL TLS profile requires a profile name, server name, CA, client certificate, and client key")
	}

	profilesMu.Lock()
	defer profilesMu.Unlock()
	if _, ok := registeredProfiles[name]; ok {
		return nil
	}

	caPEM, err := os.ReadFile(caPath)
	if err != nil {
		return fmt.Errorf("read MySQL server CA certificate: %w", err)
	}
	rootCAs := x509.NewCertPool()
	if !rootCAs.AppendCertsFromPEM(caPEM) {
		return fmt.Errorf("parse MySQL server CA certificate: no certificates found")
	}

	clientCertificate, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return fmt.Errorf("load MySQL client certificate and key: %w", err)
	}

	config := &tls.Config{
		MinVersion:   tls.VersionTLS12,
		ServerName:   serverName,
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientCertificate},
	}
	if err := mysql.RegisterTLSConfig(name, config); err != nil {
		return fmt.Errorf("register MySQL TLS profile: %w", err)
	}
	registeredProfiles[name] = struct{}{}

	return nil
}
