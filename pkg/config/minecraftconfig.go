package config

import (
	"fmt"
	"os"
)

type MinecraftConfig struct {
	WorldName string
	WorldDir  string
	ModsDir   string
	Version   string
	Port      string
	MaxMemory string
	MinMemory string
	Image     string
}

func NewMinecraftConfig() *MinecraftConfig {
	return &MinecraftConfig{
		Version:   "1.18.2",
		Port:      "25565",
		MaxMemory: "3G",
		MinMemory: "3G",
		Image:     "dm0275/minecraft-server",
	}
}

func SetupDirectories(config *MinecraftConfig) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	worldDir := fmt.Sprintf("%s/data_forge/%s/world", cwd, config.WorldName)
	modsDir := fmt.Sprintf("%s/data_forge/%s/mods", cwd, config.WorldName)

	err = os.MkdirAll(worldDir, 0o755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(modsDir, 0o755)
	if err != nil {
		return err
	}

	config.WorldDir = worldDir
	config.ModsDir = modsDir

	return nil
}
