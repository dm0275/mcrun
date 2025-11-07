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
			entries = append(entries, fmt.Sprintf("%d:%d", spec.CurseForge.ProjectID, spec.CurseForge.FileID))
		}
	}

	return strings.Join(entries, ",")
}

func (f *curseForgeModFlag) Set(value string) error {
	parts := strings.Split(value, ":")
	if len(parts) != 2 {
		return fmt.Errorf("expected format projectID:fileID, got %q", value)
	}

	projectID, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return fmt.Errorf("invalid projectID %q: %w", parts[0], err)
	}

	fileID, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return fmt.Errorf("invalid fileID %q: %w", parts[1], err)
	}

	spec := mods.Spec{
		Source: mods.SourceCurseForge,
		CurseForge: &mods.CurseForgeSpec{
			ProjectID: projectID,
			FileID:    fileID,
		},
	}

	if err := spec.Validate(); err != nil {
		return err
	}

	*f.target = append(*f.target, spec)
	return nil
}

func (f *curseForgeModFlag) Type() string {
	return "curseforge-mod"
}
