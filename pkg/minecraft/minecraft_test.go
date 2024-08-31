package minecraft

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerateComposeFile(t *testing.T) {
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
					McRunDir:  t.TempDir(),
				},
			},
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				mcconfig: &MinecraftConfig{
					WorldName: "server1",
					McRunDir:  "/invalid/dir",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenerateComposeFile(tt.args.mcconfig); (err != nil) != tt.wantErr {
				t.Errorf("GenerateComposeFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetComposeFile(t *testing.T) {
	mcRunDir := t.TempDir()

	type args struct {
		mcconfig *MinecraftConfig
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				mcconfig: &MinecraftConfig{
					WorldName: "server1",
					McRunDir:  mcRunDir,
				},
			},
			want:    fmt.Sprintf("%s/docker-compose-server1.yaml", mcRunDir),
			wantErr: false,
		},
		{
			name: "failure",
			args: args{
				mcconfig: &MinecraftConfig{
					WorldName: "server1",
					McRunDir:  "/invalid/dir",
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("MCRUN_DIR", tt.args.mcconfig.McRunDir)

			if !tt.wantErr {
				os.Create(fmt.Sprintf("%s/docker-compose-server1.yaml", tt.args.mcconfig.McRunDir))
			}

			got, err := GetComposeFile(tt.args.mcconfig)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetComposeFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetComposeFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStartServer(t *testing.T) {
	type args struct {
		mcconfig *MinecraftConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StartServer(tt.args.mcconfig); (err != nil) != tt.wantErr {
				t.Errorf("StartServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStopServer(t *testing.T) {
	type args struct {
		dockerComposeFile string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := StopServer(tt.args.dockerComposeFile); (err != nil) != tt.wantErr {
				t.Errorf("StopServer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
