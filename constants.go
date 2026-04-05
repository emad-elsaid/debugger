package main

import (
	"errors"

	"github.com/go-delve/delve/pkg/proc"
)

// Constants
const STACK_LIMIT = 50

// Errors
var ErrProcessExited = errors.New("process has exited")

// ProcLoadConfig defines how much data to load when inspecting variables
var ProcLoadConfig = proc.LoadConfig{
	FollowPointers:     true,
	MaxStringLen:       10000,
	MaxArrayValues:     100,
	MaxStructFields:    100,
	MaxMapBuckets:      100,
	MaxVariableRecurse: 2,
}
