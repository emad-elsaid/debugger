# AGENTS.md - Developer Guide for Debugger Project

This guide is for agentic coding assistants working on this Go-based GTK debugger GUI.

## Project Overview

A graphical debugger for Go programs built on top of Delve, using GTK4/Adwaita.
- **Language**: Go 1.25.0
- **Main Dependencies**: Delve, GTK4 (gotk4), Adwaita, fsnotify
- **Platform**: Linux (tested on Arch Linux)
- **Purpose**: Chrome DevTools-like debugging experience for Go applications

### Project Structure
```
.
├── main.go                 # Application entry point
├── debugger.go             # Core debugger logic
├── window.go               # Main window implementation
├── breakpoints.go          # Breakpoints panel
├── console.go              # Console panel
├── disassembly.go          # Disassembly panel
├── functions.go            # Functions panel
├── memory.go               # Memory panel
├── openfiles.go            # Open files panel
├── packages.go             # Packages panel
├── sidebar.go              # Sidebar panel
├── sources.go              # Sources panel
├── source_view.go          # Source code viewer
├── stack.go                # Stack trace panel
├── types.go                # Types panel
├── events.go               # Event definitions
├── constants.go            # Constants and errors
├── structs/                # Shared types and event system
│   ├── events.go          # Event bus implementation
│   └── types.go           # Common type definitions
├── ui/                     # GTK Builder XML UI definitions
│   ├── window.ui
│   ├── breakpoints.ui
│   └── ...
├── test_examples/          # Sample programs for testing
│   └── basic/             # Simple test program
├── Makefile               # Build system with all targets
└── AGENTS.md              # This file
```

## Build & Test Commands

### Quick Start
```bash
# Show all available make targets
make help

# Build and test
make all
```

### Building
```bash
# Build the debugger binary
make build

# Build for testing with custom name
make build-test

# Clean built binaries
make clean
```

### Running
```bash
# Build and run debugger on test target
make run

# Run interactively (manual inspection, Ctrl+C to stop)
make run-interactive

# Quick test with 10s timeout and cleanup
make quick-test

# Debug a specific package directly
./debugger run ./cmd/cli/

# Debug with runtime arguments
./debugger run . -- --secure=false

# Debug tests
./debugger test <package-path>
```

### Testing

```bash
# Run all quality checks and tests (recommended)
make test
# This runs: deadcode check, empty functions check, empty files check, 
#            linter, go vet, and all tests

# Run individual checks
make check-deadcode           # Check for unused code
make check-empty-functions    # Check for empty function bodies
make check-empty-files        # Check for files with < 5 lines
make lint                     # Run golangci-lint
make vet                      # Run go vet

# Run tests only (skip quality checks)
make test-run

# Run tests with verbose output
make test-verbose

# Run with race detector
make test-race

# Run tests with coverage report
make test-cover

# Run a single test (use go test directly)
go test -run TestEventBus
go test -run TestEventBus/Subscribe

# Run a specific test file
go test -run TestApplicationExitsOnWindowClose

# Run tests in a specific package
go test ./structs/...
```

### Linting & Formatting
```bash
# Format code
make fmt

# Vet code
make vet

# Run golangci-lint (comprehensive linting)
make lint

# Run all quality checks and tests
make check
```

### Maintenance
```bash
# Kill leftover debugger processes
make clean-processes

# Download dependencies
make deps

# Tidy go.mod
make tidy

# Install to $GOPATH/bin
make install
```

## Feature Development Workflow

When developing a new feature, follow this iterative cycle until the feature is fully delivered:

1. **Analyze** - Determine what code needs to change to implement the feature
2. **Implement** - Make the necessary code changes
3. **Build & Launch** - Run the debugger application in separate steps:
   ```bash
   # Launch debugger in background (no timeout wrapper - causes issues)
   ./debugger run test_examples/basic &
   
   # Wait for UI to settle (5-10 seconds)
   sleep 7
   ```
4. **Interact** - Exercise the UI components relevant to the feature (if automation possible)
5. **Capture** - Take a screenshot of the application window:
   ```bash
   # Start ydotool daemon if not running
   ydotoold &
   sleep 2
   
   # Trigger screenshot via keyboard shortcut (Shift+Print)
   ydotool key 42:1 99:1 99:0 42:0
   sleep 2
   
   # Screenshot saved to ~/Pictures/Screenshots/
   SCREENSHOT=$(ls -t ~/Pictures/Screenshots/ | head -1)
   ```
6. **Cleanup** - Force kill ALL debugger processes (GTK apps do not exit cleanly):
   ```bash
   killall -9 debugger 2>/dev/null
   killall -9 debug 2>/dev/null
   killall ydotoold 2>/dev/null
   sleep 1
   ```
7. **Verify** - Analyze the screenshot to confirm the feature works as expected

**Repeat** this cycle until the feature is complete and verified.

**Important Notes**:
- Always use `killall -9 debugger` to terminate - GTK applications ignore normal signals
- The debugger process will NOT terminate on its own; manual cleanup is REQUIRED after each test
- Wait at least 7 seconds after launch before capturing screenshots to ensure UI is fully rendered

This workflow ensures visual confirmation of GTK UI behavior and catches issues early in the
development process.

## Code Style Guidelines

### File Organization
- Main package in root directory
- Subpackages: `structs/` (shared types), `ui/` (UI definitions)
- UI definitions stored as `.ui` files (GTK Builder XML)
- Test files use `_test.go` suffix

### Import Conventions
```go
// Standard library imports first
import (
    "fmt"
    "log"
    "os"
)

// External dependencies grouped
import (
    "github.com/diamondburned/gotk4/pkg/gtk/v4"
    "github.com/go-delve/delve/service/api"
)

// Dot imports for common internal packages (project convention)
import (
    . "github.com/emad-elsaid/debugger/structs"
    . "github.com/emad-elsaid/types"
)
```

**Note**: This project uses dot imports (`.`) for commonly used internal packages like `structs`.
While generally discouraged, it's an established pattern here.

### Naming Conventions
- **Types**: PascalCase (`Debugger`, `EventBus`, `BreakpointsPanel`)
- **Functions/Methods**: PascalCase for exported, camelCase for unexported
- **Variables**: camelCase (`listBox`, `stackFrame`)
- **Constants**: UPPER_SNAKE_CASE (`STACK_LIMIT`, `ANALYTICS_LIMIT`)
- **Interfaces**: Implied by context (no `-er` suffix requirement)

### Type Definitions
```go
// Struct definitions
type Debugger struct {
    *debugger.Debugger  // Embedded type
    State     DebuggerState
    EventBus  *UIEventBus
    
    // Group related fields with comments
    // Project props
    Path    string
    BinName string
}

// Constants with iota
type DebuggerEvent int

const (
    EventStateChanged DebuggerEvent = iota
    EventBreakpointsChanged
    EventFileChanged
)
```

### Error Handling
```go
// Define package-level errors
var (
    ErrProcessExited            = errors.New("process has exited")
    ErrCantInspectTargetProcess = errors.New("can't inspect target process")
)

// Return errors, don't panic
func NewDebugger(bin string, args []string) (*Debugger, error) {
    if err := d.InitDebugger(); err != nil {
        return nil, err
    }
    return &d, nil
}

// Silent error handling when appropriate
f, err := os.Open(path)
if err != nil {
    return  // Early return without logging
}
defer f.Close()
```

### Logging
```go
// Use log for fatal errors
log.Fatalln("debugger needs a command: `run` or `test`")

// Use slog for structured debugging
slog.Debug("main: starting")
slog.Info("App activated, creating window")
slog.Debug("main: command", "cmd", cmd)
```

### Concurrency & UI Updates
```go
// Always update UI on main thread using glib.IdleAdd
bp.debugger.EventBus.Subscribe(EventBreakpointsChanged, func(event DebuggerEvent, data EventData) {
    glib.IdleAdd(func() bool {
        bp.updateBreakpoints()
        return false  // Return false for one-time execution
    })
})

// Use sync.RWMutex for shared data
type EventBus struct {
    handlers map[DebuggerEvent][]EventHandler
    mu       sync.RWMutex
}
```

### Testing Patterns
```go
// Use subtests for organization
func TestEventBus(t *testing.T) {
    bus := NewEventBus()
    
    t.Run("Subscribe and Emit", func(t *testing.T) {
        // Test code
    })
    
    t.Run("Multiple Subscribers", func(t *testing.T) {
        bus := NewEventBus()  // Fresh instance per subtest
        // Test code
    })
}

// Use testify for assertions (already in dependencies)
// Skip interactive tests
t.Skip("Manual test - requires interactive window closing")
```

### Comments
```go
// Package-level documentation
// DebuggerEvent represents types of debugger state changes
type DebuggerEvent int

// Document exported functions
// NewEventBus creates a new event bus
func NewEventBus() *EventBus {

// Inline comments for clarity
_ = runImmediately // Will be used by UI

// Comments for important behavior
// Note: Handlers are called synchronously. UI code should use glib.IdleAdd internally.
```

## Common Patterns

### Event-Driven Architecture
The codebase uses an event bus pattern for state changes:
```go
// Subscribe to events
d.EventBus.Subscribe(EventStateChanged, func(event DebuggerEvent, data EventData) {
    // Handle event
})

// Emit events
d.EventBus.Emit(EventBreakpointsChanged, EventData{})
```

### GTK UI Construction
```go
// Load UI from file
builder := gtk.NewBuilderFromFile("ui/breakpoints.ui")
panel := builder.GetObject("breakpoints_panel").Cast().(*gtk.Box)

// HTML escaping for GTK labels
title := html.EscapeString(fmt.Sprintf("%s:%d", file, line))
```

### Resource Management
```go
// Always defer cleanup
f, err := os.Open(path)
if err != nil {
    return
}
defer f.Close()

// Lock/unlock pattern with Delve
tg, unlock := d.Debugger.LockTargetGroup()
defer unlock()
```

## Key Files
- `main.go` - Application entry point
- `debugger.go` - Core debugger logic
- `window.go` - Main window implementation
- `structs/events.go` - Event system definitions
- `*_panel.go` files - Individual UI panels (breakpoints, stack, console, etc.)
  - `breakpoints.go`, `console.go`, `disassembly.go`, `functions.go`
  - `memory.go`, `openfiles.go`, `packages.go`, `sidebar.go`
  - `sources.go`, `source_view.go`, `stack.go`, `types.go`

## Test Examples

The `test_examples/` directory contains sample Go programs for testing the debugger locally.

### Directory Structure
```
test_examples/
├── basic/              # Simple program with various Go types
│   ├── go.mod         # Required - defines the module
│   ├── main.go        # Test program with vars, slices, maps, channels
│   └── debug          # (gitignored - built by debugger)
├── .gitignore         # Ignores compiled binaries
└── README.md          # Documentation
```

### Using Test Examples
```bash
# Run debugger on the basic example (default target)
make run

# Run with 10s timeout for quick verification
make quick-test

# Run interactively for manual testing
make run-interactive

# Debug a specific example directly
./debugger run test_examples/basic
```

### Adding New Test Examples
1. Create a new subdirectory under `test_examples/`
2. Add a `go.mod` file (required for Delve's gobuild)
3. Add your test program (`main.go` or test files)
4. Update Makefile's `TEST_TARGET` variable if needed

**Important**: Each test example must have a `go.mod` file. Delve's `gobuild` requires a proper Go module to build and debug the target.

## Development Tips
- UI files are in `ui/` directory as GTK Builder XML
- Process lifecycle managed via GTK application signals
- Debugging target specified via command line: `run` or `test`
- File watching handled by fsnotify for auto-reload
- All Makefile targets use the local `test_examples/` directory
- Compiled binaries (`debug`, `debugger`) are gitignored
