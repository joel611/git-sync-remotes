# Remote Management Feature

## ✨ New Feature: Add Remotes from TUI

The TUI now supports working with repositories that have only 1 remote and allows you to add a second remote directly from the interface!

## 🎯 What Changed

### 1. **Single Remote Support**
- TUI now launches successfully with only 1 remote configured
- Shows helpful message: "Only one remote found. Press 'a' to add a second remote."
- Gracefully handles the case where you don't have 2 remotes yet

### 2. **Add Remote Dialog** (New!)
- Press `a` to open the "Add Remote" dialog
- Interactive form with two fields:
  - **Name**: Remote name (e.g., "gitlab", "backup", "mirror")
  - **URL**: Git URL (e.g., "git@github.com:user/repo.git")
- Tab/↑/↓ to switch between fields
- Type to enter values
- Backspace to delete characters
- Enter to submit
- Esc to cancel

### 3. **Smart UI Feedback**
- Footer shows context-aware shortcuts:
  - With 1 remote: `[a]dd remote [q]uit [?]help`
  - With 2 remotes: `[f]etch [s]ync [q]uit [?]help`
- Pressing `f` or `s` with only 1 remote shows: "Need 2 remotes to [action]. Press 'a' to add a second remote."
- Success message after adding remote: "Remote 'name' added successfully! Press 'f' to fetch."

### 4. **Updated Help**
- Help overlay (press `?`) now includes:
  ```
  Actions:
    a           Add remote (when only 1 remote exists)
    f           Fetch from remotes
    s           Sync commits
  ```

## 🚀 How to Use

### Scenario 1: Starting with 1 Remote

```bash
# Your repo has only 1 remote
git remote -v
# origin  git@github.com:user/repo.git (fetch)
# origin  git@github.com:user/repo.git (push)

# Launch TUI
./git-sync-remotes-tui

# TUI shows: "Only one remote found. Press 'a' to add a second remote."
# Press 'a' to open Add Remote dialog
# Enter name: gitlab
# Tab to URL field
# Enter URL: git@gitlab.com:user/repo.git
# Press Enter to submit

# Success! Now you can use 'f' to fetch and 's' to sync
```

### Scenario 2: Adding a Backup Remote

```bash
# Already have origin, want to add a backup
./git-sync-remotes-tui

# Press 'a'
# Name: backup
# URL: git@bitbucket.org:user/repo.git
# Enter

# Remote added! Now you can sync between origin and backup
```

## 📝 Implementation Details

### Files Modified

1. **cmd/tui/main.go**
   - Changed to accept 1+ remotes (previously required exactly 2)
   - Passes remotes as pointers to handle optional second remote
   - Better error messaging for 0 remotes case

2. **internal/git/remote.go**
   - Added `AddRemote(name, url string) error` function
   - Validates remote doesn't already exist before adding
   - Uses `git remote add` command

3. **internal/models/app.go**
   - Changed `remote1` and `remote2` from values to pointers (`*git.Remote`)
   - Added `showAddRemote` flag and form fields (`addRemoteName`, `addRemoteURL`, `addRemoteField`)
   - Added `addRemoteCompleteMsg` message type
   - Updated `Init()` to skip fetch when only 1 remote exists
   - Added keyboard handling for add remote dialog (a, Tab, Enter, Esc, Backspace, typing)
   - Added `renderAddRemoteDialog()` function with highlighted input fields
   - Updated `renderFooter()` to show context-aware shortcuts
   - Updated help text to include `a` key
   - Added `addRemote()` command function for async git operation
   - Smart handling: 'f' and 's' keys now check if remote2 exists before allowing actions

### Key Design Decisions

1. **Pointers for Optional Remote**
   - Use `*git.Remote` instead of `git.Remote` so `remote2` can be `nil`
   - Clean way to represent "no second remote yet"

2. **Interactive Form**
   - Minimal but functional text input
   - Field highlighting shows which field is active
   - Tab navigation between fields
   - Placeholder text when empty: "(enter name)" and "(enter URL)"

3. **Async Operation**
   - `addRemote()` returns a command that runs in background
   - Updates model when complete via `addRemoteCompleteMsg`
   - Shows loading spinner during operation

4. **Graceful Degradation**
   - All sync/fetch operations check if `remote2 != nil`
   - Clear error messages guide user to add remote
   - No crashes or confusing behavior with 1 remote

## ✅ Testing

**Build:** ✓ Compiles successfully
**Integration:** Code changes are complete and ready for manual testing

**To test manually:**
```bash
# Build
make build

# Test with 1 remote
cd /path/to/repo/with/1/remote
/path/to/git-sync-remotes-tui

# Actions to test:
# 1. TUI launches and shows helpful message
# 2. Press 'a' - dialog appears
# 3. Type remote name
# 4. Tab to URL field
# 5. Type remote URL
# 6. Press Enter - remote is added
# 7. Press 'f' - fetch now works
# 8. Press 's' - sync now works (after fetch)
```

## 🎨 UI Preview

### Add Remote Dialog
```
┌─────────────────────────────────┐
│                                 │
│  Add Remote                     │
│                                 │
│  > Name:  gitlab                │
│    URL:   (enter URL)           │
│                                 │
│  Press Tab to switch fields     │
│  Press Enter to submit          │
│  Press Esc to cancel            │
│                                 │
└─────────────────────────────────┘
```

### Footer (1 remote)
```
[a]dd remote [q]uit [?]help
```

### Footer (2 remotes)
```
[f]etch [s]ync [q]uit [?]help
```

## 🔄 Next Steps

This feature is complete and ready to use! Future enhancements could include:

1. **Edit/Remove Remotes**: Add ability to modify or delete existing remotes
2. **Remote Selection**: When 3+ remotes exist, allow selecting which 2 to use
3. **Validation**: Add URL format validation before submitting
4. **History**: Remember previously used remote URLs
5. **Import**: Detect remotes from `.git/config` of other repos

## 📊 Impact

- **User Experience**: Much improved! No need to exit TUI to add remotes
- **Flexibility**: Can start with just 1 remote and grow organically
- **Workflow**: Seamless onboarding for new repositories
- **Code Quality**: Clean architecture with optional remotes properly handled

This feature makes the TUI more user-friendly and self-contained! 🎉
