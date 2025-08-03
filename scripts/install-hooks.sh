#!/bin/bash

# Install Git hooks for url-shortener project
# This script copies hooks from scripts/hooks/ to .git/hooks/

set -e

echo "🔧 Installing Git hooks..."

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_status() {
    echo -e "${2}${1}${NC}"
}

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Error: Not in a git repository root"
    exit 1
fi

# Check if hooks directory exists
if [ ! -d "scripts/hooks" ]; then
    echo "❌ Error: scripts/hooks directory not found"
    exit 1
fi

# Install pre-commit hook
if [ -f "scripts/hooks/pre-commit" ]; then
    cp scripts/hooks/pre-commit .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    print_status "✅ Pre-commit hook installed" "$GREEN"
else
    echo "❌ Error: scripts/hooks/pre-commit not found"
    exit 1
fi

print_status "🎉 Git hooks installation complete!" "$GREEN"
print_status "📝 The pre-commit hook will now run automatically before each commit" "$YELLOW"
print_status "🔧 To bypass hooks temporarily, use: git commit --no-verify" "$YELLOW"