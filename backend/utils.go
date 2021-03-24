package main

import (
	"bytes"
	"path/filepath"
	"strings"
)

const baseStaticURL = "http://localhost:8080/v1/MyCloud/static/"
const baseFilesURL = "http://localhost:8080/v1/MyCloud/files/"
const clientsBaseDir = "/home/orestis/MyCloud"

func getPathFromURLParam(par string) string {
	dirs := strings.Split(par, "_")
	path := ""
	for _, dirName := range dirs {
		path = filepath.Join(path, dirName)
	}
	return path
}

func getFileIconLink(filename string, isDir bool) string {
	fileIcons := map[string]string{
		".pdf": "pdf.jpg",
		".py":  "python.jpg",
	}
	const defaultFileIcon = "file.jpg"
	const directoryIcon = "folder.jpg"
	var b bytes.Buffer
	if isDir {
		b.WriteString(baseStaticURL)
		b.WriteString(directoryIcon)
		return b.String()
	}
	extension := filepath.Ext(filename)
	icon, hasKey := fileIcons[extension]
	if !hasKey {
		b.WriteString(baseStaticURL)
		b.WriteString(defaultFileIcon)
		return b.String()
	}
	b.WriteString(baseStaticURL)
	b.WriteString(icon)
	return b.String()
}

func getFileLink(filename string) string {
	var b bytes.Buffer

	b.WriteString(baseFilesURL)
	b.WriteString(filename) // append

	//fmt.Println(b.String()) // abcdef
	return b.String()
}
