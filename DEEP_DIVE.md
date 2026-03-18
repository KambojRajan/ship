# 🚢 Ship VCS — Complete Deep-Dive Documentation

> A Git-inspired, minimal Version Control System written in Go.  
> Author: KambojRajan · Version: 1.0.0 · Language: Go 1.24 · CLI: Cobra

---

## Table of Contents

1. [Project Overview](#1-project-overview)
2. [Repository Layout](#2-repository-layout)
3. [High-Level Design (HLD)](#3-high-level-design-hld)
4. [Low-Level Design (LLD)](#4-low-level-design-lld)
5. [UML — Class Diagram](#5-uml--class-diagram)
6. [UML — Sequence Diagrams (per command)](#6-uml--sequence-diagrams-per-command)
7. [Data-Flow Diagram](#7-data-flow-diagram)
8. [Object Storage Model](#8-object-storage-model)
9. [On-Disk Layout (.ship directory)](#9-on-disk-layout-ship-directory)
10. [Code Flow — Every Command Explained](#10-code-flow--every-command-explained)
11. [Core Algorithms & Strategies](#11-core-algorithms--strategies)
12. [Concurrency Infrastructure](#12-concurrency-infrastructure)
13. [CLI Reference](#13-cli-reference)
14. [Dependency Graph](#14-dependency-graph)
15. [Testing Strategy](#15-testing-strategy)
16. [Build, Install & Run](#16-build-install--run)
17. [Environment Variables](#17-environment-variables)
18. [Design Patterns Used](#18-design-patterns-used)
19. [Comparison with Git](#19-comparison-with-git)
20. [Known Limitations & Future Scope](#20-known-limitations--future-scope)

---

## 1. Project Overview

**Ship** is a ground-up reimplementation of the core primitives of Git, written in idiomatic Go. It stores file snapshots as **content-addressed objects** (SHA-1 hashed, zlib-compressed), tracks the working directory against a staging area (index), builds Merkle trees over directory structures, and chains commits into a linked-list history.

### What Ship Does

| Capability | Detail |
|---|---|
| Initialize repo | Creates `.ship/` metadata directory |
| Stage files | Hashes blobs and records them in a JSON index |
| Commit | Builds a tree object, wraps it in a commit object, updates HEAD |
| Show status | Compares live FS against staged index |
| Inspect objects | Decompress and print any stored object by hash |
| View history | Walks the commit linked list from HEAD |

### What Makes It Different From Git

- **Index is JSON** (human-readable) rather than a custom binary format.  
- **Tree nodes store hash strings** rather than raw binary SHA bytes.  
- **No branch/remote/merge/diff** support yet — core DAG primitives only.  
- Ships as a **single static binary** with no runtime dependencies.

---

## 2. Repository Layout

```
ship/
├── main.go                      # Entry-point — calls cmd.Execute()
├── go.mod / go.sum              # Module: github.com/KambojRajan/ship
├── Makefile                     # build / test / install / clean targets
├── VERSION                      # 1.0.0
│
├── cmd/                         # CLI layer — Cobra command wiring only
│   ├── root.go                  # rootCmd ("ship"), Execute()
│   ├── init.go                  # ship init
│   ├── add.go                   # ship add [files...]
│   ├── commit.go                # ship commit <message>
│   ├── status.go                # ship status [path]
│   ├── cat-file.go              # ship cat-file [-p|-t|-c|-s] <hash>
│   └── pour.go                  # ship pour  (commit log)
│
├── commands/                    # Business-logic layer — pure Go, no Cobra
│   ├── init.go                  # Init()
│   ├── add.go                   # Add()
│   ├── commit.go                # Commit()
│   ├── status.go                # Status()
│   └── cat-file.go              # CatFile()
│
├── core/
│   ├── common/                  # Shared value-types
│   │   ├── object_type.go       # ObjectType enum: BLOB / TREE / COMMIT
│   │   └── temp_dir_node.go     # IndexEntry, TempDirNode (virtual FS tree)
│   │
│   ├── entities/                # Domain model
│   │   ├── blob.go              # Blob{Data, Hash}
│   │   ├── node.go              # Node{Mode, Name, Hash}
│   │   ├── tree.go              # Tree + WriteTree + buildTempDirTree
│   │   ├── index.go             # Index{Entries} + Load/Save (JSON)
│   │   ├── Head.go              # Head{Ref, Hash} + ResolveHead / UpdateRef
│   │   ├── Commit.go            # Commit + serialize + parseCommit + LoadCommits
│   │   └── user.go              # User{Name, Email, Timestamp} + NewUserFromEnv
│   │
│   └── utils/                   # Infrastructure
│       ├── constants.go         # All string constants & ANSI colours
│       ├── error-code.go        # Named error-format strings
│       ├── hash.go              # HashBytes / HashString (SHA-1 wrappers)
│       ├── object.go            # HashObject (write=true → disk), ObjectExists, GetMode
│       ├── ship-core.go         # ShipHasBeenInitRecursive (repo discovery)
│       └── concurrent-processor.go  # Generic worker-pool
│
└── tests/
    ├── command_tests/           # Integration tests per command
    │   ├── init_test.go
    │   ├── add_test.go
    │   ├── commit_test.go
    │   ├── status_test.go
    │   ├── cat-file_test.go
    │   └── pour_test.go
    └── helpers/
        ├── setup.helper.go      # TempDir scaffolding, BurnDown
        └── assert.helper.go     # AssertNil/NotNil/Equal/Exists/FileInIndex etc.
```

---

## 3. High-Level Design (HLD)

```
┌──────────────────────────────────────────────────────────────┐
│                        User / Shell                          │
└────────────────────────────┬─────────────────────────────────┘
                             │  CLI invocation
                             ▼
┌──────────────────────────────────────────────────────────────┐
│                    CLI Layer  (cmd/)                         │
│  Cobra root command + sub-commands                           │
│  Parses flags/args, calls into commands/ layer               │
└────────────────────────────┬─────────────────────────────────┘
                             │  pure function call
                             ▼
┌──────────────────────────────────────────────────────────────┐
│                 Business Logic Layer  (commands/)            │
│  Init · Add · Commit · Status · CatFile                      │
│  Orchestrates entities + utils, owns all workflows           │
└───────────┬──────────────────────────┬───────────────────────┘
            │ uses entities            │ uses utils
            ▼                          ▼
┌────────────────────────┐   ┌─────────────────────────────────┐
│  Domain Model          │   │  Infrastructure / Utils         │
│  (core/entities/)      │   │  (core/utils/ + core/common/)   │
│                        │   │                                 │
│  Blob                  │   │  HashObject (SHA-1 + zlib)      │
│  Tree                  │   │  ObjectExists                   │
│  Commit                │   │  ShipHasBeenInitRecursive       │
│  Index                 │   │  ConcurrentProcessor            │
│  Head                  │   │  GetMode                        │
│  User                  │   │  All constants & error codes    │
│  Node                  │   │                                 │
└────────────┬───────────┘   └────────────────┬────────────────┘
             │                                │
             └────────────────┬───────────────┘
                              │ reads / writes
                              ▼
┌──────────────────────────────────────────────────────────────┐
│                  Filesystem  (.ship/)                        │
│                                                              │
│  .ship/HEAD              → "ref: refs/heads/main\n"          │
│  .ship/refs/heads/main   → "<40-char commit hash>"           │
│  .ship/index             → JSON staging area                 │
│  .ship/objects/XX/YY...  → zlib-compressed object files      │
└──────────────────────────────────────────────────────────────┘
```

---

## 4. Low-Level Design (LLD)

### 4.1 Package Dependency Graph

```
main
 └─► cmd
      └─► commands
           ├─► core/entities   (domain)
           ├─► core/common     (value objects)
           └─► core/utils      (infrastructure)
                └─► core/common
```

No circular imports. Each layer only imports **downward**.

### 4.2 Object Storage Pipeline

```
Raw file bytes
     │
     ▼
 fmt.Sprintf("%s %d\0", type, len)   ← header
     │
     ▼
 header + \0 + data  →  SHA-1  →  40-char hex hash
     │
     ▼ (if write=true)
 zlib.Compress(header + \0 + data)
     │
     ▼
 .ship/objects/<hash[0:2]>/<hash[2:]>
```

### 4.3 Index Data Model

```json
{
  "Entries": {
    "relative/path/to/file.go": {
      "Path": "relative/path/to/file.go",
      "Hash": "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3",
      "Mode": 100644
    }
  }
}
```

### 4.4 Commit Object Wire Format (serialized, before zlib)

```
commit <size>\0tree <tree-hash>\n
parent <parent-hash>\n          ← 0..N lines
author Name <email> <unix-ts> <tz>\n
committer Name <email> <unix-ts> <tz>\n
\n
<commit message>
```

### 4.5 Tree Object Wire Format

```
<mode-octal> <name>\0<hash-bytes>   ← repeated for each entry, sorted by name
```

---

## 5. UML — Class Diagram

```
┌──────────────────────────────────────────────────────────────────────────────┐
│  core/common                                                                  │
│                                                                               │
│  «enum»                     «struct»                  «struct»                │
│  ObjectType                 IndexEntry                TempDirNode             │
│  ─────────                  ──────────                ───────────             │
│  BLOB  = 0                  Path  : string            Files : map[string]     │
│  TREE  = 1                  Hash  : string                      IndexEntry    │
│  COMMIT= 2                  Mode  : uint32            Dirs  : map[string]     │
│                             ──────────                          *TempDirNode  │
│  +String() string           +Equal(IndexEntry) bool   Hash  : string         │
│                             +EqualWithoutMode() bool                          │
└──────────────────────────────────────────────────────────────────────────────┘
         ▲                          ▲                          ▲
         │ uses                     │ holds                    │ builds
         │                          │                          │
┌────────┴─────────┐   ┌────────────┴────────────┐  ┌─────────┴──────────────┐
│  «struct»        │   │  «struct»                │  │  «struct»              │
│  Blob            │   │  Index                   │  │  Tree                  │
│  ─────────       │   │  ──────────────          │  │  ──────────────        │
│  Data : []byte   │   │  Entries :               │  │  Nodes : []Node        │
│  Hash : string   │   │   map[string]IndexEntry  │  │  ──────────────        │
│                  │   │  ──────────────          │  │  +WriteTree(*Index)    │
└──────────────────┘   │  +AddIndex(IndexEntry)   │  │   (string, error)      │
                       │  +Equal(*Index) bool      │  │                        │
                       │  +Save(path) error        │  └────────────────────────┘
                       │  +LoadIndex(path)         │
                       │   (*Index, error)         │   ┌────────────────────────┐
                       └──────────────────────────-┘   │  «struct»              │
                                                        │  Node                  │
┌──────────────────────────┐   ┌────────────────────┐  │  ──────────────        │
│  «struct»                │   │  «struct»           │  │  Mode : uint32         │
│  User                    │   │  Head               │  │  Name : string         │
│  ──────────────          │   │  ──────────────     │  │  Hash : string         │
│  Name      : string      │   │  Ref  : string      │  │  ──────────────        │
│  Email     : string      │   │  Hash : string      │  │  +modeString() string  │
│  Timestamp : time.Time   │   │  ──────────────     │  └────────────────────────┘
│  ──────────────          │   │  +Write(dir) error  │
│  +String() string        │   │  +UpdateRef(dir,    │
│  +NewUserFromEnv(bool)   │   │    hash) error      │
│   User                   │   │  +ResolveHead(dir)  │
└──────────────────────────┘   │   (Head, error)     │
                               └─────────────────────┘
┌──────────────────────────────────────────────────────────────────────────────┐
│  «struct»  Commit                                                             │
│  ──────────────────────────────────────────────────────────────              │
│  TreeHash     : string                                                        │
│  ParentHashes : []string                                                      │
│  Author       : User                                                          │
│  Committer    : User                                                          │
│  Message      : string                                                        │
│  ──────────────────────────────────────────────────────────────              │
│  +NewCommit(tree, parents, author, committer, msg) *Commit                   │
│  +Commit() (string, error)           ← hashes & writes to disk               │
│  +CommitTree(tree,parents,author,msg) (string,error)                         │
│  +serialize() []byte                 ← internal                               │
│  +parseCommit([]byte) (*Commit, error)  ← package-level                     │
│  +LoadCommits(path) ([]*Commit, error)  ← package-level                     │
└──────────────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────────────┐
│  core/utils  (stateless infrastructure functions)                            │
│                                                                               │
│  HashObject(data, ObjectType, write) (string, error)                         │
│  ObjectExists(hash, basePath) bool                                            │
│  GetMode(os.FileInfo) uint32                                                  │
│  HashBytes([]byte) string                                                     │
│  HashString(string) string                                                    │
│  ShipHasBeenInitRecursive(paths...) (string, error)                          │
└──────────────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────────────┐
│  «struct»  ConcurrentProcessor                                               │
│  ──────────────────────────────────────────────────────────────              │
│  Workers    : int                                                             │
│  TaskChan   : chan Task                                                       │
│  ResultChan : chan Result                                                     │
│  wg         : sync.WaitGroup                                                 │
│  ctx        : context.Context                                                 │
│  cancelFun  : context.CancelFunc                                             │
│  processFunc: ProcessFunc                                                     │
│  ──────────────────────────────────────────────────────────────              │
│  +Start()                +Stop()      +Cancel()                               │
│  +Feed(Task)             +GetResults() <-chan Result                          │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 6. UML — Sequence Diagrams (per command)

### 6.1 `ship init`

```
User          cmd/init.go       commands/Init()        OS Filesystem
 │                │                    │                     │
 │─ ship init ───►│                    │                     │
 │                │── Init(".") ──────►│                     │
 │                │                    │── os.Getwd() ──────►│
 │                │                    │◄─ "/project" ───────│
 │                │                    │── os.Stat(path) ───►│
 │                │                    │◄─ DirEntry ─────────│
 │                │                    │── os.MkdirAll(.ship/objects) ─►│
 │                │                    │── os.MkdirAll(.ship/refs/heads)►│
 │                │                    │── os.MkdirAll(.ship/refs/tags)─►│
 │                │                    │── WriteFile(.ship/index, {}) ──►│
 │                │                    │── WriteFile(.ship/HEAD, ref:…) ►│
 │                │                    │◄─ nil ──────────────│
 │                │◄─ nil ─────────────│                     │
 │◄─ "Initialized…│                    │                     │
```

### 6.2 `ship add <path>`

```
User          cmd/add.go     commands/Add()      utils/           entities/         OS
 │                │                │               │                  │              │
 │─ ship add . ──►│                │               │                  │              │
 │                │── Add(".") ───►│               │                  │              │
 │                │                │─ShipHasBeenInitRecursive──►│     │              │
 │                │                │◄─ repoBasePath ────────────│     │              │
 │                │                │── LoadIndex(repoBasePath) ─────►│              │
 │                │                │◄─ *Index ────────────────────────│              │
 │                │                │                                                 │
 │                │                │  [for each path]                                │
 │                │                │── filepath.Walk() ─────────────────────────────►│
 │                │                │◄─ (path, FileInfo) ─────────────────────────────│
 │                │                │                                                 │
 │                │                │  [for each file]                                │
 │                │                │── os.ReadFile(path) ───────────────────────────►│
 │                │                │◄─ data []byte ──────────────────────────────────│
 │                │                │── HashObject(data, BLOB, false)─►│             │
 │                │                │◄─ hash (dry-run) ───────────────│             │
 │                │                │── ObjectExists(hash, base) ─────►│             │
 │                │                │   [if not exists]                               │
 │                │                │── HashObject(data, BLOB, true) ─►│ ───write──►│
 │                │                │◄─ hash ─────────────────────────│             │
 │                │                │── index.AddIndex(IndexEntry) ───────────────►│  │
 │                │                │                                                 │
 │                │                │── index.Save(repoBasePath) ────────────────────►│
 │                │◄─ nil ─────────│                                                 │
```

### 6.3 `ship commit <message>`

```
User        cmd/commit.go   commands/Commit()   entities/Tree   entities/Commit  Head    OS
 │               │                 │                  │                │           │      │
 │─ship commit──►│                 │                  │                │           │      │
 │               │─Commit(msg,.)──►│                  │                │           │      │
 │               │                 │─LoadIndex()──────────────────────────────────────►│  │
 │               │                 │◄─*Index ─────────────────────────────────────────│  │
 │               │                 │─WriteTree(index)─►│               │           │      │
 │               │                 │                   │─buildTempDirTree()        │      │
 │               │                 │                   │─writeTreeRecursive()      │      │
 │               │                 │                   │  └─HashObject(TREE,true)──────────►│
 │               │                 │◄─ treeHash ───────│               │           │      │
 │               │                 │─ResolveHead()────────────────────────────────►│      │
 │               │                 │◄─ Head{Ref,Hash}─────────────────────────────│      │
 │               │                 │─NewUserFromEnv(false) ← author               │      │
 │               │                 │─NewUserFromEnv(true)  ← committer            │      │
 │               │                 │─NewCommit(tree,parents,author,committer,msg)──►│    │
 │               │                 │─commit.Commit()──────────────────►│           │      │
 │               │                 │                   │    serialize()│           │      │
 │               │                 │                   │    HashObject(COMMIT,true)────────►│
 │               │                 │◄─commitHash ─────────────────────│           │      │
 │               │                 │─head.UpdateRef(base, commitHash)─────────────►│      │
 │               │◄─nil ───────────│                   │               │           │      │
```

### 6.4 `ship status`

```
User       cmd/status.go   commands/Status()    entities/Index    utils/HashObject
 │               │                 │                   │                  │
 │─ship status──►│                 │                   │                  │
 │               │─Status(".") ───►│                   │                  │
 │               │                 │─LoadIndex(base) ──►│                 │
 │               │                 │◄─*stagedIndex ─────│                 │
 │               │                 │─recalculateIndex(base)               │
 │               │                 │  └─Walk all files                    │
 │               │                 │  └─HashObject(data, BLOB, false) ───►│
 │               │                 │◄─*liveIndex                          │
 │               │                 │                                       │
 │               │                 │  [for each entry in stagedIndex]      │
 │               │                 │   compare hash with liveIndex entry   │
 │               │                 │   match  → print BLUE (staged)        │
 │               │                 │   missing/diff → print YELLOW (unstaged)
 │               │                 │  [for each entry only in liveIndex]   │
 │               │                 │   print YELLOW (untracked)            │
 │◄─coloured list│                 │                                       │
```

### 6.5 `ship cat-file -p <hash>`

```
User          cmd/cat-file.go   commands/CatFile()     OS (.ship/objects)
 │                  │                  │                       │
 │─cat-file -p h───►│                  │                       │
 │                  │─CatFile(h,"-p")─►│                       │
 │                  │                  │─ hash[0:2]=folder     │
 │                  │                  │─ hash[2:]=file        │
 │                  │                  │─ os.ReadFile(.ship/objects/XX/YY)─►│
 │                  │                  │◄─ compressed bytes ───────────────│
 │                  │                  │─ zlib.Decompress()                 │
 │                  │                  │─ split on \0 → header + body       │
 │                  │                  │─ parse header → type + size        │
 │                  │                  │  flag="-p": return "type size body"│
 │                  │◄─ result string ─│                       │
 │◄─ printed output─│                  │                       │
```

### 6.6 `ship pour`

```
User        cmd/pour.go     entities/LoadCommits()    OS
 │               │                  │                  │
 │─ ship pour ──►│                  │                  │
 │               │─LoadCommits(.)──►│                  │
 │               │                  │─getMainRef()     │
 │               │                  │  └─ReadFile(.ship/refs/heads/main)─►│
 │               │                  │◄─ commitHash ────────────────────────│
 │               │                  │─ShipHasBeenInitRecursive()           │
 │               │                  │─ReadFile(.ship/objects/XX/YY) ──────►│
 │               │                  │◄─ zlib bytes ───────────────────────│
 │               │                  │─inflateGitObject()                   │
 │               │                  │─parseCommit()                        │
 │               │                  │  ├─ split on \0                      │
 │               │                  │  ├─ parse lines: tree/parent/author   │
 │               │                  │  └─ parse message                    │
 │               │◄─ []*Commit ─────│                  │
 │◄─ (printed) ──│                  │                  │
```

---

## 7. Data-Flow Diagram

```
                          ┌──────────────────────┐
                          │   Working Directory   │
                          │  (any files / dirs)   │
                          └──────────┬───────────┘
                                     │ ship add <path>
                                     ▼
                     ┌───────────────────────────────┐
                     │      Staging Area              │
                     │   .ship/index  (JSON)          │
                     │   { path → { hash, mode } }    │
                     └──────────────┬────────────────┘
                                    │ ship commit <msg>
                        ┌───────────┴───────────┐
                        │                       │
                        ▼                       ▼
            ┌─────────────────────┐  ┌───────────────────────────┐
            │   Blob Objects      │  │   Tree Objects             │
            │ .ship/objects/XX/YY │  │ .ship/objects/XX/YY        │
            │ (zlib compressed)   │  │ (directory snapshot)       │
            └─────────────────────┘  └──────────────┬────────────┘
                                                     │
                                                     ▼
                                         ┌─────────────────────┐
                                         │   Commit Object      │
                                         │ .ship/objects/XX/YY  │
                                         │ tree, parent(s),     │
                                         │ author, committer,   │
                                         │ message              │
                                         └──────────┬──────────┘
                                                    │
                                                    ▼
                                        ┌───────────────────────┐
                                        │        HEAD            │
                                        │  .ship/HEAD            │
                                        │  ref: refs/heads/main  │
                                        │                        │
                                        │  .ship/refs/heads/main │
                                        │  <40-char commit hash> │
                                        └───────────────────────┘
```

---

## 8. Object Storage Model

Ship uses the same **content-addressable storage** (CAS) model as Git.

### How an Object is Stored

```
Step 1 — Build store payload:
  payload = "<type> <len>\0<content>"
  e.g.:    "blob 5\0hello"

Step 2 — Hash:
  sha1 = SHA1(payload)           → 20 bytes
  hash = hex(sha1)               → 40-char string
  e.g.:  "aabbcc...ff"

Step 3 — Locate bucket:
  folder = hash[0:2]             → "aa"
  file   = hash[2:]              → "bbcc...ff"
  path   = ".ship/objects/aa/bbcc...ff"

Step 4 — Compress & write (only if write=true):
  zlib.Compress(payload) → disk
```

### Object Types

| Type | Content | Hashed From |
|------|---------|-------------|
| `blob` | Raw file bytes | The file content |
| `tree` | Sorted `<mode> <name>\0<hash>` entries | Directory snapshot |
| `commit` | Text fields: tree, parent, author, committer, message | Full commit metadata |

### Content-Addressability Properties

- **Identical content always produces the same hash** — deduplication is automatic.
- **Any bit-flip changes the hash** — tamper detection is built in.
- Objects are **never deleted** (immutable append-only store).
- `ObjectExists()` checks existence before writing → **idempotent adds**.

---

## 9. On-Disk Layout (.ship directory)

After `ship init`:

```
.ship/
├── HEAD                  # "ref: refs/heads/main\n"
├── index                 # JSON staging area
├── refs/
│   ├── heads/
│   │   └── main          # latest commit hash (after first commit)
│   └── tags/             # (reserved, unused)
└── objects/
    ├── a9/
    │   └── 4a8fe5ccb...  # zlib-compressed object file
    ├── 3b/
    │   └── 1f4a2c8e9...
    └── ...
```

After `ship add file.txt` (file contents = "hello"):

```
.ship/
├── index   →  { "Entries": { "file.txt": { "Path":"file.txt","Hash":"aaf4c...","Mode":100644 } } }
└── objects/
    └── aa/
        └── f4c61dcedce3...   ← blob object for "hello"
```

After `ship commit "initial"`:

```
.ship/
├── HEAD                       → "ref: refs/heads/main\n"
├── refs/heads/main            → "e3b0c44298fc..."   ← commit hash
└── objects/
    ├── aa/f4c...              ← blob
    ├── 4b/825dc...            ← tree
    └── e3/b0c442...           ← commit
```

---

## 10. Code Flow — Every Command Explained

### `ship init`

```
main.go
 └─ cmd.Execute()
     └─ initCmd.Run()                         [cmd/init.go]
         └─ commands.Init(".")               [commands/init.go]
             ├─ os.Getwd()                   → resolve "." to absolute
             ├─ os.Stat(path)                → validate path is a directory
             ├─ check if .ship/ already exists → if yes, print "Reinitialized"
             ├─ os.MkdirAll(".ship/objects") → content store
             ├─ os.MkdirAll(".ship/refs/heads")
             ├─ os.MkdirAll(".ship/refs/tags")
             ├─ os.WriteFile(".ship/index", `{"entries":{}}`)
             └─ os.WriteFile(".ship/HEAD", "ref: refs/heads/main\n")
```

**Key decisions:**
- Initialization is **idempotent** — running `ship init` twice is safe.
- Index starts with empty JSON so `LoadIndex` never fails on a fresh repo.
- HEAD always points to `main` branch out of the box.

---

### `ship add [paths...]`

```
cmd/add.go → commands.Add(paths...)           [commands/add.go]
  │
  ├─ utils.ShipHasBeenInitRecursive(paths)    → walk up dir tree to find .ship/
  ├─ filepath.EvalSymlinks(repoBasePath)      → canonicalise symlinks
  ├─ entities.LoadIndex(repoBasePath)         → read .ship/index into *Index
  ├─ getRepoRelativePath(base, paths)         → resolve each arg to abs + rel path
  │
  └─ for each path:
      processPath(base, path, index, existingFiles)
        └─ filepath.Walk(path, func):
            ├─ skip .ship/ directory entirely
            ├─ skip directories (process only files)
            ├─ os.ReadFile(path) → raw bytes
            ├─ utils.HashObject(data, BLOB, false)  → dry-run hash
            ├─ utils.ObjectExists(hash, base)       → deduplication check
            │    └─ if already stored → skip
            ├─ utils.HashObject(data, BLOB, true)   → write blob to .ship/objects
            ├─ filepath.Rel(base, path)             → repo-relative path
            └─ index.AddIndex(IndexEntry{Path, Hash, Mode})
  │
  ├─ prune index entries for deleted files
  └─ index.Save(repoBasePath)                 → write updated JSON index
```

**Key decisions:**
- Two-phase hash (dry-run first, then write) avoids writing duplicates.
- Symlinks are resolved before path comparison.
- Deleted files are removed from the index automatically.
- File mode: `100644` (regular) or `100755` (executable).

---

### `ship commit <message>`

```
cmd/commit.go → commands.Commit(msg, path)     [commands/commit.go]
  │
  ├─ utils.ShipHasBeenInitRecursive(path)
  ├─ entities.LoadIndex(base)                  → staging area
  ├─ guard: len(index.Entries) == 0 → error "nothing to commit"
  │
  ├─ tree := entities.Tree{}
  ├─ treeHash, _ := tree.WriteTree(index)
  │    └─ buildTempDirTree(index)              → build in-memory directory tree
  │         [for each IndexEntry]
  │           split path by "/"
  │           walk TempDirNode tree, creating nodes for each dir segment
  │           place file at leaf
  │    └─ writeTreeRecursive(root)             → post-order traversal
  │         for each dir node  → recurse, get child hash
  │         serializeTree(node):
  │           collect File nodes + Dir nodes
  │           sort by Name (alphabetical)
  │           write "<mode> <name>\0<hash>" for each
  │         HashObject(treeBytes, TREE, true)  → write tree object
  │
  ├─ entities.ResolveHead(base)               → get current HEAD hash
  ├─ parents = [head.Hash]  (or [] for first commit)
  │
  ├─ author    = entities.NewUserFromEnv(false)
  ├─ committer = entities.NewUserFromEnv(true)
  │
  ├─ commit := entities.NewCommit(treeHash, parents, author, committer, msg)
  ├─ commitHash, _ := commit.Commit()
  │    └─ commit.serialize()                  → text bytes
  │    └─ utils.HashObject(content, COMMIT, true) → write commit object
  │
  └─ head.UpdateRef(base, commitHash)         → write hash to refs/heads/main
```

**Key decisions:**
- Tree is built via **recursive post-order traversal** (children first, then parent) so each parent tree hash includes all child hashes — exact Merkle tree semantics.
- The tree is serialized with **sorted entries** for deterministic hashing.
- Author/Committer are read from env vars with fallback to `$USER@localhost`.

---

### `ship status`

```
cmd/status.go → commands.Status(path)          [commands/status.go]
  │
  ├─ utils.ShipHasBeenInitRecursive(path)
  ├─ filepath.EvalSymlinks(base)
  ├─ entities.LoadIndex(base)                  → stagedIndex (what was `add`-ed)
  │
  ├─ recalculateIndex(base)                    → liveIndex (current FS state)
  │    └─ filepath.WalkDir(base)
  │         skip .ship/
  │         for each file:
  │           os.ReadFile → data
  │           HashObject(data, BLOB, false)    → hash WITHOUT writing
  │           index.AddIndex(...)
  │
  ├─ for path in stagedIndex:
  │    if NOT in liveIndex           → print YELLOW (deleted/unstaged)
  │    elif hashes match             → print BLUE   (staged, unchanged)
  │    elif hashes differ            → print YELLOW (modified/unstaged)
  │
  └─ for path in liveIndex NOT in stagedIndex:
       print YELLOW  (untracked)
```

**Key decisions:**
- Status does a **full re-hash** of the working tree (no mtime shortcuts).
- Colour output: **Blue = staged/clean**, **Yellow = unstaged/modified/deleted/untracked**.
- The comparison is strictly by hash (content-aware), not by timestamp.

---

### `ship cat-file [-p|-t|-c|-s] <hash>`

```
cmd/cat-file.go → commands.CatFile(hash, flag)  [commands/cat-file.go]
  │
  ├─ folder = hash[0:2]
  ├─ file   = hash[2:]
  ├─ os.Stat(".ship/objects/<folder>/<file>")  → verify exists
  ├─ os.ReadFile(path)                         → compressed bytes
  ├─ zlib.NewReader + io.ReadAll               → decompressed payload
  ├─ bytes.SplitN(payload, \0, 2)              → [header, body]
  ├─ strings.Split(header, " ")               → [type, size]
  │
  └─ switch flag:
      -p  → "type size body"   (pretty)
      -s  → size only
      -t  → body only          (tree/commit raw)
      default → full decompressed payload
```

---

### `ship pour`

```
cmd/pour.go → entities.LoadCommits(path)       [core/entities/Commit.go]
  │
  ├─ getMainRef(path)
  │    └─ os.ReadFile(".ship/refs/heads/main")  → commit hash
  ├─ utils.ShipHasBeenInitRecursive(path)
  ├─ os.ReadFile(".ship/objects/XX/YY")         → compressed commit object
  ├─ inflateGitObject(data)                     → zlib decompress
  └─ parseCommit(object)
       ├─ find \0 → skip header
       ├─ split on \n
       ├─ for each line: parse tree/parent/author/committer
       └─ everything after blank line = message
```

**Note:** Currently only loads the **latest commit** (HEAD), not the full chain. Full ancestry traversal is marked for future work.

---

## 11. Core Algorithms & Strategies

### 11.1 Content-Addressed Storage (CAS)

**Strategy:** SHA-1 of `header + \0 + content` → 40-hex key.  
**Why:** The hash IS the identity. Store once, reference forever. Perfect for snapshots.

```
HashObject(data, type, write):
  header = "<type> <len>"
  payload = header + "\0" + data
  hash = hex(sha1(payload))
  if write: zlib_compress(payload) → .ship/objects/<hash[0:2]>/<hash[2:]>
  return hash
```

### 11.2 Merkle Tree for Directory Snapshots

**Strategy:** Tree objects reference blob hashes for files and sub-tree hashes for directories. Any change anywhere propagates to the root tree hash.

```
root-tree-hash
  ├── src-tree-hash
  │     ├── main.go-blob-hash
  │     └── util.go-blob-hash
  └── README.md-blob-hash
```

**Implementation:** `buildTempDirTree()` converts the flat index map into an in-memory nested `TempDirNode` structure. `writeTreeRecursive()` does post-order traversal: children are hashed first, then their hashes are embedded in the parent tree serialization.

### 11.3 Repository Discovery (Recursive Parent Search)

**Strategy:** Walk up the directory tree looking for `.ship/` — identical to Git's approach.

```go
// searchBaseRepo() in ship-core.go
absPath = filepath.Abs(givenPath)
splits = absPath.split("/")
for i = len(splits); i > 0; i--:
    candidate = join(splits[:i])
    if Stat(candidate + "/.ship").IsDir():
        return candidate
```

This means `ship add src/main.go` works from any subdirectory of the repo.

### 11.4 Two-Phase Hashing in Add

Before writing a blob, `Add()` performs a **dry-run hash** (write=false) to compute the hash, then checks if the object already exists on disk. Only if missing does it call `HashObject` again with `write=true`. This prevents redundant disk I/O and preserves idempotency.

### 11.5 Index as JSON Staging Area

Instead of Git's complex binary index format, Ship stores the staging area as a JSON file:

```json
{
  "Entries": {
    "src/main.go": { "Path": "src/main.go", "Hash": "abc123", "Mode": 100644 }
  }
}
```

This trades raw performance for human readability and simplicity. `json.MarshalIndent` is used so the file is inspectable with any text editor.

### 11.6 Status via Full Re-Hash

`Status()` does **not** use timestamps or inode metadata. It re-hashes every file in the working directory and compares content hashes against the index. This is correct but O(n·filesize). Git optimises this with stat-caching in the binary index.

---

## 12. Concurrency Infrastructure

`core/utils/concurrent-processor.go` provides a **generic worker pool** — currently infrastructure for future use (e.g., parallel blob hashing during large `add` operations).

```
                  ┌─────────────────────────────────┐
 Feed(Task) ─────►│         TaskChan (buffered)      │
                  └──────────────┬──────────────────┘
                                 │
                    ┌────────────┼────────────┐
                    ▼            ▼            ▼
                 worker       worker       worker    (N goroutines)
                    │            │            │
                    └────────────┼────────────┘
                                 │
                  ┌──────────────▼──────────────────┐
 GetResults() ◄──│        ResultChan                │
                  └─────────────────────────────────┘
```

### API

```go
cp := utils.NewConcurrentProcessor(4, func(ctx context.Context, task Task) Result {
    // process task.Data
    return Result{JobID: task.Id, Data: result}
})
cp.Start()          // launch N goroutines
cp.Feed(task)       // send work
results := cp.GetResults()
cp.Stop()           // drain and close
cp.Cancel()         // context-cancel immediately
```

### Lifecycle

| Method | Behaviour |
|--------|-----------|
| `Start()` | Spawns `Workers` goroutines, each blocking on `TaskChan` |
| `Feed(Task)` | Sends task into channel; blocks if channel is full |
| `Stop()` | Closes `TaskChan`, waits for all goroutines via WaitGroup, closes `ResultChan` |
| `Cancel()` | Cancels context (immediate halt), waits, closes `ResultChan` |

---

## 13. CLI Reference

### Global

```
ship [command] [flags]
```

### Commands

#### `ship init`

Initialize a new Ship repository in the current directory.

```bash
ship init
```

- Creates `.ship/`, `.ship/objects/`, `.ship/refs/heads/`, `.ship/refs/tags/`
- Creates `.ship/index` with `{"entries":{}}`
- Creates `.ship/HEAD` with `ref: refs/heads/main`
- **Idempotent**: Running twice prints "Reinitialized existing Ship repository"

---

#### `ship add [files...]`

Stage file contents to the index.

```bash
ship add .                    # stage everything in current directory
ship add src/main.go          # stage a specific file
ship add src/ tests/          # stage multiple paths
```

- Accepts one or more file/directory paths
- Walks directories recursively
- Skips `.ship/` directory
- Hashes and stores blobs; updates index
- Removes index entries for files that no longer exist
- **Requires** repo to be initialized (walks up the tree to find `.ship/`)

---

#### `ship commit <message>`

Record staged changes as a commit.

```bash
ship commit "initial commit"
ship commit "fix: handle nil pointer in status"
```

- Requires at least one staged file
- Builds a tree object from the current index
- Reads author/committer from environment variables
- Creates a commit object pointing to the tree and any parent commits
- Updates `refs/heads/main` to the new commit hash

---

#### `ship status [path]`

Show the working tree status.

```bash
ship status             # uses current working directory
ship status /repo/path  # explicit path
```

**Output colours:**

| Colour | Meaning |
|--------|---------|
| 🔵 Blue | File is staged and unchanged |
| 🟡 Yellow | File is modified, deleted, or untracked |

---

#### `ship cat-file <flags> <hash>`

Print raw content of a stored object by its SHA-1 hash.

```bash
ship cat-file -p a94a8fe5ccb19ba61c4c0873d391e987982fbbd3  # pretty print
ship cat-file -s a94a8fe5ccb19ba61c4c0873d391e987982fbbd3  # size only
ship cat-file -t a94a8fe5ccb19ba61c4c0873d391e987982fbbd3  # tree/raw body
ship cat-file -c a94a8fe5ccb19ba61c4c0873d391e987982fbbd3  # commit body
```

| Flag | Long | Description |
|------|------|-------------|
| `-p` | `--pretty` | Pretty-print: `<type> <size> <content>` |
| `-s` | `--size` | Print object size in bytes only |
| `-t` | `--tree` | Print raw body (tree content) |
| `-c` | `--commit` | Print raw body (commit content) |

---

#### `ship pour`

Show commit history (log).

```bash
ship pour
```

- Reads HEAD → resolves to latest commit hash
- Loads and decompresses the commit object
- Prints commit metadata
- (Currently shows only HEAD commit; full ancestry walk is pending)

---

## 14. Dependency Graph

### External Dependencies

```
github.com/spf13/cobra v1.10.2       ← CLI framework
github.com/spf13/pflag v1.0.9        ← flag parsing (transitive via cobra)
github.com/inconshreveable/mousetrap ← Windows Ctrl+C helper (transitive)
```

### Standard Library Usage

| Package | Used For |
|---------|---------|
| `crypto/sha1` | Content-addressing hash |
| `compress/zlib` | Object compression/decompression |
| `encoding/json` | Index serialisation |
| `path/filepath` | Cross-platform path ops |
| `io/fs` | Directory walking |
| `os` | File I/O |
| `bytes` | Buffer building, byte scanning |
| `sync` | WaitGroup in worker pool |
| `context` | Cancellation in worker pool |
| `time` | Commit timestamps |
| `fmt` | Formatting |
| `log` | Debug logging in entities |

---

## 15. Testing Strategy

### Test Layout

```
tests/
├── command_tests/   ← integration tests (spin up real .ship repos in TempDir)
│   ├── init_test.go
│   ├── add_test.go
│   ├── commit_test.go
│   ├── status_test.go
│   ├── cat-file_test.go
│   └── pour_test.go
└── helpers/
    ├── setup.helper.go   ← TempDir scaffolding, BurnDown cleanup
    └── assert.helper.go  ← assertion helpers
```

### Test Helpers

| Helper | Description |
|--------|-------------|
| `helpers.Setup(t)` | Returns `SetupInfo{RepoDir, IncorrectDir, RandomFilePath}` using `t.TempDir()` |
| `helpers.BurnDown(t)` | Cleans up temp dirs + leftover `.ship/` |
| `helpers.WriteFile(t, dir, name, content)` | Creates a file (makes parent dirs) |
| `helpers.WriteDir(t, dir, name)` | Creates a directory |
| `helpers.DeleteFile(t, dir, name)` | Removes a file |
| `helpers.AppendToFile(path, content)` | Appends content to a file |
| `helpers.AssertNil(err)` | Panics if err ≠ nil |
| `helpers.AssertNotNil(val)` | Panics if val = nil |
| `helpers.AssertEqual(t, expected, actual)` | t.Fatalf on mismatch |
| `helpers.AssertNotEqual(t, expected, actual)` | t.Fatalf if equal |
| `helpers.AssertExists(t, path)` | Fails if file/dir doesn't exist on disk |
| `helpers.AssertFileInIndex(t, repoDir, filename)` | Fails if file not in index |
| `helpers.AssertEqualIndex(t, expected, actual)` | Deep-equals two Index objects |
| `helpers.AssertNotEqualIndex(t, expected, actual)` | Deep-not-equal two Index objects |

### Test Categories

#### Init Tests (`init_test.go`)

| Test | What it validates |
|------|-------------------|
| `TestInit_forCurrentDir_itShouldCreateBasicDirStructure` | All 4 required paths exist after init |
| `TestInit_forExistingDir_initShouldBeIdempotent` | Second init returns no error |
| `TestInit_forNonExistingDir_itShouldFail` | Non-existent path returns error |
| `TestInit_passingNonDirPath_itShouldFail` | File path returns error |
| `TestInit_forEmptyDir_itShouldPass` | Empty directory initialises cleanly |

#### Add Tests (`add_test.go`)

| Test | What it validates |
|------|-------------------|
| `TestAdd_ToEmptyInitDir_ShouldPass` | Add on init'd empty dir succeeds |
| `TestAdd_ToCurrentInitDir_ShouldPass` | Add "." works from within repo dir |
| `TestAddTwice_ToInitDir_ShouldPass` | Idempotent add |
| `TestAdd_ToUninitializedDir_ShouldFail` | Error without init |
| `TestAdd_WithNoFiles_ShouldPass` | Empty index when no files |
| `TestAdd_WithNestedDirectories_ShouldPass` | Deep directory trees |
| + more | Hash verification, deletion, update scenarios |

#### Commit Tests (`commit_test.go`)

| Test | What it validates |
|------|-------------------|
| `TestCommit_NewCommit_ShouldCreateValidCommit` | Struct fields set correctly |
| `TestCommit_WithNoParents_ShouldPass` | First commit (no parent) |
| `TestCommit_WithSingleParent_ShouldPass` | Subsequent commit with parent |
| + more | Serialization format, round-trip parse, full workflow |

#### Cat-File Tests (`cat-file_test.go`)

| Test | What it validates |
|------|-------------------|
| `TestCatFile_blobHash_ShouldPass` | `-p` flag returns `"blob 5 hello"` |
| `TestCatFile_WithSizeFlag_ShouldReturnSize` | `-s` returns size string |
| `TestCatFile_MultipleBlobsWithDifferentContent_ShouldPass` | Each file gets unique hash |
| + more | Edge cases, flag combinations |

#### Pour Tests (`pour_test.go`)

| Test | What it validates |
|------|-------------------|
| `Test_pour` | Full init→add→commit→LoadCommits round-trip |

### Running Tests

```bash
# All tests
make test

# With coverage
make test-cover

# Single package
go test -v ./tests/command_tests/...

# Single test
go test -v -run TestCommit_WithNoParents_ShouldPass ./tests/...
```

---

## 16. Build, Install & Run

### Prerequisites

- Go 1.24+
- `make` (GNU Make)

### Quick Start

```bash
# Clone / enter the project
cd ship

# Build binary → bin/ship
make build

# Install system-wide → /usr/local/bin/ship
make install

# Try it
mkdir /tmp/myrepo && cd /tmp/myrepo
ship init
echo "hello world" > hello.txt
ship add .
ship status
ship commit "first commit"
ship pour
ship cat-file -p $(cat .ship/refs/heads/main)
```

### Makefile Targets

| Target | Description |
|--------|-------------|
| `make build` | Compile → `bin/ship` with version + build-time LDFLAGS |
| `make install` | build + sudo copy to `/usr/local/bin/` |
| `make uninstall` | Remove `/usr/local/bin/ship` |
| `make dev` | build + immediately run `bin/ship` |
| `make run` | `go run .` (no build artifact) |
| `make test` | `go test -v ./tests/...` |
| `make test-cover` | Tests + HTML coverage report (`coverage.html`) |
| `make clean` | Remove `bin/`, `coverage.out`, `coverage.html` |
| `make deps` | `go mod download && go mod verify` |
| `make tidy` | `go mod tidy` |
| `make fmt` | `gofmt` all sources |

### Alternative: Shell Installer

```bash
# Local install
bash install.sh

# Remote install
curl -fsSL https://raw.githubusercontent.com/KambojRajan/ship/main/remote-install.sh | bash
```

---

## 17. Environment Variables

These environment variables control author/committer identity in commits.  
All fall back gracefully if unset.

| Variable | Fallback | Description |
|----------|---------|-------------|
| `SHIP_AUTHOR_NAME` | `$USER` or `"unknown"` | Author display name |
| `SHIP_AUTHOR_EMAIL` | `$SHIP_AUTHOR_NAME@localhost` | Author email |
| `SHIP_AUTHOR_DATE` | `time.Now()` | Author timestamp (RFC 3339) |
| `SHIP_COMMITTER_NAME` | same as author | Committer display name |
| `SHIP_COMMITTER_EMAIL` | same as author | Committer email |
| `SHIP_COMMITTER_DATE` | `time.Now()` | Committer timestamp (RFC 3339) |

**Example:**

```bash
export SHIP_AUTHOR_NAME="Alice"
export SHIP_AUTHOR_EMAIL="alice@example.com"
export SHIP_AUTHOR_DATE="2026-03-15T10:00:00Z"
ship commit "authored commit"
```

---

## 18. Design Patterns Used

### 1. Layered Architecture (Clean Architecture lite)

```
cmd → commands → entities + utils
```

Each layer has a single responsibility. `cmd/` is only about CLI wiring. `commands/` owns workflows. `entities/` owns domain state. `utils/` is pure infrastructure.

### 2. Repository Pattern

`entities/index.go` acts as a repository for staging entries:
- `LoadIndex(path)` — read from disk
- `index.Save(path)` — persist to disk
- `index.AddIndex(entry)` — in-memory mutation

### 3. Value Objects

`IndexEntry`, `User`, `Node` are value objects — compared by content, not reference. `IndexEntry.Equal()` and `Index.Equal()` implement structural equality.

### 4. Factory Functions

```go
entities.NewCommit(...)
entities.NewUserFromEnv(isCommitter)
entities.NewIndex()
utils.NewConcurrentProcessor(workers, fn)
```

All constructors are explicit factory functions, never raw struct literals in business code.

### 5. Strategy Pattern (ObjectType)

`common.ObjectType` is an enum with a `String()` method. `HashObject()` takes it as a parameter, enabling the same storage pipeline for blobs, trees, and commits:

```go
utils.HashObject(data, common.BLOB, true)
utils.HashObject(treeBytes, common.TREE, true)
utils.HashObject(commitBytes, common.COMMIT, true)
```

### 6. Template Method (serializeTree + writeTreeRecursive)

`WriteTree()` defines the algorithm skeleton:
1. Build in-memory tree (`buildTempDirTree`)
2. Recursively write nodes (`writeTreeRecursive`)
3. Serialize each node (`serializeTree`)

Each step is a separate function, following the Template Method spirit.

### 7. Worker Pool Pattern

`ConcurrentProcessor` encapsulates the goroutine lifecycle: channels, WaitGroup, context cancellation — behind a clean `Start/Feed/Stop/Cancel` API.

### 8. Idempotency by Design

- `ship init` twice → safe
- `ship add` same file twice → same blob hash, ObjectExists guard, index entry overwritten
- Object storage is **append-only** (existing objects are never overwritten)

---

## 19. Comparison with Git

| Feature | Ship | Git |
|---------|------|-----|
| Object model | SHA-1 CAS, zlib | SHA-1 CAS, zlib (identical) |
| Object types | blob, tree, commit | blob, tree, commit, tag |
| Index format | JSON | Custom binary (version 2/3/4) |
| Tree serialisation | mode+name+\0+hash-string | mode+name+\0+hash-bytes (binary) |
| Hash algorithm | SHA-1 | SHA-1 (SHA-256 in new repos) |
| Branching | Single `main` | Full ref namespace |
| Merge/Rebase | ❌ | ✅ |
| Remote/Push/Pull | ❌ | ✅ |
| Packfiles | ❌ | ✅ |
| Submodules | ❌ | ✅ |
| Hooks | ❌ | ✅ |
| Author from env | Custom `SHIP_*` vars | `GIT_AUTHOR_*` vars |
| Status colours | Blue/Yellow | Red/Green |
| Config file | ❌ (env only) | `.git/config` |
| Log command | `ship pour` | `git log` |

---

## 20. Known Limitations & Future Scope

### Current Limitations

1. **`pour` only shows HEAD commit** — parent traversal not yet implemented.
2. **No diff command** — cannot see what changed between commits.
3. **Single branch only** — always commits to `main`, no branch creation/switching.
4. **No `.shipignore`** — equivalent of `.gitignore` is absent; all files are staged.
5. **Status doesn't track truly committed state** — compares live FS against staged index, not against last commit.
6. **cat-file reads from CWD** — path to `.ship/objects` is relative, not anchored to repo root.
7. **ConcurrentProcessor unused** — built but not wired into `add`.
8. **No tag support** — `refs/tags/` directory is created but never used.

### Roadmap

| Feature | Description |
|---------|-------------|
| `ship log` | Full ancestry walk, pretty commit history |
| `ship diff` | Line-level diff between working tree and last commit |
| `ship branch` | Create/list/switch branches |
| `ship checkout` | Restore working tree from a commit |
| `.shipignore` | Pattern-based file exclusions |
| Parallel add | Wire `ConcurrentProcessor` into `Add()` for large repos |
| SHA-256 support | Upgrade to modern hash |
| Packfile | Compact object storage for large histories |
| `ship remote` | Push/pull over HTTP |

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────────┐
│                    SHIP VCS CHEAT SHEET                         │
├──────────────┬──────────────────────────────────────────────────┤
│ ship init    │ Create a new repo in current directory           │
│ ship add .   │ Stage all files                                  │
│ ship add f   │ Stage file/directory f                           │
│ ship commit  │ ship commit "your message"                       │
│ ship status  │ Show staged vs unstaged files                    │
│ ship pour    │ Show commit history                              │
│ ship cat-file│ ship cat-file -p <hash>   pretty print           │
│              │ ship cat-file -s <hash>   size                   │
│              │ ship cat-file -t <hash>   tree body              │
├──────────────┴──────────────────────────────────────────────────┤
│ On-disk:  .ship/HEAD  .ship/index  .ship/objects/               │
│ Env:      SHIP_AUTHOR_NAME  SHIP_AUTHOR_EMAIL  SHIP_AUTHOR_DATE │
│ Build:    make build   make install   make test                 │
└─────────────────────────────────────────────────────────────────┘
```

---

*Documentation generated for Ship VCS v1.0.0 — March 2026*

