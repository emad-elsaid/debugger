package main

import (
	"log"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	. "github.com/emad-elsaid/debugger/ui"
)

var (
	win  = app.NewWindow()
	ops  op.Ops
	tree W
)

func main() {
	win.Option(app.Title("Go Debugger"))

	cmdAndArgs := os.Args[1:]
	if len(cmdAndArgs) == 0 {
		log.Fatalln("debugger needs a command: `run`, `exec` or `test`")
	}

	cmd := cmdAndArgs[0]

	args := []string{"."}
	if len(cmd) > 1 {
		args = cmdAndArgs[1:]
	}

	bin := "debug"

	switch cmd {
	case "run":
	case "test":
	case "exec":
		bin = args[0]
		args = args[1:]
	default:
		log.Fatalln("invalid command: ", cmd)
	}

	debugger, err := NewDebugger(bin, args, true, cmd == "test", cmd == "exec")
	if err != nil {
		log.Fatalln(err.Error())
	}

	screen, err := NewDebugScreen(debugger)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tree = screen.Layout

	go RunWindowAndExit(debugger)
	go Refresher()
	app.Main()
}

func RunWindowAndExit(debugger *Debugger) {
	if err := EventLoop(debugger); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func EventLoop(debugger *Debugger) error {
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			debugger.Detach(true)
			return e.Err
		case system.FrameEvent:
			c := layout.NewContext(&ops, e)
			tree(c)
			e.Frame(c.Ops)
		}
	}

	return nil
}

func Refresher() {
	for {
		win.Invalidate()
		time.Sleep(time.Millisecond * 100)
	}
}
