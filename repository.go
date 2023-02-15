package milestones

import (
	"fmt"
)

func GenerateRepositoryURL(host string, owner string, repo string, p string, args ...interface{}) string {
	baseURL := fmt.Sprintf("https://%s/%s/%s", host, owner, repo)
	if p != "" {
		if path := fmt.Sprintf(p, args...); path != "" {
			return baseURL + "/" + path
		}
	}
	return baseURL
}
