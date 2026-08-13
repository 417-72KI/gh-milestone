package prompter

import (
	"github.com/cli/cli/v2/pkg/surveyext"
	ghprompter "github.com/cli/go-gh/v2/pkg/prompter"
)

// Prompter wraps the go-gh prompter and adds MarkdownEditor support.
type Prompter struct {
	gh        *ghprompter.Prompter
	editorCmd string
	stdin     ghprompter.FileReader
	stdout    ghprompter.FileWriter
	stderr    ghprompter.FileWriter
}

// New creates a new Prompter with the given I/O streams and optional editor command.
func New(editorCmd string, stdin ghprompter.FileReader, stdout ghprompter.FileWriter, stderr ghprompter.FileWriter) *Prompter {
	return &Prompter{
		gh:        ghprompter.New(stdin, stdout, stderr),
		editorCmd: editorCmd,
		stdin:     stdin,
		stdout:    stdout,
		stderr:    stderr,
	}
}

// Input prompts the user to input a single-line string.
func (p *Prompter) Input(prompt string, defaultValue string) (string, error) {
	return p.gh.Input(prompt, defaultValue)
}

// Select prompts the user to select an option from a list of options.
func (p *Prompter) Select(prompt string, defaultValue string, options []string) (int, error) {
	return p.gh.Select(prompt, defaultValue, options)
}

// Confirm prompts the user to confirm a yes/no question.
func (p *Prompter) Confirm(prompt string, defaultValue bool) (bool, error) {
	return p.gh.Confirm(prompt, defaultValue)
}

// MarkdownEditor prompts the user to edit markdown in an editor.
// If blankAllowed is true, the user can skip editing and an empty string will be returned.
func (p *Prompter) MarkdownEditor(prompt string, defaultValue string, blankAllowed bool) (string, error) {
	options := []string{
		"Launch " + surveyext.EditorName(p.editorCmd),
	}
	if blankAllowed {
		options = append(options, "Skip")
	}

	idx, err := p.gh.Select(prompt, "", options)
	if err != nil {
		return "", err
	}

	// If "Skip" was selected (only if blankAllowed), return empty string.
	if blankAllowed && idx == 1 {
		return "", nil
	}

	// Otherwise, launch the editor.
	text, err := surveyext.Edit(p.editorCmd, "*.md", defaultValue, p.stdin, p.stdout, p.stderr)
	if err != nil {
		return "", err
	}

	return text, nil
}
