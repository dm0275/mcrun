package curseforge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	defaultBaseURL = "https://www.curseforge.com/api/v1"
	widgetBaseURL  = "https://api.cfwidget.com"
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

// DownloadMod fetches metadata and downloads the target file into destDir. It
// returns the cached file path and whether the file already existed.
func (c *Client) DownloadMod(ctx context.Context, projectID, fileID int, destDir string) (string, bool, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	fileMeta, err := c.fetchFileMetadata(ctx, projectID, fileID)
	if err != nil {
		return "", false, err
	}

	if fileMeta.FileName == "" {
		return "", false, fmt.Errorf("curseforge response missing fileName for project %d file %d", projectID, fileID)
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", false, err
	}

	destPath := filepath.Join(destDir, fileMeta.FileName)
	if _, err := os.Stat(destPath); err == nil {
		return destPath, true, nil
	}

	if err := c.downloadToPath(ctx, projectID, fileID, destPath); err != nil {
		return "", false, err
	}

	return destPath, false, nil
}

func (c *Client) DownloadURL(projectID, fileID int) string {
	return fmt.Sprintf("%s/mods/%d/files/%d/download", c.baseURL, projectID, fileID)
}

// ResolveFileID tries to match a file by loader + Minecraft version similar to download.sh.
func (c *Client) ResolveFileID(ctx context.Context, projectID int, loaders []string, gameVersion string) (int, error) {
	if strings.TrimSpace(gameVersion) == "" {
		return 0, fmt.Errorf("gameVersion is required to resolve a CurseForge file")
	}

	project, err := c.fetchWidgetProject(ctx, projectID)
	if err != nil {
		return 0, err
	}

	fileID := selectWidgetFileID(project.Files, loaders, gameVersion)
	if fileID != 0 {
		return fileID, nil
	}

	return 0, fmt.Errorf("unable to find CurseForge file for project %d matching version %s", projectID, gameVersion)
}

type fileSummary struct {
	ID           int      `json:"id"`
	FileName     string   `json:"fileName"`
	GameVersions []string `json:"gameVersions"`
}

type curseForgeFile struct {
	FileName string `json:"fileName"`
}

type ModLinks struct {
	WebsiteURL string `json:"websiteUrl"`
}

type ModSummary struct {
	Name  string   `json:"name"`
	Links ModLinks `json:"links"`
}

type widgetProject struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	URLs  struct {
		CurseForge string `json:"curseforge"`
	} `json:"urls"`
	Files []widgetFile `json:"files"`
}

type widgetFile struct {
	ID         int      `json:"id"`
	URL        string   `json:"url"`
	Versions   []string `json:"versions"`
	UploadedAt string   `json:"uploaded_at"`
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
	downloadURL := c.DownloadURL(projectID, fileID)

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

// FetchModSummary returns basic information about a CurseForge project.
func (c *Client) FetchModSummary(ctx context.Context, projectID int) (*ModSummary, error) {
	project, err := c.fetchWidgetProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	summary := &ModSummary{
		Name: project.Title,
		Links: ModLinks{
			WebsiteURL: project.URLs.CurseForge,
		},
	}
	return summary, nil
}

func (c *Client) fetchWidgetProject(ctx context.Context, projectID int) (*widgetProject, error) {
	url := fmt.Sprintf("%s/%d", widgetBaseURL, projectID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 8<<10))
		return nil, fmt.Errorf("cfwidget API error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var project widgetProject
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

func selectWidgetFileID(files []widgetFile, loaders []string, gameVersion string) int {
	if len(files) == 0 {
		return 0
	}

	type fileWithTime struct {
		file widgetFile
		t    time.Time
	}

	fileList := make([]fileWithTime, 0, len(files))
	for _, f := range files {
		t, _ := time.Parse(time.RFC3339, f.UploadedAt)
		fileList = append(fileList, fileWithTime{file: f, t: t})
	}

	sort.Slice(fileList, func(i, j int) bool {
		return fileList[i].t.After(fileList[j].t)
	})

	normalizedGame := strings.ToLower(strings.TrimSpace(gameVersion))
	loaderCandidates := loaders
	if len(loaderCandidates) == 0 {
		loaderCandidates = []string{""}
	} else {
		loaderCandidates = append(loaderCandidates, "")
	}

	for _, loader := range loaderCandidates {
		normalizedLoader := strings.ToLower(strings.TrimSpace(loader))
		for _, entry := range fileList {
			if !containsIgnoreCase(entry.file.Versions, normalizedGame) {
				continue
			}
			if normalizedLoader != "" && !containsIgnoreCase(entry.file.Versions, normalizedLoader) {
				continue
			}
			return entry.file.ID
		}
	}

	return 0
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
