package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Lerma4/grimoire/internal/tui/components"
)

// View renders the whole interface.
func (m Model) View() string {
	if m.width == 0 {
		return "loading grimoire…"
	}
	if m.form != nil {
		return m.renderForm()
	}
	if m.mode == components.ModeHelp {
		return m.helpView()
	}

	headerH := 1
	footerH := 1
	bodyH := m.height - headerH - footerH
	if bodyH < 3 {
		bodyH = 3
	}

	sideW, listW, detailW := m.columns()

	sidebar := components.Sidebar(m.section, m.counts(), sideW, m.pane == components.PaneSidebar)
	center := m.renderList(listW)
	detail := m.renderDetail(detailW)

	// stack columns to equal height
	sideH, listH, detH := lineCount(sidebar), lineCount(center), lineCount(detail)
	contentH := max(max(sideH, listH), detH)
	if contentH < bodyH {
		contentH = bodyH
	}
	sidebar = padHeight(sidebar, contentH)
	center = padHeight(center, contentH)
	detail = padHeight(detail, contentH)

	var row string
	switch {
	case sideW == 0:
		row = lipgloss.JoinHorizontal(lipgloss.Top, center, detail)
	default:
		row = lipgloss.JoinHorizontal(lipgloss.Top, sidebar, center, detail)
	}

	header := components.Header(m.section.String(), m.currentFilter(), m.dbLabel(), m.width)
	footer := m.renderFooter()

	return lipgloss.JoinVertical(lipgloss.Left, header, row, footer)
}

// columns returns (sidebar, list, detail) widths for the current terminal.
func (m Model) columns() (int, int, int) {
	W := m.width
	const min = 8
	switch {
	case W < 60:
		// two columns: list + detail, sidebar hidden
		list := W * 45 / 100
		if list < min {
			list = min
		}
		detail := W - list
		if detail < min {
			detail = min
		}
		return 0, list, detail
	case W < 90:
		side := 10
		rest := W - side
		list := rest * 50 / 100
		detail := rest - list
		return side, list, detail
	default:
		side := 18
		rest := W - side
		detail := rest * 42 / 100
		list := rest - detail
		return side, list, detail
	}
}

func (m Model) renderList(width int) string {
	focused := m.pane == components.PaneList
	title := m.section.String()
	switch m.section {
	case components.SectionNotes:
		return components.NoteList(title, m.notes, m.cursor, width, focused)
	case components.SectionTags:
		names := make([]string, len(m.tags))
		for i, t := range m.tags {
			names[i] = t.Name
		}
		return components.NameList(title, names, m.cursor, width, focused)
	default:
		return components.TaskList(title, m.tasks, m.cursor, width, focused)
	}
}

func (m Model) renderDetail(width int) string {
	focused := m.pane == components.PaneDetail
	if t, ok := m.selectedTask(); ok {
		return components.TaskDetail(t, m.selTags, m.selLinkedNotes, width, focused)
	}
	if n, ok := m.selectedNote(); ok {
		preview := components.RenderMarkdown(n.Body, width-4)
		return components.NoteDetail(n, m.selTags, m.selLinkedTasks, preview, width, focused)
	}
	return emptyDetail(width, focused)
}

func emptyDetail(width int, focused bool) string {
	body := components.Styles.Muted.Render("  select an item to see its detail.")
	return components.Box("Detail", focused, width, body)
}

func (m Model) renderFooter() string {
	switch m.mode {
	case components.ModeCommand:
		return components.CommandBar(m.cmdInput.View(), m.width)
	case components.ModeSearch:
		return components.SearchBar(m.searchInput.View(), m.width)
	default:
		hint := "j/k move · h/l panes · : cmd · / search · ? help · q quit"
		msg := m.statusMsg
		if m.errMode {
			msg = m.errMsg
		}
		return components.StatusBar(m.mode.String(), hint, msg, m.errMode, m.width)
	}
}

func (m Model) currentFilter() string {
	if strings.HasPrefix(m.statusMsg, "/") {
		return m.statusMsg
	}
	if m.section == components.SectionToday {
		return "today/overdue"
	}
	if m.section == components.SectionArchive {
		return "archived"
	}
	return ""
}

func (m Model) dbLabel() string {
	if m.dbPath == "" {
		return "db: ready"
	}
	return "db: " + shortPath(m.dbPath)
}

func shortPath(p string) string {
	if i := strings.LastIndex(p, "/"); i >= 0 {
		return p[i+1:]
	}
	return p
}

func (m Model) helpView() string {
	title := components.Styles.Title.Render("✦ Grimoire — keybindings")
	rows := [][2]string{
		{"j / k", "next / previous item"},
		{"h / l", "previous / next pane"},
		{"gg / G", "first / last item"},
		{"/ + n / N", "search, next / prev result"},
		{"Enter", "focus detail"},
		{"Space", "toggle task done"},
		{"a / A", "new task / note (via :task add / :note add)"},
		{"e", "edit selected"},
		{"d / D", "archive / delete"},
		{"t m #", "Tasks / Notes / Tags"},
		{"L / U", "link / unlink task↔note (:link/:unlink)"},
		{":", "command mode"},
		{"?", "this help"},
		{"q", "quit"},
	}
	var b strings.Builder
	b.WriteString(title + "\n\n")
	for _, r := range rows {
		b.WriteString("  " + components.Styles.HelpKey.Render(padRight(r[0], 12)) + components.Styles.HelpDesc.Render(r[1]) + "\n")
	}
	b.WriteString("\n" + components.Styles.Muted.Render("press any key to close"))
	return components.BoxCenter("Help", m.width, b.String())
}

// helpers ------------------------------------------------------------

func lineCount(s string) int {
	if s == "" {
		return 0
	}
	return strings.Count(s, "\n") + 1
}

func padHeight(s string, h int) string {
	n := lineCount(s)
	if n >= h {
		return s
	}
	return s + strings.Repeat("\n", h-n)
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}
