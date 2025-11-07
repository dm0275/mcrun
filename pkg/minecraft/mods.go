package minecraft

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/dm0275/mcrun/pkg/mods"
	"github.com/dm0275/mcrun/pkg/mods/curseforge"
)

var versionRegex = regexp.MustCompile(`^\d+(\.\d+)*$`)

// SyncMods downloads and stages remote mods declared on the config.
func SyncMods(mcconfig *MinecraftConfig) error {
	if len(mcconfig.Mods) == 0 {
		return nil
	}

	var cfClient *curseforge.Client
	ctx := context.Background()

	for _, spec := range mcconfig.Mods {
		spec = spec.Normalized()
		if err := spec.Validate(); err != nil {
			return err
		}

		switch spec.SourceKey() {
		case mods.SourceCurseForge:
			if cfClient == nil {
				apiKey := strings.TrimSpace(mcconfig.CurseForgeAPIKey)
				if apiKey == "" {
					apiKey = strings.TrimSpace(os.Getenv("CURSEFORGE_API_KEY"))
				}
				client, err := curseforge.NewClient(apiKey)
				if err != nil {
					return err
				}
				cfClient = client
			}

			fileID := spec.CurseForge.FileID
			if fileID == 0 {
				gameVersion := spec.CurseForge.GameVersion
				if gameVersion == "" {
					gameVersion = deriveGameVersion(mcconfig.Version)
				}
				if gameVersion == "" {
					return fmt.Errorf("unable to determine game version for CurseForge mod %d; specify mods[].curseforge.gameVersion", spec.CurseForge.ProjectID)
				}

				loaders := candidateLoaders(spec.CurseForge.Loader, mcconfig.ServerType)
				resolvedID, err := cfClient.ResolveFileID(ctx, spec.CurseForge.ProjectID, loaders, gameVersion)
				if err != nil {
					return err
				}
				fileID = resolvedID
			}

			path, err := cfClient.DownloadMod(ctx, spec.CurseForge.ProjectID, fileID, mcconfig.ModsDir)
			if err != nil {
				return err
			}
			log.Printf("Downloaded CurseForge mod to %s", path)
		default:
			return fmt.Errorf("unsupported mod source %q", spec.Source)
		}
	}

	return nil
}

func deriveGameVersion(configVersion string) string {
	configVersion = strings.TrimSpace(configVersion)
	if configVersion == "" {
		return ""
	}

	parts := strings.Split(configVersion, "-")
	candidate := parts[len(parts)-1]

	if versionRegex.MatchString(candidate) {
		return candidate
	}

	if versionRegex.MatchString(configVersion) {
		return configVersion
	}

	return ""
}

func candidateLoaders(explicitLoader, serverType string) []string {
	if strings.TrimSpace(explicitLoader) != "" {
		return []string{explicitLoader}
	}

	switch strings.ToLower(serverType) {
	case "forge":
		return []string{"NeoForge", "Forge"}
	case "fabric":
		return []string{"Fabric"}
	default:
		return nil
	}
}
