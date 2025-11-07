package common

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dm0275/mcrun/pkg/mods"
)

type curseForgeModFlag struct {
	target *[]mods.Spec
}

func newCurseForgeModFlag(target *[]mods.Spec) *curseForgeModFlag {
	return &curseForgeModFlag{target: target}
}

func (f *curseForgeModFlag) String() string {
	if f.target == nil {
		return ""
	}

	var entries []string
	for _, spec := range *f.target {
		if spec.SourceKey() == mods.SourceCurseForge && spec.CurseForge != nil {
			entries = append(entries, formatCurseForgeSpec(spec))
		}
	}

	return strings.Join(entries, ",")
}

func (f *curseForgeModFlag) Set(value string) error {
	value = strings.TrimSpace(value)
	if value == "" {
		return fmt.Errorf("curseforge mod cannot be empty")
	}

	var (
		gameVersion string
		mainPart    string
	)

	if strings.Contains(value, "@") {
		parts := strings.SplitN(value, "@", 2)
		mainPart = parts[0]
		gameVersion = strings.TrimSpace(parts[1])
	} else {
		mainPart = value
	}

	idParts := strings.Split(mainPart, ":")
	if len(idParts) == 0 {
		return fmt.Errorf("invalid mod specification %q", value)
	}

	projectID, err := strconv.Atoi(strings.TrimSpace(idParts[0]))
	if err != nil {
		return fmt.Errorf("invalid projectID %q: %w", idParts[0], err)
	}

	var fileID int
	if len(idParts) > 1 {
		fileID, err = strconv.Atoi(strings.TrimSpace(idParts[1]))
		if err != nil {
			return fmt.Errorf("invalid fileID %q: %w", idParts[1], err)
		}
	}

	if len(idParts) > 2 {
		return fmt.Errorf("expected format projectID[:fileID][@gameVersion], got %q", value)
	}

	spec := mods.Spec{
		Source: mods.SourceCurseForge,
		CurseForge: &mods.CurseForgeSpec{
			ProjectID:   projectID,
			FileID:      fileID,
			GameVersion: gameVersion,
		},
	}

	if err := spec.Validate(); err != nil {
		return err
	}

	*f.target = append(*f.target, spec)
	return nil
}

func formatCurseForgeSpec(spec mods.Spec) string {
	if spec.CurseForge == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(strconv.Itoa(spec.CurseForge.ProjectID))
	if spec.CurseForge.FileID > 0 {
		builder.WriteString(":")
		builder.WriteString(strconv.Itoa(spec.CurseForge.FileID))
	}
	if spec.CurseForge.GameVersion != "" {
		builder.WriteString("@")
		builder.WriteString(spec.CurseForge.GameVersion)
	}
	return builder.String()
}

func (f *curseForgeModFlag) Type() string {
	return "curseforge-mod"
}
