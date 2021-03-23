package main

import (
	"fmt"
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

func getFileIconLink(filename string) string {
	fileIcons := map[string]string{
		"pdf": "foo",
		"py":  "bar",
		"c":   "baz",
	}
	fmt.Println(fileIcons)
	return "testLink"
}

func getFileLink(filename string) string {
	return "testLink"
}
