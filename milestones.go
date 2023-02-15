package milestones

import (
	"context"

	"github.com/google/go-github/v47/github"
)

func milestones(ctx context.Context, owner string, repo string, closed bool) ([]*github.Milestone, error) {
	gh, err := ghClient(ctx)
	if err != nil {
		return nil, err
	}
	state := "open"
	if closed {
		state = "closed"
	}
	opts := &github.MilestoneListOptions{
		State: state,
	}
	milestones, _, err := gh.Issues.ListMilestones(ctx, owner, repo, opts)
	return milestones, err
}
