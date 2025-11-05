# git-sync-remotes

A powerful bash script to sync commits between two git remotes automatically. Perfect for maintaining mirrors, syncing between internal and external repositories, or keeping multiple remotes in sync.

## Features

- **Auto-detect remotes**: Automatically uses both remotes if exactly 2 exist
- **Auto-detect branch**: Uses current branch if not specified
- **Smart sync**: Detects which remote is ahead and syncs accordingly
- **Branch creation**: Asks to create missing branches on remotes
- **Divergence detection**: Warns when both remotes have unique commits
- **Preview changes**: Shows commit differences before syncing
- **Auto-confirm mode**: Skip confirmations with `-y` flag
- **Color-coded output**: Easy to read status messages
- **Safe operations**: Always confirms before pushing

## Requirements

- Bash 4.0 or higher
- Git 2.0 or higher
- Unix-like environment (Linux, macOS, WSL)

## Installation

### Quick Install (Recommended)

```bash
git clone https://github.com/joel611/git-sync-remotes.git
cd git-sync-remotes
./install.sh
```

The install script will:
1. Create `~/.local/bin/` directory if it doesn't exist
2. Create a symbolic link to the script
3. Make the script executable
4. Check if `~/.local/bin` is in your PATH

### Manual Installation

```bash
# Clone the repository
git clone https://github.com/joel611/git-sync-remotes.git
cd git-sync-remotes

# Create symlink
mkdir -p ~/.local/bin
ln -s "$(pwd)/git-sync-remotes" ~/.local/bin/git-sync-remotes
chmod +x git-sync-remotes

# Ensure ~/.local/bin is in your PATH (add to ~/.zshrc or ~/.bashrc if needed)
export PATH="$PATH:$HOME/.local/bin"
```

## Uninstallation

### Quick Uninstall

```bash
cd git-sync-remotes
./uninstall.sh
```

The uninstall script will:
1. Remove the symbolic link or script from `~/.local/bin/`
2. Clean up empty directories if applicable

### Manual Uninstallation

```bash
# Remove the symbolic link or script
rm ~/.local/bin/git-sync-remotes
```

## Usage

### Basic Usage

```bash
# Sync current branch with auto-detected remotes
git-sync-remotes

# Sync current branch with auto-confirm (no prompts)
git-sync-remotes -y

# Sync specific branch
git-sync-remotes master

# Sync with specific remotes
git-sync-remotes master origin gitlab
```

### Options

- `-y`, `--yes`: Auto-confirm all prompts (useful for automation)
- `[branch-name]`: Specify branch to sync (defaults to current branch)
- `[remote1] [remote2]`: Specify which remotes to sync (auto-detects if not provided)

### Examples

#### Example 1: Sync Current Branch
```bash
$ git-sync-remotes -y
ℹ No branch specified, using current branch: main
ℹ Auto-detecting remotes...
✓ Detected 2 remotes: origin, gitlab
==========================================
Git Remote Sync Tool
==========================================

ℹ Branch: main
ℹ Remote 1: origin (git@github.com:user/repo.git)
ℹ Remote 2: gitlab (git@gitlab.com:user/repo.git)

ℹ Fetching from both remotes...
✓ Fetched from both remotes

ℹ origin has 3 commit(s) not in gitlab
ℹ Last commit: 2024-01-15 10:30:00 +0000

Latest commits on origin:
* abc1234 Add new feature
* def5678 Fix bug
* ghi9012 Update docs

ℹ Direction: origin → gitlab

ℹ Pushing to gitlab...
✓ Successfully synced to gitlab
```

#### Example 2: Missing Branch
```bash
$ git-sync-remotes feature-branch
⚠ Branch 'feature-branch' does not exist on origin
ℹ Branch exists on: gitlab

Create branch 'feature-branch' on origin? [y/n]: y
ℹ Creating and pushing branch to origin...
✓ Branch created and synced to origin
```

#### Example 3: Diverged Branches
```bash
$ git-sync-remotes main
ERROR: DIVERGED: Both remotes have unique commits!

⚠ Manual intervention required. Please resolve the divergence manually.
ℹ You may need to:
  1. Checkout the branch locally
  2. Merge or rebase the changes
  3. Push to both remotes
```

## Use Cases

### 1. Mirror Internal and External Repositories
Keep an internal GitLab and external GitHub repository in sync:
```bash
git-sync-remotes -y  # Run regularly to keep them synced
```

### 2. Backup to Multiple Remotes
Automatically backup your repository to multiple remotes:
```bash
git-sync-remotes -y main origin backup
```

### 3. Sync Feature Branches
Easily sync feature branches between remotes:
```bash
git-sync-remotes feature/new-api
```

### 4. Automated Sync with CI/CD
Add to your CI/CD pipeline:
```bash
git-sync-remotes -y ${CI_COMMIT_BRANCH}
```

## How It Works

1. **Fetch**: Downloads latest changes from both remotes
2. **Compare**: Checks commit history to see which remote is ahead
3. **Analyze**: Determines sync direction or detects divergence
4. **Preview**: Shows what will be synced
5. **Confirm**: Asks for confirmation (unless `-y` flag is used)
6. **Sync**: Pushes commits from ahead remote to behind remote

## Troubleshooting

### "Not a git repository" Error
Make sure you run the command from inside a git repository.

### "Only one remote found" Error
The script requires at least 2 remotes. Add another remote:
```bash
git remote add backup git@example.com:user/repo.git
```

### "DIVERGED" Warning
Both remotes have unique commits. Manually resolve:
```bash
git fetch origin
git fetch gitlab
git checkout main
git merge origin/main
git push origin main
git push gitlab main
```

### Permission Denied (SSH)
Ensure you have SSH keys set up for both remotes:
```bash
ssh -T git@github.com
ssh -T git@gitlab.com
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

Created by Joel Chan

## Acknowledgments

- Inspired by the need to keep multiple git remotes in sync
- Built with bash for maximum portability
- Uses git's powerful plumbing commands

## Support

If you encounter any issues or have questions:
- Open an issue on [GitHub](https://github.com/joel611/git-sync-remotes/issues)
- Check existing issues for solutions

---

**Happy Syncing! 🚀**
