# Feature: Source Code Viewing

## Overview
Central source code editor displaying Go source files with syntax highlighting,
line numbers, execution indicators, and interactive breakpoint management.
Read-only view optimized for debugging workflows.

## User Capabilities

### 1. View Current Execution Location
**User Action:** Program pauses at any location (breakpoint, step, etc.)
**Expected Behavior:**
- Source file automatically loads if not already displayed
- Current line highlighted with green background
- View auto-scrolls to show current line in center of window with animation
- Line number clearly visible
- Stack trace shows corresponding frame
- File path displayed in window title or status
- Padding to the left of the line numbers shows the breakpoint status of the current line (red = enabled, grey = disabled)

### 2. Navigate to Different File
**User Action:** Click file in Sources panel, Functions panel, or Stack frame
**Expected Behavior:**
- Previous file content cleared
- New file loads into source view
- Syntax highlighting applied
- All breakpoints for that file shown with red highlights
- EventFileChanged emitted
- Open Files panel updated with new entry
- Scroll position resets to top (unless specific line requested)

### 3. Navigate to Specific Line
**User Action:** Click stack frame, function, or breakpoint referencing specific line
**Expected Behavior:**
- File loads (if different from current)
- View scrolls to show target line
- If line is current execution: green highlight shown
- If line has breakpoint: Red highlight shown
- Line becomes visible in viewport center
- No selection/cursor (read-only mode)

### 5. View Line Numbers
**User Action:** Open any source file
**Expected Behavior:**
- Line numbers displayed in left margin
- Numbers aligned and monospace
- Padding between numbers and code
- Numbers not selectable
- Clickable for breakpoint toggling
- Clear visual separation from code area

### 6. Toggle Breakpoint on Line
**User Action:** Click on the breakpoint margin to the left of the line number
**Expected Behavior:**
- See detailed behavior in `breakpoint_management.md`
- Visual feedback immediate (red highlight appears/disappears)
- No double-click required
- Works on any executable line
- Error handling for invalid lines

### 7. View Multiple Highlights Simultaneously
**User Action:** Program paused at breakpoint
**Expected Behavior:**
- Current line: green highlight (execution indicator)
- Breakpoint line: Red highlight (if same line, green line with red highlight for the line number)
- Multiple breakpoint lines: All shown with red highlights
- Visual layering: Current line yellow > breakpoint red > normal background
- No flicker or rendering issues

### 8. Scroll Through Code
**User Action:** Use scrollbar or scroll wheel in source view
**Expected Behavior:**
- Smooth scrolling
- Syntax highlighting remains correct throughout
- Line numbers scroll with code
- Breakpoint highlights remain positioned correctly
- Current line indicator remains visible if in viewport
- No lag with large files (1000+ lines)

### 9. View File Without Debug Symbols
**User Action:** Navigate to file not in current binary's debug info
**Expected Behavior:**
- File content still displayed
- Syntax highlighting works
- Cannot set breakpoints (show dialog)
- Console logs warning
- Line numbers shown normally

## UI Elements

### Source View Widget
**Components:**
- GtkSourceView with Go language definition
- Line number margin (left)
- Breakpoint margin (left of line number margin)
- Code content area (center)
- Scrollbars (right and bottom)

**Visual Properties:**
- System usual coding font
- Font size is the default of the system
- Line spacing the default of the system
- Tab width: 4 spaces (Go convention)
- Background: default
- Text color: default

### Line Highlighting
**Current Execution Line:**
- Background: green pastel
- Full line width

**Breakpoint Line:**
- Background: `rgba(255, 0, 0, 0.2)` (semi-transparent red)
- Full line width

