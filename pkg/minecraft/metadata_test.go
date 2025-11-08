package minecraft

import (
	"os"
	"testing"
)

func TestSaveAndLoadMetadata(t *testing.T) {
	dir := t.TempDir()
	cfg := &MinecraftConfig{
		WorldName:  "test-world",
		RootDir:    dir,
		Version:    "forge-1.20.1",
		Image:      "dm0275/minecraft-server",
		ServerType: "forge",
	}

	if err := SaveServerMetadata(cfg); err != nil {
		t.Fatalf("SaveServerMetadata() error = %v", err)
	}

	meta, err := LoadServerMetadata(dir)
	if err != nil {
		t.Fatalf("LoadServerMetadata() error = %v", err)
	}

	if meta.WorldName != "test-world" || meta.Type != "forge" {
		t.Fatalf("unexpected metadata: %+v", meta)
	}

	info, err := os.Stat(dir + "/metadata.json")
	if err != nil {
		t.Fatalf("metadata file missing: %v", err)
	}

	if info.Size() == 0 {
		t.Fatalf("metadata file empty")
	}
}
