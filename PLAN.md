# Grimoire — implementation plan

A terminal TUI for tasks and markdown notes, with Vim-style bindings and
contextual linking between tasks and notes. Local-first, SQLite-backed,
pure-Go (modernc.org/sqlite) so it builds cleanly on Linux and Windows.

## Folder structure

```text
cmd/grimoire/main.go        # Cobra entrypoint
internal/
  app/                      # bootstrap: config, db init, demo seed
  domain/                   # pure entities: Task, Note, Project, Tag, TaskNoteLink
  store/                    # SQLite connection, migrations, repositories
  service/                  # application logic: task/note/project/tag/link/search
  tui/
    model.go                # top-level Bubble Tea model + Update/View
    keymap.go               # Vim-style keybindings
    styles.go               # centralized Lip Gloss palette
    messages.go             # tea.Msg types
    modes.go                # NORMAL / INSERT / COMMAND / SEARCH
    components/
      sidebar.go            # left nav: Tasks/Notes/Today/Projects/Tags/Links/Archive
      tasklist.go           # central list of tasks
      notelist.go           # central list of notes
      detail.go             # right panel: task/note detail + linked items
      command.go            # ':' command palette
      search.go             # '/' search input
      form.go               # huh-based new/edit modal
      confirm.go            # delete confirmation modal
      statusbar.go          # bottom mode + hints + messages
      header.go             # top section/filter/project/db status
```

## Database schema (migration v1)

```sql
CREATE TABLE projects (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT NOT NULL UNIQUE,
  description TEXT,
  created_at  TEXT NOT NULL,
  updated_at  TEXT NOT NULL
);

CREATE TABLE tasks (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  title        TEXT NOT NULL,
  description  TEXT,
  status       TEXT NOT NULL DEFAULT 'todo',   -- todo|doing|done|archived
  priority     TEXT NOT NULL DEFAULT 'medium', -- low|medium|high|urgent
  due_date     TEXT,
  project_id   INTEGER REFERENCES projects(id) ON DELETE SET NULL,
  created_at   TEXT NOT NULL,
  updated_at   TEXT NOT NULL,
  completed_at TEXT,
  archived_at  TEXT
);
CREATE INDEX idx_tasks_status     ON tasks(status);
CREATE INDEX idx_tasks_project_id ON tasks(project_id);

CREATE TABLE notes (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  title      TEXT NOT NULL,
  body       TEXT,
  project_id INTEGER REFERENCES projects(id) ON DELETE SET NULL,
  created_at TEXT NOT NULL,
  updated_at TEXT NOT NULL,
  archived_at TEXT
);
CREATE INDEX idx_notes_project_id ON notes(project_id);

CREATE TABLE tags (
  id    INTEGER PRIMARY KEY AUTOINCREMENT,
  name  TEXT NOT NULL UNIQUE,
  color TEXT,
  created_at TEXT NOT NULL
);

CREATE TABLE task_tags (
  task_id INTEGER NOT NULL REFERENCES tasks(id)    ON DELETE CASCADE,
  tag_id  INTEGER NOT NULL REFERENCES tags(id)     ON DELETE CASCADE,
  PRIMARY KEY (task_id, tag_id)
);

CREATE TABLE note_tags (
  note_id INTEGER NOT NULL REFERENCES notes(id)    ON DELETE CASCADE,
  tag_id  INTEGER NOT NULL REFERENCES tags(id)     ON DELETE CASCADE,
  PRIMARY KEY (note_id, tag_id)
);

CREATE TABLE task_notes (
  task_id      INTEGER NOT NULL REFERENCES tasks(id) ON DELETE CASCADE,
  note_id      INTEGER NOT NULL REFERENCES notes(id) ON DELETE CASCADE,
  relation_type TEXT NOT NULL DEFAULT 'reference',
  created_at    TEXT NOT NULL,
  PRIMARY KEY (task_id, note_id)
);
```

`relation_type` is reserved for richer semantics later; MVP uses `reference`.

## TUI components & layout

Three-column layout, collapsible on narrow terminals:

- **Sidebar** — nav sections + counts.
- **Center list** — tasks/notes/search results for the active section.
- **Detail panel** — selected item + its linked items.

Below ~90 cols the sidebar collapses to icons; below ~60 cols it switches to a
two-pane list/detail toggle. A header shows section/filter/project/db health;
a status bar shows the active mode and contextual hints.

## Vim-style keybindings

| Key        | Action                                  |
|------------|-----------------------------------------|
| `j` / `↓`  | next item                               |
| `k` / `↑`  | prev item                               |
| `h`        | focus previous pane                     |
| `l`        | focus next pane                         |
| `gg` / `G` | first / last item                       |
| `/`        | search                                  |
| `n` / `N`  | next / prev search result               |
| `Enter`    | open / focus detail                     |
| `Space`    | toggle task done/todo                   |
| `a`        | new task                                |
| `A`        | new note                                |
| `e`        | edit selected                           |
| `d`        | archive selected                        |
| `D`        | delete (confirm)                        |
| `t m p #`  | jump to Tasks / Notes / Projects / Tags |
| `L` / `U`  | link / unlink task↔note                 |
| `?`        | help                                    |
| `:`        | command mode                            |
| `q`        | quit                                    |

Command mode: `:q`/`:quit`, `:task add <t>`, `:note add <t>`, `:done`,
`:doing`, `:todo`, `:archive`, `:delete`, `:link`, `:unlink`,
`:project <name>`, `:tag <name>`.

## Task–note linking strategy

Many-to-many through `task_notes(task_id, note_id, relation_type, created_at)`.
The `LinkService` owns all link mutations and both directions of lookup
(`NotesForTask`, `TasksForNote`). The detail panel renders linked items inline
and `L`/`U` open a searchable picker so you can find a linkable item without
leaving context.

## Milestones

1. Skeleton: go.mod, cobra root + `version`/`doctor`, `go run` works.
2. Storage: SQLite open + v1 migration + repositories + demo seed.
3. Domain services: task/note/project/tag/link/search CRUD.
4. TUI shell: three columns, styles, header, status bar, basic nav.
5. Vim keymap: modes, motion, command palette, search.
6. Forms + preview: huh forms for new/edit, glamour markdown render.
7. Polish + tests: help screen, responsive layout, store/service tests.
