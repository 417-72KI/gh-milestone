package api

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/google/go-github/v50/github"
	"golang.org/x/oauth2"
)

type clientOptions struct {
	token   string
	baseURL string
}

type clientOption func(*clientOptions)

func WithToken(token string) clientOption {
	return func(ops *clientOptions) {
		ops.token = token
	}
}

func WithBaseURL(baseURL string) clientOption {
	return func(ops *clientOptions) {
		ops.baseURL = baseURL
	}
}

func FetchOwnerAndRepoFromURL(url *url.URL) (*string, *string, error) {
	path := strings.Split(url.Path, "/")
	if len(path) < 3 {
		return nil, nil, fmt.Errorf("invalid URL: %s", url)
	}
	owner := path[1]
	repo := path[2]
	return &owner, &repo, nil
}

func ghClient(ctx context.Context, ops ...clientOption) (*github.Client, error) {
	opts := clientOptions{}
	for _, op := range ops {
		op(&opts)
	}
	token := os.Getenv("GITHUB_TOKEN")
	if opts.token != "" {
		token = opts.token
	}
	if token == "" {
		return nil, errors.New("github token is missing. please use GITHUB_TOKEN environment variable")
	}
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	baseURL := os.Getenv("GITHUB_BASE_URL")
	if opts.baseURL != "" && opts.baseURL != "https://github.com" {
		baseURL = opts.baseURL
	}
	if baseURL != "" {
		var err error
		client, err = github.NewEnterpriseClient(baseURL, baseURL, tc)
		if err != nil {
			return nil, fmt.Errorf("failed to create a new github api client: %w", err)
		}
	}
	return client, nil
}
