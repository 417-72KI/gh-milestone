package milestone

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
