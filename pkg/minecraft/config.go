package minecraft

import (
	"fmt"
	"os"
	"path/filepath"
)

type MinecraftConfig struct {
	WorldName         string
	WorldDir          string
	ModsDir           string
	Version           string
	Port              string
	MaxMemory         string
	MinMemory         string
	Image             string
	Seed              string
	GameMode          string
	EnableCmdBlock    bool
	McRunDir          string
	dockerComposeFile string
}

func NewMinecraftConfig() *MinecraftConfig {
	return &MinecraftConfig{
		Version:        "1.19.3",
		Port:           "25565",
		MaxMemory:      "3G",
		MinMemory:      "3G",
		Image:          "dm0275/minecraft-server",
		GameMode:       "0",
		EnableCmdBlock: false,
	}
}

func NewMinecraftForgeConfig() *MinecraftConfig {
	return &MinecraftConfig{
		Version:        "forge-1.20.1",
		Port:           "25565",
		MaxMemory:      "3G",
		MinMemory:      "3G",
		Image:          "dm0275/minecraft-server",
		GameMode:       "0",
		EnableCmdBlock: true,
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

	return nil
}
