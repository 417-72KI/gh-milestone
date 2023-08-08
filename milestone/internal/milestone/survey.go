package milestone

import (
	"fmt"
	"strings"
)

type Action int

const (
	SubmitAction Action = iota
	PreviewAction
	MetadataAction
	CancelAction

	submitLabel   = "Submit"
	previewLabel  = "Continue in browser"
	metadataLabel = "Add metadata"
	cancelLabel   = "Cancel"
)

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
	if strings.TrimSpace(result) == "" {
		return fmt.Errorf("title can't be blank")
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

func ConfirmSubmission(p Prompt, allowPreview, allowMetadata bool) (Action, error) {
	var options []string
	options = append(options, submitLabel)

	if allowPreview {
		options = append(options, previewLabel)
	}
	if allowMetadata {
		options = append(options, metadataLabel)
	}
	options = append(options, cancelLabel)

	result, err := p.Select("What's next?", "", options)
	if err != nil {
		return -1, fmt.Errorf("could not prompt: %w", err)
	}

	switch options[result] {
	case submitLabel:
		return SubmitAction, nil
	case previewLabel:
		return PreviewAction, nil
	case metadataLabel:
		return MetadataAction, nil
	case cancelLabel:
		return CancelAction, nil
	default:
		return -1, fmt.Errorf("invalid index: %d", result)
	}
}
