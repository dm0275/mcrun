package minecraft

import (
	"fmt"
	"strings"
)

// EnsurePortAvailable checks whether another running server already uses the
// provided port. The excludeWorld parameter skips comparisons against the
// currently configured world (useful when reconfiguring an existing server).
func EnsurePortAvailable(port, excludeWorld string) error {
	if strings.TrimSpace(port) == "" {
		return fmt.Errorf("port value is required")
	}

	servers, err := ListServers()
	if err != nil {
		return err
	}

	for _, server := range servers {
		if server.Metadata == nil {
			continue
		}
		if strings.EqualFold(server.WorldName, excludeWorld) {
			continue
		}
		if server.Metadata.Port == port && server.Status == "running" {
			return fmt.Errorf("port %s is already in use by running server %s", port, server.WorldName)
		}
		if strings.TrimSpace(server.Metadata.RconPort) != "" && server.Metadata.RconPort == port && server.Status == "running" {
			return fmt.Errorf("port %s is already in use by running server %s", port, server.WorldName)
		}
	}

	return nil
}
