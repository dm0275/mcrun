package version

import (
	"fmt"
)

var (
	version string
)

func Version() string {
	return fmt.Sprintf("Version: %s\n", version)
}
