# Ship

A lightweight, Git-inspired version control system written in Go. Ship implements core VCS concepts — content-addressable object storage, a staging area, commit history, and pipeline tracing — as a clean, educational implementation.

---

## What Ship Does

Ship gives you a self-contained version control system with the following capabilities:

- **Initialize** a repository with a structured object store
- **Stage** files into a JSON-backed index
- **Commit** snapshots with full tree and metadata objects
- **Inspect** stored objects by hash
- **View** commit history
- **Check** working tree status against the index
- **Trace** the internal execution pipeline of any supported command

Every object (blob, tree, commit) is content-addressed using SHA-1 and compressed with zlib — the same foundational model as Git.

---

## Commands

### `ship init [path]`

Initializes a new Ship repository. Creates the `.ship/` directory structure:

```
.ship/
  objects/       # Content-addressable object store
  refs/
    heads/       # Branch references
    tags/        # Tag references
  HEAD           # Points to current branch (default: main)
  index          # Staging area (JSON)
```

---

### `ship add <files...>`

Stages one or more files or directories into the index.

- Recursively walks directories
- Skips the `.ship/` directory
- Computes SHA-1 hash for each file, stores a compressed blob object
- Updates the index with path, hash, and file mode
- Removes deleted files from the index automatically

```bash
ship add src/
ship add README.md go.mod
```

---

### `ship commit <message>`

Creates a commit from whatever is currently staged.

1. Loads the index
2. Builds a tree object from index entries (hierarchical, directory-aware)
3. Resolves HEAD to find the parent commit
4. Reads author/committer identity from environment variables (`GIT_AUTHOR_NAME`, `GIT_AUTHOR_EMAIL`, etc.)
5. Writes the commit object and updates HEAD

```bash
ship commit "feat: add authentication"
```

---

### `ship status [path]`

Shows the working tree status relative to the index.

- **Blue** — staged (file matches what's in the index)
- **Tan/Sand** — unstaged or untracked (differs from or absent in index)

Compares the live filesystem against the stored index and reports differences per file.

```bash
ship status
ship status ./src
```

---

### `ship pour [--oneline]`

Displays commit history (Ship's equivalent of `git log`).

- Default: full commit details (hash, author, date, message)
- `--oneline`: one commit per line — short hash + message

```bash
ship pour
ship pour --oneline
```

---

### `ship cat-file <hash>`

Inspects a stored object by its SHA-1 hash.

| Flag | Description |
|------|-------------|
| `-p, --pretty` | Pretty-print object content (default) |
| `-t, --tree` | Show as tree object |
| `-c, --commit` | Show as commit object |
| `-s, --size` | Show object size in bytes |

```bash
ship cat-file -p a3f9c1...
ship cat-file -s a3f9c1...
```

---

### `ship trace <command> [args...]`

Records and prints every internal pipeline step for a supported command. This is Ship's built-in observability layer.

**Supported commands:** `add`, `commit`, `status`

#### Output Formats

| Flag | Output |
|------|--------|
| *(default)* | Colored human-readable text with a summary footer |
| `--format json` | NDJSON — one JSON object per step, pipeable to `jq` |
| `--otel` | OpenTelemetry-compatible spans |
| `--otel-output <file>` | File to write OTel output (default: stderr, `-` for stdout) |

#### Text Output

Each step prints with a status symbol, step name, and duration:

```
✓  ResolveTargets       120µs
✓  LoadIndex            340µs
✓  WriteBlobs          1.2ms
✓  SaveIndex            210µs

Steps: 4   Duration: 1.87ms   Errors: 0
files_staged: 3
targets_resolved: 3
```

#### JSON Output

Pipeable NDJSON — one object per step:

```json
{"name":"ResolveTargets","duration_ms":0.12,"status":"ok"}
{"name":"LoadIndex","duration_ms":0.34,"status":"ok"}
```

Filter with `jq`:

```bash
ship trace add src/ --format json | jq 'select(.status=="error")'
```

#### OpenTelemetry Output

Emits OTel-compatible spans with parent/child structure, timing, and status codes — ready for ingestion into any OTel-compatible backend.

```bash
ship trace commit "fix: typo" --otel --otel-output trace.json
```

#### Examples

```bash
ship trace add README.md src/
ship trace commit "feat: init"
ship trace status
ship trace add . --format json | jq '.duration_ms'
ship trace commit "fix" --otel --otel-output spans.json
```

---

### `ship purge`

Removes temporary files created during operations. Automatically invoked after `add`, `commit`, and `status`.

---

## Object Model

Ship stores three types of objects, all content-addressed by SHA-1 and compressed with zlib:

| Type | Description |
|------|-------------|
| **Blob** | Raw file content |
| **Tree** | Directory snapshot — maps names to blob/tree hashes |
| **Commit** | Metadata + tree hash + parent reference |

Objects are stored at `.ship/objects/XX/YYYY...` where `XX` is the first two characters of the hash.

---

## Repository Layout

```
.ship/
  HEAD              # ref: refs/heads/main
  index             # JSON staging area
  objects/
    ab/
      cdef1234...   # Compressed blob/tree/commit
  refs/
    heads/
      main          # SHA-1 of latest commit
    tags/
```

---

## Tracing Architecture

Ship's trace system is built around a pluggable `Sink` interface:

```
trace.Step("StepName") → records duration + error on return
trace.Meta("key", "value") → zero-duration metadata annotation
```

Three sink implementations:

- **PrettySink** — colored terminal output with summary
- **JSONSink** — thread-safe NDJSON emission
- **OTelSink** — OpenTelemetry spans via a provider

The `trace` command swaps the active sink before executing, then restores it after — so every internal step in `add`, `commit`, or `status` is captured without modifying command logic.

---

## Installation

```bash
git clone https://github.com/KambojRajan/ship
cd ship
go build -o ship .
```

Requires Go 1.24+.
