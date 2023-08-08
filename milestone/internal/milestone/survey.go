package milestone

import "strings"

type Prompt interface {
	Input(string, string) (string, error)
	Select(string, string, []string) (int, error)
	MarkdownEditor(string, string, bool) (string, error)
	Confirm(string, bool) (bool, error)
}

func TitleSurvey(p Prompt, state *MilestoneMetadataState) error {
	result, err := p.Input("Title", state.Title)
	if err != nil {
		return err
	}

	if result != state.Title {
		state.MarkDirty()
	}

	state.Title = result

	return nil
}

func DescriptionSurvey(p Prompt, state *MilestoneMetadataState, templateContent string) error {
	if templateContent != "" {
		if state.Description != "" {
			// prevent excessive newlines between default body and template
			state.Description = strings.TrimRight(state.Description, "\n")
			state.Description += "\n\n"
		}
		state.Description += templateContent
	}

	result, err := p.MarkdownEditor("Description", state.Description, true)
	if err != nil {
		return err
	}

	if state.Description != result {
		state.MarkDirty()
	}

	state.Description = result

	return nil
}
