# AGENTS.md

Guidance for AI agents (and humans) working on **Grimoire**.

Grimoire is a local-first terminal app for tasks and markdown notes, with
Vim-style keybindings and task↔note linking. Pure-Go, SQLite-backed, built for
Linux and Windows.

- Repo: https://github.com/Lerma4/grimoire
- Module: `github.com/Lerma4/grimoire`
- Entry point: `cmd/grimoire/main.go`

## Commands

```bash
make run        # go run ./cmd/grimoire
make build      # -> bin/grimoire
make install    # go install ./cmd/grimoire  (global `grimoire`)
make test       # go test ./...
make check      # gofmt + go vet + go test   (the pre-commit gate)
make lint       # golangci-lint run (no-op if not installed)
```

**Before every commit, run the gate:**

```bash
gofmt -w .
go test ./...
go vet ./...
```

If checks fail, fix before committing. Never commit with failing tests or vet
errors. `golangci-lint` is optional — don't block on it, but run it if present.

## Architecture

Layered; dependencies point downward only.

```
cmd/grimoire   Cobra entrypoint (root, tui, version, doctor). Thin — no logic.
internal/app    Config/db-path resolution (OS-aware) + doctor diagnostics.
internal/domain Pure entities + enums + invariants. No I/O, no imports downward.
internal/store  SQLite (modernc.org/sqlite, pure Go, no CGO). Repos + embedded
                versioned SQL migrations under internal/store/migrations/.
                schema_migrations is owned by the runner, NOT the SQL files.
internal/service Application logic. Services are defined against small repo
                interfaces so they're unit-testable. Wired together by
                service.NewServices(db).
internal/tui    Bubble Tea model/update/view. Vim modes + command palette.
  components    Leaf package: centralized palette/styles, section/pane/mode
                enums, and renderers (sidebar, list, detail, header, statusbar,
                markdown). Imports only domain + lipgloss/glamour — never tui.
```

The TUI talks **only** to `service`. Services talk only to `store` interfaces.
`domain` is shared and depended on by everyone; it depends on nothing.

### Key invariants

- `internal/tui/components` must stay a leaf package (no import cycle with `tui`).
- Services take repo interfaces, not concrete `*store.XxxRepo` types. When you
  add a repo method needed by a service, extend the interface in the service file.
- SQLite TEXT timestamps are RFC3339 UTC (`domain.TimeStamp()`). Empty optional
  string columns are stored as NULL via `store.nullString` / `store.nullID`.
- `relation_type` in `task_notes` is reserved; MVP uses `reference`.
- Task state is shown with textual glyphs (`○ ◐ ● !`) — never communicate state
  by color alone.

## Storage

- Driver: `modernc.org/sqlite` (pure Go) — do not switch to a CGO driver; Windows
  portability depends on it.
- Migrations: add a new `internal/store/migrations/NNNN_name.sql`. The runner
  applies them in filename order and records each in `schema_migrations`. Never
  edit an applied migration — add a new one.
- DB path: `~/.local/share/grimoire/grimoire.db` (Linux/macOS) or
  `%LOCALAPPDATA%\grimoire\grimoire.db` (Windows); override with `GRIMOIRE_DB`.
- `store.SeedIfEmpty` plants a welcome dataset on first run only.

## Workflow

- Work in small, coherent steps. One logical change per commit.
- Run the check gate at the end of every step, then commit.
- Commit message style follows the existing history — Conventional Commits with
  Gitflow intent: `feat:`, `fix:`, `chore:`, `docs:`, `test:`, `refactor:`.
  Use clear, descriptive subjects and a short body for the "why".
- Manage the app version deliberately from now on: bump `internal/app/version.go`
  for user-visible changes before committing (`patch` for fixes, `minor` for
  features, `major` only for breaking changes).
- Push after each commit (or small coherent group). Do not accumulate large
  unpushed diffs on `main`.
- Do not commit the binary, the SQLite DB, or editor cruft (see `.gitignore`).
- Never commit secrets. Keep `GRIMOIRE_DB` test paths under `/tmp` or `t.TempDir()`.

## Testing

- `internal/store` — repo CRUD + filtering + migrations, against a temp file DB
  (`newTestDB`, single connection).
- `internal/service` — lifecycle/validation/linking/search via `newServices(t)`.
- `internal/domain` — entity invariants, glyphs, overdue logic.
- Prefer integration tests over mocks: open a real temp SQLite, migrate, and
  exercise the stack. Services are already interface-backed for targeted fakes.

## Don't

- Don't add native/CGO dependencies.
- Don't put application logic in `main.go` or rendering logic in `service`.
- Don't add abstractions with a single implementation "for later".
- Don't rely on training-data API details for the Charm stack or modernc/sqlite —
  verify against current docs (lipgloss v1 dropped `BorderTitle`, for example).
- Don't introduce global state; pass `*service.Services` / `*sql.DB` explicitly.

## TTY note

This project targets an interactive terminal. In CI/sandboxes without a TTY the
full TUI cannot be exercised interactively; cover rendering and logic with Go
tests instead (construct a `tui.Model`, set width/height, call `View()`).
