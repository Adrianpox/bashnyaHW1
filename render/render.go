package render

import (
	"fmt"
	"io"
	"text/tabwriter"

	"homework/colors"
	"homework/task"
)

const dateLayout = "07.02.2003 10:45"

func TaskList(w io.Writer, tasks []task.Task) {
	if len(tasks) == 0 {
		fmt.Fprintln(w, "Нет задач")
		return
	}

	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)
	for i, t := range tasks {
		status := colorize(colors.Yellow, "[ ]")
		if t.IsDone {
			status = colorize(colors.Green, "[x]")
		}
		id := colorize(colors.Bold, fmt.Sprintf("#%d", t.ID))
		fmt.Fprintf(tw, "%d.\t%s\t%s\t%s\t%s\n",
			i+1,
			status,
			t.CreatedAt.Format(dateLayout),
			t.Text,
			id,
		)
	}
	tw.Flush()
}

func colorize(code, s string) string {
	return code + s + colors.Reset
}
