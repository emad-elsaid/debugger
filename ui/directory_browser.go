package ui

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

var (
	IconFolder = Icon(icons.FileFolderOpen, SecondaryTextColor)
	IconFile   = Icon(icons.ActionDescription, SecondaryTextColor)
	IconUp     = Icon(icons.NavigationArrowUpward, Theme.TextColor)
)

type DirectoryBrowser struct {
	List        List
	Btns        map[string]*Clickable
	Up          Clickable
	MinItemSize int
}

func NewDirectoryBrowser() DirectoryBrowser {
	return DirectoryBrowser{
		List:        NewVerticalList(),
		Btns:        map[string]*Clickable{},
		Up:          Clickable{},
		MinItemSize: int(Theme.FontSize * 20),
	}
}

func (d *DirectoryBrowser) Layout(dist *string) W {
	stat, err := os.Stat(*dist)
	if errors.Is(err, os.ErrNotExist) {
		*dist, _ = os.UserHomeDir()
	} else if stat != nil && !stat.IsDir() {
		*dist = path.Dir(*dist)
	}

	files := []fs.DirEntry{}
	filepath.WalkDir(*dist, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == *dist {
			return nil
		}
		files = append(files, d)
		if d.IsDir() {
			return fs.SkipDir
		}

		return nil
	})

	ele := func(c C, i int) D {
		f := files[i]
		p := path.Join(*dist, f.Name())

		ico := IconFile
		if f.IsDir() {
			ico = IconFolder
		}

		doubleFont := FontSize(Theme.FontSize * 2)

		card := Border(
			Inset1(
				Columns(
					Rigid(doubleFont(ico)),
					ColSpacer1,
					Flexed(1, Label(f.Name())),
				),
			),
		)

		if !f.IsDir() {
			return Inset05(card)(c)
		}

		var btn *Clickable
		var ok bool
		if btn, ok = d.Btns[p]; !ok {
			btn = new(Clickable)
			d.Btns[p] = btn
		}

		if btn.Clicked() {
			*dist = p
		}

		if btn.Hovered() {
			card = Background(BorderColor, card)
		}

		return Inset05(func(c C) D { return btn.Layout(c, card) })(c)
	}

	return func(c C) D {
		if d.Up.Clicked() {
			*dist = path.Dir(*dist)
		}

		return Rows(
			Rigid(
				Background(ToolbarBgColor,
					Columns(
						Rigid(ToolbarButton(&d.Up, IconUp, "Up")),
						Flexed(1, Inset1(Bold(Label(*dist)))),
					),
				),
			),
			Flexed(1, Grid(&d.List, len(files), d.MinItemSize, ele)),
		)(c)
	}
}
