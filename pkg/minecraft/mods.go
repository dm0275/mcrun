package minecraft

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dm0275/mcrun/pkg/mods"
	"github.com/dm0275/mcrun/pkg/mods/curseforge"
)

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
				if apiKey == "" {
					return fmt.Errorf("curseforge mods requested but no API key supplied (set CURSEFORGE_API_KEY or --curseforge-api-key)")
				}

				client, err := curseforge.NewClient(apiKey)
				if err != nil {
					return err
				}
				cfClient = client
			}

			path, err := cfClient.DownloadMod(ctx, spec.CurseForge.ProjectID, spec.CurseForge.FileID, mcconfig.ModsDir)
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
