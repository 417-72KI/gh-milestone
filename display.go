package milestone

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/cli/cli/v2/pkg/markdown"
	"github.com/cli/cli/v2/utils"
	"github.com/cli/go-gh/pkg/text"
	"github.com/google/go-github/v47/github"
)

var whitespaceRE = regexp.MustCompile(`\s+`)

func RemoveExcessiveWhitespace(s string) string {
	return whitespaceRE.ReplaceAllString(strings.TrimSpace(s), " ")
}

func PrintMilestones(io *iostreams.IOStreams, now time.Time, prefix string, totalCount int, milestones []*github.Milestone) {
	cs := io.ColorScheme()
	table := utils.NewTablePrinter(io)
	table.AddField("Title", nil, cs.Bold)
	if !table.IsTTY() {
		table.AddField("Status", nil, cs.Bold)
	}
	table.AddField("Due date", nil, cs.Bold)
	table.AddField("Number", nil, cs.Bold)
	table.EndRow()
	for _, milestone := range milestones {
		title := milestone.Title
		table.AddField(RemoveExcessiveWhitespace(*title), nil, cs.Bold)
		if !table.IsTTY() {
			table.AddField(*milestone.State, nil, nil)
		}
		now = time.Now()
		dueOn := milestone.DueOn
		if dueOn == nil {
			table.AddField("", nil, nil)
		} else if now.Before(*dueOn) {
			table.AddField(dueOn.Format("2006/01/02"), nil, nil)
		} else if *milestone.State == "open" {
			AddTimeField(table, now, *dueOn, "", "(over)", cs.Yellow)
		} else {
			AddTimeField(table, now, *dueOn, "", "", cs.Gray)
		}
		table.AddField(strconv.Itoa(*milestone.Number), nil, nil)

		table.EndRow()
	}
	table.Render()
	remaining := totalCount - len(milestones)
	if remaining > 0 {
		fmt.Fprintf(io.Out, cs.Gray("%sAnd %d more\n"), prefix, remaining)
	}
}

func printRawMilestonePreview(out io.Writer, milestone *github.Milestone) error {
	fmt.Fprintf(out, "title:\t\t%s\n", *milestone.Title)
	fmt.Fprintf(out, "state:\t\t%s\n", *milestone.State)
	fmt.Fprintf(out, "description:\t%s\n", *milestone.Description)
	fmt.Fprintf(out, "due on:\t\t%s\n", *milestone.DueOn)
	return nil
}

func printReadableMilestonePreview(io *iostreams.IOStreams, milestone *github.Milestone) error {
	out := io.Out
	cs := io.ColorScheme()
	now := time.Now()

	fmt.Fprintf(out, "%s (%d)\n", cs.Bold(*milestone.Title), *milestone.Number)
	fmt.Fprintf(out,
		"%s • Last updated %s\n",
		milestoneStateWithColor(cs, now, milestone),
		text.RelativeTimeAgo(now, *milestone.UpdatedAt),
	)

	var (
		md  string
		err error
	)
	if *milestone.Description == "" {
		md = fmt.Sprintf("\n  %s\n\n", cs.Gray("No description provided"))
	} else {
		md, err = markdown.Render(*milestone.Description,
			markdown.WithTheme(io.TerminalTheme()),
			markdown.WithWrap(io.TerminalWidth()))
		if err != nil {
			return err
		}
	}
	fmt.Fprintf(out, "\n%s\n", md)
	fmt.Fprintf(out, cs.Gray("View this milestone on GitHub: %s\n"), *milestone.HTMLURL)

	return nil
}

func milestoneStateWithColor(cs *iostreams.ColorScheme, now time.Time, milestone *github.Milestone) string {
	if *milestone.State == "closed" {
		return fmt.Sprintf("%s %s", cs.Bold("Closed"), text.RelativeTimeAgo(now, *milestone.ClosedAt))
	}
	dueDate := milestone.DueOn
	if dueDate == nil {
		return "No due date"
	} else if now.Before(*dueDate) {
		return fmt.Sprintf("Due by %s", dueDate.Format("January 02, 2006"))
	} else {
		return fmt.Sprintf("Past due by %s", text.RelativeTimeAgo(now, *dueDate))
	}
}

func colorForMilestoneState(cs *iostreams.ColorScheme, now time.Time, milestone *github.Milestone) func(string) string {
	switch *milestone.State {
	case "open":
		if now.Before(*milestone.DueOn) {
			return cs.Green
		} else {
			return cs.Yellow
		}
	case "closed":
		return cs.Magenta
	default:
		return nil
	}
}

func AddTimeField(tp utils.TablePrinter, now, t time.Time, prefix string, suffix string, colorFunc func(string) string) {
	tf := t.Format(time.RFC3339)
	if tp.IsTTY() {
		tf = text.RelativeTimeAgo(now, t)
	}
	if len(prefix) != 0 {
		tf = prefix + " " + tf
	}
	if len(suffix) != 0 {
		tf = tf + " " + suffix
	}
	tp.AddField(tf, nil, colorFunc)
}
