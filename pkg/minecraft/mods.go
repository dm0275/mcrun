package minecraft

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

	for i := range mcconfig.Mods {
		spec := mcconfig.Mods[i].Normalized()
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
			if fileID <= 0 {
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

			if spec.Name == "" {
				if summary, err := cfClient.FetchModSummary(ctx, spec.CurseForge.ProjectID); err == nil {
					spec.Name = summary.Name
				}
			}

			cacheDir := mcconfig.ModCacheDir
			if cacheDir == "" {
				cacheDir = filepath.Join(mcconfig.McRunDir, "cache", "mods")
			}
			if err := os.MkdirAll(cacheDir, 0o755); err != nil {
				return err
			}

			cachePath, cached, err := cfClient.DownloadMod(ctx, spec.CurseForge.ProjectID, fileID, cacheDir)
			if err != nil {
				return err
			}

			if spec.Name == "" {
				spec.Name = inferNameFromPath(cachePath)
			}

			destPath, err := ensureModFromCache(cachePath, mcconfig.ModsDir)
			if err != nil {
				return err
			}

			if cached {
				log.Printf("Using cached CurseForge mod %s (source %s)", destPath, cachePath)
			} else {
				log.Printf("Downloaded CurseForge mod %s (cached at %s)", destPath, cachePath)
			}
		default:
			return fmt.Errorf("unsupported mod source %q", spec.Source)
		}

		mcconfig.Mods[i] = spec
	}

	return nil
}

func ensureModFromCache(cachePath, modsDir string) (string, error) {
	fileName := filepath.Base(cachePath)
	destPath := filepath.Join(modsDir, fileName)
	if _, err := os.Stat(destPath); err == nil {
		return destPath, nil
	}

	if err := os.MkdirAll(modsDir, 0o755); err != nil {
		return "", err
	}

	if err := os.Link(cachePath, destPath); err == nil {
		return destPath, nil
	}

	src, err := os.Open(cachePath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer func() {
		dst.Close()
		if err != nil {
			os.Remove(destPath)
		}
	}()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return destPath, nil
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

func inferNameFromPath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	if ext != "" {
		base = strings.TrimSuffix(base, ext)
	}
	return base
}
