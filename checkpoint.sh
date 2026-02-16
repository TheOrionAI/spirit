#!/bin/bash
# Orion Checkpoint Script
# Commits current state and pushes to GitHub

set -e

cd "$(dirname "$0")"

# Check if there are changes
if git diff --quiet HEAD && git diff --cached --quiet; then
    echo "No changes to checkpoint"
    exit 0
fi

# Pull latest first to avoid conflicts
git pull origin main 2>/dev/null || true

# Add all changes
git add -A

# Create checkpoint with timestamp
CHECKPOINT_TIME=$(date -u +"%Y-%m-%d %H:%M UTC")
CHECKPOINT_MSG="Checkpoint: $CHECKPOINT_TIME

Changes:
$(git status --short)

Session: $(whoami)@$(hostname)
"

git commit -m "$CHECKPOINT_MSG"

# Push to GitHub
git push origin main

echo "✅ Checkpoint created: $CHECKPOINT_TIME"
echo "🔗 Repo: https://github.com/TheOrionAI/orion-state"
