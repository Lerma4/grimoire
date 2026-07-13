package components

import (
	"fmt"
	"strings"

	"github.com/Lerma4/grimoire/internal/domain"
)

// TaskDetail renders the right panel for a selected task, including linked notes.
func TaskDetail(t domain.Task, tags []domain.Tag, linked []domain.Note, width int, focused bool) string {
	var b strings.Builder
	b.WriteString(Styles.Title.Render(truncate(t.Title, width-4)) + "\n\n")

	statusLabel := map[string]string{
		domain.StatusTodo: "○ todo", domain.StatusDoing: "◐ doing",
		domain.StatusDone: "● done", domain.StatusArchived: "□ archived",
	}[t.Status]
	b.WriteString(field("Status", statusLabel))
	b.WriteString(field("Priority", priorityLabel(t.Priority)))
	b.WriteString(field("Due", friendlyDate(t.DueDate)))
	b.WriteString(field("Tags", joinTags(tags)))
	b.WriteString("\n")
	b.WriteString(Styles.Muted.Render("Description") + "\n")
	b.WriteString(wrap(t.Description, width-4))
	b.WriteString("\n\n")
	b.WriteString(Styles.Muted.Render(fmt.Sprintf("Linked notes (%d)", len(linked))) + "\n")
	if len(linked) == 0 {
		b.WriteString(Styles.Muted.Render("  none — press L to link"))
	} else {
		for _, n := range linked {
			b.WriteString("  ✎ " + truncate(n.Title, width-6) + "\n")
		}
	}
	return boxed("Task", focused, width, strings.TrimRight(b.String(), "\n"))
}

// NoteDetail renders the right panel for a selected note, including linked tasks.
func NoteDetail(n domain.Note, tags []domain.Tag, linked []domain.Task, body string, width int, focused bool) string {
	var b strings.Builder
	b.WriteString(Styles.Title.Render(truncate(n.Title, width-4)) + "\n\n")
	b.WriteString(field("Tags", joinTags(tags)))
	b.WriteString("\n")
	b.WriteString(Styles.Muted.Render("Preview") + "\n")
	b.WriteString(body)
	b.WriteString("\n\n")
	b.WriteString(Styles.Muted.Render(fmt.Sprintf("Linked tasks (%d)", len(linked))) + "\n")
	if len(linked) == 0 {
		b.WriteString(Styles.Muted.Render("  none — press L to link"))
	} else {
		for _, t := range linked {
			b.WriteString(fmt.Sprintf("  %s %s\n", t.OverdueGlyph(), truncate(t.Title, width-6)))
		}
	}
	return boxed("Note", focused, width, strings.TrimRight(b.String(), "\n"))
}

func field(name, value string) string {
	return fmt.Sprintf("%s %s\n", Styles.Muted.Render(padRight(name, 9)), value)
}

func priorityLabel(p string) string {
	switch p {
	case domain.PriorityUrgent:
		return "urgent"
	case domain.PriorityHigh:
		return "high"
	case domain.PriorityMedium:
		return "medium"
	case domain.PriorityLow:
		return "low"
	}
	return p
}

func friendlyDate(s string) string {
	if s == "" {
		return Styles.Muted.Render("—")
	}
	return s
}

func joinTags(tags []domain.Tag) string {
	if len(tags) == 0 {
		return Styles.Muted.Render("—")
	}
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = "#" + t.Name
	}
	return strings.Join(names, " ")
}

func wrap(s string, width int) string {
	if s == "" {
		return Styles.Muted.Render("(empty)")
	}
	if width < 8 {
		return s
	}
	var out strings.Builder
	for _, line := range strings.Split(s, "\n") {
		for len(line) > width {
			out.WriteString(line[:width] + "\n")
			line = line[width:]
		}
		out.WriteString(line + "\n")
	}
	return strings.TrimRight(out.String(), "\n")
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
