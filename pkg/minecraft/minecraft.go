package minecraft

import (
	"bytes"
	"fmt"
	"github.com/dm0275/mcrun/utils"
	"os"
	"text/template"
)

func GetComposeFile(mcconfig *MinecraftConfig) (string, error) {
	mcRunDir, err := McRunHomeDir()
	if err != nil {
		return "", err
	}

	dockerComposeFilePath := fmt.Sprintf("%s/docker-compose-%s.yaml", mcRunDir, mcconfig.WorldName)
	if !utils.FileExists(dockerComposeFilePath) {
		return "", fmt.Errorf("docker compose file %s does not exist", dockerComposeFilePath)
	}

	return dockerComposeFilePath, nil
}

func GenerateComposeFile(mcconfig *MinecraftConfig) error {
	dockerComposeFilePath := fmt.Sprintf("%s/docker-compose-%s.yaml", mcconfig.McRunDir, mcconfig.WorldName)
	dockerComposeData := `version: '3.8'

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
      {{if .GameMode }}
      - GAMEMODE={{ .GameMode }}
      {{end}}
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
      {{if .Seed }}
      - LEVEL_SEED={{ .Seed }}
      {{end}}
volumes:
  world: {}
  data: {}
  mods: {}
`

	data := &bytes.Buffer{}
	tmpl := template.Must(template.New("").Parse(dockerComposeData))
	if err := tmpl.Execute(data, mcconfig); err != nil {
		return err
	}

	err := os.WriteFile(dockerComposeFilePath, data.Bytes(), 0644)
	if err != nil {
		return err
	}

	mcconfig.dockerComposeFile = dockerComposeFilePath

	return nil
}

func StartServer(mcconfig *MinecraftConfig) error {
	execCfg := utils.ExecConfig{
		Command: "docker",
		Args: []string{
			"compose",
			"-f",
			mcconfig.dockerComposeFile,
			"up",
			"-d",
		},
		Environment: map[string]string{
			"JAVA_MIN_MEM": mcconfig.MinMemory,
			"JAVA_MAX_MEM": mcconfig.MaxMemory,
		},
	}

	out, err := utils.Exec(execCfg)
	if err != nil {
		fmt.Println(out)
		return err
	}

	return nil
}

func StopServer(dockerComposeFile string) error {
	execCfg := utils.ExecConfig{
		Command: "docker",
		Args: []string{
			"compose",
			"-f",
			dockerComposeFile,
			"down",
		},
	}

	out, err := utils.Exec(execCfg)
	if err != nil {
		fmt.Println(out)
		return err
	}

	return nil
}
