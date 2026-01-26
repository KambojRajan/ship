# Ship 🚢

A lightweight, Git-inspired version control system written in Go. Ship implements core version control functionalities including repository initialization, file staging, and object storage using content-addressable storage.

## 📖 Overview

Ship is a minimal implementation of a distributed version control system that mimics Git's core architecture. It provides fundamental version control operations and stores objects using SHA-1 hashing with zlib compression, similar to how Git manages its object database.

## ✨ Features

- **Repository Initialization**: Create a new Ship repository with proper directory structure
- **File Staging**: AddIndex files to the staging area (index) for version tracking
- **Object Storage**: Content-addressable storage system using SHA-1 hashing
- **Object Inspection**: View the contents of stored objects by their hash
- **Blob Management**: Store file contents as blob objects
- **Tree Structure**: Organize files in a hierarchical tree structure
- **Commit Support**: Create commits with tree snapshots and metadata

### Key Components

- **Blob**: Represents file content
- **Tree**: Represents directory structure
- **Commit**: Captures a snapshot with metadata (author, message, timestamp)
- **Index**: Staging area that tracks files to be committed
- **Objects**: Stored in `.ship/objects/` with content-based addressing

## 🚀 Getting Started

### Prerequisites

- Go 1.24.10 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/KambojRajan/ship.git
cd ship
```

2. Build the project:
```bash
go build -o ship
```

3. (Optional) AddIndex to your PATH for system-wide access:
```bash
# For macOS/Linux
sudo mv ship /usr/local/bin/
```

## 📝 Usage

### Initialize a Repository

Initialize a new Ship repository in the current directory or a specified path:

```bash
ship init .
# or
ship init /path/to/project
```

This creates a `.ship` directory with the following structure:
```
.ship/
├── objects/           # Object database
├── refs/
│   ├── heads/        # Branch references
│   └── tags/         # Tag references
├── HEAD              # Current branch pointer
└── index             # Staging area
```

### AddIndex Files

AddIndex files to the staging area:

```bash
ship add .
# or
ship add /path/to/directory
```

The `add` command:
- Recursively walks through the directory
- Creates blob objects for each file
- Updates the index with file paths, hashes, and modes
- Automatically removes deleted files from the index

### Inspect Objects

View the contents of an object using its SHA-1 hash:

```bash
ship cat-file <hash>
```

Example:
```bash
ship cat-file a94a8fe5ccb19ba61c4c0873d391e987982fbbd3
```

This command supports:
- **blob** objects: Displays file content
- **commit** objects: Shows commit information
- **tree** objects: Currently returns an error (not yet implemented for display)

## 🔧 How It Works

### Object Storage

Ship uses content-addressable storage:

1. **Hashing**: File content is hashed using SHA-1
2. **Compression**: Content is compressed with zlib
3. **Storage**: Objects are stored in `.ship/objects/XX/YYYYYYYY...` where XX is the first 2 characters of the hash
4. **Format**: Objects are stored with a header: `<type> <size>\0<content>`

### Index (Staging Area)

The index is stored as JSON in `.ship/index` and tracks:
- File paths (relative to repository root)
- Object hashes (SHA-1)
- File modes (permissions)

### Tree Building

When creating a commit, Ship builds a tree structure:
1. Reads the index to get staged files
2. Constructs a hierarchical directory tree
3. Recursively creates tree objects from leaves to root
4. Each tree object contains entries for files and subdirectories

## 🧪 Testing

Run the test suite:

```bash
go test ./tests/command_tests/... -v
```

The test suite includes:
- **Command tests**: Integration tests for init, add, and commit commands
- **Helper utilities**: Setup helpers for test environments
- **Assertion helpers**: Custom assertions for test validation

## 🛠️ Development

### Project Structure

- **commands/**: High-level command implementations that users interact with
- **Core/common/**: Shared data structures and types
- **Core/Entities/**: Domain models representing version control objects
- **Core/utils/**: Utility functions for hashing, object manipulation, and file operations

### Key Utilities

- `HashBytes()`: Computes SHA-1 hash of byte data
- `HashObject()`: Hashes and optionally stores an object
- `StoreObject()`: Saves compressed objects to disk
- `ObjectExists()`: Checks if an object exists in the database
- `ShipHasBeenInit()`: Validates if a directory is a Ship repository

## Contributing

Contributions are welcome! Please feel free to submit issues, fork the repository, and create pull requests.

1. Fork the project
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'AddIndex some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request


## Author

**Rajan Kamboj**
- GitHub: [@KambojRajan](https://github.com/KambojRajan)

## 📚 Resources

To learn more about how version control systems work:

- [Git Internals - Git Objects](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects)
- [Building Git by James Coglan](https://shop.jcoglan.com/building-git/)
- [Write yourself a Git!](https://wyag.thb.lt/)

---

**Note**: Ship is an educational project and not intended to replace Git for production use. It's designed to help understand the fundamental concepts behind distributed version control systems.
