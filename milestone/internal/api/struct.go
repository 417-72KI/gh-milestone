package api

var MilestoneFields = []string{
	"url",
	"id",
	"number",
	"state",
	"title",
	"description",
	"creator",
	"openIssues",
	"closedIssues",
	"createdAt",
	"updatedAt",
	"closedAt",
	"dueOn",
}

type FilterOptions struct {
	Author     string
	BaseBranch string
	Fields     []string
	Repo       string
	Search     string
	State      string
}
