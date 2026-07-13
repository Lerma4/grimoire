package store

import (
	"database/sql"
	"fmt"
	"time"
)

// SeedIfEmpty populates the database with a small welcome dataset the first
// time it runs, so the app never opens on an empty screen.
func SeedIfEmpty(db *sql.DB) error {
	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&n); err != nil {
		return fmt.Errorf("count tasks: %w", err)
	}
	if n > 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin seed tx: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck

	now := time.Now().UTC().Format(time.RFC3339)
	due := now // due today

	mustExec := func(query string, args ...any) {
		if err != nil {
			return
		}
		_, err = tx.Exec(query, args...)
	}

	mustExec(`INSERT INTO tasks (title, description, status, priority, due_date, created_at, updated_at)
		VALUES (?, ?, 'todo', 'high', ?, ?, ?)`,
		"Open the help screen", "Press `?` inside the TUI to see all keybindings.", due, now, now)

	mustExec(`INSERT INTO tasks (title, description, status, priority, created_at, updated_at)
		VALUES (?, ?, 'doing', 'medium', ?, ?)`,
		"Explore task↔note linking", "Select a task and press `L` to link a note.", now, now)

	mustExec(`INSERT INTO notes (title, body, created_at, updated_at)
		VALUES (?, ?, ?, ?)`,
		"Welcome to your Grimoire",
		"# Welcome\n\n"+
			"This is a **markdown note**. Press `e` to edit it, `?` for help.\n\n"+
			"- tasks live in the center list\n"+
			"- notes render here with markdown\n"+
			"- link the two with `L`\n",
		now, now)

	if err != nil {
		return fmt.Errorf("seed insert: %w", err)
	}
	return tx.Commit()
}
