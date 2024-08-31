package version

import "testing"

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "success",
			want: "Version: 1.0.0",
		},
	}
	for _, tt := range tests {
		version = "1.0.0"

		t.Run(tt.name, func(t *testing.T) {
			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
