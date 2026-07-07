package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Lerma4/grimoire/internal/domain"
)

// Counts holds per-section item counts shown in the sidebar.
type Counts struct {
	Tasks, Notes, Today, Projects, Tags, Links, Archive int
}

// Sidebar renders the left navigation column.
func Sidebar(active Section, counts Counts, width int, focused bool) string {
	head := Styles.Title.Render("✦ Grimoire")
	items := []string{head, ""}

	render := func(s Section, n int) {
		marker := " "
		name := s.String()
		cnt := fmt.Sprintf("%d", n)
		if width < 18 {
			name = "" // collapse to glyph only
		}
		if s == active {
			line := fmt.Sprintf("%s %s", Styles.Accent.Render("▸"), s.Glyph())
			if name != "" {
				line += " " + name
			}
			line += " " + Styles.Muted.Render(cnt)
			items = append(items, Styles.Selected.Render(pad(line, width-2)))
			return
		}
		line := fmt.Sprintf("%s  %s", marker, s.Glyph())
		if name != "" {
			line += " " + name
		}
		line += " " + Styles.Muted.Render(cnt)
		items = append(items, Styles.Muted.Render(line))
	}

	render(SectionTasks, counts.Tasks)
	render(SectionNotes, counts.Notes)
	render(SectionToday, counts.Today)
	render(SectionProjects, counts.Projects)
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

func pad(s string, width int) string {
	plain := lipgloss.Width(s)
	if plain >= width {
		return s
	}
	return s + strings.Repeat(" ", width-plain)
}
