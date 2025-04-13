//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const Root = "frontend/build"

func main() {
	fmt.Println("Generating embed.go...")

	var files []string
	err := filepath.Walk(Root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("frontend/embed.go")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Fprintln(f, "// Automatically generated. DO NOT EDIT.")
	fmt.Fprintln(f, "package frontend")
	fmt.Fprintln(f, "import \"embed\"")
	fmt.Fprint(f, "//go:embed")

	for _, file := range files {
		root := strings.Split(Root, "/")[1]
		relativePath := root + strings.TrimPrefix(file, Root)
		fmt.Fprintf(f, " %s", relativePath)
	}
	fmt.Fprintln(f)
	fmt.Fprintln(f, "var Assets embed.FS")
}
