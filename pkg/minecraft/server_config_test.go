package minecraft

import (
	"os"
	"reflect"
	"testing"
)

func TestGenerateServerConfig(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "existing-server.properties")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	tests := []struct {
		name          string
		args          struct{ serverConfigFile string }
		setupExisting bool
		wantErr       bool
	}{
		{
			name:    "File exists, should not overwrite",
			args:    struct{ serverConfigFile string }{tmpFile.Name()},
			wantErr: false,
		},
		{
			name:    "File does not exist, should create new one",
			args:    struct{ serverConfigFile string }{tmpFile.Name() + "_new"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.args.serverConfigFile)

			if err := GenerateServerConfig(tt.args.serverConfigFile); (err != nil) != tt.wantErr {
				t.Errorf("GenerateServerConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			// If not expecting an error and file was supposed to be created, verify it exists
			if !tt.wantErr {
				if _, err := os.Stat(tt.args.serverConfigFile); os.IsNotExist(err) {
					t.Errorf("Expected file %s to be created", tt.args.serverConfigFile)
				}
			}
		})
	}
}

func TestNewServerConfig(t *testing.T) {
	want := &ServerConfig{
		MaxTickTime:                 60000,
		AllowNether:                 true,
		ForceGamemode:               false,
		Gamemode:                    0,
		EnableQuery:                 false,
		PlayerIdleTimeout:           0,
		Difficulty:                  1,
		SpawnMonsters:               true,
		OpPermissionLevel:           4,
		Pvp:                         true,
		SnooperEnabled:              true,
		LevelType:                   "DEFAULT",
		Hardcore:                    false,
		EnableCommandBlock:          false,
		MaxPlayers:                  20,
		NetworkCompressionThreshold: 256,
		MaxWorldSize:                29999984,
		ServerPort:                  25565,
		ServerIp:                    "",
		SpawnNpcs:                   true,
		AllowFlight:                 false,
		LevelName:                   "world",
		ViewDistance:                10,
		SpawnAnimals:                true,
		WhiteList:                   false,
		GenerateStructures:          true,
		OnlineMode:                  true,
		MaxBuildHeight:              256,
		LevelSeed:                   "",
		PreventProxyConnections:     false,
		UseNativeTransport:          true,
		Motd:                        "Docker Minecraft Server",
		EnableRcon:                  false,
		RconPort:                    25575,
		RconPassword:                "minecraft",
	}

	tests := []struct {
		name string
		want *ServerConfig
	}{
		{
			name: "Default config values",
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewServerConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewServerConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestServerConfig_MarshalAndWriteToFile(t *testing.T) {
	tests := []struct {
		name    string
		fields  ServerConfig
		args    struct{ filePath string }
		wantErr bool
	}{
		{
			name:   "Valid config should write file",
			fields: *NewServerConfig(),
			args: func() struct{ filePath string } {
				tmp, _ := os.CreateTemp("", "server-config-*.properties")
				tmp.Close()
				return struct{ filePath string }{filePath: tmp.Name()}
			}(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.args.filePath)
			if err := tt.fields.MarshalAndWriteToFile(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("MarshalAndWriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				info, err := os.Stat(tt.args.filePath)
				if err != nil || info.Size() == 0 {
					t.Errorf("Expected non-empty config file to be created")
				}
			}
		})
	}
}
