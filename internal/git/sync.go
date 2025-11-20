package git

import (
	"fmt"
	"os/exec"
)

// Push pushes a ref to a remote branch
func (r *Repository) Push(remoteName, sourceRef, branchName string) error {
	refSpec := fmt.Sprintf("%s:refs/heads/%s", sourceRef, branchName)
	cmd := exec.Command("git", "push", remoteName, refSpec)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to push to %s: %w\nOutput: %s", remoteName, err, string(output))
	}

	return nil
}

// SyncBranch syncs a branch from sourceRemote to destRemote
func (r *Repository) SyncBranch(sourceRemote, destRemote, branch string) error {
	sourceRef := fmt.Sprintf("%s/%s", sourceRemote, branch)
	return r.Push(destRemote, sourceRef, branch)
}
