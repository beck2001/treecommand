package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
)

type SortByName []os.FileInfo

func (s SortByName) Len() int {
	return len(s)
}

func (s SortByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortByName) Less(i, j int) bool {
	return s[i].Name() < s[j].Name()
}

func makePrefix(previousPrefix string, isLastParent bool, level int) string {
	if level == 0 {
		return ""
	}
	if isLastParent {
		return previousPrefix + "\t"
	}
	return previousPrefix + "|\t"
}

func makeSymbol(isLast bool) string {
	if isLast {
		return "|____"
	}
	return "|----"
}

func makeDirName(prefix string, symbol string, fi os.FileInfo) string {
	return prefix + symbol + fi.Name()
}

func makeFileName(prefix string, symbol string, fi os.FileInfo, postPrefix string) string {
	if fi.Size() > 0 {
		return prefix + symbol + fi.Name() + " (" + strconv.Itoa(int(fi.Size())) + postPrefix + ")"
	}
	return prefix + symbol + fi.Name() + " (empty)"
}

func getDirOnly(fileInfos []os.FileInfo) []os.FileInfo {
	var result []os.FileInfo
	for _, fi := range fileInfos {
		if fi.IsDir() {
			result = append(result, fi)
		}
	}
	return result
}

func normalized(fileInfos []os.FileInfo, showFiles bool) []os.FileInfo {
	var result []os.FileInfo
	if !showFiles {
		result = getDirOnly(fileInfos)
	} else {
		result = fileInfos
	}
	sort.Sort(SortByName(result))
	return result
}

func currentDirTree(writer io.Writer, path string, showFiles bool, level int, isLastParent bool, pastPrefix string) error {
	dir, err := os.Open(path)
	if err != nil {
		return err
	}
	defer dir.Close()
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	filteredFileInfos := normalized(fileInfos, showFiles)
	isLast := false
	for index, fileInfo := range filteredFileInfos {
		if index == len(filteredFileInfos) - 1 {
			isLast = true
		}
		prefix := makePrefix(pastPrefix, isLastParent, level)
		symbol := makeSymbol(isLast)
		if fileInfo.IsDir() {
			fmt.Fprintln(writer, makeDirName(prefix, symbol, fileInfo))
			currentPath := path + string(os.PathSeparator) + fileInfo.Name()
			currentDirTree(writer, currentPath, showFiles, level + 1, isLast, prefix)
		} else if !fileInfo.IsDir() && showFiles {
			fmt.Fprintln(writer, makeFileName(prefix, symbol, fileInfo, "b"))
		} else {
			continue
		}
	}
	return nil
}

func dirTree(writer io.Writer, path string, showFiles bool) error {
	err := currentDirTree(writer, path, showFiles, 0, false, "")
	if err != nil {
		return err
	}
	return nil
}

func main() {
	output := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("invalid arguments provided")
	}
	path := os.Args[1]
	showFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(output, path, showFiles)
	if err != nil {
		fmt.Printf("error occurred: %v\n", err)
	}
}