package curseforge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://www.curseforge.com/api/v1"
)

// Client exposes minimal functionality required to download mods from CurseForge.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewClient constructs a CurseForge client. The apiKey is optional and only used
// if provided.
func NewClient(apiKey string) (*Client, error) {
	return &Client{
		apiKey:  strings.TrimSpace(apiKey),
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// DownloadMod fetches metadata and downloads the target file into destDir. If the
// file already exists, the existing path is returned.
func (c *Client) DownloadMod(ctx context.Context, projectID, fileID int, destDir string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	fileMeta, err := c.fetchFileMetadata(ctx, projectID, fileID)
	if err != nil {
		return "", err
	}

	if fileMeta.FileName == "" {
		return "", fmt.Errorf("curseforge response missing fileName for project %d file %d", projectID, fileID)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	destPath := filepath.Join(destDir, fileMeta.FileName)
	if _, err := os.Stat(destPath); err == nil {
		return destPath, nil
	}

	if err := c.downloadToPath(ctx, projectID, fileID, destPath); err != nil {
		return "", err
	}

	return destPath, nil
}

// ResolveFileID tries to match a file by loader + Minecraft version similar to download.sh.
func (c *Client) ResolveFileID(ctx context.Context, projectID int, loaders []string, gameVersion string) (int, error) {
	if strings.TrimSpace(gameVersion) == "" {
		return 0, fmt.Errorf("gameVersion is required to resolve a CurseForge file")
	}

	files, err := c.listFiles(ctx, projectID)
	if err != nil {
		return 0, err
	}

	loaderCandidates := loaders
	if len(loaderCandidates) == 0 {
		loaderCandidates = []string{""}
	} else {
		loaderCandidates = append(loaderCandidates, "")
	}

	for _, loader := range loaderCandidates {
		for _, f := range files {
			if !containsIgnoreCase(f.GameVersions, gameVersion) {
				continue
			}
			if loader != "" && !containsIgnoreCase(f.GameVersions, loader) {
				continue
			}
			return f.ID, nil
		}
	}

	return 0, fmt.Errorf("unable to find CurseForge file for project %d matching version %s", projectID, gameVersion)
}

type fileSummary struct {
	ID           int      `json:"id"`
	FileName     string   `json:"fileName"`
	GameVersions []string `json:"gameVersions"`
}

func (c *Client) listFiles(ctx context.Context, projectID int) ([]fileSummary, error) {
	url := fmt.Sprintf("%s/mods/%d/files?pageIndex=0&pageSize=100&sort=dateCreated&sortDescending=true&removeAlphas=true", c.baseURL, projectID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.addHeaders(req)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("curseforge API error (%d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var payload struct {
		Data []fileSummary `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return payload.Data, nil
}

type curseForgeFile struct {
	FileName string `json:"fileName"`
}

func (c *Client) fetchFileMetadata(ctx context.Context, projectID, fileID int) (*curseForgeFile, error) {
	url := fmt.Sprintf("%s/mods/%d/files/%d", c.baseURL, projectID, fileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	c.addHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("curseforge API error (%d): %s", resp.StatusCode, strings.TrimSpace(string(b)))
	}

	var payload struct {
		Data curseForgeFile `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	return &payload.Data, nil
}

func (c *Client) downloadToPath(ctx context.Context, projectID, fileID int, destPath string) error {
	downloadURL := fmt.Sprintf("%s/mods/%d/files/%d/download", c.baseURL, projectID, fileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return err
	}
	c.addHeaders(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download mod (%d)", resp.StatusCode)
	}

	tmpFile, err := os.CreateTemp(filepath.Dir(destPath), "mcrun-mod-*")
	if err != nil {
		return err
	}
	defer func() {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}()

	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	return os.Rename(tmpFile.Name(), destPath)
}

func (c *Client) addHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "mcrun/curseforge-downloader")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}
}

func containsIgnoreCase(values []string, target string) bool {
	target = strings.ToLower(strings.TrimSpace(target))
	for _, v := range values {
		if strings.ToLower(strings.TrimSpace(v)) == target {
			return true
		}
	}
	return false
}
