# Feature: Disassembly Panel

## Overview
Low-level assembly code viewer showing machine instructions around the current
program counter (PC), useful for understanding compiled code behavior and
debugging optimized or complex code paths.

## User Capabilities

### 1. View Current Disassembly
**User Action:** Program pauses at any location (breakpoint, step, halt)
**Expected Behavior:**
- Disassembly panel automatically updates
- Shows ~1024 bytes of instructions around current PC
- Current instruction marked with `=>` prefix
- Address shown in hex format (`0x...`)
- Instruction mnemonics displayed
- Auto-scrolls to show current instruction
- Updates within 200ms of state change

**Display Format:**
```
   0x0000000000401000    MOVQ    $0x10, AX
=> 0x0000000000401007    CALLQ   0x401100
   0x000000000040100c    TESTL   AX, AX
   0x000000000040100e    JE      0x401020
   0x0000000000401010    MOVQ    8(SP), DX
```

### 2. View Instruction Details
**User Action:** Observe disassembly entries
**Expected Behavior:**
- Each line shows:
  - **Current indicator:** `=>` for current instruction, blank otherwise
  - **Address:** Hex memory address (e.g., `0x401007`)
  - **Instruction:** Assembly mnemonic and operands (e.g., `MOVQ $0x10, AX`)
- Architecture-specific assembly (matches target platform)
- Instructions shown in execution order
- Aligned columns for readability

### 3. Track Execution Through Assembly
**User Action:** Perform step operations (Step Into, Step Over)
**Expected Behavior:**
- `=>` marker moves to new instruction after each step
- Disassembly refreshes with new context
- If PC moves significantly, entire view updates
- Shows new instruction set around PC
- Useful for low-level debugging and understanding optimization

### 4. View Disassembly for Different Stack Frames
**User Action:** Select different stack frame in Stack panel
**Expected Behavior:**
- Disassembly updates to show instructions around selected frame's PC
- Current frame (top): Shows actual execution point
- Parent frames: Show call site instructions
- EventStackFrameChanged triggers update
- `=>` marker shows frame's instruction pointer

### 5. Auto-Update on State Change
**User Action:** Continue, step, or break
**Expected Behavior:**
- EventStateChanged triggers disassembly refresh
- New instructions fetched from Delve
- UI updates via glib.IdleAdd (thread-safe)
- Scroll position adjusts to show current instruction
- No manual refresh needed

### 6. View Around Breakpoints
**User Action:** Hit breakpoint
**Expected Behavior:**
- Disassembly shows instructions at breakpoint location
- Breakpoint instruction marked with `=>`
- Can see instructions before and after breakpoint
- Useful for understanding compiled breakpoint placement
- May show multiple NOPs or different instruction than source line

