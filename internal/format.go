package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
)

func PrintTable(w io.Writer, issues []*Issue) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "ID\tTYPE\tPRI\tSTATUS\tTITLE\tLABELS\n")
	for _, i := range issues {
		labels := ""
		if len(i.Labels) > 0 {
			labels = strings.Join(i.Labels, ",")
		}
		fmt.Fprintf(tw, "%s\t%s\t%d\t%s\t%s\t%s\n",
			i.ID, i.Type, i.Priority, i.Status, i.Title, labels)
	}
	tw.Flush()
}

func PrintJSON(w io.Writer, v interface{}) {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	enc.Encode(v)
}

func PrintStatusSummary(w io.Writer, issues []*Issue) {
	byStatus := map[string]int{}
	byType := map[string]int{}
	byPriority := map[int]int{}
	for _, i := range issues {
		byStatus[i.Status]++
		byType[i.Type]++
		byPriority[i.Priority]++
	}

	fmt.Fprintf(w, "Total: %d issues\n\n", len(issues))

	fmt.Fprintln(w, "By Status:")
	for _, s := range []string{"open", "in_progress", "done", "cancelled"} {
		if n := byStatus[s]; n > 0 {
			fmt.Fprintf(w, "  %-12s %d\n", s, n)
		}
	}

	fmt.Fprintln(w, "\nBy Type:")
	for _, t := range []string{"task", "bug", "feature", "epic", "chore"} {
		if n := byType[t]; n > 0 {
			fmt.Fprintf(w, "  %-12s %d\n", t, n)
		}
	}

	fmt.Fprintln(w, "\nBy Priority:")
	for p := 0; p <= 4; p++ {
		if n := byPriority[p]; n > 0 {
			fmt.Fprintf(w, "  %-12d %d\n", p, n)
		}
	}
}
