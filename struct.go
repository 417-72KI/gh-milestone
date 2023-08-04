package milestone

var MilestoneFields = []string{
	"id",
	"number",
	"state",
	"title",
	"createdAt",
	"updatedAt",
	"url",
	"openIssues",
	"closedIssues",
}

type FilterOptions struct {
	Author     string
	BaseBranch string
	Fields     []string
	Repo       string
	Search     string
	State      string
}
