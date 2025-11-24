package mods

import (
	"errors"
	"fmt"
	"strings"
)

const (
	// SourceCurseForge identifies CurseForge-backed downloads.
	SourceCurseForge = "curseforge"
)

var (
	ErrUnknownModSource = errors.New("unknown mod source")
)

// CurseForgeSpec represents the identifiers needed to download a CurseForge file.
type CurseForgeSpec struct {
	ProjectID   int    `json:"projectId"`
	FileID      int    `json:"fileId,omitempty"`
	GameVersion string `json:"gameVersion,omitempty"`
	Loader      string `json:"loader,omitempty"`
}

// Spec describes a mod reference that mcrun can resolve before launching a server.
type Spec struct {
	Source     string          `json:"source"`
	Name       string          `json:"name,omitempty"`
	CurseForge *CurseForgeSpec `json:"curseforge,omitempty"`
}

// Normalized returns a copy of the spec with normalized source naming.
func (s Spec) Normalized() Spec {
	s.Source = strings.ToLower(strings.TrimSpace(s.Source))
	return s
}

// SourceKey exposes the normalized source string.
func (s Spec) SourceKey() string {
	return strings.ToLower(strings.TrimSpace(s.Source))
}

// Validate ensures the spec contains the data required to fetch the mod.
func (s Spec) Validate() error {
	switch s.SourceKey() {
	case SourceCurseForge:
		if s.CurseForge == nil {
			return errors.New("curseforge mod missing configuration")
		}
		if s.CurseForge.ProjectID <= 0 {
			return fmt.Errorf("curseforge mod requires a positive projectId, got %d", s.CurseForge.ProjectID)
		}
	case "":
		return errors.New("mod source is required")
	default:
		return fmt.Errorf("%w: %s", ErrUnknownModSource, s.Source)
	}

	return nil
}
