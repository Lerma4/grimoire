package tui

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/Lerma4/grimoire/internal/domain"
	"github.com/Lerma4/grimoire/internal/service"
	"github.com/Lerma4/grimoire/internal/tui/components"
)

// textinputCmd is a no-op blink starter; the cursor stays steady but the input
// still receives keys via Update.
func textinputCmd(_ textinput.Model) tea.Cmd { return nil }

// runCommand parses and executes a ':' command.
func (m Model) runCommand(line string) (tea.Model, tea.Cmd) {
	m.mode = components.ModeNormal
	m.cmdInput.Blur()
	line = strings.TrimSpace(line)
	if line == "" {
		return m, nil
	}
	if err := m.execCommand(line); err != nil {
		m.setErr(err.Error())
	}
	m.refresh()
	return m, nil
}

// execCommand dispatches a single parsed command.
func (m *Model) execCommand(line string) error {
	parts := splitFields(line)
	if len(parts) == 0 {
		return nil
	}
	cmd := parts[0]
	ctx := service.Ctx

	switch cmd {
	case "q", "quit":
		m.quit = true
		return nil
	case "task":
		if len(parts) >= 3 && parts[1] == "add" {
			title := strings.Join(parts[2:], " ")
			if _, err := m.svc.Tasks.Create(ctx, service.TaskInput{Title: title}); err != nil {
				return err
			}
			m.setStatus("task added")
			return nil
		}
	case "note":
		if len(parts) >= 3 && parts[1] == "add" {
			title := strings.Join(parts[2:], " ")
			if _, err := m.svc.Notes.Create(ctx, service.NoteInput{Title: title}); err != nil {
				return err
			}
			m.setStatus("note added")
			return nil
		}
	case "done":
		if t, ok := m.selectedTask(); ok {
			return m.svc.Tasks.SetStatus(ctx, t.ID, domain.StatusDone)
		}
	case "doing":
		if t, ok := m.selectedTask(); ok {
			return m.svc.Tasks.SetStatus(ctx, t.ID, domain.StatusDoing)
		}
	case "todo":
		if t, ok := m.selectedTask(); ok {
			return m.svc.Tasks.SetStatus(ctx, t.ID, domain.StatusTodo)
		}
	case "archive":
		m.archiveSelected()
		return nil
	case "delete":
		m.deleteSelected()
		return nil
	case "link":
		return m.linkByID(parts[1:])
	case "unlink":
		return m.unlinkByID(parts[1:])
	case "project":
		name := strings.Join(parts[1:], " ")
		return m.setProjectSelected(name)
	case "tag":
		name := strings.Join(parts[1:], " ")
		return m.tagSelected(name)
	}
	return errUnknownCommand(cmd)
}

// linkByID links the selected item to the given numeric id (task→note or note→task).
func (m *Model) linkByID(args []string) error {
	if len(args) == 0 {
		return errMissingArg("id")
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return errInvalidArg("id")
	}
	ctx := service.Ctx
	if t, ok := m.selectedTask(); ok {
		if err := m.svc.Links.LinkTaskNote(ctx, t.ID, id); err != nil {
			return err
		}
		m.setStatus("linked")
		return nil
	}
	if n, ok := m.selectedNote(); ok {
		if err := m.svc.Links.LinkTaskNote(ctx, id, n.ID); err != nil {
			return err
		}
		m.setStatus("linked")
		return nil
	}
	return errNothingSelected
}

func (m *Model) unlinkByID(args []string) error {
	if len(args) == 0 {
		return errMissingArg("id")
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return errInvalidArg("id")
	}
	ctx := service.Ctx
	if t, ok := m.selectedTask(); ok {
		if err := m.svc.Links.UnlinkTaskNote(ctx, t.ID, id); err != nil {
			return err
		}
		m.setStatus("unlinked")
		return nil
	}
	if n, ok := m.selectedNote(); ok {
		if err := m.svc.Links.UnlinkTaskNote(ctx, id, n.ID); err != nil {
			return err
		}
		m.setStatus("unlinked")
		return nil
	}
	return errNothingSelected
}

func (m *Model) setProjectSelected(name string) error {
	ctx := service.Ctx
	if t, ok := m.selectedTask(); ok {
		if _, err := m.svc.Tasks.SetProject(ctx, t.ID, name); err != nil {
			return err
		}
		m.setStatus("project set")
		return nil
	}
	if n, ok := m.selectedNote(); ok {
		if _, err := m.svc.Notes.SetProject(ctx, n.ID, name); err != nil {
			return err
		}
		m.setStatus("project set")
		return nil
	}
	return errNothingSelected
}

func (m *Model) tagSelected(name string) error {
	ctx := service.Ctx
	if t, ok := m.selectedTask(); ok {
		if err := m.svc.Tasks.TagByName(ctx, t.ID, name); err != nil {
			return err
		}
		m.setStatus("tagged")
		return nil
	}
	if n, ok := m.selectedNote(); ok {
		if err := m.svc.Notes.TagByName(ctx, n.ID, name); err != nil {
			return err
		}
		m.setStatus("tagged")
		return nil
	}
	return errNothingSelected
}

// runSearch applies the query to the active list type.
func (m *Model) runSearch(q string) {
	ctx := service.Ctx
	q = strings.TrimSpace(q)
	switch m.section {
	case components.SectionNotes:
		ns, _ := m.svc.Notes.List(ctx, service.NoteFilter{Search: q})
		m.notes = ns
		m.tasks = nil
	default:
		ts, _ := m.svc.Tasks.List(ctx, service.TaskFilter{Search: q})
		m.tasks = ts
		m.notes = nil
	}
	m.cursor = 0
	m.loadDetail()
	if q != "" {
		m.setStatus("/" + q)
	}
}

// splitFields splits on whitespace but keeps it simple.
func splitFields(s string) []string {
	return strings.Fields(s)
}
