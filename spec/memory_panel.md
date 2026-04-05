# Feature: Memory Panel

## Overview
Raw memory inspector allowing manual examination of process memory at specific
addresses, displaying hex dump and ASCII representation.

## User Capabilities

### 1. Inspect Memory at Address
**User Action:** Enter memory address in hex format and click "Go" button
**Expected Behavior:**
- 256 bytes of memory read from specified address
- Memory displayed as hex dump with ASCII sidebar
- Format: 16 bytes per row
- Address, hex bytes, and ASCII shown
- Read completes within 100-200ms
- Error shown if address invalid or inaccessible

**Hex Dump Format:**
```
0x00007FF7A8B01000  48 8B 44 24 08 48 89 C3  48 8B 44 24 10 48 89 C6  |H.D$.H..H.D$.H..|
0x00007FF7A8B01010  E8 2C 01 00 00 48 89 D8  5B 5E C3 90 90 90 90 90  |.,...H..[^......|
0x00007FF7A8B01020  55 48 89 E5 48 83 EC 20  48 89 7D F8 48 89 75 F0  |UH..H.. H.}.H.u.|
```

### 2. View Hex Bytes
**User Action:** Observe memory display
**Expected Behavior:**
- Each byte shown as 2-digit hex: `00` to `FF`
- Bytes grouped in pairs for readability: `48 8B`, `44 24`
- 16 bytes per row (standard hex dump format)
- Space between byte 8 and 9 for visual grouping
- Non-printable bytes shown as hex

### 3. View ASCII Representation
**User Action:** Observe right side of memory display
**Expected Behavior:**
- Each byte converted to ASCII if printable
- Printable range: ASCII 32-126 (space to tilde)
- Non-printable bytes: Shown as `.` (period)
- Pipe delimiters: `|..text..|` around ASCII section
- Null bytes (`0x00`): Shown as `.`
- Control characters: Shown as `.`

**Examples:**
- `48 65 6C 6C 6F` → `|Hello|`
- `00 48 69 00 FF` → `|.Hi..|`
- `41 42 43 0A 44` → `|ABC.D|` (0A = newline)

### 4. Navigate Memory Regions
**User Action:** Enter different addresses sequentially
**Expected Behavior:**
- Previous memory dump cleared
- New memory dump displayed
- No history (each read is independent)
- Can inspect discontinuous memory regions
- Useful for examining:
  - Variable addresses
  - Pointer targets
  - Buffer contents
  - Stack frames
  - Heap allocations

### 5. Handle Invalid Addresses
**User Action:** Enter invalid or inaccessible address
**Expected Behavior:**
- Error message displayed below input field or in console
- Previous dump may remain visible or cleared
- Common errors:
  - "Invalid address format"
  - "Cannot read memory at 0x..."
  - "Access denied"
  - "Address out of range"
- User can retry with corrected address

### 6. Enter Address Formats
**User Action:** Type address in input field
**Expected Behavior:**
- Accepts hex format: `0x123456` or `123456`
- Leading `0x` optional but recommended
- Case-insensitive: `0xABCD` or `0xabcd`
- No spaces or separators allowed
- Full 64-bit addresses supported: `0x00007FF7A8B01000`
- Validation on submit (not real-time)

### 7. Read Fixed Size
**User Action:** Click "Go" with any address
**Expected Behavior:**
- Always reads exactly 256 bytes
- Cannot adjust read size (fixed constant)
- If less than 256 bytes available: Error or partial read
- If more needed: Must make multiple reads at different addresses
- 256 bytes = 16 rows of 16 bytes each
