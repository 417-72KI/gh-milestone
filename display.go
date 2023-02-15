package milestone

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cli/cli/v2/pkg/iostreams"
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

func ColorForMilestoneState(cs *iostreams.ColorScheme, milestone *github.Milestone) func(string) string {
	switch *milestone.State {
	case "open":
		return cs.Green
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
