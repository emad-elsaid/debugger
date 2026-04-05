package structs

// FunctionInfo represents metadata about a function in the debugged binary
type FunctionInfo struct {
	Name string
	File string
	Line int
}

// DebuggerState represents the runtime state of the debugger
type DebuggerState uint64
