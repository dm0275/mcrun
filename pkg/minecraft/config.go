package minecraft

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dm0275/mcrun/pkg/mods"
)

type MinecraftConfig struct {
	WorldName         string
	WorldDir          string
	ModsDir           string
	ModCacheDir       string
	RootDir           string
	Version           string
	Port              string
	MaxMemory         string
	MinMemory         string
	Image             string
	Seed              string
	GameMode          string
	EnableCmdBlock    bool
	McRunDir          string
	ServerConfigFile  string
	LocalServerConfig bool
	MountDirs         []string
	Mods              []mods.Spec
	ExistingMods      []mods.Spec
	CurseForgeAPIKey  string
	ServerType        string
	dockerComposeFile string
}

func NewMinecraftConfig() *MinecraftConfig {
	return &MinecraftConfig{
		Version:           "1.19.3",
		Port:              "25565",
		MaxMemory:         "3G",
		MinMemory:         "3G",
		Image:             "dm0275/minecraft-server",
		GameMode:          "0",
		EnableCmdBlock:    false,
		LocalServerConfig: true,
		CurseForgeAPIKey:  os.Getenv("CURSEFORGE_API_KEY"),
		ServerType:        "vanilla",
	}
}

func NewMinecraftForgeConfig() *MinecraftConfig {
	return &MinecraftConfig{
		Version:           "forge-1.20.1",
		Port:              "25565",
		MaxMemory:         "3G",
		MinMemory:         "3G",
		Image:             "dm0275/minecraft-server",
		GameMode:          "0",
		EnableCmdBlock:    true,
		LocalServerConfig: true,
		CurseForgeAPIKey:  os.Getenv("CURSEFORGE_API_KEY"),
		ServerType:        "forge",
	}
}

func NewMinecraftFabricConfig() *MinecraftConfig {
	return &MinecraftConfig{
		Version:           "fabric-1.20.1",
		Port:              "25565",
		MaxMemory:         "3G",
		MinMemory:         "3G",
		Image:             "dm0275/minecraft-server",
		GameMode:          "0",
		EnableCmdBlock:    true,
		LocalServerConfig: true,
		CurseForgeAPIKey:  os.Getenv("CURSEFORGE_API_KEY"),
		ServerType:        "fabric",
	}
}

func McRunHomeDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	defaultHomeDir := filepath.Join(homeDir, ".mcrun")

	var mcRunDir string
	if os.Getenv("MCRUN_DIR") != "" {
		mcRunDir = os.Getenv("MCRUN_DIR")
	} else {
		mcRunDir = defaultHomeDir
	}

	return mcRunDir, nil
}

func SetupDirectories(mcconfig *MinecraftConfig) error {
	mcrunDir, err := McRunHomeDir()
	if err != nil {
		return err
	}

	worldDir := fmt.Sprintf("%s/%s/world", mcrunDir, mcconfig.WorldName)
	modsDir := fmt.Sprintf("%s/%s/mods", mcrunDir, mcconfig.WorldName)

	err = os.MkdirAll(worldDir, 0o755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(modsDir, 0o755)
	if err != nil {
		return err
	}

	mcconfig.McRunDir = mcrunDir
	mcconfig.WorldDir = worldDir
	mcconfig.ModsDir = modsDir
	mcconfig.ModCacheDir = filepath.Join(mcrunDir, "cache", "mods")
	mcconfig.RootDir = fmt.Sprintf("%s/%s", mcrunDir, mcconfig.WorldName)

	if mcconfig.LocalServerConfig {
		serverConfig := fmt.Sprintf("%s/%s/server.properties", mcrunDir, mcconfig.WorldName)
		mcconfig.ServerConfigFile = serverConfig
	}

	return nil
}
