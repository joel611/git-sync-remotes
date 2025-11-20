package git

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Push pushes a ref to a remote branch with a 30 second timeout
func (r *Repository) Push(remoteName, sourceRef, branchName string) error {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	refSpec := fmt.Sprintf("%s:refs/heads/%s", sourceRef, branchName)
	cmd := exec.CommandContext(ctx, "git", "push", remoteName, refSpec)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("push to %s timed out after 30 seconds", remoteName)
		}
		return fmt.Errorf("failed to push to %s: %w\nOutput: %s", remoteName, err, string(output))
	}

	return nil
}

// SyncBranch syncs a branch from sourceRemote to destRemote
func (r *Repository) SyncBranch(sourceRemote, destRemote, branch string) error {
	sourceRef := fmt.Sprintf("%s/%s", sourceRemote, branch)
	return r.Push(destRemote, sourceRef, branch)
}
