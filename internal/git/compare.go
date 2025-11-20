package git

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// SyncStatus represents the synchronization status between two remotes
type SyncStatus int

const (
	InSync       SyncStatus = iota // Both remotes have identical commits
	Remote1Ahead                    // Remote 1 is ahead of remote 2
	Remote2Ahead                    // Remote 2 is ahead of remote 1
	Diverged                        // Both remotes have unique commits
	BranchMissing                   // Branch doesn't exist on one or both remotes
)

func (s SyncStatus) String() string {
	switch s {
	case InSync:
		return "in sync"
	case Remote1Ahead:
		return "remote1 ahead"
	case Remote2Ahead:
		return "remote2 ahead"
	case Diverged:
		return "diverged"
	case BranchMissing:
		return "branch missing"
	default:
		return "unknown"
	}
}

// Commit represents a git commit
type Commit struct {
	SHA       string
	ShortSHA  string
	Message   string
	Author    string
	Date      time.Time
	DateStr   string
}

// CompareResult contains the comparison between two remotes for a branch
type CompareResult struct {
	Status           SyncStatus
	Remote1Ahead     int
	Remote2Ahead     int
	Remote1SHA       string
	Remote2SHA       string
	Remote1Commits   []Commit
	Remote2Commits   []Commit
	Remote1HasBranch bool
	Remote2HasBranch bool
}

// CompareBranch compares a branch between two remotes
func (r *Repository) CompareBranch(remote1, remote2, branch string) (*CompareResult, error) {
	remote1Ref := fmt.Sprintf("%s/%s", remote1, branch)
	remote2Ref := fmt.Sprintf("%s/%s", remote2, branch)

	result := &CompareResult{}

	// Check if branches exist on remotes
	sha1, err1 := r.getCommitSHA(remote1Ref)
	result.Remote1HasBranch = (err1 == nil)

	sha2, err2 := r.getCommitSHA(remote2Ref)
	result.Remote2HasBranch = (err2 == nil)

	// If branch doesn't exist on one or both remotes
	if !result.Remote1HasBranch || !result.Remote2HasBranch {
		result.Status = BranchMissing
		if result.Remote1HasBranch {
			result.Remote1SHA = sha1
		}
		if result.Remote2HasBranch {
			result.Remote2SHA = sha2
		}
		return result, nil
	}

	result.Remote1SHA = sha1
	result.Remote2SHA = sha2

	// If SHAs are identical, remotes are in sync
	if sha1 == sha2 {
		result.Status = InSync
		return result, nil
	}

	// Count commits ahead
	ahead1, err := r.countCommitsAhead(remote2Ref, remote1Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to count commits ahead for %s: %w", remote1, err)
	}

	ahead2, err := r.countCommitsAhead(remote1Ref, remote2Ref)
	if err != nil {
		return nil, fmt.Errorf("failed to count commits ahead for %s: %w", remote2, err)
	}

	result.Remote1Ahead = ahead1
	result.Remote2Ahead = ahead2

	// Determine status
	if ahead1 > 0 && ahead2 > 0 {
		result.Status = Diverged
	} else if ahead1 > 0 {
		result.Status = Remote1Ahead
		// Get commits unique to remote1
		commits, err := r.getCommitsBetween(remote2Ref, remote1Ref, 50)
		if err != nil {
			return nil, err
		}
		result.Remote1Commits = commits
	} else if ahead2 > 0 {
		result.Status = Remote2Ahead
		// Get commits unique to remote2
		commits, err := r.getCommitsBetween(remote1Ref, remote2Ref, 50)
		if err != nil {
			return nil, err
		}
		result.Remote2Commits = commits
	}

	return result, nil
}

// getCommitSHA returns the commit SHA for a ref
func (r *Repository) getCommitSHA(ref string) (string, error) {
	cmd := exec.Command("git", "rev-parse", ref)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get SHA for %s: %w", ref, err)
	}

	return strings.TrimSpace(string(output)), nil
}

// countCommitsAhead counts how many commits ref2 is ahead of ref1
func (r *Repository) countCommitsAhead(ref1, ref2 string) (int, error) {
	revRange := fmt.Sprintf("%s..%s", ref1, ref2)
	cmd := exec.Command("git", "rev-list", "--count", revRange)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("failed to count commits: %w", err)
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, fmt.Errorf("invalid count output: %w", err)
	}

	return count, nil
}

// getCommitsBetween returns commits between two refs
func (r *Repository) getCommitsBetween(ref1, ref2 string, limit int) ([]Commit, error) {
	revRange := fmt.Sprintf("%s..%s", ref1, ref2)
	args := []string{"log", "--format=%H|%h|%s|%an|%ai", fmt.Sprintf("-n%d", limit), revRange}
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []Commit{}, nil
	}

	commits := make([]Commit, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "|", 5)
		if len(parts) < 5 {
			continue
		}

		date, err := time.Parse("2006-01-02 15:04:05 -0700", parts[4])
		if err != nil {
			date = time.Now() // Fallback
		}

		commits = append(commits, Commit{
			SHA:      parts[0],
			ShortSHA: parts[1],
			Message:  parts[2],
			Author:   parts[3],
			Date:     date,
			DateStr:  parts[4],
		})
	}

	return commits, nil
}

// GetCommitDetails returns detailed information about a commit
func (r *Repository) GetCommitDetails(sha string) (*Commit, error) {
	cmd := exec.Command("git", "show", "--format=%H|%h|%s|%an|%ai|%b", "--no-patch", sha)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit details: %w", err)
	}

	line := strings.TrimSpace(string(output))
	parts := strings.SplitN(line, "|", 6)
	if len(parts) < 5 {
		return nil, fmt.Errorf("invalid commit output")
	}

	date, err := time.Parse("2006-01-02 15:04:05 -0700", parts[4])
	if err != nil {
		date = time.Now()
	}

	commit := &Commit{
		SHA:      parts[0],
		ShortSHA: parts[1],
		Message:  parts[2],
		Author:   parts[3],
		Date:     date,
		DateStr:  parts[4],
	}

	return commit, nil
}

// GetCommitDiff returns the diff for a commit
func (r *Repository) GetCommitDiff(sha string) (string, error) {
	cmd := exec.Command("git", "show", "--format=", sha)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get commit diff: %w", err)
	}

	return string(output), nil
}
