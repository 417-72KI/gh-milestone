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
	"github.com/cli/go-gh/v2/pkg/text"
	"github.com/google/go-github/v53/github"
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
	table.AddField("Completed", nil, cs.Bold)
	table.AddField("Open", nil, cs.Bold)
	table.AddField("Closed", nil, cs.Bold)
	table.EndRow()
	for _, milestone := range milestones {
		title := milestone.Title
		table.AddField(RemoveExcessiveWhitespace(*title), nil, colorForMilestoneState(cs, now, milestone))
		if !table.IsTTY() {
			table.AddField(*milestone.State, nil, nil)
		}
		dueOn := milestone.DueOn
		if dueOn == nil {
			table.AddField("No due date", nil, cs.Gray)
		} else if now.Before(dueOn.Time) {
			table.AddField(dueOn.Format("2006/01/02"), nil, nil)
		} else if *milestone.State == "open" {
			AddTimeField(table, now, dueOn.Time, "", "(over)", cs.Yellow)
		} else {
			AddTimeField(table, now, dueOn.Time, "", "", cs.Gray)
		}
		table.AddField(strconv.Itoa(*milestone.Number), nil, nil)

		table.AddField(fmt.Sprintf("%d%%", completionRate(milestone)), nil, nil)
		table.AddField(fmt.Sprintf("%d", *milestone.OpenIssues), nil, nil)
		table.AddField(fmt.Sprintf("%d", *milestone.ClosedIssues), nil, nil)

		table.EndRow()
	}
	table.Render()
	remaining := totalCount - len(milestones)
	if remaining > 0 {
		fmt.Fprintf(io.Out, cs.Gray("%sAnd %d more\n"), prefix, remaining)
	}
}

func PrintRawMilestonePreview(out io.Writer, milestone *github.Milestone) error {
	fmt.Fprintf(out, "title:\t\t%s\n", *milestone.Title)
	fmt.Fprintf(out, "state:\t\t%s\n", *milestone.State)
	description := "(No description provided.)"
	if milestone.Description != nil && *milestone.Description != "" {
		description = *milestone.Description
	}
	fmt.Fprintf(out, "description:\t%s\n", description)
	dueOn := "(No due date)"
	if milestone.DueOn != nil {
		dueOn = milestone.DueOn.Format("2006-01-02")
	}
	fmt.Fprintf(out, "due on:\t\t%s\n", dueOn)
	fmt.Fprintf(out, "completed:\t%d%%\n", completionRate(milestone))
	fmt.Fprintf(out, "open:\t\t%d\n", *milestone.OpenIssues)
	fmt.Fprintf(out, "closed:\t\t%d\n", *milestone.ClosedIssues)
	return nil
}

func PrintReadableMilestonePreview(io *iostreams.IOStreams, milestone *github.Milestone) error {
	out := io.Out
	cs := io.ColorScheme()
	now := time.Now()

	fmt.Fprintf(out, "%s (%d)\n", cs.Bold(*milestone.Title), *milestone.Number)
	fmt.Fprintf(out,
		"%s • %s complete (%s Open %d Closed) • Last updated %s\n",
		milestoneStateWithColor(cs, now, milestone),
		cs.Boldf("%d%%", completionRate(milestone)),
		cs.Boldf("%d", *milestone.OpenIssues),
		*milestone.ClosedIssues,
		text.RelativeTimeAgo(now, milestone.UpdatedAt.Time),
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
		return fmt.Sprintf("%s %s", cs.Bold("Closed"), text.RelativeTimeAgo(now, milestone.ClosedAt.Time))
	}
	dueDate := milestone.DueOn
	if dueDate == nil {
		return "No due date"
	} else if now.Before(dueDate.Time) {
		return fmt.Sprintf("Due by %s", dueDate.Format("January 02, 2006"))
	} else {
		return fmt.Sprintf("Past due by %s", text.RelativeTimeAgo(now, dueDate.Time))
	}
}

func colorForMilestoneState(cs *iostreams.ColorScheme, now time.Time, milestone *github.Milestone) func(string) string {
	switch *milestone.State {
	case "open":
		if milestone.DueOn == nil || now.Before(milestone.DueOn.Time) {
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

func completionRate(milestone *github.Milestone) int {
	open := *milestone.OpenIssues
	closed := *milestone.ClosedIssues
	total := open + closed
	if total == 0 {
		return 0
	}
	return int(float64(closed) / float64(total) * 100)
}
