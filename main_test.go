package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadDotEnvAllowsEnvironmentOnlyStartup(t *testing.T) {
	previousDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("get current directory: %v", err)
	}
	directory := t.TempDir()
	if err := os.Chdir(directory); err != nil {
		t.Fatalf("change to temporary directory: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(previousDir); err != nil {
			t.Errorf("restore current directory: %v", err)
		}
	})

	if err := loadDotEnv(); err != nil {
		t.Fatalf("loadDotEnv() without a .env file = %v, want nil", err)
	}

	if err := os.WriteFile(filepath.Join(directory, ".env"), []byte(`INVALID="unterminated`), 0o600); err != nil {
		t.Fatalf("write invalid .env: %v", err)
	}
	if err := loadDotEnv(); err == nil {
		t.Fatal("loadDotEnv() accepted an invalid .env file")
	}
}
