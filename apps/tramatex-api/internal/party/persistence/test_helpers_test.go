package persistence

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestReadEnvRemote_ReadsFile(t *testing.T) {
	_, currentFile, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", ".."))
	path := filepath.Join(root, ".env.remote")

	original, err := os.ReadFile(path)
	hadOriginal := err == nil

	content := []byte("DB_HOST=remotehost\nDB_PORT=5432\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to write .env.remote: %v", err)
	}
	defer func() {
		if hadOriginal {
			_ = os.WriteFile(path, original, 0644)
		} else {
			_ = os.Remove(path)
		}
	}()

	env, err := readEnvRemote()
	if err != nil {
		t.Fatalf("readEnvRemote should not error: %v", err)
	}
	if env["DB_HOST"] != "remotehost" {
		t.Fatalf("expected DB_HOST to match, got %s", env["DB_HOST"])
	}
	if env["DB_PORT"] != "5432" {
		t.Fatalf("expected DB_PORT to match, got %s", env["DB_PORT"])
	}
}

func TestReadEnvFile_ParsesEntries(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "env-test-*.env")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	path := tmpFile.Name()
	_ = tmpFile.Close()
	t.Cleanup(func() {
		_ = os.Remove(path)
	})

	content := []byte("# comment\nDB_HOST=localhost\nDB_PORT=\"5432\"\nINVALID_LINE\nDB_NAME='tramatex'\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}

	env, err := readEnvFile(path)
	if err != nil {
		t.Fatalf("readEnvFile should not error: %v", err)
	}
	if env["DB_HOST"] != "localhost" {
		t.Fatalf("expected DB_HOST to match, got %s", env["DB_HOST"])
	}
	if env["DB_PORT"] != "5432" {
		t.Fatalf("expected DB_PORT to match, got %s", env["DB_PORT"])
	}
	if env["DB_NAME"] != "tramatex" {
		t.Fatalf("expected DB_NAME to match, got %s", env["DB_NAME"])
	}
}

func TestLoadTestDBConfig_EnvOverrides(t *testing.T) {
	t.Setenv("TRAMATEX_TEST_DB_HOST", "db-host")
	t.Setenv("TRAMATEX_TEST_DB_PORT", "5544")
	t.Setenv("TRAMATEX_TEST_DB_USER", "db-user")
	t.Setenv("TRAMATEX_TEST_DB_PASSWORD", "db-pass")
	t.Setenv("TRAMATEX_TEST_DB_NAME", "db-name")
	t.Setenv("TRAMATEX_TEST_DB_SSLMODE", "require")

	config := loadTestDBConfig()
	if config.Host != "db-host" || config.Port != "5544" || config.User != "db-user" || config.Password != "db-pass" || config.Name != "db-name" || config.SSLMode != "require" {
		t.Fatalf("expected env overrides to be applied")
	}
}

func TestLoadTestDBConfig_EnvLocalOverrides(t *testing.T) {
	t.Setenv("TRAMATEX_TEST_DB_HOST", "")
	t.Setenv("TRAMATEX_TEST_DB_PORT", "")
	t.Setenv("TRAMATEX_TEST_DB_USER", "")
	t.Setenv("TRAMATEX_TEST_DB_PASSWORD", "")
	t.Setenv("TRAMATEX_TEST_DB_NAME", "")
	t.Setenv("TRAMATEX_TEST_DB_SSLMODE", "")

	_, currentFile, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "..", ".."))
	path := filepath.Join(root, ".env.local")

	original, err := os.ReadFile(path)
	hadOriginal := err == nil

	content := []byte("DB_HOST=postgres\nDB_PORT=5432\nDB_USER=local\nDB_PASSWORD=localpass\nDB_NAME=localdb\nDB_SSLMODE=disable\nSSH_HOST=ssh-host\n")
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatalf("failed to write .env.local: %v", err)
	}
	t.Cleanup(func() {
		if hadOriginal {
			_ = os.WriteFile(path, original, 0644)
		} else {
			_ = os.Remove(path)
		}
	})

	config := loadTestDBConfig()
	if config.Host != "ssh-host" {
		t.Fatalf("expected SSH_HOST to override postgres host, got %s", config.Host)
	}
	if config.User != "local" || config.Password != "localpass" || config.Name != "localdb" || config.Port != "5432" || config.SSLMode != "disable" {
		t.Fatalf("expected local env overrides to be applied")
	}
}
