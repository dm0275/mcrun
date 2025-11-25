package minecraft

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dm0275/mcrun/utils"
	"os"
	"path/filepath"
	"text/template"
)

func GetComposeFile(mcconfig *MinecraftConfig) (string, error) {
	mcRunDir, err := McRunHomeDir()
	if err != nil {
		return "", err
	}

	dockerComposeFilePath := fmt.Sprintf("%s/docker-compose-%s.yaml", mcRunDir, mcconfig.WorldName)
	if !utils.FileExists(dockerComposeFilePath) {
		return "", fmt.Errorf("docker compose file %s does not exist: %w", dockerComposeFilePath, os.ErrNotExist)
	}

	return dockerComposeFilePath, nil
}

func GenerateComposeFile(mcconfig *MinecraftConfig) error {
	dockerComposeFilePath := fmt.Sprintf("%s/docker-compose-%s.yaml", mcconfig.McRunDir, mcconfig.WorldName)
	dockerComposeData := `services:
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
      {{- if .LocalServerConfig }}
      - {{ .ServerConfigFile }}:/opt/minecraft/server.properties
      {{- end }}
      {{- if .MountDirs }}
      {{- range $index, $dir := .MountDirs }}
      - {{ $.RootDir }}/{{ $dir }}:/opt/minecraft/{{ $dir }}
      {{- end }}
      {{- end }}
    environment:
      - JAVA_MIN_MEM={{ .MinMemory }}
      - JAVA_MAX_MEM={{ .MaxMemory }}
      {{- if .LocalServerConfig }}
      - LOAD_PROPERTY_FILE=true
      {{- end }}
      {{- if .GameMode }}
      - GAMEMODE={{ .GameMode }}
      {{- end }}
      {{- if .EnableCmdBlock }}
      - ENABLE_CMD_BLOCK={{ .EnableCmdBlock }}
      {{- end}}
      - MAX_PLAYERS
      - DIFFICULTY
      - MOTD
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
      {{- if .Seed }}
      - LEVEL_SEED={{ .Seed }}
      {{- end}}
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

func StartServerFromCompose(worldName string) error {
	if strings.TrimSpace(worldName) == "" {
		return fmt.Errorf("world name is required to start server")
	}

	cfg := NewMinecraftConfig()
	cfg.WorldName = worldName
	composeFile, err := GetComposeFile(cfg)
	if err != nil {
		return err
	}

	execCfg := utils.ExecConfig{
		Command: "docker",
		Args: []string{
			"compose",
			"-f",
			composeFile,
			"up",
			"-d",
		},
	}

	if mcRunDir, err := McRunHomeDir(); err == nil {
		if meta, metaErr := LoadServerMetadata(filepath.Join(mcRunDir, worldName)); metaErr == nil {
			if strings.TrimSpace(meta.MinMemory) != "" {
				cfg.MinMemory = meta.MinMemory
			}
			if strings.TrimSpace(meta.MaxMemory) != "" {
				cfg.MaxMemory = meta.MaxMemory
			}
		}
	}

	execCfg.Environment = map[string]string{
		"JAVA_MIN_MEM": cfg.MinMemory,
		"JAVA_MAX_MEM": cfg.MaxMemory,
	}

	out, execErr := utils.Exec(execCfg)
	if execErr != nil {
		fmt.Println(out)
		return execErr
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

func DeleteServerResources(worldName string) error {
	if worldName == "" {
		return fmt.Errorf("world name is required to delete resources")
	}

	mcRunDir, err := McRunHomeDir()
	if err != nil {
		return err
	}

	rootDir := filepath.Join(mcRunDir, worldName)
	composeFile := fmt.Sprintf("%s/docker-compose-%s.yaml", mcRunDir, worldName)

	if utils.FileExists(composeFile) {
		if err := os.Remove(composeFile); err != nil {
			return err
		}
	}

	if utils.FileExists(rootDir) {
		return os.RemoveAll(rootDir)
	}

	return fmt.Errorf("server directory %s does not exist: %w", rootDir, os.ErrNotExist)
}

// GetServerStatus inspects the docker container backing the specified world and
// returns the container state (running, exited, etc.). Missing containers result
// in a "stopped" status.
func GetServerStatus(worldName string) string {
	if strings.TrimSpace(worldName) == "" {
		return "unknown"
	}

	containerName := fmt.Sprintf("%s-minecraft", worldName)
	execCfg := utils.ExecConfig{
		Command: "docker",
		Args:    []string{"inspect", "-f", "{{.State.Status}}", containerName},
	}

	out, err := utils.Exec(execCfg)
	if err != nil {
		if strings.Contains(out, "No such object") {
			return "stopped"
		}
		return "unknown"
	}

	status := strings.TrimSpace(out)
	if status == "" {
		status = "unknown"
	}
	return status
}
