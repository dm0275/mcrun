package api

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dm0275/mcrun/pkg/minecraft"
	"github.com/dm0275/mcrun/pkg/mods"
)

var errWorldNameRequired = errors.New("worldName is required")

// CreateServerRequest describes the expected JSON payload for provisioning servers.
type CreateServerRequest struct {
	Type              string      `json:"type"`
	WorldName         string      `json:"worldName"`
	Version           string      `json:"version"`
	Port              string      `json:"port"`
	MaxMemory         string      `json:"maxMemory"`
	MinMemory         string      `json:"minMemory"`
	Image             string      `json:"image"`
	Seed              string      `json:"seed"`
	GameMode          string      `json:"gameMode"`
	EnableCmdBlock    *bool       `json:"enableCmdBlock"`
	LocalServerConfig *bool       `json:"localServerConfig"`
	MountDirs         []string    `json:"mountDirs"`
	Mods              []mods.Spec `json:"mods"`
	CurseForgeAPIKey  string      `json:"curseForgeApiKey"`
	EnableRcon        *bool       `json:"enableRcon"`
	RconPort          string      `json:"rconPort"`
	RconPassword      string      `json:"rconPassword"`
}

// TypeOrDefault ensures a stable value in responses.
func (r CreateServerRequest) TypeOrDefault() string {
	if strings.TrimSpace(r.Type) == "" {
		return "vanilla"
	}
	return strings.ToLower(r.Type)
}

// ToMinecraftConfig merges the request payload with the default Minecraft config.
func (r CreateServerRequest) ToMinecraftConfig() (*minecraft.MinecraftConfig, error) {
	worldName := strings.TrimSpace(r.WorldName)
	if worldName == "" {
		return nil, errWorldNameRequired
	}
	if strings.ContainsAny(worldName, `/\`) {
		return nil, fmt.Errorf("worldName %q contains invalid path characters", worldName)
	}

	var cfg *minecraft.MinecraftConfig
	switch strings.ToLower(strings.TrimSpace(r.Type)) {
	case "", "vanilla":
		cfg = minecraft.NewMinecraftConfig()
	case "forge":
		cfg = minecraft.NewMinecraftForgeConfig()
	case "fabric":
		cfg = minecraft.NewMinecraftFabricConfig()
	default:
		return nil, fmt.Errorf("unsupported server type %q", r.Type)
	}

	cfg.WorldName = worldName

	if strings.TrimSpace(r.Version) != "" {
		cfg.Version = r.Version
	}
	if strings.TrimSpace(r.Port) != "" {
		cfg.Port = r.Port
	}
	if strings.TrimSpace(r.MaxMemory) != "" {
		cfg.MaxMemory = r.MaxMemory
	}
	if strings.TrimSpace(r.MinMemory) != "" {
		cfg.MinMemory = r.MinMemory
	}
	if strings.TrimSpace(r.Image) != "" {
		cfg.Image = r.Image
	}
	if strings.TrimSpace(r.Seed) != "" {
		cfg.Seed = r.Seed
	}
	if strings.TrimSpace(r.GameMode) != "" {
		cfg.GameMode = r.GameMode
	}
	if r.MountDirs != nil {
		cfg.MountDirs = r.MountDirs
	}
	if r.EnableCmdBlock != nil {
		cfg.EnableCmdBlock = *r.EnableCmdBlock
	}
	if r.LocalServerConfig != nil {
		cfg.LocalServerConfig = *r.LocalServerConfig
	}
	if r.EnableRcon != nil {
		cfg.EnableRcon = *r.EnableRcon
	}
	if strings.TrimSpace(r.RconPort) != "" {
		cfg.RconPort = r.RconPort
	}
	if strings.TrimSpace(r.RconPassword) != "" {
		cfg.RconPassword = r.RconPassword
	}
	if strings.TrimSpace(r.CurseForgeAPIKey) != "" {
		cfg.CurseForgeAPIKey = strings.TrimSpace(r.CurseForgeAPIKey)
	}
	for _, modSpec := range r.Mods {
		modSpec = modSpec.Normalized()
		if err := modSpec.Validate(); err != nil {
			return nil, err
		}
		cfg.Mods = append(cfg.Mods, modSpec)
	}

	return cfg, nil
}
