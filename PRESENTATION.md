# Technical Interview Presentation -- Intuit Demo Round
## Rajan Kamboj

---

### Slide 1 -- Introduction

**Rajan Kamboj | Back-End Engineer**

- Back-end engineer specializing in high-throughput distributed systems and financial infrastructure
- Currently building a cash reconciliation engine processing daily transactions for 65,000+ SBI ATMs
- Final-year Computer Engineering, NIT Kurukshetra (GPA 8.71/10)
- Core strengths: Java/Spring Boot, Go, system design, data integrity at scale
- Built production systems handling billions in daily transaction volume
- Open-source contributor -- authored a Git-style VCS from scratch in Go

---

### Slide 2 -- Why Intuit

**Mission Alignment & Engineering Culture**

- Intuit powers the financial backbone for 100M+ customers -- that mission resonates with my work in financial infrastructure
- Engineering culture of "Assessing for Awesome" aligns with how I approach problems: first principles, deep ownership
- Intuit's investment in developer platforms and internal tooling matches my passion for building foundational systems
- Opportunity to work at the intersection of scale, reliability, and real financial impact
- Strong track record of turning complex financial domains into simple, powerful products (TurboTax, QuickBooks)

---

### Slide 3 -- Why This Role

**Skills Match & Growth Trajectory**

- Direct experience building financial processing engines -- I understand reconciliation, idempotency, and data consistency
- Proven ability to design systems from scratch (VCS, recon engine) -- not just glue existing tools together
- Immediate contribution: Spring Boot, distributed systems, database optimization, event-driven architecture
- Growth goal: deepen expertise in large-scale platform engineering and observability at Intuit's scale
- I want to build systems that millions of people depend on without ever thinking about

---

### Slide 4 -- Project Overview: Ship VCS

**A Git-Style Version Control System, Built From Scratch**

- Ground-up reimplementation of Git's core internals in Go -- not a wrapper, not a fork
- Implements content-addressable storage, Merkle trees, staging area, commit chains, and object inspection
- Serves as a deep exploration of how distributed version control actually works at the data structure level
- Solo project: designed the architecture, wrote every line, built the CI/CD pipeline, and packaged for multi-platform release
- 2,500+ lines of Go, 42 integration tests, cross-platform builds for Linux/macOS/Windows

---

### Slide 5 -- Defining the Problem

**Why Build a VCS From Scratch?**

- Most engineers use Git daily but have no mental model of its internals -- this limits debugging and architectural thinking
- Existing educational resources explain Git conceptually but rarely implement the actual data structures
- No lightweight, readable VCS implementation existed in Go with production-grade tooling (CI, packaging, testing)
- Goal: build something real enough to expose every design decision Git makes -- hashing, compression, tree construction
- Constraint: zero external dependencies for core logic -- only stdlib (crypto, compress, encoding)

---

### Slide 6 -- Key Goals

**Non-Negotiable Engineering Requirements**

- **Correctness**: SHA-1 content addressing must guarantee data integrity -- any bit flip must change the hash
- **Determinism**: identical working trees must always produce identical commit hashes across runs and platforms
- **Idempotency**: every operation (init, add, commit) must be safe to repeat without side effects
- **Maintainability**: clean layered architecture -- CLI parsing, business logic, and domain model strictly separated
- **Testability**: every command must be verifiable through integration tests with real filesystem state
- **Portability**: static binaries with zero runtime dependencies for Linux, macOS, and Windows

---

### Slide 7 -- Tech Stack

**Minimal Dependencies, Maximum Control**

- **Language**: Go 1.24 -- chosen for static typing, fast compilation, excellent stdlib, and native concurrency primitives
- **CLI Framework**: Cobra -- industry-standard Go CLI library (only external dependency for core logic)
- **Storage**: Filesystem-based content-addressable store -- no database, objects stored as zlib-compressed files
- **Hashing**: crypto/sha1 (stdlib) | **Compression**: compress/zlib (stdlib) | **Serialization**: encoding/json (stdlib)
- **CI/CD**: GitHub Actions (matrix: 3 OS x Go 1.24) + GoReleaser (tar.gz, deb, rpm, Homebrew tap)
- **Testing**: Go's native testing framework + custom assertion helpers | Race detector enabled in CI

---

### Slide 8 -- System Architecture

**Clean Layered Architecture With Strict Dependency Flow**

- **Layer 1 -- CLI** (`cmd/`): Pure Cobra wiring -- parses args/flags, zero business logic, delegates downward
- **Layer 2 -- Commands** (`commands/`): Workflow orchestration -- pure Go functions with no CLI dependency, callable from any context
- **Layer 3 -- Domain** (`core/entities/`): Rich domain objects -- Blob, Tree, Commit, Index, Head with serialization and persistence
- **Layer 4 -- Infrastructure** (`core/utils/`): Stateless utilities -- hashing, object storage, repo discovery, concurrent worker pool
- Dependency flow is strictly unidirectional: `main -> cmd -> commands -> entities + utils` -- no circular imports
- Mirrors Git's own internal architecture: porcelain (cmd) vs. plumbing (core)

---

### Slide 9 -- Data Modeling & Core Logic

**Content-Addressable Storage + Merkle Trees**

- **Object format**: `<type> <size>\0<content>` -> SHA-1 hash -> zlib compress -> store at `.ship/objects/<hash[0:2]>/<hash[2:]>`
- **Three object types**: Blob (file content), Tree (directory snapshot with sorted entries), Commit (metadata + tree hash + parent chain)
- **Merkle tree construction**: flat index -> virtual directory tree -> post-order recursive serialization -> deterministic hashing
- **Index format**: JSON (`{path: {Path, Hash, Mode}}`) -- deliberately human-readable vs. Git's binary format (debugging tradeoff)
- **Deduplication is automatic**: identical content produces identical hashes -- `ObjectExists()` prevents redundant writes
- **Commit chain**: each commit stores parent hash(es), forming a linked-list history anchored at `.ship/refs/heads/main`

---

### Slide 10 -- Security & Reliability

**Integrity by Design, Not by Afterthought**

- **Content integrity**: SHA-1 of `header + content` means any corruption -- even a single bit flip -- changes the hash (tamper detection)
- **Immutable object store**: objects are append-only, never modified or deleted -- eliminates an entire class of data corruption bugs
- **Symlink resolution**: `filepath.EvalSymlinks()` called before path operations to prevent symlink-based path traversal attacks
- **Internal metadata isolation**: `.ship/` directory explicitly skipped during add/status to prevent indexing internal state
- **Idempotent operations**: init, add, and object writes all check existing state before acting -- safe to retry on failure
- **Error chain propagation**: centralized error format strings with `%w` wrapping for proper Go error chain debugging

---

### Slide 11 -- Performance Optimizations

**Deliberate Tradeoffs With Clear Rationale**

- **Two-phase hashing in `add`**: dry-run hash first (write=false), then existence check, then write only if new -- avoids redundant disk I/O
- **Object existence check**: `os.Stat()` called before `os.Create()` to skip already-stored objects -- critical for large repos
- **Static binaries**: `CGO_ENABLED=0` produces fully static executables -- zero runtime dependencies, instant startup
- **Binary stripping**: `-s -w` ldflags remove debug info and DWARF tables -- smaller binary footprint for distribution
- **Built but not yet wired**: `ConcurrentProcessor` worker pool with channels + WaitGroup -- designed for future parallel blob hashing
- **Known limitation I'd fix next**: `status` does a full re-hash of every file (no mtime cache) -- O(n * filesize) vs. Git's O(1) stat check

---

### Slide 12 -- Engineering Challenges

**Three Problems That Shaped the Design**

- **Challenge 1 -- Deterministic Tree Hashing**
  - Problem: tree objects must hash identically regardless of file insertion order
  - Solution: alphabetical sort of entries before serialization -- same approach Git uses
  - Result: identical working trees always produce identical commit hashes across platforms

- **Challenge 2 -- Merkle Tree From Flat Index**
  - Problem: index stores flat file paths but trees require nested directory structure
  - Solution: built `TempDirNode` virtual filesystem -- split paths on `/`, construct tree, then post-order recursive serialization
  - Result: clean separation between index representation and tree object construction

- **Challenge 3 -- Repository Discovery**
  - Problem: commands must work from any subdirectory, not just repo root
  - Solution: walk up the directory tree checking for `.ship/` at each level -- identical to Git's approach
  - Result: seamless UX -- `ship status` works from any nested directory

---

### Slide 13 -- Code Quality & Testing

**42 Integration Tests, Zero Mocks**

- **Testing philosophy**: every test creates a real `.ship/` repo in `t.TempDir()`, runs actual commands, and verifies filesystem state
- **Zero mocking**: all tests exercise the full code path end-to-end -- if it passes, it works
- **Edge case coverage**: empty files, 1MB files, unicode filenames, emoji filenames, corrupted index, binary content, deeply nested dirs
- **Custom assertion helpers**: `AssertNil`, `AssertExists`, `AssertFileInIndex` -- domain-specific test vocabulary
- **CI pipeline**: GitHub Actions matrix (Ubuntu + macOS + Windows) with race detector (`-race`) enabled on every run
- **Code quality**: `gofmt` + `go vet` enforced via `make lint` | GoReleaser automates multi-platform release packaging

---

### Slide 14 -- Key Lessons Learned

**What Building a VCS Taught Me**

- **Technical**: content-addressable storage is one of the most elegant patterns in computer science -- once you understand it, you see it everywhere (Docker, IPFS, blockchain)
- **Technical**: Merkle trees aren't just theory -- implementing one forced me to think deeply about determinism, ordering, and hash propagation
- **Architecture**: strict layered architecture pays off immediately -- I could swap Cobra for any CLI framework without touching business logic
- **Testing**: integration tests with real filesystem state catch bugs that unit tests with mocks never would
- **Leadership**: writing DEEP_DIVE.md and CONTRIBUTING.md forced me to articulate decisions I'd made intuitively -- documentation is design review

---

### Slide 15 -- Future Improvements

**What I'd Build Next**

- **Branching & checkout**: implement ref management, HEAD detachment, and working tree reconstruction from any commit
- **Delta compression & packfiles**: current object-per-file model doesn't scale -- packfiles with delta encoding would reduce storage 10-100x
- **Observability**: wire up OpenTelemetry traces (already planned) to instrument every command with span-level profiling via Jaeger
- **Parallel operations**: activate the `ConcurrentProcessor` worker pool for `add` -- benchmark against sequential for repos with 10K+ files
- **Remote protocol**: implement `ship push` / `ship pull` over HTTP -- the hardest unsolved problem in distributed VCS

---

### Slide 16 -- Closing Summary

**Why I'm the Right Fit**

- **Engineering capability**: I build systems from first principles -- content-addressable stores, Merkle trees, reconciliation engines, event-driven architectures. I don't just use tools, I understand how they work underneath.

- **Culture fit**: I write thorough documentation, I test edge cases obsessively, I make deliberate architectural decisions and can explain the tradeoffs. I care about code that other engineers want to work on.

- **Why I'm excited about Intuit**: I've spent the last year building financial infrastructure that processes real money for real banks. Intuit does this at a scale that few companies in the world match -- and I want to be part of that.

---

## Suggested Diagrams

### Diagram 1 -- System Architecture Diagram

**What to include:**

```
+------------------+
|    CLI Layer     |  cmd/  (Cobra commands)
|  ship init/add/  |
|  commit/status   |
+--------+---------+
         |
         v
+------------------+
|  Command Layer   |  commands/  (Pure Go functions)
|  Init, Add,      |
|  Commit, Status  |
+--------+---------+
         |
    +----+----+
    |         |
    v         v
+--------+ +----------+
| Domain | | Infra    |
| Layer  | | Layer    |
|--------| |----------|
| Blob   | | Hash     |
| Tree   | | Object   |
| Commit | | Store    |
| Index  | | ShipCore |
| Head   | | Worker   |
+--------+ | Pool     |
           +----------+
                |
                v
        +---------------+
        |  Filesystem   |
        |  .ship/       |
        |  objects/     |
        |  refs/        |
        |  index        |
        |  HEAD         |
        +---------------+
```

**Key callouts on the diagram:**
- Unidirectional dependency arrows (no circular imports)
- Label each layer with its responsibility
- Highlight that Domain and Infra are siblings, not parent-child
- Show `.ship/` directory structure as the "database"

---

### Diagram 2 -- Data Flow Diagram: `ship commit` Pipeline

**What to include:**

```
Working Directory
       |
       v
[1] Read .ship/index (JSON)
       |
       v
[2] Build TempDirNode virtual tree
    (split paths on "/", nest into directory structure)
       |
       v
[3] Post-order recursive traversal
    For each directory (leaves first):
       |
       +---> [3a] Collect file entries (blob hashes from index)
       +---> [3b] Collect subdirectory entries (tree hashes from children)
       +---> [3c] Sort all entries alphabetically by name
       +---> [3d] Serialize: "<mode> <name>\0<hash>" for each entry
       +---> [3e] HashObject(serialized, TREE, write=true)
       |           -> Prepend "tree <size>\0"
       |           -> SHA-1 hash
       |           -> zlib compress
       |           -> Store at .ship/objects/<hash[0:2]>/<hash[2:]>
       |
       v
[4] Root tree hash = final tree object hash
       |
       v
[5] Resolve HEAD -> get parent commit hash (if any)
       |
       v
[6] Build commit object:
    "tree <root-tree-hash>\n"
    "parent <parent-hash>\n"     (if not first commit)
    "author <name> <email> <ts> <tz>\n"
    "committer <name> <email> <ts> <tz>\n"
    "\n"
    "<commit message>"
       |
       v
[7] HashObject(commit-content, COMMIT, write=true)
       |
       v
[8] Update .ship/refs/heads/main with new commit hash
       |
       v
[9] Commit complete. Object store now contains:
    - All blob objects (from prior `add`)
    - All tree objects (directories)
    - New commit object
    - Updated ref pointer
```

**Key callouts on the diagram:**
- Show the Merkle tree property: changing one file changes every tree hash up to root
- Highlight the sorting step (determinism guarantee)
- Show that blobs were already stored during `ship add` -- commit only creates trees + commit object
- Label the zlib compression + SHA-1 hashing at the HashObject step
