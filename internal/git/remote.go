package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// Remote represents a git remote
type Remote struct {
	Name string
	URL  string
}

// ListRemotes returns all remotes configured in the repository
func (r *Repository) ListRemotes() ([]Remote, error) {
	cmd := exec.Command("git", "remote")
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to list remotes: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []Remote{}, nil
	}

	remotes := make([]Remote, 0, len(lines))
	for _, name := range lines {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		url, err := r.GetRemoteURL(name)
		if err != nil {
			return nil, err
		}

		remotes = append(remotes, Remote{
			Name: name,
			URL:  url,
		})
	}

	return remotes, nil
}

// GetRemoteURL returns the URL for a given remote name
func (r *Repository) GetRemoteURL(name string) (string, error) {
	cmd := exec.Command("git", "remote", "get-url", name)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get URL for remote %s: %w", name, err)
	}

	return strings.TrimSpace(string(output)), nil
}

// RemoteExists checks if a remote with the given name exists
func (r *Repository) RemoteExists(name string) (bool, error) {
	remotes, err := r.ListRemotes()
	if err != nil {
		return false, err
	}

	for _, remote := range remotes {
		if remote.Name == name {
			return true, nil
		}
	}

	return false, nil
}

// Fetch fetches from the specified remote
func (r *Repository) Fetch(remoteName string) error {
	cmd := exec.Command("git", "fetch", remoteName)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to fetch from %s: %w\nOutput: %s", remoteName, err, string(output))
	}

	return nil
}

// AddRemote adds a new remote to the repository
func (r *Repository) AddRemote(name, url string) error {
	// Check if remote already exists
	exists, err := r.RemoteExists(name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("remote '%s' already exists", name)
	}

	// Add the remote
	cmd := exec.Command("git", "remote", "add", name, url)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to add remote %s: %w\nOutput: %s", name, err, string(output))
	}

	return nil
}
