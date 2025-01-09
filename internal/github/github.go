package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/soerenschneider/gollum/internal/metrics"
	"golang.org/x/exp/slices"
)

var ErrUnauthorized = errors.New("unauthorized. either token is invalid, expired or missing the correct scope")

type GithubClient struct {
	httpClient *http.Client
	token      *string

	// unauthorized is as bool that is true when the system detects we lack permissions to call the packages API.
	// this is used to prevent wasting further calls to the API in order to save quota.
	unauthorized atomic.Bool

	rateLimitMutex   sync.RWMutex
	rateLimitedUntil *RateLimitError
}

func NewGithubClient(client *http.Client, token *string) (*GithubClient, error) {
	if client == nil {
		client = &http.Client{
			Timeout: 5 * time.Second,
		}
	}

	ret := &GithubClient{
		httpClient: client,
		token:      token,
	}

	return ret, nil
}

func (g *GithubClient) isRateLimited() error {
	g.rateLimitMutex.RLock()
	defer g.rateLimitMutex.RUnlock()
	if g.rateLimitedUntil != nil && time.Now().Before(g.rateLimitedUntil.Info.Reset) {
		return g.rateLimitedUntil
	}
	return nil
}

func (g *GithubClient) getReleases(ctx context.Context, queryParams RepoQuery) ([]Release, error) {
	metrics.GithubRequestsTotal.WithLabelValues(queryParams.Owner, queryParams.Repo).Inc()
	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases", queryParams.Owner, queryParams.Repo)
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	page := 1
	hasNextPage := true
	var ret []Release

	for hasNextPage {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			err := func() error {
				if err := g.isRateLimited(); err != nil {
					return err
				}

				params := url.Values{}
				params.Add("per_page", "100")
				params.Add("page", strconv.Itoa(page))

				parsedURL.RawQuery = params.Encode()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
				if err != nil {
					return fmt.Errorf("failed to create request: %w", err)
				}

				if g.token != nil && *g.token != "" {
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *g.token))
				}

				resp, err := g.httpClient.Do(req)
				if err != nil {
					return fmt.Errorf("failed to send request: %w", err)
				}

				defer func() {
					_ = resp.Body.Close()
				}()

				if resp.StatusCode != http.StatusOK {
					return g.evaluateAndTransformError(resp)
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return fmt.Errorf("failed to read response body: %w", err)
				}

				var parsed []Release
				if err := json.Unmarshal(body, &parsed); err != nil {
					return fmt.Errorf("failed to parse JSON: %w", err)
				}

				ret = append(ret, parsed...)
				linkHeader := resp.Header.Get("Link")
				hasNextPage = linkHeader != "" && strings.Contains(linkHeader, "rel=\"next\"")
				page++

				return nil
			}()

			if err != nil {
				return nil, err
			}
		}
	}

	if len(ret) == 0 {
		return nil, errors.New("no releases found for the repository")
	}

	return ret, nil
}

func GetRateLimitInfo(resp *http.Response) *RateLimitError {
	limitStr := resp.Header.Get("X-RateLimit-Limit")
	remainingStr := resp.Header.Get("X-RateLimit-Remaining")
	resetStr := resp.Header.Get("X-RateLimit-Reset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return nil
	}

	remaining, err := strconv.Atoi(remainingStr)
	if err != nil {
		return nil
	}

	resetUnix, err := strconv.ParseInt(resetStr, 10, 64)
	if err != nil {
		return nil
	}

	resetTime := time.Unix(resetUnix, 0)

	return &RateLimitError{RateLimitInfo{
		Limit:     limit,
		Remaining: remaining,
		Reset:     resetTime,
	}}
}

func (g *GithubClient) GetReleases(ctx context.Context, params RepoQuery) ([]Release, error) {
	releases, err := g.getReleases(ctx, params)
	if err != nil {
		metrics.GithubRequestErrors.WithLabelValues(params.Owner, params.Repo, "releases").Inc()
		return nil, err
	}

	var ret []Release
	for _, event := range releases {
		// ignore releases that are already built
		if !slices.Contains(params.IgnoreReleases, event.TagName) {
			ret = append(ret, event)
		}
	}

	return ret, nil
}

func (g *GithubClient) evaluateAndTransformError(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}

	genericErr := fmt.Errorf("got status code %d", resp.StatusCode)

	if resp.StatusCode == http.StatusForbidden {
		rateLimitErr := GetRateLimitInfo(resp)
		if rateLimitErr != nil && rateLimitErr.Info.Remaining <= 0 {
			g.rateLimitMutex.Lock()
			defer g.rateLimitMutex.Unlock()
			g.rateLimitedUntil = rateLimitErr
			return rateLimitErr
		}
		return genericErr
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorized
	}

	return genericErr
}

func (g *GithubClient) GetAssets(ctx context.Context, query ArtifactQuery) ([]ReleaseAsset, error) {
	if err := g.isRateLimited(); err != nil {
		return nil, err
	}

	endpoint := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d/assets", query.Owner, query.Repo, query.Release.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		metrics.GithubRequestErrors.WithLabelValues(query.Owner, query.Repo, "assets").Inc()
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if g.token != nil && *g.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *g.token))
	}

	resp, err := g.httpClient.Do(req)
	if err != nil {
		metrics.GithubRequestErrors.WithLabelValues(query.Owner, query.Repo, "assets").Inc()
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, g.evaluateAndTransformError(resp)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		metrics.GithubRequestErrors.WithLabelValues(query.Owner, query.Repo, "assets").Inc()
		return nil, err
	}

	var parsed []ReleaseAsset
	if err := json.Unmarshal(data, &parsed); err != nil {
		metrics.GithubRequestErrors.WithLabelValues(query.Owner, query.Repo, "assets").Inc()
	}

	return parsed, err
}

func (g *GithubClient) GetPackages(ctx context.Context, query ArtifactQuery) ([]Package, error) {
	if g.unauthorized.Load() {
		return nil, ErrUnauthorized
	}

	pckgs, err := g.getPackages(ctx, query.Owner, query.Repo, "container")
	if err != nil {
		if errors.Is(err, ErrUnauthorized) {
			g.unauthorized.Store(true)
		}
		return nil, err
	}

	ret := make([]Package, 0, len(pckgs))
	for _, p := range pckgs {
		if p.Tag == query.Release.TagName {
			ret = append(ret, p)
		}
	}

	return ret, nil
}

func (g *GithubClient) getPackages(ctx context.Context, owner, repo, packageType string) ([]Package, error) {
	endpoint := fmt.Sprintf("https://api.github.com/users/%s/packages/%s/%s/versions", owner, repo, packageType)
	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	page := 1
	hasNextPage := true
	var ret []Package

	for hasNextPage {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			err := func() error {
				if err := g.isRateLimited(); err != nil {
					return err
				}

				params := url.Values{}
				params.Add("per_page", "100")
				params.Add("page", strconv.Itoa(page))

				parsedURL.RawQuery = params.Encode()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
				if err != nil {
					return err
				}

				if g.token != nil && *g.token != "" {
					req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *g.token))
				}

				resp, err := g.httpClient.Do(req)
				if err != nil {
					metrics.GithubRequestErrors.WithLabelValues(owner, repo, "packages").Inc()
					return err
				}

				defer func() {
					_ = resp.Body.Close()
				}()

				if resp.StatusCode != http.StatusOK {
					return g.evaluateAndTransformError(resp)
				}

				data, err := io.ReadAll(resp.Body)
				if err != nil {
					metrics.GithubRequestErrors.WithLabelValues(owner, repo, "packages").Inc()
					return err
				}

				var parsed []Package
				if err := json.Unmarshal(data, &parsed); err != nil {
					metrics.GithubRequestErrors.WithLabelValues(owner, repo, "packages").Inc()
					return err
				}

				ret = append(ret, parsed...)
				linkHeader := resp.Header.Get("Link")
				hasNextPage = linkHeader != "" && strings.Contains(linkHeader, "rel=\"next\"")
				page++

				return nil
			}()

			if err != nil {
				return nil, err
			}
		}
	}

	return ret, nil
}
