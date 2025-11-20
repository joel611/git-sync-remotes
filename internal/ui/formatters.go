package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/joel611/git-sync-remotes/internal/git"
)

var (
	// Color styles for sync states
	inSyncStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))  // Green
	aheadStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("220")) // Yellow
	divergedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // Red
)

// FormatSyncStatus formats the sync status for display with color coding
func FormatSyncStatus(result *git.CompareResult, remote1, remote2 string) string {
	switch result.Status {
	case git.InSync:
		return inSyncStyle.Render(fmt.Sprintf("✓ %s and %s are in sync", remote1, remote2))
	case git.Remote1Ahead:
		return aheadStyle.Render(fmt.Sprintf("→ %s has %d commit(s) ahead of %s",
			remote1, result.Remote1Ahead, remote2))
	case git.Remote2Ahead:
		return aheadStyle.Render(fmt.Sprintf("← %s has %d commit(s) ahead of %s",
			remote2, result.Remote2Ahead, remote1))
	case git.Diverged:
		return divergedStyle.Render(fmt.Sprintf("⚠ Diverged: %s has %d, %s has %d unique commits",
			remote1, result.Remote1Ahead, remote2, result.Remote2Ahead))
	case git.BranchMissing:
		// Determine which remote(s) are missing the branch
		if !result.Remote1HasBranch && !result.Remote2HasBranch {
			return aheadStyle.Render(fmt.Sprintf("⚠ Branch doesn't exist on %s or %s. Push to create it.", remote1, remote2))
		} else if !result.Remote1HasBranch {
			return aheadStyle.Render(fmt.Sprintf("⚠ Branch missing on %s. Press 'b' then 'c' to create it.", remote1))
		} else {
			return aheadStyle.Render(fmt.Sprintf("⚠ Branch missing on %s. Press 'b' then 'c' to create it.", remote2))
		}
	default:
		return "Unknown status"
	}
}
