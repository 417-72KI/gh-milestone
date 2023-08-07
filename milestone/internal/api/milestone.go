package api

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/go-github/v53/github"
)

func Milestones(ctx context.Context, owner string, repo string, filterOpts FilterOptions) ([]*github.Milestone, error) {
	gh, err := ghClient(ctx)
	if err != nil {
		return nil, err
	}
	opts := &github.MilestoneListOptions{
		Direction: "desc",
		State:     filterOpts.State,
	}
	milestones, _, err := gh.Issues.ListMilestones(ctx, owner, repo, opts)
	return milestones, err
}

func GetMilestone(ctx context.Context, owner string, repo string, number int) (*github.Milestone, error) {
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
	owner, repo, err := FetchOwnerAndRepoFromURL(url)
	if err != nil {
		return nil, err
	}

	if len(path) != 5 || path[3] != "milestone" {
		return nil, fmt.Errorf("invalid URL: %s", url)
	}
	number, err := strconv.Atoi(path[4])
	if err != nil {
		return nil, err
	}

	gh, err := ghClient(ctx, WithBaseURL(scheme+"://"+host))
	if err != nil {
		return nil, err
	}

	milestone, _, err := gh.Issues.GetMilestone(ctx, *owner, *repo, number)
	return milestone, err
}

type CloseMilestoneOptions struct {
	IO        *iostreams.IOStreams
	Owner     string
	Repo      string
	Milestone *github.Milestone
}

func CloseMilestone(ctx context.Context, opts CloseMilestoneOptions) (*github.Milestone, error) {
	cs := opts.IO.ColorScheme()
	milestone := opts.Milestone

	if *milestone.State == "closed" {
		fmt.Fprintf(opts.IO.ErrOut, cs.Yellow("%s has already closed.\n"), *milestone.HTMLURL)
		return nil, nil
	}
	number := *milestone.Number

	gh, err := ghClient(ctx)
	if err != nil {
		return nil, err
	}

	editedMilestone := &github.Milestone{
		ClosedAt: new(github.Timestamp),
		State:    new(string),
	}
	*editedMilestone.ClosedAt = github.Timestamp{Time: time.Now()}
	*editedMilestone.State = "closed"

	result, _, err := gh.Issues.EditMilestone(ctx, opts.Owner, opts.Repo, number, editedMilestone)
	return result, err
}

type ReopenMilestoneOptions struct {
	IO        *iostreams.IOStreams
	Owner     string
	Repo      string
	Milestone *github.Milestone
}

func ReopenMilestone(ctx context.Context, opts ReopenMilestoneOptions) (*github.Milestone, error) {
	cs := opts.IO.ColorScheme()
	milestone := opts.Milestone

	if *milestone.State != "closed" {
		fmt.Fprintf(opts.IO.ErrOut, cs.Yellow("%s has not closed.\n"), *milestone.HTMLURL)
		return nil, nil
	}
	number := *milestone.Number

	gh, err := ghClient(ctx)
	if err != nil {
		return nil, err
	}

	editedMilestone := &github.Milestone{
		ClosedAt: new(github.Timestamp),
		State:    new(string),
	}
	*editedMilestone.ClosedAt = github.Timestamp{Time: time.Now()}
	*editedMilestone.State = "open"

	result, _, err := gh.Issues.EditMilestone(ctx, opts.Owner, opts.Repo, number, editedMilestone)
	return result, err
}
