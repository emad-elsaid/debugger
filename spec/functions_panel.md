# Feature: Functions Panel

## Overview
Searchable list of all functions from the debugged binary, enabling quick navigation
to function definitions regardless of current execution state.

## User Capabilities

### 1. View All Functions
**User Action:** Click "Functions" tab in bottom panel
**Expected Behavior:**
- List displays all functions from binary's debug information
- Each entry shows:
  - **Primary Text:** Function name (e.g., `main.processRequest`)
  - **Secondary Text:** Location (`filename.go:line`)
- Functions sorted alphabetically by name
- Includes:
  - Package functions
  - Methods (with receiver notation)
  - Anonymous functions/closures
  - Init functions
  - Main function
  - Standard library functions (if debug symbols present)
- Scrollable for large function counts (100s to 1000s)

**Function Name Formats:**
- Simple function: `main.helperFunction`
- Method: `(*Server).HandleRequest` (pointer receiver)
- Method: `(Server).GetConfig` (value receiver)
- Package function: `github.com/user/pkg.ExportedFunc`
- Anonymous: `main.main.func1`, `main.init.func2`

### 2. Search Functions
**User Action:** Type in search entry field at top of panel
**Expected Behavior:**
- Function list filters in real-time
- Case-insensitive matching
- Matches against BOTH:
  - Function name
  - File location
- Matching functions remain visible
- Non-matching functions hidden
- Clear button (X) appears when text present
- Empty search shows all functions

**Search Examples:**
- `"process"` → Shows ProcessData, processRequest, ProcessingHandler
- `"handler"` → Shows all Handler methods, handler.go functions
- `"main.go"` → Shows all functions defined in main.go
- `"(*Server)"` → Shows all methods on Server pointer receiver
- `"init"` → Shows init functions

### 3. Navigate to Function
**User Action:** Click function in list
**Expected Behavior:**
- EventFileChanged emitted with:
  - File path containing function
  - Line number of function definition (first line)
- Source view loads file and scrolls to function
- Function signature visible at top of viewport
- Function body visible below
- Sources panel remains visible
- Can immediately set breakpoint in function

### 4. Clear Search
**User Action:** Click X button in search field or clear text
**Expected Behavior:**
- Search field clears
- All functions become visible again
- List returns to full alphabetical order
- Scroll position may reset
- Cursor remains in search field

### 5. Keyboard Navigation
**User Action:** Use keyboard in functions panel
**Expected Behavior:**
- **Tab:** Focus search field
- **Arrow keys:** Navigate function list
- **Enter:** Open selected function in source view
- **Escape:** Clear search and focus list
- **Type-ahead:** Typing focuses search field automatically

### 6. Inspect Function Location
**User Action:** Hover over function entry (if tooltips implemented)
**Expected Behavior:**
- Tooltip shows full file path
- Shows complete function signature (if available)
- Useful for distinguishing functions with same name

### 7. Refresh After Restart
**User Action:** Restart debugger process
**Expected Behavior:**
- Functions panel clears
- Re-fetches function list from new binary
- New functions appear if added
- Deleted functions removed
- List re-sorts alphabetically
- Search filter clears

