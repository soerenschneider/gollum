package github

import (
	"fmt"
	"time"
)

type ArtifactQuery struct {
	Owner   string
	Repo    string
	Release Release
}

type RepoQuery struct {
	Owner          string
	Repo           string
	Since          *time.Time
	IgnoreReleases []string
}

type Release struct {
	ID      int64  `json:"id"`
	TagName string `json:"tag_name"`

	HasAssets *bool `json:"has_assets"`
}

type ReleaseAsset struct {
	URL                string    `json:"url"`
	ID                 int64     `json:"id"`
	NodeID             string    `json:"node_id"`
	Name               string    `json:"name"`
	Label              string    `json:"label"`
	ContentType        string    `json:"content_type"`
	State              string    `json:"state"`
	Size               int64     `json:"size"`
	DownloadCount      int64     `json:"download_count"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	BrowserDownloadURL string    `json:"browser_download_url"`
}

type Package struct {
	Name       string `json:"name"`
	Tag        string `json:"tag"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Digest     string `json:"digest"`
	Repository string `json:"repository"`
}

type RateLimitInfo struct {
	Limit     int
	Remaining int
	Reset     time.Time
}

type RateLimitError struct {
	Info RateLimitInfo
}

// Error implements the error interface for RateLimitError.
func (e *RateLimitError) Error() string {
	return fmt.Sprintf("Rate limit exceeded. Limit: %d, Remaining: %d, Reset at: %s",
		e.Info.Limit, e.Info.Remaining, e.Info.Reset)
}
