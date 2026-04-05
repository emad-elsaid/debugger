package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
)

var dbg *Debugger

func main() {
	log.Println("main: starting (before slog setup)")

	// Set log level to DEBUG
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))

	log.Println("main: slog setup complete")
	slog.Debug("main: starting")

	cmdAndArgs := os.Args[1:]
	if len(cmdAndArgs) == 0 {
		log.Fatalln("debugger needs a command: `run` or `test`")
	}

	cmd := cmdAndArgs[0]
	slog.Debug("main: command", "cmd", cmd)

	args := []string{"."}
	if len(cmdAndArgs) > 1 {
		args = cmdAndArgs[1:]
	}
	slog.Debug("main: args", "args", args)

	app := adw.NewApplication("com.github.emad-elsaid.debugger", gio.ApplicationFlagsNone)
	slog.Debug("main: app created")

	var err error
	slog.Debug("main: creating debugger")
	dbg, err = NewDebugger("debug", args, true, cmd == "test")
	if err != nil {
		log.Fatalln(err.Error())
	}
	slog.Debug("main: debugger created")

	app.ConnectActivate(func() {
		slog.Info("App activated, creating window")
		win := NewDebugWindow(app, dbg)
		win.Present()
		slog.Info("Window presented")
	})

	// Handle shutdown to clean up debugger
	app.ConnectShutdown(func() {
		slog.Info("App shutdown signal received")
		if dbg.Debugger != nil {
			// Detach in background - don't block shutdown
			go func() {
				slog.Info("Detaching debugger")
				dbg.Detach(true)
				slog.Info("Debugger detached")
			}()
		}
		slog.Info("Shutdown complete")
	})

	slog.Info("Starting app.Run()")
	// Only pass application name to GTK, not debugger arguments
	code := app.Run([]string{os.Args[0]})
	slog.Info("app.Run() exited", "code", code)
	os.Exit(code)
}
