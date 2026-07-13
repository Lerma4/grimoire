package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Lerma4/grimoire/internal/service"
	"github.com/Lerma4/grimoire/internal/tui/components"
)

// Update routes messages to the active mode.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok && km.String() == "ctrl+c" {
		m.quit = true
		return m, tea.Quit
	}
	// huh form takes over when active.
	if m.form != nil {
		return m.updateForm(msg)
	}
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.MouseMsg:
		return m.updateMouse(msg)
	case tea.KeyMsg:
		if m.mode == components.ModeHelp {
			m.mode = components.ModeNormal
			return m, nil
		}
		switch m.mode {
		case components.ModeNormal:
			return m.updateNormal(msg)
		case components.ModeCommand, components.ModeSearch:
			return m.updatePrompt(msg)
		}
	}
	return m, nil
}

func (m Model) updateNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q":
		m.quit = true
		return m, tea.Quit
	case "esc":
		m.errMsg = ""
		m.statusMsg = ""
		// ponytail: Esc backs out of the focused pane (reverse of Tab/Enter).
		m.pane = panePrev(m.pane)
	case "j", "down":
		if m.pane == components.PaneSidebar {
			m.section = nextSection(m.section, 1)
			m.cursor = 0
			m.refresh()
		} else if m.cursor < m.listLen()-1 {
			m.cursor++
			m.loadDetail()
		}
	case "k", "up":
		if m.pane == components.PaneSidebar {
			m.section = nextSection(m.section, -1)
			m.cursor = 0
			m.refresh()
		} else if m.cursor > 0 {
			m.cursor--
			m.loadDetail()
		}
	case "g":
		if m.pendingG {
			m.cursor = 0
			m.pendingG = false
			m.loadDetail()
		} else {
			m.pendingG = true
		}
		return m, nil
	case "G":
		m.cursor = max(m.listLen()-1, 0)
		m.loadDetail()
	case "h":
		m.pane = panePrev(m.pane)
	case "l":
		m.pane = paneNext(m.pane)
	case "tab":
		m.pane = paneNext(m.pane)
	case "enter":
		m.pane = components.PaneDetail
	case "t":
		m.setSection(components.SectionTasks)
	case "m":
		m.setSection(components.SectionNotes)
	case "#":
		m.setSection(components.SectionTags)
	case " ":
		m.toggleDone()
	case "n":
		m.setStatus("next result: use / then n")
	case "?":
		m.mode = components.ModeHelp
	case ":":
		return m.enterCommand("")
	case "/":
		return m.enterSearch()
	case "a":
		return m.openTaskForm(true)
	case "A":
		return m.openNoteForm(true)
	case "e":
		return m.openEditForm()
	case "d":
		m.archiveSelected()
	case "D":
		return m.openDeleteConfirm()
	case "L":
		return m.enterCommand("link ")
	case "U":
		return m.enterCommand("unlink ")
	}
	m.pendingG = false
	return m, nil
}

// updateMouse maps clicks/scroll to panes and list rows. Only acts in normal
// mode so command/search/form input keeps the keyboard.
func (m Model) updateMouse(msg tea.MouseMsg) (tea.Model, tea.Cmd) {
	if m.mode != components.ModeNormal {
		return m, nil
	}
	// ponytail: geometry tied to boxed() chrome (border+title+sep=3) + 1-line header.
	const headerH, boxChrome = 1, 3
	sideW, listW, _ := m.columns()
	switch {
	case sideW > 0 && msg.X < sideW:
		// sidebar: sections start 2 rows into the body (Grimoire title + blank).
		idx := msg.Y - (headerH + boxChrome) - 2
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft &&
			idx >= 0 && idx < len(components.AllSections) {
			m.setSection(components.AllSections[idx])
		}
	case msg.X < sideW+listW:
		row := msg.Y - (headerH + boxChrome)
		n := m.listLen()
		switch {
		case msg.Button == tea.MouseButtonWheelUp && m.cursor > 0:
			m.cursor--
			m.loadDetail()
		case msg.Button == tea.MouseButtonWheelDown && m.cursor < n-1:
			m.cursor++
			m.loadDetail()
		case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft &&
			row >= 0 && row < n:
			m.cursor = row
			m.pane = components.PaneList
			m.loadDetail()
		}
	default:
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			m.pane = components.PaneDetail
		}
	}
	return m, nil
}

func (m Model) updatePrompt(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		m.mode = components.ModeNormal
		m.cmdInput.Blur()
		m.searchInput.Blur()
		return m, nil
	}
	switch m.mode {
	case components.ModeCommand:
		if msg.String() == "enter" {
			return m.runCommand(m.cmdInput.Value())
		}
		var cmd tea.Cmd
		m.cmdInput, cmd = m.cmdInput.Update(msg)
		return m, cmd
	case components.ModeSearch:
		if msg.String() == "enter" {
			m.runSearch(m.searchInput.Value())
			m.mode = components.ModeNormal
			return m, nil
		}
		var cmd tea.Cmd
		m.searchInput, cmd = m.searchInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

// enterCommand switches to command mode, optionally pre-filled.
func (m Model) enterCommand(prefill string) (tea.Model, tea.Cmd) {
	m.mode = components.ModeCommand
	m.cmdInput.Reset()
	m.cmdInput.SetValue(prefill)
	m.cmdInput.CursorEnd()
	m.cmdInput.Focus()
	return m, m.cmdInput.Cursor.BlinkCmd() //nolint:staticcheck
}

// enterSearch switches to search mode.
func (m Model) enterSearch() (tea.Model, tea.Cmd) {
	m.mode = components.ModeSearch
	m.searchInput.Reset()
	m.searchInput.Focus()
	return m, m.searchInput.Cursor.BlinkCmd() //nolint:staticcheck
}

func (m *Model) setSection(s components.Section) {
	m.section = s
	m.pane = components.PaneList
	m.cursor = 0
	m.refresh()
}

func (m *Model) toggleDone() {
	t, ok := m.selectedTask()
	if !ok {
		return
	}
	if _, err := m.svc.Tasks.ToggleDone(service.Ctx, t.ID); err != nil {
		m.setErr("toggle: " + err.Error())
		return
	}
	m.setStatus("toggled")
	m.refresh()
}

func (m *Model) archiveSelected() {
	if t, ok := m.selectedTask(); ok {
		if err := m.svc.Tasks.Archive(service.Ctx, t.ID); err != nil {
			m.setErr(err.Error())
			return
		}
		m.setStatus("archived")
		m.refresh()
		return
	}
	if n, ok := m.selectedNote(); ok {
		if err := m.svc.Notes.Archive(service.Ctx, n.ID); err != nil {
			m.setErr(err.Error())
			return
		}
		m.setStatus("archived")
		m.refresh()
	}
}

func (m *Model) deleteSelected() {
	if t, ok := m.selectedTask(); ok {
		if err := m.svc.Tasks.Delete(service.Ctx, t.ID); err != nil {
			m.setErr(err.Error())
			return
		}
		m.setStatus("deleted")
		m.refresh()
		return
	}
	if n, ok := m.selectedNote(); ok {
		if err := m.svc.Notes.Delete(service.Ctx, n.ID); err != nil {
			m.setErr(err.Error())
			return
		}
		m.setStatus("deleted")
		m.refresh()
	}
}

func (m *Model) setStatus(s string) { m.statusMsg = s; m.errMsg = ""; m.errMode = false }
func (m *Model) setErr(s string)    { m.errMsg = s; m.errMode = true; m.statusMsg = "" }

func nextSection(s components.Section, delta int) components.Section {
	idx := 0
	for i, x := range components.AllSections {
		if x == s {
			idx = i
			break
		}
	}
	idx += delta
	n := len(components.AllSections)
	if idx < 0 {
		idx = 0
	}
	if idx >= n {
		idx = n - 1
	}
	return components.AllSections[idx]
}

func panePrev(p components.Pane) components.Pane {
	if p == components.PaneDetail {
		return components.PaneList
	}
	return components.PaneSidebar
}

func paneNext(p components.Pane) components.Pane {
	switch p {
	case components.PaneSidebar:
		return components.PaneList
	case components.PaneList:
		return components.PaneDetail
	}
	return components.PaneDetail
}
