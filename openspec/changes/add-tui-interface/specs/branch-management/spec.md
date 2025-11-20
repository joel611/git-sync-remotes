# Branch Management Capability

## ADDED Requirements

### Requirement: Branch Selection
The TUI SHALL allow users to switch between branches to view their sync status across remotes.

#### Scenario: Open branch selector
- **WHEN** user presses `b` key
- **THEN** a branch selection dialog SHALL be displayed
- **AND** dialog SHALL list all branches that exist on at least one remote
- **AND** branches SHALL be grouped by existence (both remotes, remote1 only, remote2 only)

#### Scenario: Switch to different branch
- **WHEN** branch selection dialog is open
- **AND** user selects a branch that exists on both remotes
- **THEN** TUI SHALL switch to that branch
- **AND** commit lists SHALL be refreshed for the selected branch
- **AND** header SHALL update to show new branch name

#### Scenario: Display current branch indicator
- **WHEN** branch selection dialog is displayed
- **THEN** the currently viewed branch SHALL be highlighted
- **AND** a checkmark or indicator SHALL be shown next to it

### Requirement: Branch Existence Indicators
The TUI SHALL clearly indicate which remotes have each branch.

#### Scenario: Show branch remote availability
- **WHEN** branch selection dialog is displayed
- **THEN** each branch SHALL show indicators for:
  - Remote 1 presence (icon or color)
  - Remote 2 presence (icon or color)
- **AND** branches missing from one remote SHALL be visually distinct

#### Scenario: Filter branches by remote
- **WHEN** branch selection dialog is open
- **AND** user presses `1` key
- **THEN** only branches existing on remote 1 SHALL be displayed
- **WHEN** user presses `2` key
- **THEN** only branches existing on remote 2 SHALL be displayed
- **WHEN** user presses `a` key
- **THEN** all branches SHALL be displayed

### Requirement: Create Missing Branch
The TUI SHALL allow users to create a branch on a remote where it is missing.

#### Scenario: Offer to create missing branch
- **WHEN** user switches to a branch that exists on only one remote
- **THEN** an info message SHALL be displayed
- **AND** message SHALL indicate the branch is missing from one remote
- **AND** message SHALL offer option to create it (press `c` to create)

#### Scenario: Create branch on remote
- **WHEN** viewing a branch that exists on remote 1 but not remote 2
- **AND** user presses `c` key
- **THEN** a confirmation dialog SHALL be displayed
- **AND** dialog SHALL show:
  - Branch name to be created
  - Target remote
  - Source commit (from existing remote)
- **WHEN** user confirms
- **THEN** branch SHALL be created on the target remote
- **AND** a success message SHALL be displayed

#### Scenario: Branch creation failure
- **WHEN** branch creation fails (e.g., permissions, network error)
- **THEN** an error dialog SHALL be displayed
- **AND** error message SHALL include the reason for failure
- **AND** branch state SHALL remain unchanged

### Requirement: Branch Information Display
The TUI SHALL display detailed information about the selected branch.

#### Scenario: Show branch metadata
- **WHEN** a branch is selected in the branch dialog
- **AND** user presses `i` key
- **THEN** a branch info overlay SHALL be displayed
- **AND** overlay SHALL show:
  - Branch name
  - Last commit on each remote
  - Commit authors and dates
  - Sync status between remotes

### Requirement: Branch Search
The TUI SHALL provide a search function to quickly find branches by name.

#### Scenario: Search branches by name
- **WHEN** branch selection dialog is open
- **AND** user presses `/` key
- **THEN** a search input SHALL appear
- **WHEN** user types characters
- **THEN** branch list SHALL filter to show only matching branches
- **AND** matching SHALL be case-insensitive
- **AND** partial matches SHALL be included

#### Scenario: Clear branch search
- **WHEN** branch search is active
- **AND** user presses `Esc` key
- **THEN** search input SHALL be cleared
- **AND** full branch list SHALL be displayed again

### Requirement: Branch Creation Validation
The TUI SHALL validate branch creation operations before execution.

#### Scenario: Prevent duplicate branch creation
- **WHEN** user attempts to create a branch on a remote where it already exists
- **THEN** an error message SHALL be displayed
- **AND** no creation operation SHALL be performed

#### Scenario: Validate branch name format
- **WHEN** creating a new branch
- **THEN** branch name SHALL be validated against git naming rules
- **AND** if invalid, an error message SHALL explain the naming requirements
- **AND** creation SHALL be blocked until a valid name is provided

### Requirement: Refresh Branch List
The TUI SHALL allow users to refresh the branch list to detect new branches on remotes.

#### Scenario: Refresh branch list
- **WHEN** branch selection dialog is open
- **AND** user presses `r` key
- **THEN** a fetch operation SHALL be performed
- **AND** branch list SHALL be refreshed with latest remote data
- **AND** a loading indicator SHALL be displayed during fetch
