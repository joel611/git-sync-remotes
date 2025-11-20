# TUI Core Capability

## ADDED Requirements

### Requirement: TUI Application Launch
The system SHALL provide a terminal user interface mode that can be launched independently of the CLI mode.

#### Scenario: Launch TUI with default settings
- **WHEN** user executes `git-sync-remotes-tui` without arguments
- **THEN** the TUI SHALL launch in the current directory
- **AND** verify it is a git repository
- **AND** display the main interface

#### Scenario: Launch TUI with auto-detected remotes
- **WHEN** TUI launches in a repository with exactly 2 remotes
- **THEN** both remotes SHALL be automatically selected
- **AND** the current branch SHALL be detected and displayed

#### Scenario: Launch fails in non-git directory
- **WHEN** TUI is launched outside a git repository
- **THEN** an error message SHALL be displayed
- **AND** the application SHALL exit with a non-zero status

#### Scenario: Launch with 3+ remotes requires selection
- **WHEN** TUI launches in a repository with more than 2 remotes
- **THEN** a remote selection interface SHALL be displayed
- **AND** user SHALL select exactly 2 remotes before proceeding

### Requirement: Three-Pane Layout
The TUI SHALL display a three-pane layout consisting of a header, dual commit views, and a detail pane.

#### Scenario: Display header pane with status
- **WHEN** TUI is active
- **THEN** the header pane SHALL display:
  - Current branch name
  - Remote 1 name and URL
  - Remote 2 name and URL
  - Sync status (in sync, remote1 ahead, remote2 ahead, diverged)
  - Number of commits ahead for each remote

#### Scenario: Display dual commit lists
- **WHEN** TUI is active
- **THEN** two side-by-side commit list panes SHALL be displayed
- **AND** left pane SHALL show commits from remote 1
- **AND** right pane SHALL show commits from remote 2
- **AND** commits SHALL be displayed with short SHA, message, and relative date

#### Scenario: Display detail pane for selected commit
- **WHEN** user selects a commit in either list
- **THEN** the detail pane SHALL display:
  - Full commit SHA
  - Author name and email
  - Commit date
  - Full commit message
  - Diff preview (file changes)

### Requirement: Keyboard Navigation
The TUI SHALL support keyboard-based navigation using vim-like keybindings.

#### Scenario: Navigate between panes
- **WHEN** user presses `Tab` key
- **THEN** focus SHALL move to the next pane in order: commit-list-1 → commit-list-2 → detail → commit-list-1

#### Scenario: Navigate within commit list
- **WHEN** commit list is focused
- **AND** user presses `j` or Down arrow
- **THEN** selection SHALL move to the next commit
- **WHEN** user presses `k` or Up arrow
- **THEN** selection SHALL move to the previous commit

#### Scenario: Scroll detail pane
- **WHEN** detail pane is focused
- **AND** user presses `d` or Page Down
- **THEN** detail view SHALL scroll down one page
- **WHEN** user presses `u` or Page Up
- **THEN** detail view SHALL scroll up one page

#### Scenario: Exit application
- **WHEN** user presses `q` or `Ctrl+C`
- **THEN** a confirmation dialog SHALL be displayed
- **AND** if confirmed, the application SHALL exit gracefully

### Requirement: Visual Indicators
The TUI SHALL provide visual indicators for commit state and sync status.

#### Scenario: Highlight commits unique to one remote
- **WHEN** displaying commit lists
- **THEN** commits that exist only in remote 1 SHALL be highlighted in one color
- **AND** commits that exist only in remote 2 SHALL be highlighted in another color
- **AND** commits that exist in both remotes SHALL be displayed in neutral color

#### Scenario: Display sync status icons
- **WHEN** header shows sync status
- **THEN** "in sync" state SHALL show a checkmark icon
- **AND** "ahead" state SHALL show an arrow pointing to the behind remote
- **AND** "diverged" state SHALL show a warning icon

### Requirement: Responsive Layout
The TUI SHALL adapt to different terminal sizes while maintaining usability.

#### Scenario: Handle small terminal size
- **WHEN** terminal width is less than 80 columns
- **THEN** commit lists SHALL stack vertically instead of side-by-side
- **AND** commit messages SHALL be truncated with ellipsis

#### Scenario: Handle terminal resize
- **WHEN** terminal is resized during TUI operation
- **THEN** layout SHALL automatically adjust to new dimensions
- **AND** selected commit SHALL remain in view

### Requirement: Real-Time Status Updates
The TUI SHALL update status information when git operations complete.

#### Scenario: Refresh after fetch
- **WHEN** a fetch operation completes
- **THEN** commit lists SHALL be refreshed
- **AND** header status SHALL be updated
- **AND** visual indicators SHALL reflect new state

#### Scenario: Show loading state during operations
- **WHEN** a long-running git operation is in progress
- **THEN** a loading spinner SHALL be displayed in the header
- **AND** the operation name SHALL be shown
- **AND** user input SHALL be blocked until operation completes

### Requirement: Help System
The TUI SHALL provide in-application help for keyboard shortcuts and features.

#### Scenario: Display help overlay
- **WHEN** user presses `?`
- **THEN** a help overlay SHALL be displayed over the main interface
- **AND** help SHALL list all available keyboard shortcuts
- **AND** pressing `?` again or `Esc` SHALL close the help overlay

#### Scenario: Show contextual shortcuts in footer
- **WHEN** TUI is active
- **THEN** footer SHALL display the most relevant keyboard shortcuts for current context
- **AND** shortcuts SHALL update based on focused pane
