# Feature: Watch Expressions

## Overview
User-defined expression monitoring that evaluates arbitrary Go expressions
at the current program state, displaying live values in the sidebar.

## User Capabilities

### 1. Add Watch Expression
**User Action:** Enter expression in watch input field and submit
**Expected Behavior:**
- Expression added to watch list
- Immediate evaluation at current program state
- Result displayed next to expression
- Empty input ignored (no watch created)
- Expression stored in watch list
- EventWatchesChanged emitted
- Input field cleared for next entry
- New watch appears at bottom of list

**Supported Expressions:**
- Variable names: `username`, `count`, `config`
- Struct fields: `user.Name`, `config.Port`
- Pointer dereferences: `*ptr`, `user.Address.City`
- Array/slice elements: `items[0]`, `matrix[1][2]`
- Map lookups: `cache["key"]`, `data[id]`
- Function calls: `len(items)`, `cap(buffer)`
- Arithmetic: `count * 2`, `x + y`
- Comparisons: `x > 10`, `name == "test"`
- Complex expressions: `user.Orders[0].Total * 1.1`

### 2. View Watch Value
**User Action:** Observe watch list in sidebar
**Expected Behavior:**
- Expression shown on left
- Current value shown on right
- Format: `expression = value`
- Value updates automatically on state changes
- Long values truncated to 50 characters with "..." suffix
- Different value types formatted appropriately:
  - Strings: `"hello"` (with quotes)
  - Numbers: `42`, `3.14`
  - Booleans: `true`, `false`
  - Pointers: `0x123456` or `<nil>`
  - Structs: `{field1: value1, field2: value2}` (truncated)
  - Arrays/Slices: `[1, 2, 3]` or `[1, 2, 3, ...]`
  - Maps: `map[key1:value1 key2:value2]` (truncated)

### 3. View Watch Evaluation Error
**User Action:** Watch expression becomes invalid or cannot be evaluated
**Expected Behavior:**
- Value shows: "Error evaluating"
- Red text color for error indication
- Expression remains in list
- Can be deleted or kept for later
- Console may log detailed error
- Tooltip shows full error message (if hover implemented)

**Common Error Causes:**
- Variable out of scope: `localVar` from different frame
- Undefined variable: `nonexistent`
- Type mismatch: `"string" + 123`
- Nil pointer dereference: `(*NilPointer).Field`
- Invalid syntax: `x +`

### 4. Delete Watch Expression
**User Action:** Click trash icon next to watch
**Expected Behavior:**
- Watch immediately removed from list
- No confirmation dialog
- EventWatchesChanged emitted
- List updates to remove entry
- Remaining watches remain unchanged
- Can re-add same expression later

### 5. Watch Updates on State Change
**User Action:** Step, continue, or any operation changing program state
**Expected Behavior:**
- All watches re-evaluated automatically
- Values update to reflect new state
- Updates happen within 200ms of state change
- No flicker during update
- Thread-safe update via glib.IdleAdd
- EventStateChanged triggers re-evaluation

### 6. Watch Scoped Variables
**User Action:** Add watch for local variable, then navigate stack frames
**Expected Behavior:**
- Watch evaluates in context of CURRENT stack frame
- Selecting frame #0: Shows local vars of current function
- Selecting frame #1: Shows local vars of caller function
- If variable doesn't exist in selected frame: "Error evaluating"
- Package-level variables always accessible regardless of frame
- Watch list doesn't change, but values update per frame context

### 7. Watch Complex Types
**User Action:** Add watch for struct, map, slice, or array
**Expected Behavior:**
- Value truncated to 50 characters
- Nested structures shown as much as fits
- Format prioritizes readability
- Full value NOT shown (use Variables panel for deep inspection)
- Examples:
  - `user = {Name: "Alice", Age: 30, ...}` (if more fields)
  - `items = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, ...]` (if > 50 chars)
  - `data = map[a:1 b:2 c:3 d:4 ...]` (if > 50 chars)

### 8. Watch Expressions Persist Across Steps
**User Action:** Add watch, then perform multiple step operations
**Expected Behavior:**
- Watch list remains intact
- Values update after each step
- Watches don't disappear
- Can track variable changes over time
- Useful for observing loop variables, accumulator changes, etc.

## UI Elements

### Sidebar Watch Widget
**Location:** Right sidebar, top section (above breakpoints)

**Components:**
1. **Header:** "Watches" label
2. **Input Field:** Text entry for new expressions
3. **Add Button:** Submit new watch (or Enter key)
4. **Watch List:** Scrollable vertical list

**Watch Entry Format:**
```
expression = value        [🗑️]
```
- Expression: Left-aligned, regular font
- Equals sign: Separator
- Value: Right side, monospace font
- Delete icon: Trash button, far right

**Empty State:**
- "No watches" message
- Instruction: "Add expressions to monitor their values"

### Input Field
- Placeholder text: "Enter expression to watch"
- Single-line text entry
- Enter key: Submit and add watch
- Escape key: Clear input (optional)
- Auto-focus after adding watch (for rapid entry)

### Value Display
**Normal Value:**
- Monospace font for precise alignment
- Black text color
- Truncated with "..." if > 50 characters

**Error Value:**
- "Error evaluating" text
- Red text color (#cc0000 or similar)
- No truncation (fixed message)

## Evaluation Configuration

### Load Config (Delve API)
When evaluating expressions, the following limits apply:
- **Follow pointers:** Enabled (dereferences pointers)
- **Max string length:** 10,000 characters
- **Max array values:** 100 elements
- **Max struct fields:** 100 fields
- **Max variable recursion:** 2 levels deep
- **Max map buckets:** 100 buckets

### Evaluation Context
- Evaluates in currently selected stack frame
- Uses current goroutine context
- Respects variable scoping rules
- Cannot evaluate if process is running
- Cannot evaluate if process has exited

## Event Interactions

### EventWatchesChanged
**Emitted When:**
- Watch expression added
- Watch expression deleted
- (Not emitted when values update, only when list changes)

**Data Payload:**
- None (subscribers refresh full watch list)

**Subscribers:**
- Sidebar watch widget (refresh display)

### EventStateChanged (Subscriber)
**Trigger:** Process state changes (step, continue, breakpoint)
**Watch Widget Response:**
- Re-evaluate all watch expressions
- Update displayed values
- Mark errors if evaluation fails
- Refresh UI on main thread
- Complete within 200ms

### EventStackFrameChanged (Subscriber)
**Trigger:** User selects different stack frame
**Watch Widget Response:**
- Re-evaluate all watches in new frame context
- Update values (may change or error if out of scope)
- UI update via glib.IdleAdd

## Performance

### Evaluation Time
- Simple variable: < 10ms
- Complex struct: < 50ms
- Map/slice iteration: < 100ms
- 10 watches: < 200ms total evaluation

### UI Responsiveness
- Value update: < 50ms after evaluation
- Add/delete: Immediate (< 20ms)
- State change update: < 300ms total (eval + render)

### Memory Usage
- ~1KB per watch expression
- 50 watches: ~50KB (negligible)

## Error Handling

### Invalid Expression Syntax
**Trigger:** User enters malformed expression
**Behavior:**
- Watch added to list anyway
- Value shows "Error evaluating"
- Console logs syntax error
- User can delete and retry

### Variable Out of Scope
**Trigger:** Local variable from different function
**Behavior:**
- Value shows "Error evaluating"
- Watch remains in list
- Becomes valid if navigating to correct frame

### Nil Pointer Dereference
**Trigger:** Expression like `(*nilPtr).Field`
**Behavior:**
- Value shows "Error evaluating"
- Doesn't crash debugger
- Safe error handling

### Process Not Paused
**Trigger:** Evaluation attempt while process running
**Behavior:**
- Value shows previous value or "Error evaluating"
- Re-evaluates when process pauses
- No crashes

## Accessibility
- Keyboard navigation: Tab to input, Enter to submit
- Arrow keys: Navigate watch list
- Delete key: Remove selected watch (optional)
- Screen reader: Announces expression and value

## Known Limitations
- Value truncated to 50 characters (no expand option)
- Cannot edit existing watch (must delete and re-add)
- No watch history or value change highlighting
- Cannot copy watch value (without full clipboard support)
- No conditional watches (no "break when value changes")
- No watch grouping or organization
- Maximum practical watches: ~50 (performance reasons)
- Cannot watch CPU registers or memory addresses directly
- Evaluation errors don't show detailed messages in UI (console only)
