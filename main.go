package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
)

func formatOutput(file os.FileInfo, prefix string, isLast bool) string {
	glyph := "├"
	if isLast {
		glyph = "└"
	}

	return prefix + glyph + "───" + file.Name() + "\n"
}

func scanDir(path string, getFiles bool, prefix string) (res string) {
	dir, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}

	if items, err := dir.Readdir(0); err == nil {
		sort.Slice(items, func(i, j int) bool { return items[i].Name() < items[j].Name() })

		for i := 0; i < len(items); i++ {
			currFile := items[i]
			if i == len(items)-1 {
				res += formatOutput(currFile, prefix, true)
				if items[i].IsDir() {
					res += scanDir(filepath.Join(path, currFile.Name()), getFiles, prefix+"\t")
				}
			} else {
				res += formatOutput(currFile, prefix, false)
				if items[i].IsDir() {
					res += scanDir(filepath.Join(path, currFile.Name()), getFiles, prefix+"│\t")
				}
			}
		}
	}

	return res
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	fmt.Println(scanDir(path, printFiles, ""))

	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
