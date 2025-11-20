package models

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/joel611/git-sync-remotes/internal/git"
	"github.com/joel611/git-sync-remotes/internal/ui"
)

// Pane represents which pane is currently focused
type Pane int

const (
	CommitList1 Pane = iota
	CommitList2
	DetailPane
)

// Model represents the main TUI application state
type Model struct {
	repo          *git.Repository
	remote1       *git.Remote
	remote2       *git.Remote
	currentBranch string

	// State
	compareResult *git.CompareResult
	focusedPane   Pane
	selectedIndex int
	loading       bool
	err           error
	message       string

	// UI components
	spinner spinner.Model

	// Dimensions
	width  int
	height int

	// Flags
	showHelp       bool
	showBranches   bool
	showSyncDialog bool
	showAddRemote  bool
	quitting       bool

	// Add remote form fields
	addRemoteName  string
	addRemoteURL   string
	addRemoteField int // 0 = name, 1 = url

	// Branch selector
	branches           []git.Branch
	selectedBranchIdx  int
	branchSearchQuery  string
	branchSearchActive bool

	// Branch creation
	showBranchCreate bool
	createBranchName string
	createOnRemote   string // Which remote to create the branch on

	// Branch info overlay
	showBranchInfo bool
}

// NewModel creates a new TUI model
func NewModel(repo *git.Repository, remote1, remote2 *git.Remote, branch string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		repo:          repo,
		remote1:       remote1,
		remote2:       remote2,
		currentBranch: branch,
		focusedPane:   CommitList1,
		spinner:       s,
	}
}

// Messages for async operations
type fetchCompleteMsg struct {
	err error
}

type compareCompleteMsg struct {
	result *git.CompareResult
	err    error
}

type syncCompleteMsg struct {
	err error
}

type addRemoteCompleteMsg struct {
	remote *git.Remote
	err    error
}

type branchesLoadedMsg struct {
	branches []git.Branch
	err      error
}

type branchSwitchMsg struct {
	branch string
	err    error
}

type branchCreateMsg struct {
	branchName string
	remoteName string
	err        error
}

// Messages for initialization
type initFetchMsg struct{}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	// Start spinner and auto-fetch in background if we have 2 remotes
	if m.remote2 != nil {
		return tea.Batch(
			m.spinner.Tick,
			fetchRemotes(m.repo, m.remote1.Name, m.remote2.Name),
		)
	}

	// If only 1 remote, just start spinner
	return m.spinner.Tick
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case fetchCompleteMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Fetch failed: %v. Press 'f' to retry.", msg.err)
			return m, nil
		}
		m.message = ""
		// Only compare if we have 2 remotes
		if m.remote2 != nil {
			return m, compareBranch(m.repo, m.remote1.Name, m.remote2.Name, m.currentBranch)
		}
		return m, nil

	case compareCompleteMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Compare failed: %v", msg.err)
			return m, nil
		}
		m.compareResult = msg.result
		m.message = ""
		return m, nil

	case syncCompleteMsg:
		m.loading = false
		m.showSyncDialog = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Sync failed: %v", msg.err)
			return m, nil
		}
		m.message = "Sync successful!"
		// Re-fetch and compare after sync
		return m, tea.Batch(
			fetchRemotes(m.repo, m.remote1.Name, m.remote2.Name),
		)

	case addRemoteCompleteMsg:
		m.loading = false
		m.showAddRemote = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Failed to add remote: %v", msg.err)
			return m, nil
		}
		// Update remote2 with the newly added remote
		m.remote2 = msg.remote
		m.message = fmt.Sprintf("Remote '%s' added successfully! Press 'f' to fetch.", msg.remote.Name)
		// Clear form fields
		m.addRemoteName = ""
		m.addRemoteURL = ""
		return m, nil

	case branchesLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Failed to load branches: %v", msg.err)
			m.showBranches = false
			return m, nil
		}
		m.branches = msg.branches
		m.selectedBranchIdx = 0
		// Find current branch in the list
		for i, branch := range m.branches {
			if branch.Name == m.currentBranch {
				m.selectedBranchIdx = i
				break
			}
		}
		return m, nil

	case branchSwitchMsg:
		m.loading = false
		m.showBranches = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Failed to switch branch: %v", msg.err)
			return m, nil
		}
		m.currentBranch = msg.branch
		m.message = fmt.Sprintf("Switched to branch '%s'. Press 'f' to fetch.", msg.branch)
		// Reset state for new branch
		m.compareResult = nil
		m.selectedIndex = 0
		return m, nil

	case branchCreateMsg:
		m.loading = false
		m.showBranchCreate = false
		if msg.err != nil {
			m.err = msg.err
			m.message = fmt.Sprintf("Failed to create branch: %v", msg.err)
			return m, nil
		}
		m.message = fmt.Sprintf("Branch '%s' created on %s. Press 'b' to refresh branches.", msg.branchName, msg.remoteName)
		// Refresh branches after creation
		return m, loadBranches(m.repo, m.remote1.Name, m.remote2.Name)
	}

	return m, nil
}

// handleKeyPress handles keyboard input
func (m Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle dialog-specific keys first
	if m.showHelp {
		if msg.String() == "?" || msg.String() == "esc" {
			m.showHelp = false
		}
		return m, nil
	}

	if m.showAddRemote {
		switch msg.String() {
		case "esc":
			m.showAddRemote = false
			m.addRemoteName = ""
			m.addRemoteURL = ""
			m.addRemoteField = 0
			return m, nil
		case "tab", "down":
			// Switch between fields
			m.addRemoteField = (m.addRemoteField + 1) % 2
			return m, nil
		case "up":
			m.addRemoteField = (m.addRemoteField + 1) % 2
			return m, nil
		case "enter":
			// Submit form
			if m.addRemoteName != "" && m.addRemoteURL != "" {
				m.loading = true
				return m, addRemote(m.repo, m.addRemoteName, m.addRemoteURL)
			}
			return m, nil
		case "backspace":
			// Delete character from current field
			if m.addRemoteField == 0 && len(m.addRemoteName) > 0 {
				m.addRemoteName = m.addRemoteName[:len(m.addRemoteName)-1]
			} else if m.addRemoteField == 1 && len(m.addRemoteURL) > 0 {
				m.addRemoteURL = m.addRemoteURL[:len(m.addRemoteURL)-1]
			}
			return m, nil
		default:
			// Add character to current field
			if len(msg.String()) == 1 {
				if m.addRemoteField == 0 {
					m.addRemoteName += msg.String()
				} else {
					m.addRemoteURL += msg.String()
				}
			}
			return m, nil
		}
	}

	if m.showBranchCreate {
		switch msg.String() {
		case "y", "Y":
			m.showBranchCreate = false
			m.loading = true
			// Determine source ref - use the remote that HAS the branch
			var sourceRef string
			if m.selectedBranchIdx < len(m.branches) {
				selectedBranch := m.branches[m.selectedBranchIdx]
				if selectedBranch.ExistsRemote1 {
					sourceRef = fmt.Sprintf("%s/%s", m.remote1.Name, selectedBranch.Name)
				} else if selectedBranch.ExistsRemote2 {
					sourceRef = fmt.Sprintf("%s/%s", m.remote2.Name, selectedBranch.Name)
				} else {
					// Branch doesn't exist on either remote - use current branch as source
					sourceRef = fmt.Sprintf("%s/%s", m.remote1.Name, m.currentBranch)
				}
			} else {
				// Fallback to current branch
				sourceRef = fmt.Sprintf("%s/%s", m.remote1.Name, m.currentBranch)
			}
			return m, createBranch(m.repo, m.createOnRemote, m.createBranchName, sourceRef)
		case "n", "N", "esc":
			m.showBranchCreate = false
		}
		return m, nil
	}

	if m.showBranches {
		// Handle branch search mode
		if m.branchSearchActive {
			switch msg.String() {
			case "esc":
				m.branchSearchActive = false
				m.branchSearchQuery = ""
				return m, nil
			case "enter":
				m.branchSearchActive = false
				return m, nil
			case "backspace":
				if len(m.branchSearchQuery) > 0 {
					m.branchSearchQuery = m.branchSearchQuery[:len(m.branchSearchQuery)-1]
				}
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.branchSearchQuery += msg.String()
				}
				return m, nil
			}
		}

		// Handle branch selector navigation
		switch msg.String() {
		case "esc":
			m.showBranches = false
			m.branchSearchQuery = ""
			return m, nil
		case "j", "down":
			if m.selectedBranchIdx < len(m.branches)-1 {
				m.selectedBranchIdx++
			}
			return m, nil
		case "k", "up":
			if m.selectedBranchIdx > 0 {
				m.selectedBranchIdx--
			}
			return m, nil
		case "enter":
			// Switch to selected branch
			if m.selectedBranchIdx < len(m.branches) {
				selectedBranch := m.branches[m.selectedBranchIdx]
				if selectedBranch.Name != m.currentBranch {
					m.loading = true
					return m, switchBranch(m.repo, selectedBranch.Name, m.remote1.Name, m.remote2.Name)
				} else {
					m.showBranches = false
					m.message = "Already on this branch"
				}
			}
			return m, nil
		case "/":
			m.branchSearchActive = true
			return m, nil
		case "r":
			// Refresh branches
			m.loading = true
			return m, loadBranches(m.repo, m.remote1.Name, m.remote2.Name)
		case "c":
			// Create branch on missing remote
			if m.selectedBranchIdx < len(m.branches) {
				selectedBranch := m.branches[m.selectedBranchIdx]
				// Check if branch is missing on one remote
				if !selectedBranch.ExistsRemote1 && selectedBranch.ExistsRemote2 {
					// Branch missing on remote1
					m.createBranchName = selectedBranch.Name
					m.createOnRemote = m.remote1.Name
					m.showBranchCreate = true
				} else if selectedBranch.ExistsRemote1 && !selectedBranch.ExistsRemote2 {
					// Branch missing on remote2
					m.createBranchName = selectedBranch.Name
					m.createOnRemote = m.remote2.Name
					m.showBranchCreate = true
				} else {
					m.message = "Branch exists on both remotes or neither remote"
				}
			}
			return m, nil
		case "i":
			// Show branch info
			if m.selectedBranchIdx < len(m.branches) {
				m.showBranchInfo = true
			}
			return m, nil
		}
		return m, nil
	}

	if m.showBranchInfo {
		// Close info overlay with any key
		m.showBranchInfo = false
		return m, nil
	}

	if m.showSyncDialog {
		switch msg.String() {
		case "y", "Y":
			m.showSyncDialog = false
			m.loading = true
			return m, performSync(m.repo, m.compareResult, m.remote1.Name, m.remote2.Name, m.currentBranch)
		case "n", "N", "esc":
			m.showSyncDialog = false
		}
		return m, nil
	}

	// Global keys
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		return m, tea.Quit

	case "?":
		m.showHelp = true
		return m, nil

	case "a":
		// Only allow adding remote if we don't have 2 remotes yet
		if m.remote2 == nil {
			m.showAddRemote = true
			m.addRemoteField = 0
		} else {
			m.message = "Already have 2 remotes configured"
		}
		return m, nil

	case "b":
		// Open branch selector (only if we have 2 remotes)
		if m.remote2 != nil && !m.loading {
			m.showBranches = true
			m.loading = true
			return m, loadBranches(m.repo, m.remote1.Name, m.remote2.Name)
		} else if m.remote2 == nil {
			m.message = "Need 2 remotes for branch management. Press 'a' to add a second remote."
		}
		return m, nil

	case "f":
		if !m.loading && m.remote2 != nil {
			m.loading = true
			m.message = "Fetching..."
			return m, fetchRemotes(m.repo, m.remote1.Name, m.remote2.Name)
		} else if m.remote2 == nil {
			m.message = "Need 2 remotes to fetch. Press 'a' to add a second remote."
		}

	case "s":
		if m.remote2 == nil {
			m.message = "Need 2 remotes to sync. Press 'a' to add a second remote."
		} else if m.compareResult != nil && !m.loading {
			if m.compareResult.Status == git.InSync {
				m.message = "Remotes are already in sync"
			} else if m.compareResult.Status == git.Diverged {
				m.message = "Remotes have diverged - manual intervention required"
			} else {
				m.showSyncDialog = true
			}
		}

	case "tab":
		m.focusedPane = (m.focusedPane + 1) % 3

	case "j", "down":
		if m.focusedPane == CommitList1 || m.focusedPane == CommitList2 {
			m.selectedIndex++
		}

	case "k", "up":
		if (m.focusedPane == CommitList1 || m.focusedPane == CommitList2) && m.selectedIndex > 0 {
			m.selectedIndex--
		}
	}

	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if m.width == 0 {
		return "Loading..."
	}

	var content string

	// Show help overlay if active
	if m.showHelp {
		return m.renderHelp()
	}

	// Show add remote dialog if active
	if m.showAddRemote {
		return m.renderAddRemoteDialog()
	}

	// Show branch selector if active
	if m.showBranches {
		return m.renderBranchSelector()
	}

	// Show branch creation dialog if active
	if m.showBranchCreate {
		return m.renderBranchCreationDialog()
	}

	// Show branch info overlay if active
	if m.showBranchInfo {
		return m.renderBranchInfo()
	}

	// Show sync dialog if active
	if m.showSyncDialog {
		return m.renderSyncDialog()
	}

	// Render header
	header := m.renderHeader()

	// Render main content area
	mainContent := m.renderMainContent()

	// Render footer
	footer := m.renderFooter()

	content = lipgloss.JoinVertical(lipgloss.Left,
		header,
		mainContent,
		footer,
	)

	return content
}

// renderHeader renders the header pane
func (m Model) renderHeader() string {
	var status string

	// Show appropriate status
	if m.remote2 == nil {
		status = "Only one remote found. Press 'a' to add a second remote."
	} else if m.loading {
		status = fmt.Sprintf("%s Fetching from remotes...", m.spinner.View())
	} else if m.compareResult != nil {
		status = ui.FormatSyncStatus(m.compareResult, m.remote1.Name, m.remote2.Name)
	} else {
		status = "Ready. Press 'f' to fetch from remotes."
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true)

	header := fmt.Sprintf("Branch: %s | %s",
		m.currentBranch,
		status,
	)

	if m.message != "" {
		header += "\n" + m.message
	}

	return headerStyle.Render(header)
}

// renderMainContent renders the main content area
func (m Model) renderMainContent() string {
	// Show two panes even if no compare result yet
	var remote1Commits, remote2Commits []git.Commit

	if m.compareResult != nil {
		remote1Commits = m.compareResult.Remote1Commits
		remote2Commits = m.compareResult.Remote2Commits
	}

	// Render commit lists (will show empty if no commits)
	var left, right string
	if m.remote2 != nil {
		left = m.renderCommitList(m.remote1.Name, remote1Commits, m.focusedPane == CommitList1)
		right = m.renderCommitList(m.remote2.Name, remote2Commits, m.focusedPane == CommitList2)
	} else {
		// Only one remote - show single pane
		left = m.renderCommitList(m.remote1.Name, remote1Commits, true)
		right = m.renderPlaceholder()
	}

	columns := lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		right,
	)

	return columns
}

// renderPlaceholder renders a placeholder pane
func (m Model) renderPlaceholder() string {
	style := lipgloss.NewStyle().
		Width(m.width / 2).
		Padding(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderRight(true)

	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("SECOND REMOTE"))
	lines = append(lines, "")
	lines = append(lines, "  Press 'a' to add a second remote")
	lines = append(lines, "")

	return style.Render(strings.Join(lines, "\n"))
}

// renderCommitList renders a commit list
func (m Model) renderCommitList(remoteName string, commits []git.Commit, focused bool) string {
	style := lipgloss.NewStyle().
		Width(m.width / 2).
		Padding(1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderLeft(true).
		BorderRight(true)

	if focused {
		style = style.BorderForeground(lipgloss.Color("205"))
	}

	// Color for commits unique to this remote
	commitStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("117")) // Light blue

	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("COMMITS (%s)", remoteName)))
	lines = append(lines, "")

	if len(commits) == 0 {
		if m.compareResult == nil {
			lines = append(lines, "  (waiting for fetch...)")
		} else {
			lines = append(lines, "  (no unique commits)")
		}
	} else {
		for i, commit := range commits {
			prefix := "  "
			selector := ""
			if focused && i == m.selectedIndex {
				prefix = ""
				selector = "> "
			}

			// Format commit line with color
			commitLine := fmt.Sprintf("%s %s", commit.ShortSHA, commit.Message)
			if len(commits) > 0 {
				// Commits in this list are unique to this remote, so color them
				commitLine = commitStyle.Render(commitLine)
			}

			lines = append(lines, fmt.Sprintf("%s%s%s", prefix, selector, commitLine))
		}
	}

	return style.Render(strings.Join(lines, "\n"))
}

// renderFooter renders the footer with keyboard shortcuts
func (m Model) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder()).
		BorderTop(true)

	var shortcuts string
	if m.remote2 == nil {
		shortcuts = "[a]dd remote [q]uit [?]help"
	} else {
		shortcuts = "[f]etch [s]ync [q]uit [?]help"
	}
	return footerStyle.Render(shortcuts)
}

// renderHelp renders the help overlay
func (m Model) renderHelp() string {
	helpStyle := lipgloss.NewStyle().
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205"))

	help := `
Keyboard Shortcuts:

  Navigation:
    ↑/k         Move up in list
    ↓/j         Move down in list
    Tab         Switch between panes

  Actions:
    a           Add remote (when only 1 remote exists)
    b           Branch selector (switch/manage branches)
    f           Fetch from remotes
    s           Sync commits

  Branch Management (in branch selector):
    c           Create branch on missing remote
    i           Show branch information

  Other:
    ?           Toggle this help
    q/Ctrl+C    Quit

Press ? or Esc to close this help.
`

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		helpStyle.Render(help))
}

// renderSyncDialog renders the sync confirmation dialog
func (m Model) renderSyncDialog() string {
	dialogStyle := lipgloss.NewStyle().
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205"))

	var direction string
	var count int
	if m.compareResult.Status == git.Remote1Ahead {
		direction = fmt.Sprintf("%s → %s", m.remote1.Name, m.remote2.Name)
		count = m.compareResult.Remote1Ahead
	} else {
		direction = fmt.Sprintf("%s → %s", m.remote2.Name, m.remote1.Name)
		count = m.compareResult.Remote2Ahead
	}

	dialog := fmt.Sprintf(`
Sync Confirmation

Direction: %s
Commits to sync: %d

Continue? [y/n]
`, direction, count)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogStyle.Render(dialog))
}

// renderAddRemoteDialog renders the add remote dialog
func (m Model) renderAddRemoteDialog() string {
	dialogStyle := lipgloss.NewStyle().
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205"))

	// Highlight current field
	nameStyle := lipgloss.NewStyle()
	urlStyle := lipgloss.NewStyle()

	if m.addRemoteField == 0 {
		nameStyle = nameStyle.Foreground(lipgloss.Color("205")).Bold(true)
	} else {
		urlStyle = urlStyle.Foreground(lipgloss.Color("205")).Bold(true)
	}

	nameValue := m.addRemoteName
	if nameValue == "" {
		nameValue = "(enter name)"
	}

	urlValue := m.addRemoteURL
	if urlValue == "" {
		urlValue = "(enter URL)"
	}

	dialog := fmt.Sprintf(`Add Remote

%s Name:  %s
%s URL:   %s

Press Tab to switch fields
Press Enter to submit
Press Esc to cancel
`,
		func() string {
			if m.addRemoteField == 0 {
				return ">"
			}
			return " "
		}(),
		nameStyle.Render(nameValue),
		func() string {
			if m.addRemoteField == 1 {
				return ">"
			}
			return " "
		}(),
		urlStyle.Render(urlValue),
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogStyle.Render(dialog))
}

// renderBranchSelector renders the branch selector dialog
func (m Model) renderBranchSelector() string {
	dialogStyle := lipgloss.NewStyle().
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205"))

	if m.loading {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			dialogStyle.Render("Loading branches..."))
	}

	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render("Select Branch"))
	lines = append(lines, "")

	// Filter branches if search is active
	visibleBranches := m.branches
	if m.branchSearchQuery != "" {
		filtered := []git.Branch{}
		for _, branch := range m.branches {
			if strings.Contains(strings.ToLower(branch.Name), strings.ToLower(m.branchSearchQuery)) {
				filtered = append(filtered, branch)
			}
		}
		visibleBranches = filtered
	}

	if len(visibleBranches) == 0 {
		lines = append(lines, "  No branches found")
	} else {
		// Show branches with indicators
		for i, branch := range visibleBranches {
			prefix := "  "
			if i == m.selectedBranchIdx {
				prefix = "> "
			}

			// Add indicators for remote existence
			indicators := ""
			if branch.ExistsRemote1 && branch.ExistsRemote2 {
				indicators = " [both]"
			} else if branch.ExistsRemote1 {
				indicators = " [remote1]"
			} else if branch.ExistsRemote2 {
				indicators = " [remote2]"
			}

			// Highlight current branch
			branchText := branch.Name + indicators
			if branch.Name == m.currentBranch {
				branchText = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("* " + branchText)
			}

			lines = append(lines, prefix+branchText)
		}
	}

	lines = append(lines, "")
	if m.branchSearchActive {
		lines = append(lines, fmt.Sprintf("Search: %s_", m.branchSearchQuery))
		lines = append(lines, "Press Esc to exit search")
	} else if m.branchSearchQuery != "" {
		lines = append(lines, fmt.Sprintf("Filter: %s (Press / to search again)", m.branchSearchQuery))
	} else {
		lines = append(lines, "↑/↓/j/k: Navigate  Enter: Switch  c: Create  i: Info  /: Search  r: Refresh  Esc: Close")
	}

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogStyle.Render(strings.Join(lines, "\n")))
}

// renderBranchCreationDialog renders the branch creation confirmation dialog
func (m Model) renderBranchCreationDialog() string {
	dialogStyle := lipgloss.NewStyle().
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205"))

	dialog := fmt.Sprintf(`
Create Branch

Branch: %s
Remote: %s

This will create the branch on %s from the existing branch on the other remote.

Continue? [y/n]
`, m.createBranchName, m.createOnRemote, m.createOnRemote)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogStyle.Render(dialog))
}

// renderBranchInfo renders the branch info overlay
func (m Model) renderBranchInfo() string {
	dialogStyle := lipgloss.NewStyle().
		Padding(2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("205"))

	if m.selectedBranchIdx >= len(m.branches) {
		return lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			dialogStyle.Render("No branch selected"))
	}

	selectedBranch := m.branches[m.selectedBranchIdx]

	var infoLines []string
	infoLines = append(infoLines, lipgloss.NewStyle().Bold(true).Render("Branch Information"))
	infoLines = append(infoLines, "")
	infoLines = append(infoLines, fmt.Sprintf("Name: %s", selectedBranch.Name))
	infoLines = append(infoLines, "")
	infoLines = append(infoLines, "Availability:")

	// Show which remotes have this branch
	remote1Status := "✗ Not found"
	if selectedBranch.ExistsRemote1 {
		remote1Status = "✓ Available"
	}
	remote2Status := "✗ Not found"
	if selectedBranch.ExistsRemote2 {
		remote2Status = "✓ Available"
	}

	infoLines = append(infoLines, fmt.Sprintf("  %s: %s", m.remote1.Name, remote1Status))
	infoLines = append(infoLines, fmt.Sprintf("  %s: %s", m.remote2.Name, remote2Status))
	infoLines = append(infoLines, "")

	// Show sync status if branch exists on both remotes
	if selectedBranch.ExistsRemote1 && selectedBranch.ExistsRemote2 {
		if selectedBranch.Name == m.currentBranch && m.compareResult != nil {
			statusLine := "Status: "
			switch m.compareResult.Status {
			case git.InSync:
				statusLine += "In sync ✓"
			case git.Remote1Ahead:
				statusLine += fmt.Sprintf("%s ahead by %d commits", m.remote1.Name, m.compareResult.Remote1Ahead)
			case git.Remote2Ahead:
				statusLine += fmt.Sprintf("%s ahead by %d commits", m.remote2.Name, m.compareResult.Remote2Ahead)
			case git.Diverged:
				statusLine += "Diverged (manual merge required)"
			}
			infoLines = append(infoLines, statusLine)
		} else {
			infoLines = append(infoLines, "Status: Switch to this branch to see sync status")
		}
	} else if selectedBranch.ExistsRemote1 || selectedBranch.ExistsRemote2 {
		infoLines = append(infoLines, "Status: Branch only exists on one remote")
		infoLines = append(infoLines, "Press 'c' to create on the other remote")
	} else {
		infoLines = append(infoLines, "Status: Branch doesn't exist on either remote")
	}

	infoLines = append(infoLines, "")
	infoLines = append(infoLines, "Press any key to close")

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		dialogStyle.Render(strings.Join(infoLines, "\n")))
}

// Command functions for async operations

func fetchRemotes(repo *git.Repository, remote1, remote2 string) tea.Cmd {
	return func() tea.Msg {
		err1 := repo.Fetch(remote1)
		if err1 != nil {
			return fetchCompleteMsg{err: err1}
		}

		err2 := repo.Fetch(remote2)
		if err2 != nil {
			return fetchCompleteMsg{err: err2}
		}

		return fetchCompleteMsg{err: nil}
	}
}

func compareBranch(repo *git.Repository, remote1, remote2, branch string) tea.Cmd {
	return func() tea.Msg {
		result, err := repo.CompareBranch(remote1, remote2, branch)
		return compareCompleteMsg{result: result, err: err}
	}
}

func performSync(repo *git.Repository, result *git.CompareResult, remote1, remote2, branch string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if result.Status == git.Remote1Ahead {
			err = repo.SyncBranch(remote1, remote2, branch)
		} else if result.Status == git.Remote2Ahead {
			err = repo.SyncBranch(remote2, remote1, branch)
		}
		return syncCompleteMsg{err: err}
	}
}

func addRemote(repo *git.Repository, name, url string) tea.Cmd {
	return func() tea.Msg {
		err := repo.AddRemote(name, url)
		if err != nil {
			return addRemoteCompleteMsg{remote: nil, err: err}
		}

		// Get the newly added remote
		remoteURL, err := repo.GetRemoteURL(name)
		if err != nil {
			return addRemoteCompleteMsg{remote: nil, err: err}
		}

		return addRemoteCompleteMsg{
			remote: &git.Remote{Name: name, URL: remoteURL},
			err:    nil,
		}
	}
}

func loadBranches(repo *git.Repository, remote1, remote2 string) tea.Cmd {
	return func() tea.Msg {
		branches, err := repo.ListAllBranches(remote1, remote2)
		return branchesLoadedMsg{branches: branches, err: err}
	}
}

func switchBranch(repo *git.Repository, branch, remote1, remote2 string) tea.Cmd {
	return func() tea.Msg {
		// Note: We don't actually checkout the branch locally
		// We just change which branch we're viewing for comparison
		// This is intentional - the tool works with remote refs only
		return branchSwitchMsg{branch: branch, err: nil}
	}
}

func createBranch(repo *git.Repository, remoteName, branchName, sourceRef string) tea.Cmd {
	return func() tea.Msg {
		err := repo.CreateBranchOnRemote(remoteName, branchName, sourceRef)
		return branchCreateMsg{
			branchName: branchName,
			remoteName: remoteName,
			err:        err,
		}
	}
}
