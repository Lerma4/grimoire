package tui

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"

	"github.com/Lerma4/grimoire/internal/domain"
	"github.com/Lerma4/grimoire/internal/service"
	"github.com/Lerma4/grimoire/internal/tui/components"
)

// Deps are the runtime dependencies injected into the TUI.
type Deps struct {
	DB     *sql.DB
	DBPath string
}

// Model is the top-level Bubble Tea model for Grimoire.
type Model struct {
	svc *service.Services

	width, height int

	section components.Section
	pane    components.Pane
	mode    components.Mode

	pendingG bool // waiting for second 'g' in 'gg'

	// displayed lists, filtered by the active section
	tasks    []domain.Task
	notes    []domain.Note
	projects []domain.Project
	tags     []domain.Tag
	cursor   int

	// context for the detail panel of the current selection
	selTags        []domain.Tag
	selLinkedNotes []domain.Note
	selLinkedTasks []domain.Task

	cmdInput    textinput.Model
	searchInput textinput.Model

	form      *huh.Form
	formApply func(*Model)

	statusMsg string
	errMsg    string
	errMode   bool

	dbPath string

	quit bool
}

// NewModel constructs the model and loads the initial dataset.
func NewModel(deps Deps) Model {
	ti := textinput.New()
	ti.Prompt = ""
	ti.CharLimit = 200

	si := textinput.New()
	si.Prompt = ""
	si.CharLimit = 200

	m := Model{
		svc:         service.NewServices(deps.DB),
		mode:        components.ModeNormal,
		pane:        components.PaneList,
		section:     components.SectionTasks,
		cmdInput:    ti,
		searchInput: si,
		dbPath:      deps.DBPath,
	}
	m.refresh()
	return m
}

// Init starts the program (no async work needed; data is loaded sync).
func (m Model) Init() tea.Cmd {
	return nil
}

// refresh reloads the lists appropriate to the current section.
func (m *Model) refresh() {
	ctx := service.Ctx
	switch m.section {
	case components.SectionTasks:
		t, _ := m.svc.Tasks.List(ctx, service.TaskFilter{})
		m.tasks = t
		m.notes = nil
	case components.SectionNotes:
		n, _ := m.svc.Notes.List(ctx, service.NoteFilter{})
		m.notes = n
		m.tasks = nil
	case components.SectionToday:
		t, _ := m.svc.Tasks.List(ctx, service.TaskFilter{Overdue: true})
		// include tasks due today
		m.tasks = t
		m.notes = nil
	case components.SectionArchive:
		t, _ := m.svc.Tasks.List(ctx, service.TaskFilter{Status: domain.StatusArchived})
		n, _ := m.svc.Notes.List(ctx, service.NoteFilter{IncludeArchived: true})
		m.tasks = t
		m.notes = filterArchivedNotes(n)
	case components.SectionProjects:
		p, _ := m.svc.Projects.List(ctx)
		m.projects = p
	case components.SectionTags:
		tg, _ := m.svc.Tags.List(ctx)
		m.tags = tg
	case components.SectionLinks:
		// links section: show tasks that have linked notes as a jumping-off point
		t, _ := m.svc.Tasks.List(ctx, service.TaskFilter{})
		m.tasks = t
		m.notes = nil
	}
	if m.cursor > 0 && m.cursor >= m.listLen() {
		m.cursor = max(m.listLen()-1, 0)
	}
	m.loadDetail()
}

func filterArchivedNotes(n []domain.Note) []domain.Note {
	var out []domain.Note
	for _, x := range n {
		if x.ArchivedAt != "" {
			out = append(out, x)
		}
	}
	return out
}

// loadDetail fetches tags and linked items for the current selection.
func (m *Model) loadDetail() {
	m.selTags = nil
	m.selLinkedNotes = nil
	m.selLinkedTasks = nil
	ctx := service.Ctx
	if t, ok := m.selectedTask(); ok {
		m.selTags, _ = m.svc.Tasks.Tags(ctx, t.ID)
		m.selLinkedNotes, _ = m.svc.Links.NotesForTask(ctx, t.ID)
		return
	}
	if n, ok := m.selectedNote(); ok {
		m.selTags, _ = m.svc.Notes.Tags(ctx, n.ID)
		m.selLinkedTasks, _ = m.svc.Links.TasksForNote(ctx, n.ID)
	}
}

// listLen returns the number of rows in the active center list.
func (m Model) listLen() int {
	switch m.section {
	case components.SectionNotes:
		return len(m.notes)
	case components.SectionProjects:
		return len(m.projects)
	case components.SectionTags:
		return len(m.tags)
	default:
		return len(m.tasks)
	}
}

func (m Model) selectedTask() (domain.Task, bool) {
	if m.section != components.SectionTasks && m.section != components.SectionToday &&
		m.section != components.SectionArchive && m.section != components.SectionLinks {
		return domain.Task{}, false
	}
	if m.cursor < 0 || m.cursor >= len(m.tasks) {
		return domain.Task{}, false
	}
	return m.tasks[m.cursor], true
}

func (m Model) selectedNote() (domain.Note, bool) {
	if m.section != components.SectionNotes && m.section != components.SectionArchive {
		return domain.Note{}, false
	}
	if m.cursor < 0 || m.cursor >= len(m.notes) {
		return domain.Note{}, false
	}
	return m.notes[m.cursor], true
}

func (m Model) counts() components.Counts {
	ctx := service.Ctx
	c := components.Counts{}
	if ts, err := m.svc.Tasks.List(ctx, service.TaskFilter{}); err == nil {
		c.Tasks = len(ts)
	}
	if ns, err := m.svc.Notes.List(ctx, service.NoteFilter{}); err == nil {
		c.Notes = len(ns)
	}
	if ts, err := m.svc.Tasks.List(ctx, service.TaskFilter{Overdue: true}); err == nil {
		c.Today = len(ts)
	}
	if ps, err := m.svc.Projects.List(ctx); err == nil {
		c.Projects = len(ps)
	}
	if tg, err := m.svc.Tags.List(ctx); err == nil {
		c.Tags = len(tg)
	}
	c.Links = countLinks(m)
	if at, err := m.svc.Tasks.List(ctx, service.TaskFilter{Status: domain.StatusArchived}); err == nil {
		c.Archive = len(at)
	}
	return c
}

func countLinks(m Model) int {
	// ponytail: O(tasks) scan for a count; fine for a local single-user DB.
	n := 0
	for _, t := range m.tasks {
		if ns, err := m.svc.Links.NotesForTask(service.Ctx, t.ID); err == nil {
			n += len(ns)
		}
	}
	return n
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
