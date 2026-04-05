# Test Examples

This directory contains sample Go programs used for testing the debugger.

## Directory Structure

- `basic/` - Simple program demonstrating various Go types (strings, ints, slices, maps, channels, etc.)

## Purpose

These examples are used by the Makefile targets:
- `make run` - Runs the debugger on the basic example
- `make quick-test` - Fast test with timeout using the basic example
- `make run-interactive` - Interactive debugging session

## Adding New Examples

To add a new test example:
1. Create a new subdirectory under `test_examples/`
2. Add a `go.mod` and `main.go` (or test files)
3. Update Makefile if needed to reference the new example
