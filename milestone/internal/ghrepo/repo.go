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

func GenerateRepoURL(repo Interface, p string, args ...interface{}) string {
	baseURL := fmt.Sprintf("%s%s/%s", hostPrefix(repo.RepoHost()), repo.RepoOwner(), repo.RepoName())
	if p != "" {
		if path := fmt.Sprintf(p, args...); path != "" {
			return baseURL + "/" + path
		}
	}
	return baseURL
}

const localhost = "github.localhost"

func hostPrefix(hostname string) string {
	if strings.EqualFold(hostname, localhost) {
		return fmt.Sprintf("http://%s/", hostname)
	}
	return fmt.Sprintf("https://%s/", hostname)
}
