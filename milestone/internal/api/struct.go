package api

import (
	"github.com/google/go-github/v62/github"
)

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

func ConvertMilestoneToMap(milestone *github.Milestone, fields []string) map[string]any {
	output := map[string]any{}
	for _, field := range fields {
		switch field {
		case "url":
			output[field] = *milestone.URL
		case "id":
			output[field] = *milestone.ID
		case "number":
			output[field] = *milestone.Number
		case "state":
			output[field] = *milestone.State
		case "title":
			output[field] = *milestone.Title
		case "description":
			output[field] = *milestone.Description
		case "creator":
			output[field] = *milestone.Creator.Login
		case "openIssues":
			output[field] = *milestone.OpenIssues
		case "closedIssues":
			output[field] = *milestone.ClosedIssues
		case "createdAt":
			output[field] = milestone.CreatedAt
		case "updatedAt":
			output[field] = milestone.UpdatedAt
		case "closedAt":
			if milestone.ClosedAt == nil {
				output[field] = "(null)"
			} else {
				output[field] = milestone.ClosedAt
			}
		case "dueOn":
			if milestone.DueOn == nil {
				output[field] = "(null)"
			} else {
				output[field] = milestone.DueOn.Format("2006/01/02")
			}
		}
	}
	return output
}
