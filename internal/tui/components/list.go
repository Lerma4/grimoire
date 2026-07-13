package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Lerma4/grimoire/internal/domain"
)

// Counts holds per-section item counts shown in the sidebar.
type Counts struct {
	Tasks, Notes, Today, Tags, Links, Archive int
}

// Sidebar renders the left navigation column.
func Sidebar(active Section, counts Counts, width int, focused bool) string {
	head := Styles.Title.Render("✦ Grimoire")
	items := []string{head, ""}

	render := func(s Section, n int) {
		name := s.String()
		cnt := fmt.Sprintf("%d", n)
		if width < 18 {
			name = "" // collapse to glyph only
		}
		var left string
		style := Styles.Muted
		if s == active {
			left = fmt.Sprintf("%s %s", Styles.Accent.Render("▸"), s.Glyph())
			style = Styles.Selected
		} else {
			left = fmt.Sprintf("   %s", s.Glyph())
		}
		if name != "" {
			left += " " + name
		}
		gap := width - 2 - lipgloss.Width(left) - lipgloss.Width(cnt)
		if gap < 1 {
			gap = 1
		}
		line := left + strings.Repeat(" ", gap) + Styles.Muted.Render(cnt)
		items = append(items, style.Render(line))
	}

	render(SectionTasks, counts.Tasks)
	render(SectionNotes, counts.Notes)
	render(SectionToday, counts.Today)
	render(SectionTags, counts.Tags)
	render(SectionLinks, counts.Links)
	render(SectionArchive, counts.Archive)

	body := strings.Join(items, "\n")
	return boxed("Sections", focused, width, body)
}

// TaskList renders the central column when the active section shows tasks.
func TaskList(title string, tasks []domain.Task, cursor, width int, focused bool) string {
	if len(tasks) == 0 {
		return boxed(title, focused, width, Styles.Muted.Render("  no tasks — press a to add one"))
	}
	rows := make([]string, 0, len(tasks))
	for i, t := range tasks {
		glyph := Styles.Warning.Render(t.OverdueGlyph())
		switch t.Status {
		case domain.StatusDone:
			glyph = Styles.Success.Render(t.OverdueGlyph())
		case domain.StatusDoing:
			glyph = Styles.Accent.Render(t.OverdueGlyph())
		}
		mark := " "
		prio := priorityMark(t.Priority)
		label := truncate(t.Title, width-8)
		line := fmt.Sprintf("%s %s %s %s", mark, glyph, prio, label)
		if i == cursor {
			line = fmt.Sprintf("%s %s %s %s",
				Styles.Accent.Render("▸"), glyph, prio, Styles.CursorRow.Render(label))
		}
		rows = append(rows, line)
	}
	return boxed(title, focused, width, strings.Join(rows, "\n"))
}

// NoteList renders the central column when the active section shows notes.
func NoteList(title string, notes []domain.Note, cursor, width int, focused bool) string {
	if len(notes) == 0 {
		return boxed(title, focused, width, Styles.Muted.Render("  no notes — press A to add one"))
	}
	rows := make([]string, 0, len(notes))
	for i, n := range notes {
		mark := " "
		archived := ""
		if n.ArchivedAt != "" {
			archived = Styles.Muted.Render(" ▽")
		}
		label := truncate(n.Title, width-6)
		line := fmt.Sprintf("%s %s%s", mark, label, archived)
		if i == cursor {
			line = fmt.Sprintf("%s %s%s", Styles.Accent.Render("▸"), Styles.CursorRow.Render(label), archived)
		}
		rows = append(rows, line)
	}
	return boxed(title, focused, width, strings.Join(rows, "\n"))
}

// NameList renders a generic list of named items (e.g. tags).
func NameList(title string, names []string, cursor, width int, focused bool) string {
	if len(names) == 0 {
		return boxed(title, focused, width, Styles.Muted.Render("  nothing here yet"))
	}
	rows := make([]string, 0, len(names))
	for i, name := range names {
		label := truncate(name, width-6)
		line := "  " + label
		if i == cursor {
			line = Styles.Accent.Render("▸") + " " + Styles.CursorRow.Render(label)
		}
		rows = append(rows, line)
	}
	return boxed(title, focused, width, strings.Join(rows, "\n"))
}

func priorityMark(p string) string {
	switch p {
	case domain.PriorityUrgent:
		return Styles.Danger.Render("!!!")
	case domain.PriorityHigh:
		return Styles.Danger.Render("!!")
	case domain.PriorityMedium:
		return Styles.Warning.Render("!")
	case domain.PriorityLow:
		return Styles.Muted.Render("·")
	}
	return " "
}

func truncate(s string, max int) string {
	if max <= 3 {
		return s
	}
	if len(s) > max {
		return s[:max-1] + "…"
	}
	return s
}
