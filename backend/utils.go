package main

import (
	"bytes"
	"math"
	"os"
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

//Return the size of path Directory in Gb
func getDirSize(path string) (float64, error) {
	var size int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	sizeGB := float64(size) / 1024.0 / 1024.0 / 1024.0
	return toFixed(sizeGB, 3), err
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
