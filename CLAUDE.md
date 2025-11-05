# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

git-sync-remotes is a bash script that syncs commits between two git remotes automatically. It's designed for maintaining mirrors, syncing between internal/external repositories, or keeping multiple remotes in sync.

## Core Architecture

### Main Script: `git-sync-remotes`

Single bash script (git-sync-remotes:1-329) that performs the entire sync workflow:

1. **Argument Parsing** (git-sync-remotes:74-182)
   - Handles `-y`/`--yes` flag for auto-confirmation
   - Auto-detects remotes (if exactly 2 exist) or accepts explicit remote arguments
   - Auto-detects current branch or accepts branch argument
   - Supports these usage patterns:
     - `git-sync-remotes` - auto-detect everything
     - `git-sync-remotes -y` - auto-detect with auto-confirm
     - `git-sync-remotes branch` - specify branch, auto-detect remotes
     - `git-sync-remotes remote1 remote2` - specify remotes, use current branch
     - `git-sync-remotes branch remote1 remote2` - specify everything

2. **Fetch Phase** (git-sync-remotes:193-205)
   - Fetches from both remotes to get latest state
   - Validates remotes are accessible

3. **Branch Validation** (git-sync-remotes:207-245)
   - Checks if branch exists on both remotes
   - Offers to create missing branches with user confirmation
   - Uses `git ls-remote` for remote branch detection

4. **Commit Comparison** (git-sync-remotes:247-290)
   - Uses `git rev-parse` to get commit SHAs
   - Uses `git rev-list --count` to count commits ahead/behind
   - Displays commit logs and dates to help user understand differences

5. **Sync Logic** (git-sync-remotes:292-325)
   - **In Sync**: Both remotes at same commit - exit successfully
   - **One Ahead**: Push commits from ahead remote to behind remote
   - **Diverged**: Both have unique commits - requires manual intervention
   - Uses `git push remote ref:refs/heads/branch` for syncing

### Installation Scripts

- **install.sh** - Creates symlink in `~/.local/bin` with path verification
- **uninstall.sh** - Removes installation and offers to clean up empty directories

Both scripts share common color output functions (print_error, print_success, print_warning, print_info) and confirm() function for user prompts.

## Development Commands

### Testing the Script

```bash
# Test with auto-detection
./git-sync-remotes

# Test with auto-confirm
./git-sync-remotes -y

# Test with specific branch
./git-sync-remotes main

# Test with specific remotes
./git-sync-remotes main origin gitlab
```

### Installation/Uninstallation

```bash
# Install locally
./install.sh

# Uninstall
./uninstall.sh
```

## Key Technical Details

### Remote Detection Logic
- If 0 remotes: Error
- If 1 remote: Error (need at least 2)
- If 2 remotes: Auto-use both
- If 3+ remotes: Error, require explicit specification

### Sync Direction Algorithm
The script determines sync direction by counting commits:
- `REMOTE1_AHEAD=$(git rev-list --count $REMOTE2_REF..$REMOTE1_REF)`
- `REMOTE2_AHEAD=$(git rev-list --count $REMOTE1_REF..$REMOTE2_REF)`

If both counts > 0, the branches have diverged and require manual merge.

### Important Script Behavior
- Uses `set -e` to exit on errors
- Changes to current working directory on start to support running from anywhere
- Always confirms before pushing (unless `-y` flag is used)
- All git operations use remote refs directly (e.g., `origin/main`) rather than checking out branches locally

## Requirements

- Bash 4.0+
- Git 2.0+
- Unix-like environment (Linux, macOS, WSL)
- SSH keys configured for all remotes being synced
