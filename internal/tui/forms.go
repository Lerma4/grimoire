package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/Lerma4/grimoire/internal/domain"
	"github.com/Lerma4/grimoire/internal/service"
	"github.com/Lerma4/grimoire/internal/tui/components"
)

// priorityOptions returns the huh options for task priority.
func priorityOptions(selected string) []huh.Option[string] {
	opts := []huh.Option[string]{}
	for _, p := range []string{domain.PriorityLow, domain.PriorityMedium, domain.PriorityHigh, domain.PriorityUrgent} {
		opt := huh.NewOption(p, p)
		if p == selected {
			opt = opt.Selected(true)
		}
		opts = append(opts, opt)
	}
	return opts
}

// openEditForm opens the appropriate form for the current selection.
func (m Model) openEditForm() (tea.Model, tea.Cmd) {
	if _, ok := m.selectedTask(); ok {
		return m.openTaskForm(false)
	}
	if _, ok := m.selectedNote(); ok {
		return m.openNoteForm(false)
	}
	m.setErr("nothing selected to edit")
	return m, nil
}

// openTaskForm opens a huh form to create or edit a task.
func (m Model) openTaskForm(create bool) (tea.Model, tea.Cmd) {
	var title, due, desc, prio string
	var targetID int64
	prio = domain.PriorityMedium
	if !create {
		t, ok := m.selectedTask()
		if !ok {
			m.setErr("no task selected")
			return m, nil
		}
		title, desc, prio, targetID = t.Title, t.Description, t.Priority, t.ID
		due = prettyDate(t.DueDate)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Title").Value(&title).Prompt(""),
			huh.NewSelect[string]().Title("Priority").Options(priorityOptions(prio)...).Value(&prio),
			huh.NewInput().Title("Due (YYYY-MM-DD)").Value(&due).Prompt(""),
			huh.NewText().Title("Description (markdown)").Value(&desc).Lines(5),
		),
	).WithWidth(formWidth(m.width)).WithShowHelp(true)

	m.form = form
	m.formApply = func(mm *Model) {
		d, _ := service.ParseDueDate(due)
		in := service.TaskInput{Title: title, Description: desc, Priority: prio, DueDate: d}
		if create {
			if _, err := mm.svc.Tasks.Create(service.Ctx, in); err != nil {
				mm.setErr(err.Error())
				return
			}
			mm.setStatus("task created")
		} else {
			if _, err := mm.svc.Tasks.Update(service.Ctx, targetID, in); err != nil {
				mm.setErr(err.Error())
				return
			}
			mm.setStatus("task updated")
		}
		mm.refresh()
	}
	return m, form.Init()
}

// openNoteForm opens a huh form to create or edit a note.
func (m Model) openNoteForm(create bool) (tea.Model, tea.Cmd) {
	var title, body string
	var targetID int64
	if !create {
		n, ok := m.selectedNote()
		if !ok {
			m.setErr("no note selected")
			return m, nil
		}
		title, body, targetID = n.Title, n.Body, n.ID
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Title").Value(&title).Prompt(""),
			huh.NewText().Title("Body (markdown)").Value(&body).Lines(12),
		),
	).WithWidth(formWidth(m.width)).WithShowHelp(true)

	m.form = form
	m.formApply = func(mm *Model) {
		in := service.NoteInput{Title: title, Body: body}
		if create {
			if _, err := mm.svc.Notes.Create(service.Ctx, in); err != nil {
				mm.setErr(err.Error())
				return
			}
			mm.setStatus("note created")
		} else {
			if _, err := mm.svc.Notes.Update(service.Ctx, targetID, in); err != nil {
				mm.setErr(err.Error())
				return
			}
			mm.setStatus("note updated")
		}
		mm.refresh()
	}
	return m, form.Init()
}

// openDeleteConfirm asks for confirmation before deleting the selected item.
func (m Model) openDeleteConfirm() (tea.Model, tea.Cmd) {
	var confirm bool
	label := "selected item"
	if t, ok := m.selectedTask(); ok {
		label = "task: " + t.Title
	} else if n, ok := m.selectedNote(); ok {
		label = "note: " + n.Title
	} else {
		m.setErr("nothing selected")
		return m, nil
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().Title(fmt.Sprintf("Delete %s?", label)).Value(&confirm),
		),
	).WithWidth(formWidth(m.width))

	m.form = form
	m.formApply = func(mm *Model) {
		if !confirm {
			mm.setStatus("cancelled")
			return
		}
		if t, ok := mm.selectedTask(); ok {
			if err := mm.svc.Tasks.Delete(service.Ctx, t.ID); err != nil {
				mm.setErr(err.Error())
				return
			}
		} else if n, ok := mm.selectedNote(); ok {
			if err := mm.svc.Notes.Delete(service.Ctx, n.ID); err != nil {
				mm.setErr(err.Error())
				return
			}
		}
		mm.setStatus("deleted")
		mm.refresh()
	}
	return m, form.Init()
}

// updateForm delegates to the active huh form and applies results on completion.
func (m Model) updateForm(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.form == nil {
		return m, nil
	}
	model, cmd := m.form.Update(msg)
	form, ok := model.(*huh.Form)
	if !ok {
		return m, cmd
	}
	m.form = form
	switch form.State {
	case huh.StateCompleted:
		apply := m.formApply
		m.formApply = nil
		m.form = nil
		m.mode = components.ModeNormal
		if apply != nil {
			apply(&m)
		}
	case huh.StateAborted:
		m.form = nil
		m.formApply = nil
		m.mode = components.ModeNormal
		m.setStatus("cancelled")
	}
	return m, cmd
}

// renderForm overlays the active form centered on screen.
func (m Model) renderForm() string {
	if m.form == nil {
		return ""
	}
	body := m.form.View()
	w := lipgloss.Width(body)
	maxW := m.width - 4
	if w > maxW {
		w = maxW
	}
	box := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(components.ColorAccent).
		BorderBackground(components.ColorBg).
		Background(components.ColorBg).
		Padding(0, 1).
		Width(w)
	rendered := box.Render(body)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, rendered)
}

func formWidth(w int) int {
	if w < 50 {
		return w - 2
	}
	if w > 80 {
		return 80
	}
	return w - 10
}

func prettyDate(s string) string {
	if len(s) >= 10 {
		return s[:10]
	}
	return s
}
