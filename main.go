package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	commentRegex = regexp.MustCompile(`%.*$`)
	envRegex     = regexp.MustCompile(`\\begin\{.*?}[\s\S]*?\\end\{.*?}`)
	chineseRegex = regexp.MustCompile(`\p{Han}`)
)

type FileInfo struct {
	Path      string
	WordCount int
}

func countChineseCharacters(filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("錯誤: %v\n", err)
			os.Exit(1)
		}
	}(file)

	count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = commentRegex.ReplaceAllString(line, "")
		line = envRegex.ReplaceAllString(line, "")
		count += len(chineseRegex.FindAllString(line, -1))
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func processDirectory(dirPath string) ([]FileInfo, error) {
	var fileInfos []FileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".tex") {
			count, err := countChineseCharacters(path)
			if err != nil {
				return err
			}
			fileInfos = append(fileInfos, FileInfo{Path: path, WordCount: count})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].WordCount > fileInfos[j].WordCount
	})

	return fileInfos, nil
}

func printTable(fileInfos []FileInfo) {
	totalCount := 0
	fmt.Println("+----------------------+-----------+")
	fmt.Println("| File                 | Word Count |")
	fmt.Println("+----------------------+-----------+")
	for _, fi := range fileInfos {
		filename := filepath.Base(fi.Path)
		if len(filename) > 20 {
			filename = filename[:17] + "..."
		}
		fmt.Printf("| %-20s | %9d |\n", filename, fi.WordCount)
		totalCount += fi.WordCount
	}
	fmt.Println("+----------------------+-----------+")
	fmt.Printf("| Total                | %9d |\n", totalCount)
	fmt.Println("+----------------------+-----------+")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("使用方法: latex-directory-counter <LaTeX目錄路徑>")
		os.Exit(1)
	}

	dirPath := os.Args[1]
	fileInfos, err := processDirectory(dirPath)
	if err != nil {
		fmt.Printf("錯誤: %v\n", err)
		os.Exit(1)
	}

	printTable(fileInfos)
}
