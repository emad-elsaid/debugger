# Feature: Sources Panel

## Overview
Hierarchical source file browser displaying all source files from the debugged binary,
organized by package/directory with search and navigation capabilities.

## User Capabilities

### 1. View All Source Files
**User Action:** Click "Sources" tab in bottom panel
**Expected Behavior:**
- Tree view displays all source files included in binary
- Files organized hierarchically by package/directory structure
- Each package node is expandable/collapsible
- All packages expanded by default on first load
- File count visible (implicitly by tree size)
- Updates when process restarts (binary changes)

**Organization Example:**
```
📁 main
  └─ main.go
  └─ server.go
📁 github.com/user/project/handlers
  └─ api.go
  └─ auth.go
📁 github.com/user/project/models
  └─ user.go
  └─ database.go
📁 runtime
  └─ asm.go
  └─ proc.go
```

### 2. Expand/Collapse Package
**User Action:** Click on package/directory node
**Expected Behavior:**
- Package expands to show contained files (if collapsed)
- Package collapses to hide files (if expanded)
- Smooth animation (if GTK provides)
- State independent of other packages
- Folder icon changes: 📁 (closed) ↔ 📂 (open)
- Child files indent under parent

### 3. Open Source File
**User Action:** Click on file name in tree
**Expected Behavior:**
- EventFileChanged emitted with file path
- Source view loads and displays clicked file
- Syntax highlighting applied
- Scroll position resets to top
- Breakpoints for that file shown (if any)
- Open Files panel updated with new entry
- Sources panel remains visible (doesn't auto-hide)
- Selected file highlighted in tree

### 4. Search/Filter Files
**User Action:** Type in search entry field at top of panel
**Expected Behavior:**
- Tree filters in real-time as user types
- Matching files remain visible
- Non-matching files hidden
- Package nodes shown if any child matches
- Empty packages hidden
- Case-insensitive matching
- Matches anywhere in file path (not just filename)
- Clear button (X) appears in search field when text present

**Search Examples:**
- `"main"` → Shows main.go, main_test.go, any file with "main" in path
- `"handler"` → Shows files in handlers/ package, handler.go, etc.
- `"user.go"` → Shows user.go, user_test.go
- `""` (empty) → Shows all files (filter cleared)

### 5. Clear Search Filter
**User Action:** Click X button in search field or clear text manually
**Expected Behavior:**
- All files and packages become visible again
- Tree returns to pre-filter state
- Expansion states preserved (previously expanded nodes stay expanded)
- Search field clears
- Cursor remains in search field for new search

### 6. Refresh After Restart
**User Action:** Restart debugger process 
**Expected Behavior:**
- Sources panel clears
- Re-fetches source files from new binary
- Re-populates tree with potentially updated file list
- New files appear if added to project
- Deleted files removed from tree
- Search filter clears
- All packages re-expanded
