package main

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func WatchAddRecursive(watcher *fsnotify.Watcher, p string) error {
	return filepath.WalkDir(p, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasPrefix(info.Name(), ".") {
			return fs.SkipDir
		}

		if info.IsDir() {
			return watcher.Add(path)
		}

		return err
	})
}
