# Feature: Window Layout

## Overview
Three-panel split interface design providing organized workspace for debugging with header
controls, main content area, and bottom panel system.

## User Interface Structure

### Overall Layout

#### Window Sections
The main window divided into three primary areas:
1. **Header Bar** (top) - Debugging controls
2. **Main Content Area** (middle) - Source view and sidebar
3. **Bottom Panel** (bottom) - Tabbed information panels

```
┌─────────────────────────────────────────────────┐
│           Header Bar (Controls)                 │
├──────────────────────────┬──────────────────────┤
│                          │                      │
│   Source View            │   Sidebar            │
│   (Center)               │   (Right)            │
│                          │                      │
│                          │                      │
├──────────────────────────┴──────────────────────┤
│                                                 │
│   Bottom Tabbed Panel                           │
│                                                 │
└─────────────────────────────────────────────────┘
```

### Header Bar Section

#### Location and Purpose
- **Position**: Top of window
- **Height**: Single row, compact
- **Purpose**: Primary debugging controls
- **Always Visible**: Fixed position, doesn't scroll

#### Controls and Elements
Includes:
- **Continue** button (▶) + F5 tooltip
- **Pause** button (⏸)
- **Restart** button (↻) + Ctrl+Shift+F5 tooltip
- **Step Over** button + F10 tooltip
- **Step Into** button + F11 tooltip
- **Step Out** button + Shift+F11 tooltip

#### Visual Design
- **Icons**: Standard debugging icons
- **Tooltips**: Show keyboard shortcuts
- **State**: Buttons enable/disable based on process state
- **Spacing**: Logical grouping of related controls
- **Style**: Follows GTK/Adwaita design guidelines

### Main Content Area

#### Horizontal Split Layout
The content area uses a horizontal paned divider:
- **Left/Center Side**: Source code view (expandable)
- **Right Side**: Sidebar with watches and breakpoints
- **Divider**: Draggable splitter between sides
- **Resizable**: User can adjust split ratio

#### Source View (Left/Center)

**Characteristics**:
- **Primary Focus**: Largest screen real estate
- **Expandable**: Can expand to nearly full width
- **Content**: Syntax-highlighted source code
- **Scrolling**: Vertical and horizontal scrolling
- **Features**: Line numbers, current line highlight, breakpoint highlights

**Default Size**:
- **Width**: Approximately 70-80% of content area
- **Adjustable**: User can resize via splitter

#### Sidebar (Right)

**Characteristics**:
- **Compact**: Narrower than source view
- **Fixed Widgets**: Contains specific UI elements
- **Vertical Layout**: Widgets stacked vertically
- **Scrolling**: May scroll if content exceeds height

**Contents**:
1. **Watch Expressions Widget** (top)
   - Entry field for adding watches
   - List of watch expressions and values
   - Delete buttons per watch

2. **Breakpoints Widget** (below watches)
   - Compact breakpoint list
   - Toggle switches for enable/disable
   - Delete buttons per breakpoint

**Default Size**:
- **Width**: Approximately 20-30% of content area. Minimum width so that content doesn't get hidden when sidebar is too narrow.
- **Adjustable**: User can resize via splitter

#### Pane Divider

**Features**:
- **Draggable**: Click and drag to resize
- **Visual Indicator**: Line or handle showing drag area
- **Constraints**: Minimum sizes prevent collapsing completely
- **State**: Position persists during session (may not persist between sessions)

### Bottom Panel Section

#### Tabbed Interface
- **Layout**: Horizontal tab bar with panel below
- **Tabs**: 10 different information panels
- **Non-Closable**: Tabs cannot be closed
- **Always Present**: All tabs always available

#### Tab List
In order (may vary):
1. **Console** - Event logging
2. **Stack** - Stack trace
3. **Types** - Type browser
4. **Functions** - Function list
5. **Packages** - Package list
6. **Sources** - File browser
7. **Disassembly** - Assembly code
8. **Breakpoints** - Breakpoints panel
9. **Open Files** - Recent files
10. **Memory** - Memory inspector

#### Panel Resizing

**Vertical Resizing**:
- **Divider**: Between content area and bottom panel
- **Draggable**: Click and drag vertically to resize
- **Collapse**: Can minimize bottom panel to just tab bar
- **Expand**: Can expand bottom panel to take significant height
- **Default**: Moderate height showing ~10-15 rows of content

**Horizontal**:
- **Full Width**: Bottom panel uses entire window width
- **No Horizontal Split**: No left/right division in bottom panel

### Window Properties

#### Window Size
- **Minimum**: Likely has minimum width/height constraints
- **Default**: Opens at reasonable size for debugging
- **Resizable**: User can resize window freely
- **Maximizable**: Can maximize to full screen
- **State**: Window size/position may persist (TBD)

#### Window Title
- **Content**: Likely shows application name or debugged binary name
- **Format**: "Debugger - <binary-name>" or similar
- **Updates**: May update to reflect current state

### Responsive Behavior

#### Panel Visibility
- **Small Windows**: Panels remain accessible but may require scrolling
- **Large Windows**: More content visible without scrolling
- **Minimum Sizes**: Panels don't collapse below usable size

#### Scrolling
- **Source View**: Scrolls independently
- **Sidebar**: May scroll if widgets exceed height
- **Bottom Panel**: Content scrolls within panel area
- **Independent**: Each area scrolls independently

### Visual Theming

#### GTK/Adwaita Integration
- **Native Look**: Uses GTK4 native widgets
- **Theme Support**: Follows system GTK theme
- **Light/Dark**: May support system light/dark theme
- **Consistent**: All panels use consistent styling

#### Colors
- **Standard**: GTK standard colors for UI elements
- **Highlights**: Custom colors for debugging features:
  - Yellow: Current execution line
  - Red: Breakpoint lines
- **Contrast**: Readable syntax highlighting

### Keyboard Navigation

#### Focus Management
- **Tab Order**: Logical tab order through UI elements
- **Focus Indicators**: Visual indication of focused element
- **Keyboard Control**: Navigate without mouse

#### Efficient Workflow
- **Primary Focus**: Source view largest for code focus
- **Quick Access**: Sidebar for frequently needed info
- **Detailed Info**: Bottom panels for deep inspection
- **Controls**: Always visible in header

#### Fixed Tab Set
- **No Add/Remove**: Cannot add or remove tabs
- **No Reorder**: Cannot reorder tabs
- **No Detach**: Cannot detach panels to separate windows
- **Fixed Layout**: Overall layout structure is fixed

