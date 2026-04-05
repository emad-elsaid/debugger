# Feature: Stack Trace Inspection

## Overview
Displays call stack frames for the current goroutine, allowing navigation through
execution context and inspection of different stack levels.

## User Capabilities

### 1. View Current Stack Trace
**User Action:** Program pauses (breakpoint, step, or halt)
**Expected Behavior:**
- Stack panel automatically updates
- Shows frames from current (top) to bottom (main/goroutine start)
- Each frame displays:
  - Frame number: `#0`, `#1`, `#2`, etc. (top is #0)
  - Function name: Fully qualified (e.g., `main.processRequest`)
  - Location: `filename:line` (basename only for display)
  - Full path: Available in tooltip on hover
- Topmost frame (current execution) highlighted/selected
- Limited to 50 frames maximum (STACK_LIMIT)
- Auto-scrolls to show current frame

### 2. Select Stack Frame
**User Action:** Click on any frame in stack trace list
**Expected Behavior:**
- Selected frame becomes highlighted
- Source view navigates to selected frame's file and line
- Variables panel (if implemented) shows:
  - Local variables for selected frame
  - Function arguments for selected frame
- Disassembly view updates to selected frame's assembly
- EventStackFrameChanged emitted
- Current execution line (yellow) remains at actual PC location
- Selected frame marked with distinct visual indicator

### 3. Navigate Up the Stack
**User Action:** Click frame closer to bottom (#0 → #1 → #2)
**Expected Behavior:**
- Shows who called the current function
- Source view shows caller's code at call site
- Call site line becomes visible (not highlighted as current)
- Variables show caller's scope
- Useful for understanding how execution reached current point
- Frame selection highlights new frame

### 4. Navigate Down the Stack
**User Action:** Click frame closer to top (from #2 → #1 → #0)
**Expected Behavior:**
- Returns toward current execution point
- Source view shows callee's code
- Selecting #0 returns to actual execution location
- Variables show callee's scope
- Frame selection highlights new frame

### 5. View Stack on Goroutine
**User Action:** Program pauses in any goroutine
**Expected Behavior:**
- Stack shows frames for current/selected goroutine only
- Goroutine ID visible (if UI displays it)
- Other goroutines' stacks not shown (single goroutine view)
- Switching goroutines (if supported) updates stack
- Main goroutine shows stack from main.main downward

### 6. Inspect Deep Call Stacks
**User Action:** Program pauses in deeply nested function call (50+ levels)
**Expected Behavior:**
- First 50 frames displayed (STACK_LIMIT = 50)
- Warning/indicator if stack truncated: "Showing 50 of N frames"
- Oldest frames (#50+) not accessible
- Performance remains responsive
- No crash or hang with deep stacks

### 7. View Stack After Crash/Panic
**User Action:** Program panics, debugger catches panic
**Expected Behavior:**
- Stack shows frames at panic point
- Topmost frame is where panic occurred
- Stack includes all frames up to panic site
- Can navigate frames to understand panic cause
- Error message or panic value shown in console

### 8. Refresh Stack After Step
**User Action:** Perform Step Over, Step Into, or Step Out
**Expected Behavior:**
- Stack updates immediately
- Step Over: Same stack depth, #0 frame advances
- Step Into: Stack grows (new #0 frame), previous #0 becomes #1
- Step Out: Stack shrinks (current #0 removed), previous #1 becomes #0
- Selected frame resets to #0 (current execution)
- Source view tracks new #0 location

## UI Elements

### Stack Trace Panel
**Location:** Bottom tabbed panel, "Stack" tab

**List View:**
- Vertical scrollable list
- Single-column layout
- Each row represents one frame

**Frame Row Format:**
```
#0 main.processRequest
   /path/to/file.go:42

#1 main.handleConnection
   /path/to/file.go:103

#2 main.startServer
   /path/to/file.go:25
```

**Visual Hierarchy:**
- Frame number: Small, gray, left-aligned
- Function name: Bold, larger font, primary text color
- File location: Smaller, gray, indented or second line
- Selected frame: Background highlight (light blue/gray)
- Current frame (#0): Possibly bold or distinct color

**Tooltips:**
- Hover over frame: Show full file path
- Hover over function: Show full signature (if available)

### Empty State
**When Shown:**
- Process not started
- Process exited
- Process running (not paused)

**Message:**
- "Stack trace unavailable - process is running"
- "Stack trace unavailable - process has exited"

### Error State
**When Shown:**
- Cannot retrieve stack (Delve error)
- Invalid goroutine state

**Message:**
- "Error retrieving stack trace: <reason>"
- Red text or error icon

## Stack Frame Information

### Frame Number
- Zero-indexed from top (#0 is current execution)
- Sequential numbering (#0, #1, #2, ... #49)
- No gaps in numbering

### Function Name
**Format:** `package.FunctionName`
**Examples:**
- `main.main` - main function
- `main.processRequest` - package function
- `(*Server).HandleRequest` - method on pointer receiver
- `github.com/user/pkg.Function` - external package function

**Special Cases:**
- Anonymous functions: `main.main.func1`, `main.main.func2`
- Closures: Same as anonymous functions
- Goroutine start: `runtime.goexit`
- Init functions: `main.init`

### Location Display
**Short Format:** `filename.go:42`
- Basename only (no directory path)
- Colon separator
- Line number precise

**Full Path (Tooltip):** `/full/path/to/file.go:42`
- Absolute path from filesystem
- Available on hover
- Used for source view navigation

## Event Interactions

### EventStackFrameChanged
**Emitted When:**
- User clicks stack frame
- Stack frame selection changes

**Data Payload:**
- Selected frame index (0-49)
- File path
- Line number
- Function name

**Subscribers:**
- Source view (navigate to frame location)
- Variables panel (show frame's locals)
- Disassembly view (show frame's assembly)

### EventStateChanged (Subscriber)
**Trigger:** Process state changes (continue, step, breakpoint hit)
**Stack Panel Response:**
- Fetch new stack trace from Delve
- Clear previous stack display
- Render new stack frames
- Reset selection to #0
- Scroll to top
- Update UI on main thread (glib.IdleAdd)

## Goroutine Context

### Current Goroutine
- Stack always shows frames for "current" goroutine
- Current goroutine: The one that triggered pause/breakpoint
- Goroutine ID tracked internally
- If no goroutine selected, first available goroutine auto-selected

### Switching Goroutines (Future)
- Not currently implemented
- Would require goroutine selector UI
- Stack would refresh when goroutine changed

## Performance

### Stack Retrieval
- Delve API call: ~10-50ms
- 50 frames: < 100ms total retrieval
- UI rendering: < 50ms for 50 frames
- Total update latency: < 200ms

### Large Stacks
- 50+ frames: Only first 50 fetched (STACK_LIMIT)
- No performance degradation
- Memory usage: ~5KB for 50 frames

### Frequent Updates
- Step operations: Stack updates every step
- No noticeable lag
- Event bus ensures thread-safe updates

## Error Handling

### Stack Unavailable
**Causes:**
- Process running (not paused)
- Process exited
- Delve API error
- Invalid goroutine

**Behavior:**
- Clear stack display
- Show error message in panel
- Log to console
- Don't crash or hang

### Invalid Frame Selection
**Causes:**
- Frame index out of range
- Frame no longer exists after step

**Behavior:**
- Reset selection to #0
- Log warning
- Don't emit EventStackFrameChanged
- Graceful degradation

### File Not Found
**Causes:**
- Frame references file not on disk
- Build artifacts moved

**Behavior:**
- Frame still shown in stack
- Source view shows error when navigating
- Frame remains selectable
- Location shown in stack even if file missing

## Accessibility
- Keyboard navigation: Up/Down arrows to select frames
- Enter key: Select frame (if not already default)
- Screen readers: Announce frame number, function, location
- High contrast: Selected frame clearly visible

## Known Limitations
- Stack limited to 50 frames (STACK_LIMIT constant)
- Single goroutine view (no multi-goroutine stack comparison)
- No inlined function frame detection (Delve limitation)
- No recursive call cycle detection/highlighting
- Cannot edit frame variables from stack panel
- Cannot restart from arbitrary frame (no "jump to frame")
