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

- Vim-style modal TUI (Normal / Insert / Command / Search)
- Three-column layout: sidebar · list · detail
- Tasks with status, priority, due dates, projects, tags
- Markdown notes with live preview
- Many-to-many task ↔ note linking
- SQLite storage (pure-Go driver, no CGO)
- Cross-platform: Linux & Windows, one `go install`

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
| `a`         | New task                                        |
| `A`         | New note                                        |
| `e`         | Edit selected                                   |
| `d`         | Archive                                         |
| `D`         | Delete (confirm)                                |
| `t m p #`   | Jump to Tasks / Notes / Projects / Tags         |
| `L` / `U`   | Link / unlink task ↔ note                       |
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
:project <name>      set project
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
make fmt      # gofmt -w .
make vet      # go vet ./...
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
