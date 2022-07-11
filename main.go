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
	tree = NewOpenScreen().Layout

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
