# Project Context

## Purpose

git-sync-remotes is a bash script that syncs commits between two git remotes automatically. It's designed for:
- Maintaining repository mirrors across different hosting platforms (GitHub, GitLab, Bitbucket, etc.)
- Syncing between internal and external repositories
- Keeping multiple remotes in sync for backup/redundancy
- Automated CI/CD workflows that require multi-remote synchronization

**Key goals:**
- Simple, portable, single-script solution
- Intelligent auto-detection of remotes and branches
- Safe operation with user confirmations and dry-run mode
- Clear, color-coded output for easy understanding
- Handle common edge cases (missing branches, divergence)

## Tech Stack

- **Bash 4.0+** - Shell scripting language for portability across Unix-like systems
- **Git 2.0+** - Version control system, using both porcelain and plumbing commands
- **Unix-like environment** - Linux, macOS, WSL, or any POSIX-compliant system
- **ANSI color codes** - For terminal output formatting

## Project Conventions

### Code Style

- **Error handling**: Use `set -e` to exit on errors
- **Variable naming**:
  - UPPERCASE for global configuration variables (e.g., `REMOTE1`, `BRANCH`, `AUTO_CONFIRM`)
  - lowercase for local function variables
- **Color output**: Standardized helper functions for consistent messaging
  - `print_error()` - Red text for errors
  - `print_success()` - Green text with checkmark for success
  - `print_warning()` - Yellow text with warning symbol
  - `print_info()` - Blue text with info symbol
- **User interaction**: `confirm()` function for yes/no prompts with auto-confirm support
- **Comments**: Include comprehensive usage documentation in script header
- **Formatting**:
  - Clear section separators (commented "Step 1:", "Step 2:", etc.)
  - Consistent indentation (4 spaces)
  - Line breaks between logical sections

### Architecture Patterns

**Single-script modular architecture** with five distinct phases:

1. **Argument Parsing** (git-sync-remotes:77-189)
   - Parse command-line flags (`-y`, `--dry-run`)
   - Auto-detect or validate remotes and branch
   - Support multiple usage patterns

2. **Fetch Phase** (git-sync-remotes:207-219)
   - Fetch from both remotes
   - Validate remote accessibility

3. **Branch Validation** (git-sync-remotes:221-267)
   - Check branch existence on both remotes
   - Offer to create missing branches

4. **Commit Comparison** (git-sync-remotes:269-312)
   - Use git plumbing commands to compare commit SHAs
   - Count commits ahead/behind
   - Display commit logs for user review

5. **Sync Logic** (git-sync-remotes:314-358)
   - Determine sync direction or detect divergence
   - Execute push or show dry-run preview
   - Confirm before making changes (unless auto-confirmed)

**Key architectural decisions:**
- Work directly with remote refs (e.g., `origin/main`) - never checkout branches locally
- Use git plumbing commands (`rev-parse`, `rev-list`, `ls-remote`) for reliability
- Fail-fast with clear error messages
- Always operate in current working directory (supports running from anywhere via alias)

### Testing Strategy

**Manual testing approach** with documented test scenarios:

```bash
# Test auto-detection
./git-sync-remotes

# Test with flags
./git-sync-remotes -y
./git-sync-remotes --dry-run

# Test with specific parameters
./git-sync-remotes main
./git-sync-remotes main origin gitlab
```

**Edge cases to test:**
- Both remotes in sync
- One remote ahead
- Both remotes diverged
- Branch missing on one remote
- Branch missing on both remotes
- Only one remote configured
- More than two remotes configured
- Invalid remote names
- Not in a git repository

**Testing installation:**
```bash
./install.sh    # Test installation
./uninstall.sh  # Test cleanup
```

### Git Workflow

- **Main branch**: `main` (single primary branch)
- **Commit messages**: Clear, descriptive commits following conventional style
  - `feat:` for new features
  - `fix:` for bug fixes
  - `docs:` for documentation changes
- **No local branching**: Simple linear history on main
- **Distribution**: Users install via symlink, so changes are immediately reflected

## Domain Context

### Git Remote Synchronization Concepts

- **Remote refs**: References like `origin/main` that track remote branch state
- **Commit comparison**: Using `git rev-list --count` to determine ahead/behind status
- **Divergence**: State where both remotes have unique commits not in the other
  - Requires manual merge/rebase - cannot be auto-resolved
- **Fast-forward push**: Pushing when one remote is strictly ahead of the other
- **Branch creation**: Using `git push remote ref:refs/heads/branch` syntax

### Sync Direction Algorithm

```bash
REMOTE1_AHEAD=$(git rev-list --count $REMOTE2_REF..$REMOTE1_REF)
REMOTE2_AHEAD=$(git rev-list --count $REMOTE1_REF..$REMOTE2_REF)
```

- If both counts = 0: Remotes are in sync
- If only one count > 0: That remote is ahead, sync from it
- If both counts > 0: Branches have diverged, requires manual intervention

## Important Constraints

### Technical Constraints

- **Remote requirement**: Must have at least 2 remotes configured
- **Auto-detection limit**: Only works automatically with exactly 2 remotes
  - With 0-1 remotes: Error
  - With 3+ remotes: Must explicitly specify which two to use
- **Divergence handling**: Cannot automatically resolve when both remotes have unique commits
- **Environment**: Requires Unix-like environment with bash 4.0+
- **Git version**: Requires Git 2.0+ for modern plumbing commands
- **Authentication**: Requires pre-configured SSH keys or credentials for all remotes

### Operational Constraints

- **Must run from git repository**: Script validates it's in a git repo before proceeding
- **Network dependency**: Requires network access to fetch from remotes
- **Permissions**: User must have push access to remotes for sync to work
- **Line endings**: Script must maintain LF line endings (not CRLF) for Unix compatibility

### Design Philosophy

- **Safety first**: Always confirm before destructive operations (unless `-y` flag)
- **Transparency**: Show users what will happen before it happens
- **Portability**: Single bash script, no external dependencies beyond git
- **User-friendly**: Clear error messages, helpful suggestions for resolution

## External Dependencies

### Required

- **Git CLI** (v2.0+)
  - Commands used: `fetch`, `push`, `rev-parse`, `rev-list`, `ls-remote`, `remote`, `log`
  - Both porcelain (user-facing) and plumbing (scripting) commands
- **Bash** (v4.0+)
  - Array support for remote detection
  - Parameter expansion
  - Process substitution for reading command output
- **SSH** (for remote access)
  - Assumes SSH key-based authentication is configured
  - Users must set up keys for each remote before using the tool

### System Utilities

- Standard Unix tools (implicitly available on target systems):
  - `wc` - Count lines for remote branch existence check
  - `tr` - Trim whitespace
  - `grep` - Pattern matching for remote validation
  - `read` - User input for confirmations

### Installation Dependencies

- **install.sh** requires:
  - `mkdir` - Create installation directory
  - `ln` - Create symbolic link
  - `chmod` - Make script executable
  - `readlink` - Verify symlink creation

### Remote Services

- **Git hosting platforms** (external):
  - GitHub, GitLab, Bitbucket, or any Git-compatible hosting service
  - Must support SSH or HTTPS authentication
  - Must support standard Git protocol operations
