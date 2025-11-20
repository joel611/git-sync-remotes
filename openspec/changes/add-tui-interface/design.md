# Design: Interactive TUI Interface

## Context

git-sync-remotes currently provides a bash CLI for syncing commits between two git remotes. While functional for automation, it lacks:
- Visual feedback during operation
- Easy exploration of commit differences
- Interactive branch management
- Intuitive navigation of remote state

The TUI will address these gaps while preserving the CLI for scripts and simple workflows.

**Stakeholders**: CLI users wanting better visibility, users managing multiple repositories regularly

**Constraints**:
- Must maintain bash CLI for backward compatibility and automation
- Must work with existing git infrastructure (no server-side changes)
- Should be installable without complex dependencies
- Must handle 2 remotes (matching current scope)

## Goals / Non-Goals

**Goals**:
- Provide lazygit-style TUI for visualizing and syncing git remotes
- Support interactive sync operations with visual confirmation
- Enable branch switching and creation from TUI
- Maintain CLI functionality for automation/scripting
- Single-binary distribution for easy installation

**Non-Goals**:
- Supporting more than 2 remotes (future enhancement)
- Local commit operations (staging, committing) - this is a remote sync tool
- Replacing bash CLI entirely (keep both modes)
- Web-based or GUI interface
- Editing files or resolving conflicts within TUI

## Decisions

### 1. Framework: Go + Bubbletea

**Decision**: Use Go with the bubbletea TUI framework (same stack as lazygit)

**Rationale**:
- Bubbletea is battle-tested (powers lazygit, glow, soft-serve)
- Go compiles to single static binary (easy distribution)
- Excellent cross-platform support (Linux, macOS, Windows/WSL)
- Strong git integration libraries available (go-git)
- Active community and ecosystem

**Alternatives considered**:
- **Bash + whiptail**: Too limited for rich UI, poor navigation experience
- **Python + textual**: Requires Python runtime, harder distribution, slower startup
- **Rust + ratatui**: Excellent performance but slower compile times, smaller ecosystem for git ops

### 2. Architecture: Dual-Mode Binary

**Decision**: Create a Go binary `git-sync-remotes-tui` that includes both CLI and TUI modes

**Structure**:
```
git-sync-remotes           # Original bash script (kept for portability)
git-sync-remotes-tui       # Go binary with TUI + CLI reimplementation
```

**Mode detection**:
- `git-sync-remotes-tui` (no args) → Launch TUI
- `git-sync-remotes-tui --cli <args>` → Use CLI mode (Go reimplementation)
- `git-sync-remotes <args>` → Original bash CLI (unchanged)

**Rationale**:
- Preserves bash script for users without Go binary
- Allows gradual migration
- Bash script can delegate to TUI if binary exists
- No breaking changes to existing workflows

**Alternatives considered**:
- **Replace bash entirely**: Too disruptive, loses portability
- **TUI-only binary**: Requires keeping bash and Go in sync
- **Make TUI the default**: Confusing for automation users

### 3. UI Layout: Three-Pane Design

**Decision**: Use a three-pane layout similar to lazygit

```
┌─────────────────────────────────────────────────────────────┐
│ Status │ Branch: main  │ Remote1: origin (2 ahead)          │
│        │               │ Remote2: gitlab (0 ahead)          │
├─────────────────────────────────────────────────────────────┤
│ COMMITS (Remote1)      │ COMMITS (Remote2)                  │
│ > abc123 feat: add TUI │   def456 fix: typo                 │
│   def456 fix: typo     │   ghi789 docs: update README       │
│   ghi789 docs: README  │                                    │
├─────────────────────────────────────────────────────────────┤
│ DIFF PREVIEW                                                │
│ commit abc123                                               │
│ Author: Joel Chan <...>                                     │
│ Date: 2025-01-20                                            │
│                                                             │
│ +++ new file content                                        │
└─────────────────────────────────────────────────────────────┘
[s]ync [b]ranch [f]etch [q]uit [?]help
```

**Panes**:
1. **Header** - Status bar with branch, remote info, sync state
2. **Commit Lists** - Side-by-side commit lists for both remotes
3. **Detail View** - Commit details/diff for selected commit
4. **Footer** - Keyboard shortcuts

**Rationale**:
- Familiar to lazygit users
- Clear visual comparison of both remotes
- Efficient use of terminal space
- Supports keyboard navigation (vim-like)

### 4. Git Operations: go-git Library

**Decision**: Use go-git library for git operations

**Rationale**:
- Pure Go implementation (no git binary dependency)
- Faster than shelling out to git CLI
- Type-safe git object manipulation
- Well-maintained (used by GitHub, GitLab tooling)

**Fallback**: Shell out to git CLI for operations not supported by go-git (rare)

**Alternatives considered**:
- **git2go (libgit2 bindings)**: Requires CGO, harder to build/distribute
- **Shell out only**: Slower, harder to parse output, no type safety

### 5. State Management: Elm Architecture (TEA)

**Decision**: Use Bubbletea's Elm architecture for state management

**Model**:
```go
type Model struct {
    repo         *git.Repository
    remote1      RemoteState
    remote2      RemoteState
    currentBranch string
    focusedPane  Pane
    selectedCommit *Commit
    // ... UI state
}
```

**Updates**: Commands trigger state updates (fetch, sync, branch switch)

**Rationale**:
- Natural fit for bubbletea
- Predictable state transitions
- Easy to test and reason about
- Handles async operations (git fetch) cleanly

## Risks / Trade-offs

### Risk: Binary Size and Distribution
- **Risk**: Go binary will be larger than bash script (~5-10MB vs <10KB)
- **Mitigation**: Provide pre-built binaries for common platforms, keep bash script as lightweight option
- **Trade-off**: Accepted for better UX

### Risk: Build Complexity
- **Risk**: Users need Go toolchain to build from source
- **Mitigation**: Provide releases on GitHub with pre-built binaries for Linux/macOS/Windows
- **Trade-off**: Accepted - most users will download binary

### Risk: Feature Parity
- **Risk**: Go implementation might diverge from bash script behavior
- **Mitigation**:
  - Extract bash script logic into shared test cases
  - Test both implementations against same scenarios
  - Document differences clearly
- **Trade-off**: Minor behavioral differences acceptable if documented

### Risk: Dependency on go-git
- **Risk**: go-git may not support all git operations
- **Mitigation**: Fall back to git CLI for unsupported operations
- **Trade-off**: Slightly slower for edge cases, but acceptable

## Migration Plan

### Phase 1: Initial TUI (MVP)
1. Implement Go binary with basic TUI
2. Support: fetch, view commits, view diffs, sync
3. Keep bash script as primary, TUI as experimental
4. Release as `git-sync-remotes-tui` (separate binary)

### Phase 2: Feature Parity
1. Add branch management (create, switch)
2. Add all bash script features to TUI
3. Polish UI/UX based on feedback
4. Document migration path

### Phase 3: Integration
1. Update install.sh to optionally install TUI binary
2. Make bash script delegate to TUI when available
3. Promote TUI as recommended interface
4. Keep bash for automation

### Rollback Plan
- If TUI has critical issues, users can continue using bash script
- Remove TUI binary, update install.sh
- Archive this change without merging TUI into main script

## Open Questions

1. **Keyboard shortcuts**: Should we match lazygit's keybindings or create our own?
   - **Proposal**: Start with lazygit-inspired, document differences

2. **Color scheme**: Dark mode only or support themes?
   - **Proposal**: Start with single dark theme, add themes later

3. **Performance**: How to handle repos with 1000s of commits?
   - **Proposal**: Implement pagination (show 50 commits at a time), lazy load diffs

4. **Error handling**: How to display git errors in TUI?
   - **Proposal**: Modal dialog for errors, log to file for debugging

5. **Configuration**: Should TUI have its own config file?
   - **Proposal**: Start with sensible defaults, add config in Phase 2
