package minecraft

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
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
	Seed      string
	McRunDir  string
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

func GenerateComposeFile(mcconfig *MinecraftConfig) error {
	compose := `version: '3.8'

services:
  {{ .WorldName }}:
    tty: true
    stdin_open: true
    image: "{{ .Image }}:{{ .Version }}"
    container_name: "{{ .WorldName }}-minecraft"
    ports:
      - "{{ .Port }}:25565"
    volumes:
      - {{ .WorldDir }}:/opt/minecraft/world
      - {{ .ModsDir }}:/opt/minecraft/mods
    environment:
      - GAMEMODE
      - MAX_PLAYERS
      - DIFFICULTY
      - MOTD
      - ENABLE_CMD_BLOCK
      - MAX_TICK_TIME
      - GENERATOR_SETTINGS
      - ALLOW_NETHER
      - FORCE_GAMEMODE
      - ENABLE_QUERY
      - PLAYER_IDLE_TIMEOUT
      - SPAWN_MONSTERS
      - OP_PERMISSION_LEVEL
      - PVP
      - LEVEL_TYPE
      - HARDCORE
      - NETWORK_COMPRESSION_THRESHOLD
      - RESOURCE_PACK_SHA1
      - MAX_WORLD_SIZE
      - LEVEL_SEED
volumes:
  world: {}
  data: {}
  mods: {}
`

	data := &bytes.Buffer{}
	tmpl := template.Must(template.New("").Parse(compose))
	if err := tmpl.Execute(data, mcconfig); err != nil {
		return err
	}

	err := os.WriteFile(fmt.Sprintf("%s/docker-compose-%s.yaml", mcconfig.McRunDir, mcconfig.WorldName), data.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
