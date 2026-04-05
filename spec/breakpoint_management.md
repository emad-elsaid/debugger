# Feature: Breakpoint Management

## Overview
Comprehensive breakpoint creation, modification, and inspection across multiple UI interfaces.
Users can set breakpoints via source view clicks, manage them in dedicated panels,
and control execution with precise line-level granularity.

## User Capabilities

### 1. Create or enable or delete Breakpoint (Source View)
**User Action:** Left-click on gutter to the left line number in source view
**Expected Behavior:**
- If no breakpoint exists: Create new breakpoint
- If break point exists and enabled: disable it
- If break point exists and disabled: delete it
- Line background changes to reflect the breakpoint status (none = no breakpoint, red = enabled, grey = disabled)
- Breakpoint appears in Breakpoints panel
- Breakpoint appears in Sidebar breakpoints widget
- Red/grey point appears in the breakpoint gutter
- Console logs: "[HH:MM:SS] Breakpoint (operation) at file.go:42"
- EventBreakpointsChanged emitted
- All breakpoint displays update immediately

### 2. Enable All Breakpoints
**User Action:** Click "Enable All" button in Breakpoints panel
**Expected Behavior:**
- All breakpoints set to enabled state
- All toggle switches turn ON
- All source view highlights restored
- Console logs: "[HH:MM:SS] All breakpoints enabled"
- EventBreakpointsChanged emitted

### 3. Disable All Breakpoints
**User Action:** Click "Disable All" button in Breakpoints panel
**Expected Behavior:**
- All breakpoints set to disabled state
- All toggle switches turn OFF
- Source view highlights removed
- Program runs without stopping at any breakpoint
- Console logs: "[HH:MM:SS] All breakpoints disabled"
- EventBreakpointsChanged emitted

### 4. Clear All Breakpoints
**User Action:** Click "Clear All" button in Breakpoints panel
**Expected Behavior:**
- All breakpoints deleted permanently
- Breakpoints panel shows empty state
- Sidebar breakpoints widget shows empty state
- All source view highlights removed
- Console logs: "[HH:MM:SS] All breakpoints cleared"
- EventBreakpointsChanged emitted
- Confirmation dialog NOT shown (direct clear)

### 5. Breakpoint Hit
**User Action:** Program execution reaches enabled breakpoint
**Expected Behavior:**
- Process pauses immediately
- Current line indicator (yellow) shown on breakpoint line
- Hit count increments by 1 in breakpoint display
- Stack trace updates to show current call stack
- Console logs: "[HH:MM:SS] Breakpoint hit at file.go:42"
- Source view scrolls to breakpoint location
- EventStateChanged emitted (paused state)

### 6. View Breakpoint Hit Count
**User Action:** Observe breakpoint entry in panel
**Expected Behavior:**
- "Hit count: N" displayed below location
- Count starts at 0 for new breakpoints
- Increments each time breakpoint is reached
- Persists across pause/continue cycles
- Resets to 0 on process restart
- Updated in real-time when breakpoint hit

