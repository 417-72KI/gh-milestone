package milestone

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
)

func MilestoneNumberAndRepoFromArg(arg string) (int, ghrepo.Interface, error) {
	milestoneNumber, baseRepo := milestoneMetadataFromURL(arg)

	if milestoneNumber == 0 {
		var err error
		milestoneNumber, err = strconv.Atoi(strings.TrimPrefix(arg, "#"))
		if err != nil {
			return 0, nil, fmt.Errorf("invalid milestone format: %q", arg)
		}
	}

	return milestoneNumber, baseRepo, nil
}

var milestoneURLRE = regexp.MustCompile(`^/([^/]+)/([^/]+)/milestone/(\d+)`)

func milestoneMetadataFromURL(s string) (int, ghrepo.Interface) {
	u, err := url.Parse(s)
	if err != nil {
		return 0, nil
	}

	if u.Scheme != "https" && u.Scheme != "http" {
		return 0, nil
	}

	m := milestoneURLRE.FindStringSubmatch(u.Path)
	if m == nil {
		return 0, nil
	}

	repo := ghrepo.NewWithHost(m[1], m[2], u.Hostname())
	milestoneNumber, _ := strconv.Atoi(m[3])
	return milestoneNumber, repo
}
