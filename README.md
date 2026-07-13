```
         ___
       ╱     ╲
      │  ✦ ◇  │     Grimoire
      │   ◆   │     a terminal grimoire for tasks & markdown notes
       ╲___╱
```

> Bound in SQLite. Inked in Vim. Opened from any shell.

**Grimoire** is a local-first terminal app for managing tasks and markdown
notes — and for linking the two contextually. Think of it as a small arcane
tome you keep in your pocket terminal: fast, keyboard-driven, offline, and
with zero external services.

- Vim-style modal TUI (Normal / Command / Search)
- Three-column layout: sidebar · list · detail (collapses on narrow terminals)
- Tasks with status, priority, due dates, tags
- Markdown notes with live preview (glamour)
- Many-to-many task ↔ note linking
- SQLite storage (pure-Go driver, no CGO)
- Cross-platform: Linux & Windows, one `go install`

```
┌─✦ Grimoire───────────────────────┬─Tasks───────────────┬─Task──────────────────────┐
│ ▸ ○ Tasks    4                   │ ▸ ○ !! Ship v0.1     │  Ship v0.1                │
│   ✎ Notes    3                   │   ◐ !  Wire the TUI  │                           │
│   ★ Today    2                   │   ●    Seed the db   │  Status     ◐ doing       │
│   # Tags     5                   │   ○    Write docs    │  Priority   high          │
│   ↔ Links    2                   │                      │  Due        2026-07-10    │
│   ▽ Archive  7                   │                      │  Tags       #release      │
│                                  │                      │                           │
│                                  │                      │  Linked notes (1)         │
│                                  │                      │   ✎ Release checklist     │
├──────────────────────────────────┴──────────────────────┴───────────────────────────┤
│ [NORMAL]  j/k move · h/l panes · : cmd · / search · ? help · q quit                  │
└──────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Install

### Linux / macOS

```bash
go install github.com/Lerma4/grimoire/cmd/grimoire@latest
```

Make sure `$HOME/go/bin` is on your `PATH`:

```bash
# bash / zsh
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.profile
# fish
fish_add_path "$HOME/go/bin"
```

### Windows (PowerShell)

```powershell
go install github.com/Lerma4/grimoire/cmd/grimoire@latest
```

Add the Go bin dir to your `PATH`:

```powershell
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";$env:USERPROFILE\go\bin", "User")
```

### Build manually

```bash
git clone https://github.com/Lerma4/grimoire
cd grimoire
go build -o grimoire ./cmd/grimoire
```

After install, run from anywhere:

```bash
grimoire
```

---

## Commands

| Command           | Description                                  |
|-------------------|----------------------------------------------|
| `grimoire`        | Open the TUI (default)                       |
| `grimoire tui`    | Open the TUI                                 |
| `grimoire version`| Print version                                |
| `grimoire doctor` | Check database, config paths, terminal env   |

---

## Keybindings (Vim-style)

### Normal mode

| Key         | Action                                          |
|-------------|-------------------------------------------------|
| `j` / `↓`   | Next item                                       |
| `k` / `↑`   | Previous item                                   |
| `h` / `l`   | Focus previous / next pane                      |
| `gg` / `G`  | First / last item                               |
| `/`         | Search                                          |
| `n` / `N`   | Next / previous search result                   |
| `Enter`     | Open / focus detail                             |
| `Space`     | Toggle task done ↔ todo                         |
| `a`         | New task (form)                                |
| `A`         | New note (form)                                |
| `e`         | Edit selected (form)                           |
| `d`         | Archive                                        |
| `D`         | Delete (confirm)                               |
| `t m #`   | Jump to Tasks / Notes / Tags                    |
| `L` / `U`   | Link / unlink task ↔ note (`:link`/`:unlink`)   |
| `?`         | Help                                            |
| `:`         | Command mode                                    |
| `q`         | Quit                                            |

### Command mode (`:`)

```
:quit  :q            quit
:task add <title>    create a task
:note add <title>    create a note
:done  :doing  :todo set task status
:archive             archive selected
:delete              delete selected (confirm)
:link  :unlink       link / unlink task ↔ note
:tag <name>          add tag
```

---

## Task status glyphs

| Glyph | Status   |
|-------|----------|
| `○`   | todo     |
| `◐`   | doing    |
| `●`   | done     |
| `!`   | overdue  |

---

## Data location

The SQLite database is created automatically on first run.

- **Linux/macOS:** `~/.local/share/grimoire/grimoire.db`
- **Windows:**       `%LOCALAPPDATA%\grimoire\grimoire.db`

Override with the `GRIMOIRE_DB` environment variable.

---

## Development

```bash
make run      # go run ./cmd/grimoire
make build    # build ./bin/grimoire
make install  # go install ./cmd/grimoire
make test     # go test ./...
make check    # gofmt + go vet + go test (the pre-commit gate)
make lint     # golangci-lint run (skipped if not installed)
```

### Optional: golangci-lint

```bash
# install (one way)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# run
golangci-lint run
```

---

## License

MIT
