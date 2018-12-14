package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

func filter(vs []os.FileInfo, f func(os.FileInfo) bool) []os.FileInfo {
	res := make([]os.FileInfo, 0)
	for _, v := range vs {
		if f(v) {
			res = append(res, v)
		}
	}

	return res
}

func formatOutput(file os.FileInfo, prefix string, isLast bool) string {
	name := file.Name()
	glyph := "├"
	if isLast {
		glyph = "└"
	}

	if !file.IsDir() {
		size := strconv.FormatInt(file.Size(), 10) + "b"
		if size == "0b" {
			size = "empty"
		}

		name += " (" + size + ")"
	}

	return prefix + glyph + "───" + name + "\n"
}

func scanDir(path string, printFiles bool, prefix string) (res string) {
	dir, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}

	if items, err := dir.Readdir(0); err == nil {
		sort.Slice(items, func(i, j int) bool { return items[i].Name() < items[j].Name() })
		if !printFiles {
			items = filter(items, func(item os.FileInfo) bool { return item.IsDir() })
		}

		for i := 0; i < len(items); i++ {
			currFile := items[i]
			if i == len(items)-1 {
				res += formatOutput(currFile, prefix, true)
				if currFile.IsDir() {
					res += scanDir(filepath.Join(path, currFile.Name()), printFiles, prefix+"\t")
				}
			} else {
				res += formatOutput(currFile, prefix, false)
				if currFile.IsDir() {
					res += scanDir(filepath.Join(path, currFile.Name()), printFiles, prefix+"│\t")
				}
			}
		}
	}

	return res
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	fmt.Fprint(out, scanDir(path, printFiles, ""))

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
