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

// Init initializes the model
func (m Model) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}

	// Only fetch if we have 2 remotes
	if m.remote2 != nil {
		cmds = append(cmds, fetchRemotes(m.repo, m.remote1.Name, m.remote2.Name))
	} else {
		// Show message that user needs to add a second remote
		m.message = "Only one remote found. Press 'a' to add a second remote."
	}

	return tea.Batch(cmds...)
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
			m.message = fmt.Sprintf("Fetch failed: %v", msg.err)
			return m, nil
		}
		m.message = "Fetch complete"
		return m, compareBranch(m.repo, m.remote1.Name, m.remote2.Name, m.currentBranch)

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
	if m.loading {
		status = fmt.Sprintf("%s Loading...", m.spinner.View())
	} else if m.compareResult != nil {
		status = ui.FormatSyncStatus(m.compareResult, m.remote1.Name, m.remote2.Name)
	} else {
		status = "Initializing..."
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
	if m.compareResult == nil {
		return "Loading..."
	}

	// For now, simple commit list view
	left := m.renderCommitList(m.remote1.Name, m.compareResult.Remote1Commits, m.focusedPane == CommitList1)
	right := m.renderCommitList(m.remote2.Name, m.compareResult.Remote2Commits, m.focusedPane == CommitList2)

	columns := lipgloss.JoinHorizontal(lipgloss.Top,
		left,
		right,
	)

	return columns
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

	var lines []string
	lines = append(lines, lipgloss.NewStyle().Bold(true).Render(fmt.Sprintf("COMMITS (%s)", remoteName)))
	lines = append(lines, "")

	if len(commits) == 0 {
		lines = append(lines, "  (no unique commits)")
	} else {
		for i, commit := range commits {
			prefix := "  "
			if focused && i == m.selectedIndex {
				prefix = "> "
			}
			lines = append(lines, fmt.Sprintf("%s%s %s", prefix, commit.ShortSHA, commit.Message))
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
    f           Fetch from remotes
    s           Sync commits

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
