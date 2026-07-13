package tui

import (
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/Lerma4/grimoire/internal/store"
	"github.com/Lerma4/grimoire/internal/tui/components"
)

func newModel(t *testing.T) Model {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "tui.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	db.SetMaxOpenConns(1)
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := store.SeedIfEmpty(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	m := NewModel(Deps{DB: db, DBPath: "test.db"})
	m.width, m.height = 120, 40
	return m
}

func key(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func press(m Model, s string) Model {
	mm, _ := m.Update(key(s))
	return mm.(Model)
}

func pressKey(m Model, msg tea.KeyMsg) Model {
	mm, _ := m.Update(msg)
	return mm.(Model)
}

// TestSidebarNavDirection checks that j moves the selection DOWN the sidebar
// (Tasks -> Notes -> Today) and k moves it back UP.
func TestSidebarNavDirection(t *testing.T) {
	m := newModel(t)
	m = press(m, "h")
	if m.pane != components.PaneSidebar {
		t.Fatalf("expected sidebar focus, pane=%d", m.pane)
	}
	m = press(m, "j")
	if m.section != components.SectionNotes {
		t.Fatalf("j: expected SectionNotes, got %d", m.section)
	}
	m = press(m, "j")
	if m.section != components.SectionToday {
		t.Fatalf("j: expected SectionToday, got %d", m.section)
	}
	m = press(m, "k")
	if m.section != components.SectionNotes {
		t.Fatalf("k: expected SectionNotes, got %d", m.section)
	}
	m = pressKey(m, tea.KeyMsg{Type: tea.KeyDown})
	if m.section != components.SectionToday {
		t.Fatalf("down: expected SectionToday, got %d", m.section)
	}
	m = pressKey(m, tea.KeyMsg{Type: tea.KeyUp})
	if m.section != components.SectionNotes {
		t.Fatalf("up: expected SectionNotes, got %d", m.section)
	}
}

// TestNotesSectionRenders lands on Notes via the sidebar and renders.
func TestNotesSectionRenders(t *testing.T) {
	m := newModel(t)
	m = press(m, "h")
	m = press(m, "j")
	if m.section != components.SectionNotes {
		t.Fatalf("expected notes section, got %d", m.section)
	}
	if out := m.View(); out == "" {
		t.Fatal("empty view on notes section")
	}
}

// TestEscBacksOutOfPanes checks that Esc exits the detail pane (opened via
// Enter) back to the list, and the list back to the sidebar.
func TestEscBacksOutOfPanes(t *testing.T) {
	m := newModel(t)
	m = press(m, "enter")
	if m.pane != components.PaneDetail {
		t.Fatalf("enter: expected detail pane, got %d", m.pane)
	}
	m = press(m, "esc")
	if m.pane != components.PaneList {
		t.Fatalf("esc from detail: expected list pane, got %d", m.pane)
	}
	m = press(m, "esc")
	if m.pane != components.PaneSidebar {
		t.Fatalf("esc from list: expected sidebar pane, got %d", m.pane)
	}
}

// TestAllSectionsRender walks every section and renders it.
func TestAllSectionsRender(t *testing.T) {
	m := newModel(t)
	for _, s := range components.AllSections {
		m.section = s
		m.cursor = 0
		m.refresh()
		_ = m.View()
	}
}
