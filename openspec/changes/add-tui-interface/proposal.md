# Change: Add Interactive TUI Interface with Bubbletea

## Why

The current CLI tool requires users to understand command-line arguments and provides limited visibility into the state of multiple remotes. A Terminal User Interface (TUI) similar to lazygit would provide:
- Real-time visual comparison of commits across remotes
- Interactive navigation through commit history and diffs
- Intuitive branch switching and management
- Streamlined sync operations with visual confirmation

This enhances the user experience while keeping the existing CLI functional for automation and simple usage scenarios.

## What Changes

- Add Go-based TUI implementation using bubbletea framework
- Implement multi-pane interface showing:
  - Remote status overview (commit counts, sync state)
  - Commit list for each remote with visual indicators
  - Diff viewer for selected commits
  - Branch selector/manager
- Support interactive sync operations with visual feedback
- Enable branch creation and switching from TUI
- Maintain backward compatibility with existing CLI script
- Add new `git-sync-remotes --tui` flag or make TUI the default with `--cli` fallback
- Provide single-binary distribution via Go build

## Impact

- **Affected specs**:
  - New: `tui-core` (core TUI framework and layout)
  - New: `interactive-sync` (sync operations from TUI)
  - New: `branch-management` (branch operations from TUI)
- **Affected code**:
  - New: `cmd/tui/` - Go TUI implementation
  - New: `internal/git/` - Git operations library (extracted from bash script)
  - New: `internal/ui/` - UI components and models
  - Modified: `git-sync-remotes` - Add TUI launcher and mode detection
  - New: `Makefile` or `build.sh` - Go build scripts
  - New: `go.mod`, `go.sum` - Go dependencies
- **Distribution changes**:
  - Users can choose between bash CLI (portable, no build) or TUI (requires Go binary)
  - Installation will need to build Go binary or download pre-built releases
