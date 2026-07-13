package store

import (
	"context"
	"testing"

	"github.com/Lerma4/grimoire/internal/domain"
)

func TestMigrate_AndSeed(t *testing.T) {
	db := newTestDB(t)

	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations`).Scan(&n); err != nil {
		t.Fatalf("read migrations: %v", err)
	}
	if n == 0 {
		t.Fatal("expected at least one applied migration")
	}

	// seed inserts exactly once.
	if err := SeedIfEmpty(db); err != nil {
		t.Fatalf("seed: %v", err)
	}
	var tasks int
	if err := db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&tasks); err != nil {
		t.Fatalf("count tasks: %v", err)
	}
	if tasks == 0 {
		t.Fatal("seed should create tasks")
	}
	before := tasks
	if err := SeedIfEmpty(db); err != nil {
		t.Fatalf("re-seed: %v", err)
	}
	if err := db.QueryRow(`SELECT COUNT(*) FROM tasks`).Scan(&tasks); err != nil {
		t.Fatalf("recount tasks: %v", err)
	}
	if tasks != before {
		t.Fatalf("seed should be idempotent; got %d want %d", tasks, before)
	}
}

func TestTaskRepo_CRUDAndFilter(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	r := NewTaskRepo(db)

	t1, err := r.Create(ctx, domain.Task{Title: "write tests", Priority: domain.PriorityHigh})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if t1.ID == 0 || t1.Status != domain.StatusTodo {
		t.Fatalf("unexpected task: %+v", t1)
	}

	got, err := r.Get(ctx, t1.ID)
	if err != nil || got.Title != "write tests" {
		t.Fatalf("get: %v %+v", err, got)
	}

	if _, err := r.Create(ctx, domain.Task{Title: "another", Status: domain.StatusDoing}); err != nil {
		t.Fatalf("create 2: %v", err)
	}

	all, err := r.List(ctx, TaskFilter{})
	if err != nil || len(all) != 2 {
		t.Fatalf("list all: %v len=%d", err, len(all))
	}

	doing, err := r.List(ctx, TaskFilter{Status: domain.StatusDoing})
	if err != nil || len(doing) != 1 {
		t.Fatalf("list doing: %v len=%d", err, len(doing))
	}

	found, err := r.List(ctx, TaskFilter{Search: "tests"})
	if err != nil || len(found) != 1 {
		t.Fatalf("search: %v len=%d", err, len(found))
	}

	if err := r.SetStatus(ctx, t1.ID, domain.StatusDone); err != nil {
		t.Fatalf("setstatus: %v", err)
	}
	got, _ = r.Get(ctx, t1.ID)
	if got.Status != domain.StatusDone || got.CompletedAt == "" {
		t.Fatalf("done not stamped: %+v", got)
	}
	// reopening clears completed_at
	if err := r.SetStatus(ctx, t1.ID, domain.StatusTodo); err != nil {
		t.Fatalf("reopen: %v", err)
	}
	got, _ = r.Get(ctx, t1.ID)
	if got.CompletedAt != "" {
		t.Fatalf("completed_at should be cleared on reopen: %+v", got)
	}

	count, err := r.Count(ctx)
	if err != nil || count != 2 {
		t.Fatalf("count: %v %d", err, count)
	}

	if err := r.Delete(ctx, t1.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}
	count, _ = r.Count(ctx)
	if count != 1 {
		t.Fatalf("after delete count=%d", count)
	}
}

func TestNoteRepo_CRUDAndArchive(t *testing.T) {
	db := newTestDB(t)
	ctx := context.Background()
	r := NewNoteRepo(db)

	n, err := r.Create(ctx, domain.Note{Title: "ideas", Body: "# hello"})
	if err != nil {
		t.Fatalf("create: %v", err)
	}

	all, _ := r.List(ctx, NoteFilter{})
	if len(all) != 1 {
		t.Fatalf("list len=%d", len(all))
	}

	if hits, _ := r.List(ctx, NoteFilter{Search: "hello"}); len(hits) != 1 {
		t.Fatalf("search body miss")
	}

	if err := r.Archive(ctx, n.ID); err != nil {
		t.Fatalf("archive: %v", err)
	}
	active, _ := r.List(ctx, NoteFilter{})
	if len(active) != 0 {
		t.Fatalf("archived note should be hidden from default list")
	}
	archived, _ := r.List(ctx, NoteFilter{IncludeArchived: true})
	if len(archived) != 1 {
		t.Fatalf("archived note should appear with IncludeArchived")
	}

	if err := r.Unarchive(ctx, n.ID); err != nil {
		t.Fatalf("unarchive: %v", err)
	}
	active, _ = r.List(ctx, NoteFilter{})
	if len(active) != 1 {
		t.Fatalf("unarchived note should be visible again")
	}
}
