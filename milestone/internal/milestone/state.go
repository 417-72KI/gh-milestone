package milestone

import (
	"time"

	"github.com/google/go-github/v68/github"
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

func ParseTime(t string) (*github.Timestamp, error) {
	location, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}
	dueOn, err := time.ParseInLocation("2006/01/02", t, location)
	if err != nil {
		return nil, err
	}
	return &github.Timestamp{Time: dueOn}, nil
}
