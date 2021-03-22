package main

import (
	"path/filepath"
	"strings"
)

func getPathFromURLParam(par string) string {
	dirs := strings.Split(par, "_")
	path := ""
	for _, dirName := range dirs {
		path = filepath.Join(path, dirName)
	}
	return path
}
