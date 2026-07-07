package domain

import (
	"strings"
	"testing"
	"time"
)

func TestTask_Validate(t *testing.T) {
	cases := []struct {
		name string
		task Task
		ok   bool
	}{
		{"empty title", Task{Status: StatusTodo}, false},
		{"empty status", Task{Title: "x"}, false},
		{"bad status", Task{Title: "x", Status: "weird"}, false},
		{"bad priority", Task{Title: "x", Status: StatusTodo, Priority: "critical"}, false},
		{"valid todo", Task{Title: "x", Status: StatusTodo, Priority: PriorityLow}, true},
		{"empty priority ok", Task{Title: "x", Status: StatusDoing}, true},
	}
	for _, c := range cases {
		err := c.task.Validate()
		if c.ok && err != nil {
			t.Errorf("%s: expected ok, got %v", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("%s: expected error, got nil", c.name)
		}
	}
}

func TestStatusGlyph(t *testing.T) {
	want := map[string]string{
		StatusTodo: "○", StatusDoing: "◐", StatusDone: "●", StatusArchived: "□",
	}
	for status, glyph := range want {
		if g := StatusGlyph(status); g != glyph {
			t.Errorf("glyph(%s) = %q want %q", status, g, glyph)
		}
	}
}

func TestTask_IsOverdue(t *testing.T) {
	past := time.Now().Add(-24 * time.Hour).UTC().Format(time.RFC3339)
	future := time.Now().Add(24 * time.Hour).UTC().Format(time.RFC3339)

	cases := []struct {
		name string
		task Task
		want bool
	}{
		{"no due date", Task{Status: StatusTodo}, false},
		{"past due todo", Task{Status: StatusTodo, DueDate: past}, true},
		{"past due but done", Task{Status: StatusDone, DueDate: past}, false},
		{"past due but archived", Task{Status: StatusArchived, DueDate: past}, false},
		{"future due", Task{Status: StatusTodo, DueDate: future}, false},
		{"garbage date", Task{Status: StatusTodo, DueDate: "not-a-date"}, false},
	}
	for _, c := range cases {
		if got := c.task.IsOverdue(); got != c.want {
			t.Errorf("%s: IsOverdue = %v want %v", c.name, got, c.want)
		}
	}
}

func TestTask_OverdueGlyph(t *testing.T) {
	past := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	task := Task{Status: StatusTodo, DueDate: past}
	if g := task.OverdueGlyph(); g != "!" {
		t.Errorf("overdue glyph = %q want !", g)
	}
	done := Task{Status: StatusDone}
	if g := done.OverdueGlyph(); g != "●" {
		t.Errorf("done glyph = %q want ●", g)
	}
}

func TestValidStatusesAndPriorities(t *testing.T) {
	if len(ValidStatuses) != 4 {
		t.Errorf("expected 4 statuses, got %d", len(ValidStatuses))
	}
	if len(ValidPriorities) != 4 {
		t.Errorf("expected 4 priorities, got %d", len(ValidPriorities))
	}
	// sanity: every status string is lowercase alphanumeric
	for _, s := range ValidStatuses {
		if s == "" || strings.ToLower(s) != s {
			t.Errorf("bad status token %q", s)
		}
	}
}
