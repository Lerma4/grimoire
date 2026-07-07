// Package domain defines the pure business entities of Grimoire.
// Entities are persistence-agnostic: they carry data and invariants only.
package domain

import "time"

// now is the time source; package tests may override it.
var now = func() time.Time { return time.Now().UTC() }

// TimeStamp returns the current UTC time as an RFC3339 string, matching the
// TEXT columns used in SQLite.
func TimeStamp() string { return now().Format(time.RFC3339) }

// Project groups related tasks and notes.
type Project struct {
	ID          int64
	Name        string
	Description string
	CreatedAt   string
	UpdatedAt   string
}
