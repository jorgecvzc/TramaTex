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
