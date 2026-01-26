# Ship Deployment Guide

This guide explains how to properly set up and deploy Ship for easy user installation.

## Current Setup

Your project now has **three installation methods**:

### 1. 🚀 One-Line Remote Installation (Recommended for Users)

```bash
curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/master/remote-install.sh | bash
```

**How it works:**
- Downloads the `remote-install.sh` script
- Detects OS and architecture
- Tries to download pre-built binary from GitHub releases
- Falls back to building from source if binary not available
- Installs to `/usr/local/bin/ship`

### 2. 📦 Go Install Method

```bash
go install github.com/KambojRajan/ship@latest
```

**How it works:**
- Requires Go to be installed
- Downloads and builds from source
- Installs to `$GOPATH/bin`

### 3. 🔧 Manual Installation

```bash
git clone https://github.com/KambojRajan/ship.git
cd ship
./install.sh
# OR
make install
```

## How to Deploy

### Step 1: Commit and Push Your Changes

```bash
cd /Users/rajankamboj/Desktop/something/ship
git add .
git commit -m "feat: add remote installation script"
git push origin master
```

### Step 2: Create a Release (Triggers Binary Build)

Create a new release to trigger GoReleaser:

```bash
# Tag the release
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

This will:
- ✅ Trigger the GitHub Actions workflow
- ✅ Run tests
- ✅ Build binaries for:
  - Linux (amd64, arm64, 386)
  - macOS (amd64, arm64)
  - Windows (amd64, 386)
- ✅ Create archives (tar.gz for Unix, zip for Windows)
- ✅ Generate checksums
- ✅ Create GitHub release with all assets

### Step 3: Test the Installation

After the release is created, test the one-line install:

```bash
curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/master/remote-install.sh | bash
```

## Verification Checklist

Before creating a release:

- [ ] All tests pass locally: `go test ./tests/...`
- [ ] README is updated with installation instructions
- [ ] VERSION file is updated
- [ ] All changes are committed

## File Overview

### Installation Scripts

1. **`remote-install.sh`** (NEW)
   - For remote curl installation
   - Handles binary downloads and source builds
   - Should be in GitHub repository

2. **`install.sh`** (EXISTING)
   - For local installation after cloning
   - Builds from source in current directory
   - Used by developers

### Configuration Files

1. **`.goreleaser.yml`**
   - Configures binary builds for multiple platforms
   - Defines release artifacts
   - Already configured ✅

2. **`.github/workflows/release.yml`**
   - GitHub Actions workflow
   - Triggers on tag push
   - Runs GoReleaser
   - Already configured ✅

## Common Issues and Solutions

### Issue: 404 Error on Curl

**Problem:** 
```bash
curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/main/remote-install.sh | bash
# 404: Not Found
```

**Solution:**
Your default branch is `master`, not `main`. Use:
```bash
curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/master/remote-install.sh | bash
```

### Issue: No Pre-built Binaries

**Problem:** Remote install can't find pre-built binaries

**Solution:** 
- Create a release with `git tag v1.0.0 && git push origin v1.0.0`
- Wait for GitHub Actions to complete
- Pre-built binaries will be available at: `https://github.com/KambojRajan/ship/releases`

### Issue: Go Version Mismatch

**Problem:** Build fails due to Go version

**Solution:**
- Update your Go to 1.24.10 or higher
- Or update `go.mod` to match your Go version

## Next Steps

1. **Commit the new files:**
   ```bash
   git add remote-install.sh README.md
   git commit -m "feat: add one-line remote installation support"
   git push origin master
   ```

2. **Create your first release:**
   ```bash
   git tag -a v1.0.0 -m "Release v1.0.0 - First stable release"
   git push origin v1.0.0
   ```

3. **Wait for GitHub Actions** to build binaries (5-10 minutes)

4. **Test the installation:**
   ```bash
   curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/master/remote-install.sh | bash
   ```

5. **Update your documentation** to showcase the one-line install

## Marketing Your Installation

Update your README header with:

```markdown
# Ship 🚢

> A lightweight, Git-inspired version control system written in Go

## Quick Start

```bash
curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/master/remote-install.sh | bash
```

That's it! Ship is now installed. Try:
```bash
ship init .
ship add .
ship commit -m "Initial commit"
```
```

This creates a **great first impression** for users! 🎉
