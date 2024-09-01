package minecraft

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMcRunHomeDir(t *testing.T) {
	homeDir := t.TempDir()
	customDir := t.TempDir()

	tests := []struct {
		name      string
		want      string
		customDir bool
		wantErr   bool
	}{
		{
			name:    "default dir",
			want:    fmt.Sprintf("%s/.mcrun", homeDir),
			wantErr: false,
		},
		{
			name:      "custom dir",
			customDir: true,
			want:      customDir,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("HOME", homeDir)

			if tt.customDir {
				t.Setenv("MCRUN_DIR", customDir)
			}

			got, err := McRunHomeDir()
			if (err != nil) != tt.wantErr {
				t.Errorf("McRunHomeDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("McRunHomeDir() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMinecraftConfig(t *testing.T) {
	tests := []struct {
		name string
		want *MinecraftConfig
	}{
		{
			name: "success",
			want: &MinecraftConfig{
				Version:        "1.19.3",
				Port:           "25565",
				MaxMemory:      "3G",
				MinMemory:      "3G",
				Image:          "dm0275/minecraft-server",
				GameMode:       "0",
				EnableCmdBlock: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMinecraftConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMinecraftConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewMinecraftForgeConfig(t *testing.T) {
	tests := []struct {
		name string
		want *MinecraftConfig
	}{
		{
			name: "success",
			want: &MinecraftConfig{
				Version:        "forge-1.20.1",
				Port:           "25565",
				MaxMemory:      "3G",
				MinMemory:      "3G",
				Image:          "dm0275/minecraft-server",
				GameMode:       "0",
				EnableCmdBlock: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMinecraftForgeConfig(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMinecraftForgeConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetupDirectories(t *testing.T) {
	type args struct {
		mcconfig *MinecraftConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				mcconfig: &MinecraftConfig{
					WorldName: "server1",
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				mcconfig: &MinecraftConfig{
					WorldName: "server1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			homeDir := t.TempDir()
			if tt.wantErr {
				t.Setenv("MCRUN_DIR", "/invalid/dir")
			} else {
				t.Setenv("MCRUN_DIR", homeDir)
			}

			if err := SetupDirectories(tt.args.mcconfig); (err != nil) != tt.wantErr {
				t.Errorf("SetupDirectories() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
