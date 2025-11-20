package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/joel611/git-sync-remotes/internal/git"
	"github.com/joel611/git-sync-remotes/internal/models"
)

func main() {
	// Get working directory
	wd, err := git.GetWorkingDirectory()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Open repository
	repo, err := git.OpenRepository(wd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Get remotes
	remotes, err := repo.ListRemotes()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing remotes: %v\n", err)
		os.Exit(1)
	}

	// Validate remote count
	if len(remotes) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No git remotes found\n")
		fmt.Fprintf(os.Stderr, "You need at least one remote. Add one with:\n")
		fmt.Fprintf(os.Stderr, "  git remote add origin <url>\n")
		os.Exit(1)
	} else if len(remotes) > 2 {
		fmt.Fprintf(os.Stderr, "Error: Found %d remotes\n", len(remotes))
		fmt.Fprintf(os.Stderr, "Please specify which two remotes to sync (support for 3+ remotes coming soon)\n")
		for _, r := range remotes {
			fmt.Fprintf(os.Stderr, "  - %s: %s\n", r.Name, r.URL)
		}
		os.Exit(1)
	}

	// Get current branch
	currentBranch, err := repo.GetCurrentBranch()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current branch: %v\n", err)
		os.Exit(1)
	}

	// Initialize TUI model with available remotes
	var remote1, remote2 *git.Remote
	remote1 = &remotes[0]
	if len(remotes) > 1 {
		remote2 = &remotes[1]
	}

	m := models.NewModel(repo, remote1, remote2, currentBranch)

	// Run TUI
	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
