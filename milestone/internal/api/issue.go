package api

import (
	"context"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/go-github/v70/github"

	"github.com/417-72KI/gh-milestone/milestone/internal/ghrepo"
)

type GetIssueListOfMilestoneOptions struct {
	IO        *iostreams.IOStreams
	Repo      ghrepo.Interface
	Milestone *github.Milestone
}

func GetIssueListOfMilestone(ctx context.Context, opts GetIssueListOfMilestoneOptions) ([]*github.Issue, error) {
	gh, err := ghClient(ctx, WithBaseURL(ghrepo.HostWithScheme(opts.Repo)))
	if err != nil {
		return nil, err
	}

	issues, _, err := gh.Issues.ListByRepo(ctx, opts.Repo.RepoOwner(), opts.Repo.RepoName(), &github.IssueListByRepoOptions{
		Milestone: *opts.Milestone.Title,
	})

	return issues, err
}
