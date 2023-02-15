package milestone

import (
	"context"

	"github.com/google/go-github/v47/github"
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
