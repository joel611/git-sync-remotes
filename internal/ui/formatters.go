package ui

import (
	"fmt"

	"github.com/joel611/git-sync-remotes/internal/git"
)

// FormatSyncStatus formats the sync status for display
func FormatSyncStatus(result *git.CompareResult, remote1, remote2 string) string {
	switch result.Status {
	case git.InSync:
		return fmt.Sprintf("✓ %s and %s are in sync", remote1, remote2)
	case git.Remote1Ahead:
		return fmt.Sprintf("→ %s has %d commit(s) ahead of %s",
			remote1, result.Remote1Ahead, remote2)
	case git.Remote2Ahead:
		return fmt.Sprintf("← %s has %d commit(s) ahead of %s",
			remote2, result.Remote2Ahead, remote1)
	case git.Diverged:
		return fmt.Sprintf("⚠ Diverged: %s has %d, %s has %d unique commits",
			remote1, result.Remote1Ahead, remote2, result.Remote2Ahead)
	default:
		return "Unknown status"
	}
}
