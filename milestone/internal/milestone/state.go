package milestone

import (
	"github.com/google/go-github/v53/github"
)

type MilestoneMetadataState struct {
	Title       string
	Description string
	DueOn       *github.Timestamp

	dirty bool // whether user i/o has modified this
}

func (m *MilestoneMetadataState) MarkDirty() {
	m.dirty = true
}

func (m *MilestoneMetadataState) ConvertToMilestone() github.Milestone {
	return github.Milestone{
		Title:       &m.Title,
		Description: &m.Description,
		DueOn:       m.DueOn,
	}
}
