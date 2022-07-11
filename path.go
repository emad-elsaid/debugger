package main

import (
	"os"
	"path"
)

func ProjectDir(p string) string {
	files := []string{
		"go.mod",
		".git",
	}

	for cp := p; cp != path.Dir(cp); cp = path.Dir(cp) {
		for _, f := range files {
			indicator := path.Join(cp, f)
			if _, err := os.Stat(indicator); err == nil {
				return cp
			}
		}
	}

	return p
}
