package syncer

import "fmt"

type Summary struct {
	Create int
	Update int
	Delete int
	Noop   int
}

func Summarize(plan Plan) Summary {
	var s Summary
	for _, action := range plan.Actions {
		switch action.Type {
		case ActionCreate:
			s.Create++
		case ActionUpdate:
			s.Update++
		case ActionDelete:
			s.Delete++
		case ActionNoop:
			s.Noop++
		}
	}
	return s
}

func FormatSummary(summary Summary) string {
	return fmt.Sprintf("create: %d, update: %d, delete: %d, noop: %d", summary.Create, summary.Update, summary.Delete, summary.Noop)
}

func FormatActions(plan Plan) []string {
	var lines []string
	for _, action := range plan.Actions {
		if action.Type == ActionNoop {
			continue
		}
		line := fmt.Sprintf("%s %s -> %s", action.Type, action.Source, action.Dest)
		if action.Note != "" {
			line = fmt.Sprintf("%s (%s)", line, action.Note)
		}
		lines = append(lines, line)
	}
	if len(lines) == 0 {
		lines = append(lines, "no changes")
	}
	return lines
}
