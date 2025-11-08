package minecraft

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListServers(t *testing.T) {
	temp := t.TempDir()
	// prepare directories
	must := func(path string) {
		t.Helper()
		if err := os.MkdirAll(path, 0o755); err != nil {
			t.Fatalf("failed to create dir: %v", err)
		}
	}

	must(filepath.Join(temp, "cache"))
	must(filepath.Join(temp, "world-one"))
	must(filepath.Join(temp, "world-two"))

	compose := filepath.Join(temp, "docker-compose-world-one.yaml")
	if err := os.WriteFile(compose, []byte("version: '3'"), 0o644); err != nil {
		t.Fatalf("failed to write compose file: %v", err)
	}

	t.Setenv("MCRUN_DIR", temp)

	servers, err := ListServers()
	if err != nil {
		t.Fatalf("ListServers() error = %v", err)
	}

	if len(servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(servers))
	}

	if servers[0].WorldName != "world-one" || !servers[0].HasCompose {
		t.Fatalf("unexpected first server: %+v", servers[0])
	}

	if servers[1].WorldName != "world-two" || servers[1].HasCompose {
		t.Fatalf("unexpected second server: %+v", servers[1])
	}
}
