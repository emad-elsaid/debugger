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
		log.Fatalln("debugger needs a command: `run` or `test`")
	}

	cmd := cmdAndArgs[0]

	args := []string{"."}
	if len(cmd) > 1 {
		args = cmdAndArgs[1:]
	}

	debugger, err := NewDebugger("debug", args, true, cmd == "test")
	if err != nil {
		log.Fatalln(err.Error())
	}

	screen, err := NewDebugScreen(debugger)
	if err != nil {
		log.Fatalln(err.Error())
	}

	tree = screen.Layout

	go RunWindowAndExit()
	go Refresher()
	app.Main()
}

func RunWindowAndExit() {
	if err := EventLoop(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func EventLoop() error {
	for e := range win.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
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
