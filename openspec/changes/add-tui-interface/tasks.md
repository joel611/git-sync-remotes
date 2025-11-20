# Implementation Tasks

## 1. Project Setup
- [ ] 1.1 Initialize Go module (`go.mod`) with project name and dependencies
- [ ] 1.2 Add bubbletea, bubbles, lipgloss dependencies
- [ ] 1.3 Add go-git library for git operations
- [ ] 1.4 Create project directory structure (`cmd/`, `internal/git/`, `internal/ui/`, `internal/models/`)
- [ ] 1.5 Create Makefile with build, install, test, and clean targets
- [ ] 1.6 Add .gitignore for Go binaries and build artifacts

## 2. Git Operations Layer
- [ ] 2.1 Implement git repository initialization and validation (`internal/git/repo.go`)
- [ ] 2.2 Implement remote detection and listing (`internal/git/remote.go`)
- [ ] 2.3 Implement branch listing and validation (`internal/git/branch.go`)
- [ ] 2.4 Implement fetch operations for remotes (`internal/git/fetch.go`)
- [ ] 2.5 Implement commit comparison logic (ahead/behind counting) (`internal/git/compare.go`)
- [ ] 2.6 Implement push operations for sync (`internal/git/sync.go`)
- [ ] 2.7 Implement branch creation on remote (`internal/git/branch.go`)
- [ ] 2.8 Add error handling and git operation logging
- [ ] 2.9 Write unit tests for git operations (using test fixtures)

## 3. Core TUI Framework
- [ ] 3.1 Create main bubbletea Model struct (`internal/models/app.go`)
- [ ] 3.2 Implement Init() function for initial state and commands
- [ ] 3.3 Implement Update() function for handling messages and state transitions
- [ ] 3.4 Implement View() function for rendering main layout
- [ ] 3.5 Create message types for git operations (FetchMsg, SyncMsg, etc.)
- [ ] 3.6 Implement application startup and repository validation

## 4. UI Components - Layout
- [ ] 4.1 Create header component with branch and remote info (`internal/ui/header.go`)
- [ ] 4.2 Create dual commit list component (`internal/ui/commits.go`)
- [ ] 4.3 Create detail pane component for commit preview (`internal/ui/detail.go`)
- [ ] 4.4 Create footer component with keyboard shortcuts (`internal/ui/footer.go`)
- [ ] 4.5 Implement responsive layout logic for terminal resizing
- [ ] 4.6 Create color scheme and styling with lipgloss

## 5. UI Components - Dialogs
- [ ] 5.1 Create confirmation dialog component (`internal/ui/dialog.go`)
- [ ] 5.2 Create error dialog component
- [ ] 5.3 Create help overlay component (`internal/ui/help.go`)
- [ ] 5.4 Create branch selector dialog (`internal/ui/branch_selector.go`)
- [ ] 5.5 Create loading/progress indicator component (`internal/ui/spinner.go`)
- [ ] 5.6 Create sync log viewer overlay (`internal/ui/log.go`)

## 6. Keyboard Navigation
- [ ] 6.1 Implement pane focus management and Tab navigation
- [ ] 6.2 Implement vim-like navigation (j/k/h/l) in commit lists
- [ ] 6.3 Implement scrolling in detail pane (d/u, PageDown/PageUp)
- [ ] 6.4 Implement help overlay toggle (?)
- [ ] 6.5 Implement quit confirmation (q, Ctrl+C)
- [ ] 6.6 Implement search in branch selector (/)
- [ ] 6.7 Add keyboard shortcut reference to help overlay

## 7. Fetch Feature
- [ ] 7.1 Implement fetch command trigger (f key)
- [ ] 7.2 Implement auto-fetch on TUI launch
- [ ] 7.3 Add fetch progress indicator
- [ ] 7.4 Handle fetch errors with user-friendly messages
- [ ] 7.5 Refresh commit lists after successful fetch
- [ ] 7.6 Update header status after fetch

## 8. Interactive Sync Feature
- [ ] 8.1 Implement sync command trigger (s key)
- [ ] 8.2 Create sync confirmation dialog with commit details
- [ ] 8.3 Implement sync direction detection (which remote is ahead)
- [ ] 8.4 Implement push operation with progress feedback
- [ ] 8.5 Handle sync success and update UI state
- [ ] 8.6 Handle sync errors (permissions, network, conflicts)
- [ ] 8.7 Implement sync preview mode (Shift+S)
- [ ] 8.8 Block sync when remotes are in sync or diverged
- [ ] 8.9 Implement sync log tracking (l key to view)

## 9. Branch Management Feature
- [ ] 9.1 Implement branch selector trigger (b key)
- [ ] 9.2 Create branch list with remote existence indicators
- [ ] 9.3 Implement branch switching and UI refresh
- [ ] 9.4 Implement branch creation trigger (c key when branch is missing)
- [ ] 9.5 Create branch creation confirmation dialog
- [ ] 9.6 Handle branch creation success/failure
- [ ] 9.7 Implement branch search (/) in branch selector
- [ ] 9.8 Implement branch info overlay (i key)
- [ ] 9.9 Implement branch list refresh (r key)
- [ ] 9.10 Add branch name validation

## 10. Visual Indicators
- [ ] 10.1 Implement commit highlighting (unique to remote1, remote2, or both)
- [ ] 10.2 Add sync status icons in header (checkmark, arrows, warning)
- [ ] 10.3 Add branch existence indicators in branch selector
- [ ] 10.4 Implement loading spinner for async operations
- [ ] 10.5 Add color coding for sync states (in sync = green, diverged = red)

## 11. Testing
- [ ] 11.1 Write unit tests for git operations layer
- [ ] 11.2 Write unit tests for UI components (rendering)
- [ ] 11.3 Create test repository fixtures for integration testing
- [ ] 11.4 Test TUI with 2 remotes in sync
- [ ] 11.5 Test TUI with one remote ahead
- [ ] 11.6 Test TUI with diverged remotes
- [ ] 11.7 Test branch creation workflow
- [ ] 11.8 Test error handling (network failures, permission errors)
- [ ] 11.9 Test terminal resize handling
- [ ] 11.10 Test with large commit histories (1000+ commits)

## 12. CLI Mode (Go Reimplementation)
- [ ] 12.1 Implement CLI argument parsing with flags package
- [ ] 12.2 Implement CLI mode detection (--cli flag)
- [ ] 12.3 Reimplement bash script logic in Go for CLI mode
- [ ] 12.4 Add color output helpers for CLI mode
- [ ] 12.5 Implement --dry-run flag for CLI
- [ ] 12.6 Implement -y/--yes auto-confirm flag for CLI
- [ ] 12.7 Ensure CLI mode output matches bash script format
- [ ] 12.8 Test CLI mode against bash script test cases

## 13. Build and Distribution
- [ ] 13.1 Create Makefile with cross-platform build targets (Linux, macOS, Windows)
- [ ] 13.2 Test binary on Linux x86_64
- [ ] 13.3 Test binary on macOS (Intel and ARM)
- [ ] 13.4 Test binary on Windows/WSL
- [ ] 13.5 Create install script to download and install TUI binary
- [ ] 13.6 Update main install.sh to optionally install TUI binary
- [ ] 13.7 Set up GitHub Actions for automated builds (optional for MVP)
- [ ] 13.8 Create release artifacts (tar.gz, checksums)

## 14. Documentation
- [ ] 14.1 Update README.md with TUI installation instructions
- [ ] 14.2 Add TUI usage guide with screenshots/GIFs
- [ ] 14.3 Document keyboard shortcuts in README
- [ ] 14.4 Add troubleshooting section for TUI
- [ ] 14.5 Update CLAUDE.md with TUI architecture details
- [ ] 14.6 Create CONTRIBUTING.md with development setup for Go
- [ ] 14.7 Document differences between bash CLI and Go CLI modes
- [ ] 14.8 Add migration guide from bash-only to TUI

## 15. Integration and Polish
- [ ] 15.1 Modify bash script to detect and delegate to TUI binary if available
- [ ] 15.2 Add version command (--version) to show TUI version
- [ ] 15.3 Implement graceful degradation if go-git doesn't support an operation
- [ ] 15.4 Add configuration file support for TUI preferences (optional)
- [ ] 15.5 Implement session state persistence (remember last branch, etc.)
- [ ] 15.6 Performance optimization for large repositories
- [ ] 15.7 Final UI polish (colors, spacing, alignment)
- [ ] 15.8 Conduct user acceptance testing with sample workflows

## Dependencies Between Tasks
- Tasks 1.x must complete before all other tasks (setup)
- Tasks 2.x (git layer) must complete before 7.x, 8.x, 9.x (features)
- Tasks 3.x and 4.x (core framework) must complete before 6.x, 7.x, 8.x, 9.x
- Tasks 5.x (dialogs) needed for 8.x and 9.x
- Task 12.x (CLI mode) can be done in parallel with TUI development
- Tasks 13.x (build) can start after 8.x and 9.x are complete
- Tasks 14.x (docs) should be done near completion
- Tasks 15.x (integration) should be done last

## Parallelizable Work
- Section 2 (git layer) and Section 4 (UI layout) can be developed in parallel
- Section 5 (dialogs) can be developed in parallel with Section 6 (navigation)
- Section 12 (CLI mode) can be developed independently
- Testing (Section 11) can be written incrementally alongside features
