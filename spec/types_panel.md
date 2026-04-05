# Feature: Types Panel

## Overview
Alphabetically sorted browser of all Go types from the debugged binary,
with search capability for quick type lookup.

## User Capabilities

### 1. View All Types
**User Action:** Click "Types" tab in bottom panel
**Expected Behavior:**
- List displays all type definitions from binary's debug information
- Types sorted alphabetically by name
- Each entry shows type name only (no additional details)
- Includes:
  - Primitive types (int, string, bool, etc.)
  - Struct types
  - Interface types
  - Custom type aliases
  - Pointer types
  - Channel types
  - Function types
  - Map types
  - Array/slice types
- Scrollable for potentially thousands of types

**Type Name Examples:**
- `string`
- `int`, `int32`, `uint64`
- `User`, `Server`, `Config`
- `*http.Request` (pointer types)
- `map[string]int`
- `[]byte`, `[10]string`
- `chan int`
- `interface{}`
- `github.com/user/pkg.CustomType`

### 2. Search Types
**User Action:** Type in search entry field
**Expected Behavior:**
- Type list filters in real-time
- Case-insensitive substring matching
- Matches anywhere in type name
- Matching types remain visible
- Non-matching types hidden
- Clear button (X) appears when text present
- Empty search shows all types

**Search Examples:**
- `"user"` → Shows User, UserData, userInternal
- `"int"` → Shows int, int32, interface{}, pointer types
- `"*http"` → Shows all pointer types from http package
- `"map"` → Shows all map types
- `"Server"` → Shows Server, ServerConfig, TCPServer

### 3. View Type (Read-Only)
**User Action:** Click on type in list
**Expected Behavior:**
- Type becomes selected (highlighted)
- NO navigation or detailed view (read-only panel)
- No definition shown
- No fields or methods displayed
- Simply marks selection visually

**Note:** Current implementation is read-only browser.
Future enhancement could navigate to type definition in source.

### 4. Clear Search
**User Action:** Click X button or clear text manually
**Expected Behavior:**
- Search field clears
- All types visible again
- Alphabetical order restored
- Selection cleared or maintained
- Cursor remains in search field

### 5. Keyboard Navigation
**User Action:** Use keyboard in types panel
**Expected Behavior:**
- Tab: Focus search field
- Arrow keys: Navigate type list
- Enter: Select type (no action currently)
- Escape: Clear search
- Type-ahead: First letter jumps to types starting with that letter

## UI Elements

### Types Panel
**Location:** Bottom tabbed panel, "Types" tab

**Components:**
1. **Search Field:** Text entry, placeholder "Search types..."
2. **Type List:** Scrollable single-column list
3. **Scrollbar:** Vertical scroll

### Search Field
- **Placeholder:** "Search types..."
- **Position:** Top of panel, full width
- **Clear Button:** X icon when text present
- **Auto-focus:** Optional keyboard shortcut

### List View
**Each Row:**
- Type name (single line)
- Monospace or regular font
- Standard list row height

**Visual States:**
- Normal: Default appearance
- Hover: Light highlight
- Selected: Darker highlight
- Focused: Focus ring for keyboard nav

### Empty State
**When Shown:**
- Process not started
- Binary has no type info (rare)
- Search matches nothing

**Messages:**
- Not started: "No types available"
- No matches: "No types match '<search term>'"

## Type Discovery

### Source
- Delve API: Type information from DWARF debug symbols
- All types encountered during compilation
- Includes standard library types used in program

### Included Types
**Built-in Types:**
- bool
- string
- int, int8, int16, int32, int64
- uint, uint8, uint16, uint32, uint64
- float32, float64
- complex64, complex128
- byte, rune

**User-Defined Types:**
- Struct definitions
- Interface definitions
- Type aliases (`type MyInt int`)
- Named types

**Composite Types:**
- Pointer types (`*User`)
- Slice types (`[]string`)
- Array types (`[10]int`)
- Map types (`map[string]interface{}`)
- Channel types (`chan int`, `<-chan bool`)
- Function types (`func(int) string`)

**Package Types:**
- Project types
- Standard library types (http.Request, io.Reader, etc.)
- Third-party types (if used in binary)

### Type Representation
- Full package path included for disambiguiation
- Examples:
  - `main.User`
  - `net/http.Request`
  - `github.com/user/project/models.Database`

## Search Algorithm

### Matching
- **Case-insensitive** substring match
- Matches anywhere in type name
- Includes package path in match
- OR logic: Single search term matches any part

### Performance
- Real-time (no debounce)
- 1,000 types: < 50ms
- 10,000 types: < 200ms
- 100,000 types: < 1 second (edge case)

## Event Interactions

### EventProcessRestarted (Subscriber)
**Trigger:** Process restarts
**Types Panel Response:**
- Clear type list
- Re-fetch types from new binary
- Rebuild list UI
- Clear search filter
- Re-sort alphabetically
- Update on main thread

### No Events Emitted
- Types panel is purely informational
- Clicking types doesn't navigate or emit events
- Read-only browsing only

## Performance

### Load Time
- 100 types: < 100ms
- 1,000 types: < 500ms
- 10,000 types: 1-2 seconds
- 100,000 types: 5-10 seconds (very large binaries)

### Scroll Performance
- Smooth for < 1,000 types
- May lag with 10,000+ types (GTK limitation)

### Memory Usage
- ~200 bytes per type entry
- 1,000 types: ~200KB
- 10,000 types: ~2MB
- Negligible impact

## Accessibility
- Keyboard navigation: Arrow keys, Page Up/Down
- Screen reader: Announces type names
- Focus indicators: Clear visual feedback
- Tab order: Search → List

## Error Handling

### No Types Available
**Trigger:** Binary without debug symbols
**Behavior:**
- Empty list
- Message: "No types found"
- Console logs warning
- Panel remains functional

### Delve API Error
**Trigger:** Cannot retrieve type information
**Behavior:**
- Empty list
- Error in console
- Error dialog may appear
- Retry by restarting

## Known Limitations
- **Read-only:** Cannot navigate to type definitions
- **No details:** No fields, methods, or structure shown
- **No filtering:** Cannot filter by category (struct, interface, etc.)
- **No sorting options:** Alphabetical only
- **No grouping:** All types in flat list (not grouped by package)
- **No type relationships:** Cannot see which types embed others
- **No usage info:** Cannot see where type is used
- **No method listing:** Cannot see type's methods
- **Large binaries:** May have slow load with 10k+ types
- **No hierarchy:** Cannot see type inheritance/embedding tree
- **No export status:** Cannot distinguish exported vs. unexported
- **No custom actions:** Clicking type does nothing useful currently

## Future Enhancement Ideas
(Not currently implemented)
- Navigate to type definition in source
- Show type details (fields, methods) in detail pane
- Group types by package
- Filter by type category (struct/interface/etc.)
- Show type size and alignment
- Jump to type usage locations
- Show type dependency graph
- Export type information
