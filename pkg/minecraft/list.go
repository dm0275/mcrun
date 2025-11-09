package minecraft

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dm0275/mcrun/utils"
)

// ServerInfo describes minimal metadata for a provisioned server directory.
type ServerInfo struct {
	WorldName  string          `json:"worldName"`
	Path       string          `json:"path"`
	HasCompose bool            `json:"hasCompose"`
	Status     string          `json:"status"`
	Metadata   *ServerMetadata `json:"metadata,omitempty"`
}

// ListServers enumerates the directories under the mcrun home directory and
// reports basic metadata for each potential server.
func ListServers() ([]ServerInfo, error) {
	mcRunDir, err := McRunHomeDir()
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(mcRunDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []ServerInfo{}, nil
		}
		return nil, err
	}

	servers := make([]ServerInfo, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if name == "cache" || strings.HasPrefix(name, ".") {
			continue
		}

		info := ServerInfo{
			WorldName: name,
			Path:      filepath.Join(mcRunDir, name),
		}

		composeFile := fmt.Sprintf("%s/docker-compose-%s.yaml", mcRunDir, name)
		info.HasCompose = utils.FileExists(composeFile)
		info.Status = GetServerStatus(name)

		if meta, err := LoadServerMetadata(info.Path); err == nil {
			info.Metadata = meta
		}

		servers = append(servers, info)
	}

	sort.Slice(servers, func(i, j int) bool {
		return servers[i].WorldName < servers[j].WorldName
	})

	return servers, nil
}
