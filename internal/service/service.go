// Package service contains the application logic of Grimoire. Services wrap
// repositories and enforce invariants and validation; they are the only layer
// the TUI talks to.
package service

import "context"

// Ctx is the default background context used by service calls.
var Ctx = context.Background()
