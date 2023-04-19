package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gioui.org/widget"
	. "github.com/emad-elsaid/debugger/ui"
	. "github.com/emad-elsaid/types"
)

type OpenFilesPanel struct {
	List   List
	Filter widget.Editor
}

func NewOpenFilesPanel() OpenFilesPanel {
	return OpenFilesPanel{
		List:   NewVerticalList(),
		Filter: LineEditor(),
	}
}

func (o *OpenFilesPanel) Layout(d *Debugger, a *Analytics) W {
	filter := o.Filter.Text()
	files := o.Files(d.Target().Pid())
	filtered := Slice[string](files).Select(func(s string) bool { return strings.Contains(s, filter) })
	ele := func(c C, i int) D { return Inset1(Label(filtered[i]))(c) }
	openFiles := Slice[int](a.OpenFiles).Fetch(len(a.OpenFiles)-1, 0)

	return Columns(
		Flexed(1,
			Rows(
				RowSpacer1,
				Rigid(
					FormRow(
						Rigid(Label(fmt.Sprintf(" %d", len(filtered)))),
						ColSpacer1,
						Flexed(1, TextInput(&o.Filter, "Search open files...")),
					),
				),
				Rigid(HSpacer1),
				Flexed(1, ZebraList(&o.List, len(filtered), ele)),
			),
		),
		Rigid(
			Inset1(
				Panel(
					fmt.Sprintf("Open files (%d)", openFiles),
					Constraint(300, 100, Chart(a.OpenFiles, 100)),
				),
			),
		),
	)
}

func (o *OpenFilesPanel) Files(pid int) []string {
	fd := fmt.Sprintf("/proc/%d/fd", pid)
	fs, err := os.ReadDir(fd)
	if err != nil {
		return []string{}
	}

	paths := make([]string, 0, len(fs))
	for _, v := range fs {
		dest, _ := os.Readlink(path.Join(fd, v.Name()))
		paths = append(paths, dest)
	}

	return paths
}
