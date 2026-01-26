#!/bin/bash
# Quick deployment script - Run this to push changes and create a release

set -e

echo "🚢 Ship Deployment Helper"
echo ""

# Check if there are uncommitted changes
if [[ -n $(git status -s) ]]; then
    echo "📝 Uncommitted changes detected."
    echo ""
    git status -s
    echo ""
    read -p "Commit these changes? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        read -p "Enter commit message: " commit_msg
        git add .
        git commit -m "$commit_msg"
        echo "✅ Changes committed"
    else
        echo "❌ Please commit your changes first"
        exit 1
    fi
fi

# Push to remote
echo ""
echo "📤 Pushing to remote..."
git push origin master || git push origin main
echo "✅ Pushed successfully"

# Ask about creating a release
echo ""
read -p "Create a new release? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    # Get current version
    CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
    echo ""
    echo "Current version: $CURRENT_VERSION"
    read -p "Enter new version (e.g., v1.0.0): " new_version

    if [[ -z "$new_version" ]]; then
        echo "❌ Version cannot be empty"
        exit 1
    fi

    read -p "Enter release message: " release_msg

    echo ""
    echo "Creating release $new_version..."
    git tag -a "$new_version" -m "$release_msg"
    git push origin "$new_version"

    echo ""
    echo "✅ Release $new_version created!"
    echo ""
    echo "🔄 GitHub Actions is now building binaries..."
    echo "Monitor progress at: https://github.com/KambojRajan/ship/actions"
    echo ""
    echo "When complete, users can install with:"
    echo "  curl -sSL https://raw.githubusercontent.com/KambojRajan/ship/master/remote-install.sh | bash"
else
    echo "Skipping release creation"
fi

echo ""
echo "✅ Done!"
