# Feature: Process Control

## Overview
Core debugging process lifecycle management including execution control, pausing,
restarting, and stepping through code.

## User Capabilities

### 1. Continue Execution
**User Action:** Click "Continue" button or press F5
**Expected Behavior:**
- Resume program execution from current pause point
- Button becomes disabled while running
- Console logs: "[HH:MM:SS] Process continued"
- Program runs until next breakpoint, error, or completion
- UI updates when process halts at breakpoint
- Current line indicator moves to new location
- Current line is highlighted with green color

### 2. Pause/Halt Execution
**User Action:** Click "Pause" button while program is running
**Expected Behavior:**
- Immediately suspend program execution
- UI updates to show paused state
- Current execution line highlighted in source view
- Stack trace updates to show current call stack
- All panels refresh with current state
- Control buttons re-enable for stepping operations

### 3. Restart Process
**User Action:** Click "Restart" button
**Expected Behavior:**
- Console logs: "[HH:MM:SS] Restarting process..."
- Terminates current debugging session
- Rebuilds binary with latest source changes
- Preserves all breakpoints
- Starts fresh debugging session
- Remains in the current file + line
- Console logs: "[HH:MM:SS] ✓ Process restarted"
- Updates all panels with new process state

### 4. Step Over
**User Action:** Click "Step Over" button
**Expected Behavior:**
- Execute next line of code without entering function calls
- Current line indicator advances to next line
- If line contains function call, entire function executes
- Stack trace remains at current frame level
- All variable values update to reflect new state
- Disassembly view updates to new instruction

### 5. Step Into
**User Action:** Click "Step Into" button
**Expected Behavior:**
- Execute next line and enter function calls
- If current line calls a function, debugger enters that function
- Current line moves to first line of called function
- Stack trace adds new frame for entered function
- Source view switches to file containing function definition
- Variables panel shows local variables of new frame

### 6. Step Out
**User Action:** Click "Step Out" button
**Expected Behavior:**
- Execute until current function returns
- Current line moves to caller's next line
- Stack trace removes current frame
- Source view returns to calling file
- Variables panel shows caller's local variables
- If at top-level, behaves like Continue

## UI Elements

### Header Bar Buttons
- **Continue** (Play icon)
- **Pause** (Pause icon)
- **Restart** (Refresh icon)
- **Step Over** (Step over icon)
- **Step Into** (Step into icon)
- **Step Out** (Step out icon)

### Button State Management
- Buttons disabled when process is running (except Pause)
- Buttons enabled when process is paused
- Buttons disabled when process has exited
- Visual feedback through sensitivity changes

## State Transitions

```
[Initial] -> Continue -> [Running] -> Breakpoint Hit -> [Paused]
[Paused] -> Step Over -> [Paused at next line]
[Paused] -> Step Into -> [Paused in function]
[Paused] -> Step Out -> [Paused in caller]
[Any] -> Restart -> [Paused at main.main]
[Running/Paused] -> Process Exit -> [Terminated]
```

## Event Emissions
- `EventStateChanged` - After every process state change
- `EventProcessRestarted` - After successful restart
- `EventProcessExited` - When program terminates
- `EventError` - On operation failures

## Error Handling

### Process Exited
**Trigger:** Any operation when process has terminated
**Behavior:**
- Show error dialog: "Process has exited"
- Console logs error message
- Buttons remain disabled
- User must restart to continue debugging

### Operation Failures
**Trigger:** Delve API errors (network, target issues)
**Behavior:**
- Show error dialog with failure reason
- Console logs: "[HH:MM:SS] ERROR: <reason>"
- Process state remains unchanged
- User can retry operation
