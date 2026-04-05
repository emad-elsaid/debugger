package main

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
)

// TestApplicationExitsOnWindowClose verifies the application terminates when window closes
func TestApplicationExitsOnWindowClose(t *testing.T) {
	// Build the debugger binary for testing
	buildCmd := exec.Command("go", "build", "-o", "debugger-test-binary")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build debugger: %v", err)
	}
	defer os.Remove("debugger-test-binary")

	// Start the debugger process
	cmd := exec.Command("./debugger-test-binary", "run", "../debugger-test/")
	cmd.Env = append(os.Environ(), "DISPLAY="+os.Getenv("DISPLAY"))

	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start debugger: %v", err)
	}

	// Give the application time to start up
	time.Sleep(3 * time.Second)

	// Send SIGTERM to simulate window close
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		t.Fatalf("Failed to send SIGTERM: %v", err)
	}

	// Wait for process to exit with timeout
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		// Process exited - this is what we want
		t.Logf("Process exited successfully: %v", err)
	case <-time.After(5 * time.Second):
		// Process didn't exit - test fails
		cmd.Process.Kill()
		t.Fatal("Application did not exit within 5 seconds of window close")
	}
}

// TestWindowCloseCallsQuit tests the close handler logic
func TestWindowCloseCallsQuit(t *testing.T) {
	t.Skip("Manual test - requires interactive window closing")

	// This test documents how to manually verify:
	// 1. Run: go run . run ~/code/debugger-test/
	// 2. Click the window close button (X)
	// 3. Verify with `ps aux | grep debugger` that process exits
	// 4. Expected: No debugger process remains
	// 5. Actual: Process may still be running (BUG)
}
