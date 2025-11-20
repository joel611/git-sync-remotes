package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Repository represents a git repository
type Repository struct {
	Path string
}

// OpenRepository opens a git repository at the given path
func OpenRepository(path string) (*Repository, error) {
	// Verify we're in a git repository
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("not a git repository: %w", err)
	}

	gitDir := strings.TrimSpace(string(output))
	if gitDir == "" {
		return nil, fmt.Errorf("not a git repository")
	}

	return &Repository{Path: path}, nil
}

// GetCurrentBranch returns the name of the current branch
func (r *Repository) GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" || branch == "HEAD" {
		return "", fmt.Errorf("detached HEAD or invalid branch")
	}

	return branch, nil
}

// GetWorkingDirectory returns the current working directory
func GetWorkingDirectory() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get working directory: %w", err)
	}
	return wd, nil
}
