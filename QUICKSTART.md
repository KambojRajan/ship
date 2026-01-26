# 🚀 Quick Start Guide

Get started with Ship in under 2 minutes!

## Installation

### One-Command Install

```bash
curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/main/install.sh | bash
```

Or clone and install manually:

```bash
git clone https://github.com/KambojRajan/ship.git
cd ship
./install.sh
```

## First Steps

### 1. Initialize a Repository

```bash
# Create a new directory for your project
mkdir my-project
cd my-project

# Initialize Ship
ship init .
```

### 2. Create Some Files

```bash
echo "Hello, Ship!" > hello.txt
echo "# My Project" > README.md
mkdir src
echo "package main" > src/main.go
```

### 3. Stage Your Files

```bash
# Stage all files
ship add .

# Or stage specific files
ship add hello.txt
```

### 4. Check What's Staged

```bash
# View the index (staging area)
cat .ship/index
```

### 5. Inspect Objects

```bash
# Get the hash from the index and view the object
ship cat-file <hash>
```

## Common Commands

| Command | Description | Example |
|---------|-------------|---------|
| `ship init <path>` | Initialize a new repository | `ship init .` |
| `ship add <path>` | Stage files for commit | `ship add .` |
| `ship cat-file <hash>` | View object content | `ship cat-file abc123...` |
| `ship --help` | Show help | `ship --help` |

## Directory Structure

After initialization, Ship creates:

```
.ship/
├── objects/          # All versioned content
│   └── XX/           # First 2 chars of hash
│       └── XXXXXX... # Rest of hash
├── refs/
│   ├── heads/       # Branch references
│   └── tags/        # Tag references
├── HEAD             # Current branch
└── index            # Staging area (JSON)
```

## Tips

- **Staging**: Files must be staged before committing
- **Objects**: All content is stored as SHA-1 hashed objects
- **Index**: The staging area tracks what will go into the next commit
- **Paths**: Use absolute or relative paths with commands

## Next Steps

1. Explore the [full documentation](README.md)
2. Learn about [object storage](#how-it-works) 
3. Contribute to the [project](CONTRIBUTING.md)

## Need Help?

```bash
# General help
ship --help

# Command-specific help
ship init --help
ship add --help
ship cat-file --help
```

## Uninstall

```bash
make uninstall
# or
sudo rm /usr/local/bin/ship
```

Happy Shipping! 🚢
