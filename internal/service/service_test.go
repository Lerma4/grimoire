package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Lerma4/grimoire/internal/domain"
	"github.com/Lerma4/grimoire/internal/store"

	_ "modernc.org/sqlite"
)

// newServices builds the full service layer over a migrated temp database.
func newServices(t *testing.T) *Services {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "svc.db"))
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	db.SetMaxOpenConns(1)
	if err := store.Migrate(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	t.Cleanup(func() { db.Close() })
	return NewServices(db)
}

func TestTaskService_Lifecycle(t *testing.T) {
	s := newServices(t)
	ctx := context.Background()

	if _, err := s.Tasks.Create(ctx, TaskInput{Title: "  "}); err == nil {
		t.Fatal("expected empty-title validation error")
	}

	t1, err := s.Tasks.Create(ctx, TaskInput{Title: "ship it", Priority: domain.PriorityHigh})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if t1.Status != domain.StatusTodo || t1.Priority != domain.PriorityHigh {
		t.Fatalf("defaults wrong: %+v", t1)
	}

	if _, err := s.Tasks.Update(ctx, t1.ID, TaskInput{Title: "ship it now"}); err != nil {
		t.Fatalf("update: %v", err)
	}

	if err := s.Tasks.SetStatus(ctx, t1.ID, "bogus"); err == nil {
		t.Fatal("expected invalid-status error")
	}
	if err := s.Tasks.SetStatus(ctx, t1.ID, domain.StatusDoing); err != nil {
		t.Fatalf("setstatus doing: %v", err)
	}

	t1, _ = s.Tasks.Get(ctx, t1.ID)
	if t1.Status != domain.StatusDoing {
		t.Fatalf("status not applied: %+v", t1)
	}

	// toggle to done, then back.
	_, _ = s.Tasks.ToggleDone(ctx, t1.ID)
	t1, _ = s.Tasks.Get(ctx, t1.ID)
	if t1.Status != domain.StatusDone || t1.CompletedAt == "" {
		t.Fatalf("toggle to done failed: %+v", t1)
	}
	_, _ = s.Tasks.ToggleDone(ctx, t1.ID)
	t1, _ = s.Tasks.Get(ctx, t1.ID)
	if t1.Status != domain.StatusTodo || t1.CompletedAt != "" {
		t.Fatalf("toggle back failed: %+v", t1)
	}

	if _, err := s.Tasks.SetProject(ctx, t1.ID, "Release"); err != nil {
		t.Fatalf("setproject: %v", err)
	}
	t1, _ = s.Tasks.Get(ctx, t1.ID)
	projs, _ := s.Projects.List(ctx)
	if t1.ProjectID == 0 || len(projs) != 1 {
		t.Fatalf("project not set/created: %+v projs=%d", t1, len(projs))
	}

	if err := s.Tasks.TagByName(ctx, t1.ID, "urgent"); err != nil {
		t.Fatalf("tag: %v", err)
	}
	tags, _ := s.Tasks.Tags(ctx, t1.ID)
	if len(tags) != 1 || tags[0].Name != "urgent" {
		t.Fatalf("tags = %+v", tags)
	}

	if err := s.Tasks.Archive(ctx, t1.ID); err != nil {
		t.Fatalf("archive: %v", err)
	}
	active, _ := s.Tasks.List(ctx, TaskFilter{})
	if len(active) != 0 {
		t.Fatalf("archived task should be hidden, got %d", len(active))
	}
	archived, _ := s.Tasks.List(ctx, TaskFilter{Status: domain.StatusArchived})
	if len(archived) != 1 {
		t.Fatalf("archived list len=%d", len(archived))
	}
}

func TestNoteService_Lifecycle(t *testing.T) {
	s := newServices(t)
	ctx := context.Background()

	if _, err := s.Notes.Create(ctx, NoteInput{Title: ""}); err == nil {
		t.Fatal("expected empty-title error")
	}
	n, err := s.Notes.Create(ctx, NoteInput{Title: "plan", Body: "# body"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := s.Notes.Update(ctx, n.ID, NoteInput{Title: "plan v2", Body: "## x"}); err != nil {
		t.Fatalf("update: %v", err)
	}
	if err := s.Notes.Archive(ctx, n.ID); err != nil {
		t.Fatalf("archive: %v", err)
	}
	active, _ := s.Notes.List(ctx, NoteFilter{})
	if len(active) != 0 {
		t.Fatalf("archived note should be hidden")
	}
	if err := s.Notes.Unarchive(ctx, n.ID); err != nil {
		t.Fatalf("unarchive: %v", err)
	}
	if _, err := s.Notes.SetProject(ctx, n.ID, "Docs"); err != nil {
		t.Fatalf("setproject: %v", err)
	}
	if err := s.Notes.TagByName(ctx, n.ID, "draft"); err != nil {
		t.Fatalf("tag: %v", err)
	}
}

func TestLinkService_AndSearch(t *testing.T) {
	s := newServices(t)
	ctx := context.Background()

	task, _ := s.Tasks.Create(ctx, TaskInput{Title: "investigate"})
	note, _ := s.Notes.Create(ctx, NoteInput{Title: "research notes", Body: "see #investigate"})

	// linking to a nonexistent note must fail validation.
	if err := s.Links.LinkTaskNote(ctx, task.ID, note.ID+9999); err == nil {
		t.Fatal("expected error linking to missing note")
	}
	if err := s.Links.LinkTaskNote(ctx, task.ID, note.ID); err != nil {
		t.Fatalf("link: %v", err)
	}

	linked, _ := s.Links.NotesForTask(ctx, task.ID)
	if len(linked) != 1 || linked[0].ID != note.ID {
		t.Fatalf("notes for task = %+v", linked)
	}
	back, _ := s.Links.TasksForNote(ctx, note.ID)
	if len(back) != 1 || back[0].ID != task.ID {
		t.Fatalf("tasks for note = %+v", back)
	}

	// search finds both.
	res, err := s.Search.Search(ctx, "investigate")
	if err != nil {
		t.Fatalf("search: %v", err)
	}
	if len(res.Tasks) != 1 || len(res.Notes) != 1 {
		t.Fatalf("search buckets = tasks:%d notes:%d", len(res.Tasks), len(res.Notes))
	}

	if err := s.Links.UnlinkTaskNote(ctx, task.ID, note.ID); err != nil {
		t.Fatalf("unlink: %v", err)
	}
	linked, _ = s.Links.NotesForTask(ctx, task.ID)
	if len(linked) != 0 {
		t.Fatalf("expected unlinked, got %d", len(linked))
	}
}

func TestProjectService_FindOrCreate(t *testing.T) {
	s := newServices(t)
	ctx := context.Background()

	p1, err := s.Projects.FindOrCreate(ctx, "Solo")
	if err != nil {
		t.Fatalf("findorcreate: %v", err)
	}
	p2, err := s.Projects.FindOrCreate(ctx, "Solo")
	if err != nil {
		t.Fatalf("findorcreate 2: %v", err)
	}
	if p1.ID != p2.ID {
		t.Fatalf("findorcreate should return same project: %d != %d", p1.ID, p2.ID)
	}
}

func TestParseDueDate(t *testing.T) {
	cases := []struct{ in, wantPrefix string }{
		{"", ""},
		{"2026-07-07", "2026-07-07"},
	}
	for _, c := range cases {
		got, err := ParseDueDate(c.in)
		if err != nil {
			t.Fatalf("parse %q: %v", c.in, err)
		}
		if c.wantPrefix == "" {
			if got != "" {
				t.Fatalf("expected empty, got %q", got)
			}
			continue
		}
		if len(got) < len(c.wantPrefix) || got[:len(c.wantPrefix)] != c.wantPrefix {
			t.Fatalf("parse %q -> %q want prefix %q", c.in, got, c.wantPrefix)
		}
	}
	if _, err := ParseDueDate("not-a-date"); err == nil {
		t.Fatal("expected error for invalid date")
	}
}
