package curseforge

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://api.curseforge.com/v1"
)

// Client exposes minimal functionality required to download mods from CurseForge.
type Client struct {
	apiKey     string
	baseURL    string
	userAgent  string
	httpClient *http.Client
}

// NewClient constructs a CurseForge client using the provided API key.
func NewClient(apiKey string) (*Client, error) {
	if strings.TrimSpace(apiKey) == "" {
		return nil, errors.New("curseforge api key is required")
	}

	return &Client{
		apiKey:    strings.TrimSpace(apiKey),
		baseURL:   defaultBaseURL,
		userAgent: "mcrun-curser/1.0",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// DownloadMod fetches the file metadata then downloads the mod into destDir.
// If the file already exists it is returned without downloading again.
func (c *Client) DownloadMod(ctx context.Context, projectID, fileID int, destDir string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	fileMeta, err := c.fetchFileMetadata(ctx, projectID, fileID)
	if err != nil {
		return "", err
	}

	if fileMeta.FileName == "" || fileMeta.DownloadURL == "" {
		return "", fmt.Errorf("curseforge response missing file data for project %d file %d", projectID, fileID)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	destPath := filepath.Join(destDir, fileMeta.FileName)
	if _, err := os.Stat(destPath); err == nil {
		return destPath, nil
	}

	if err := c.downloadToPath(ctx, fileMeta.DownloadURL, destPath); err != nil {
		return "", err
	}

	return destPath, nil
}

type curseForgeFile struct {
	FileName    string `json:"fileName"`
	DownloadURL string `json:"downloadUrl"`
}

func (c *Client) fetchFileMetadata(ctx context.Context, projectID, fileID int) (*curseForgeFile, error) {
	url := fmt.Sprintf("%s/mods/%d/files/%d", c.baseURL, projectID, fileID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

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

func (c *Client) downloadToPath(ctx context.Context, downloadURL, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, downloadURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", c.userAgent)

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
