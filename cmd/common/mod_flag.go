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

	parts := strings.SplitN(value, "@", 2)
	projectStr := strings.TrimSpace(parts[0])
	gameVersion := ""
	if len(parts) == 2 {
		gameVersion = strings.TrimSpace(parts[1])
	}

	projectID, err := strconv.Atoi(projectStr)
	if err != nil {
		return fmt.Errorf("invalid projectID %q: %w", projectStr, err)
	}

	spec := mods.Spec{
		Source: mods.SourceCurseForge,
		CurseForge: &mods.CurseForgeSpec{
			ProjectID:   projectID,
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
	if spec.CurseForge.GameVersion != "" {
		builder.WriteString("@")
		builder.WriteString(spec.CurseForge.GameVersion)
	}
	return builder.String()
}

func (f *curseForgeModFlag) Type() string {
	return "curseforge-mod"
}
