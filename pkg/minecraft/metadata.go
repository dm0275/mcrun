package minecraft

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/dm0275/mcrun/pkg/mods"
)

const metadataFileName = "metadata.json"

// ServerMetadata stores basic information about a provisioned server.
type ServerMetadata struct {
	WorldName string      `json:"worldName"`
	Type      string      `json:"type"`
	Version   string      `json:"version"`
	Image     string      `json:"image"`
	Mods      []mods.Spec `json:"mods,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
}

// SaveServerMetadata writes metadata for the provided configuration, preserving the
// original creation timestamp if a metadata file already exists.
func SaveServerMetadata(cfg *MinecraftConfig) error {
	if cfg.RootDir == "" {
		return fmt.Errorf("server root directory is required to write metadata")
	}

	metaPath := filepath.Join(cfg.RootDir, metadataFileName)
	meta := &ServerMetadata{
		WorldName: cfg.WorldName,
		Type:      cfg.ServerType,
		Version:   cfg.Version,
		Image:     cfg.Image,
		Mods:      cfg.Mods,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if existing, err := LoadServerMetadata(cfg.RootDir); err == nil {
		meta.CreatedAt = existing.CreatedAt
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(metaPath, data, 0o644); err != nil {
		return err
	}

	return nil
}

// LoadServerMetadata loads metadata for the given world directory.
func LoadServerMetadata(rootDir string) (*ServerMetadata, error) {
	metaPath := filepath.Join(rootDir, metadataFileName)
	data, err := os.ReadFile(metaPath)
	if err != nil {
		return nil, err
	}

	var meta ServerMetadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}
