package main

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func WatchAddRecursive(watcher *fsnotify.Watcher, p string) error {
	return filepath.WalkDir(p, func(parent string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(d.Name(), ".") {
			return fs.SkipDir
		}

		if d.IsDir() {
			return watcher.Add(parent)
		}

		return err
	})
}
