package mods

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// Hash computes a stable fingerprint for the mod spec, ignoring display-only fields.
func (s Spec) HashValue() string {
	source := strings.ToLower(strings.TrimSpace(s.Source))
	projectID := 0
	fileID := 0
	gameVersion := ""
	loader := ""
	url := strings.TrimSpace(s.URL)

	if s.CurseForge != nil {
		projectID = s.CurseForge.ProjectID
		fileID = s.CurseForge.FileID
		gameVersion = strings.ToLower(strings.TrimSpace(s.CurseForge.GameVersion))
		loader = strings.ToLower(strings.TrimSpace(s.CurseForge.Loader))
	}

	raw := fmt.Sprintf("%s:%d:%d:%s:%s:%s", source, projectID, fileID, gameVersion, loader, url)
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// Key returns an identifier to match specs across runs.
func (s Spec) Key() string {
	source := strings.ToLower(strings.TrimSpace(s.Source))
	projectID := 0
	fileID := 0
	gameVersion := ""
	loader := ""

	if s.CurseForge != nil {
		projectID = s.CurseForge.ProjectID
		fileID = s.CurseForge.FileID
		gameVersion = strings.ToLower(strings.TrimSpace(s.CurseForge.GameVersion))
		loader = strings.ToLower(strings.TrimSpace(s.CurseForge.Loader))
	}

	return fmt.Sprintf("%s:%d:%d:%s:%s", source, projectID, fileID, gameVersion, loader)
}
