package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ValidateBranchName checks if a branch name is valid according to git naming conventions
func ValidateBranchName(name string) error {
	if name == "" {
		return fmt.Errorf("branch name cannot be empty")
	}

	// Check for invalid characters and patterns
	invalidPatterns := []struct {
		pattern string
		message string
	}{
		{".", "cannot start with a dot"},
		{"..", "cannot contain '..'"},
		{"~", "cannot contain '~'"},
		{"^", "cannot contain '^'"},
		{":", "cannot contain ':'"},
		{"?", "cannot contain '?'"},
		{"*", "cannot contain '*'"},
		{"[", "cannot contain '['"},
		{" ", "cannot contain spaces"},
		{"@{", "cannot contain '@{'"},
	}

	// Check if starts with dot
	if strings.HasPrefix(name, ".") {
		return fmt.Errorf("branch name cannot start with a dot")
	}

	// Check if ends with .lock
	if strings.HasSuffix(name, ".lock") {
		return fmt.Errorf("branch name cannot end with '.lock'")
	}

	// Check for invalid patterns
	for _, invalid := range invalidPatterns {
		if strings.Contains(name, invalid.pattern) {
			return fmt.Errorf("branch name %s", invalid.message)
		}
	}

	// Check for control characters
	for _, ch := range name {
		if ch < 32 || ch == 127 {
			return fmt.Errorf("branch name cannot contain control characters")
		}
	}

	return nil
}

// Branch represents a git branch
type Branch struct {
	Name          string
	Remote        string
	ExistsRemote1 bool
	ExistsRemote2 bool
}

// ListRemoteBranches returns all branches on the specified remote
func (r *Repository) ListRemoteBranches(remoteName string) ([]string, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", remoteName)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches on %s: %w", remoteName, err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}

	branches := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Format: <sha> refs/heads/<branch-name>
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		ref := parts[1]
		if !strings.HasPrefix(ref, "refs/heads/") {
			continue
		}

		branchName := strings.TrimPrefix(ref, "refs/heads/")
		branches = append(branches, branchName)
	}

	return branches, nil
}

// BranchExistsOnRemote checks if a branch exists on the specified remote
func (r *Repository) BranchExistsOnRemote(remoteName, branchName string) (bool, error) {
	cmd := exec.Command("git", "ls-remote", "--heads", remoteName, fmt.Sprintf("refs/heads/%s", branchName))
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to check branch on %s: %w", remoteName, err)
	}

	result := strings.TrimSpace(string(output))
	return result != "", nil
}

// CreateBranchOnRemote creates a branch on the specified remote from a source ref with a 30 second timeout
func (r *Repository) CreateBranchOnRemote(remoteName, branchName, sourceRef string) error {
	// Validate branch name
	if err := ValidateBranchName(branchName); err != nil {
		return fmt.Errorf("invalid branch name: %w", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Push the source ref to the remote as the new branch
	refSpec := fmt.Sprintf("%s:refs/heads/%s", sourceRef, branchName)
	cmd := exec.CommandContext(ctx, "git", "push", remoteName, refSpec)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("push to %s timed out after 30 seconds", remoteName)
		}
		return fmt.Errorf("failed to create branch %s on %s: %w\nOutput: %s",
			branchName, remoteName, err, string(output))
	}

	return nil
}

// ListAllBranches returns all branches across both remotes
func (r *Repository) ListAllBranches(remote1, remote2 string) ([]Branch, error) {
	branches1, err := r.ListRemoteBranches(remote1)
	if err != nil {
		return nil, err
	}

	branches2, err := r.ListRemoteBranches(remote2)
	if err != nil {
		return nil, err
	}

	// Create a map to track unique branches
	branchMap := make(map[string]*Branch)

	for _, name := range branches1 {
		branchMap[name] = &Branch{
			Name:          name,
			ExistsRemote1: true,
			ExistsRemote2: false,
		}
	}

	for _, name := range branches2 {
		if b, exists := branchMap[name]; exists {
			b.ExistsRemote2 = true
		} else {
			branchMap[name] = &Branch{
				Name:          name,
				ExistsRemote1: false,
				ExistsRemote2: true,
			}
		}
	}

	// Convert map to slice
	result := make([]Branch, 0, len(branchMap))
	for _, branch := range branchMap {
		result = append(result, *branch)
	}

	return result, nil
}
