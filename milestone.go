package milestone

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/go-github/v50/github"
)

func milestones(ctx context.Context, owner string, repo string, state string) ([]*github.Milestone, error) {
	gh, err := ghClient(ctx)
	if err != nil {
		return nil, err
	}
	opts := &github.MilestoneListOptions{
		Direction: "desc",
		State:     state,
	}
	milestones, _, err := gh.Issues.ListMilestones(ctx, owner, repo, opts)
	return milestones, err
}

func getMilestone(ctx context.Context, owner string, repo string, number int) (*github.Milestone, error) {
	gh, err := ghClient(ctx)
	if err != nil {
		return nil, err
	}

	milestone, _, err := gh.Issues.GetMilestone(ctx, owner, repo, number)
	return milestone, err
}

func GetMilestoneByURL(ctx context.Context, url *url.URL) (*github.Milestone, error) {
	scheme := url.Scheme
	host := url.Hostname()
	path := strings.Split(url.Path, "/")

	if len(path) != 5 || path[3] != "milestone" {
		return nil, fmt.Errorf("invalid URL: %s", url)
	}
	owner := path[1]
	repo := path[2]
	number, err := strconv.Atoi(path[4])
	if err != nil {
		return nil, err
	}

	gh, err := ghClient(ctx, WithBaseURL(scheme+"://"+host))
	if err != nil {
		return nil, err
	}

	milestone, _, err := gh.Issues.GetMilestone(ctx, owner, repo, number)
	return milestone, err
}

type closeMilestoneOptions struct {
	ctx       context.Context
	IO        *iostreams.IOStreams
	owner     string
	repo      string
	milestone *github.Milestone
}

func closeMilestone(opts closeMilestoneOptions) (*github.Milestone, error) {
	cs := opts.IO.ColorScheme()
	milestone := opts.milestone

	if *milestone.State == "closed" {
		fmt.Fprintf(opts.IO.ErrOut, cs.Yellow("%s has already closed."), *milestone.HTMLURL)
		return nil, nil
	}

	number := *milestone.Number

	gh, err := ghClient(opts.ctx)
	if err != nil {
		return nil, err
	}

	editedMilestone := &github.Milestone{
		ClosedAt: new(github.Timestamp),
		State:    new(string),
	}
	*editedMilestone.ClosedAt = github.Timestamp{Time: time.Now()}
	*editedMilestone.State = "closed"

	result, _, err := gh.Issues.EditMilestone(opts.ctx, opts.owner, opts.repo, number, editedMilestone)
	return result, err
}
