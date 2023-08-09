package ghrepo

import (
	"fmt"
	"strings"
)

type Interface interface {
	RepoName() string
	RepoOwner() string
	RepoHost() string
}

func NewWithHost(owner, repo, hostname string) Interface {
	return &ghRepo{
		owner:    owner,
		name:     repo,
		hostname: normalizeHostname(hostname),
	}
}

func GenerateRepoURL(repo Interface, p string, args ...interface{}) string {
	baseURL := fmt.Sprintf("%s%s/%s", HostWithScheme(repo), repo.RepoOwner(), repo.RepoName())
	if p != "" {
		if path := fmt.Sprintf(p, args...); path != "" {
			return baseURL + "/" + path
		}
	}
	return baseURL
}

func HostWithScheme(repo Interface) string {
	return hostPrefix(repo.RepoHost())
}

const localhost = "github.localhost"

func hostPrefix(hostname string) string {
	if strings.EqualFold(hostname, localhost) {
		return fmt.Sprintf("http://%s/", hostname)
	}
	return fmt.Sprintf("https://%s/", hostname)
}

func normalizeHostname(h string) string {
	return strings.ToLower(strings.TrimPrefix(h, "www."))
}

type ghRepo struct {
	owner    string
	name     string
	hostname string
}

func (r ghRepo) RepoOwner() string {
	return r.owner
}

func (r ghRepo) RepoName() string {
	return r.name
}

func (r ghRepo) RepoHost() string {
	return r.hostname
}
