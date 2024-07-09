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
	commentRegex  = regexp.MustCompile(`%.*$`)
	envFirstRegex = regexp.MustCompile(`\\begin\{.*?}`)
	envEndRegex   = regexp.MustCompile(`\\end\{.*?}`)
	decoRegex     = regexp.MustCompile(`\\[a-zA-Z]+(\{.*?})`)
	syntaxRegex   = regexp.MustCompile(`\\[a-zA-Z]+`)
	chineseRegex  = regexp.MustCompile(`\p{Han}`)
	englishRegex  = regexp.MustCompile(`[a-zA-Z]+`)
)

func removeEnvCommands(text string) string {
	text = envFirstRegex.ReplaceAllString(text, "")
	text = envEndRegex.ReplaceAllString(text, "")
	return text
}

func removeDecoCommands(text string) string {
	return decoRegex.ReplaceAllStringFunc(text, func(match string) string {
		result := decoRegex.FindStringSubmatch(match)
		if len(result) > 1 {
			// 返回捕獲組的內容，去掉大括號
			return result[1][1 : len(result[1])-1]
		}
		return match
	})
}

type FileInfo struct {
	Path         string
	ChineseCount int
	EnglishCount int
}

func countWords(filePath string) (int, int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, 0, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("錯誤: %v\n", err)
			os.Exit(1)
		}
	}(file)

	chineseCount := 0
	englishCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 以下屏蔽順序很重要
		line = removeEnvCommands(line)
		line = commentRegex.ReplaceAllString(line, "")
		line = removeDecoCommands(line)
		line = syntaxRegex.ReplaceAllString(line, "")
		chineseCount += len(chineseRegex.FindAllString(line, -1))
		englishCount += len(englishRegex.FindAllString(line, -1))
	}

	if err := scanner.Err(); err != nil {
		return 0, 0, err
	}

	return chineseCount, englishCount, nil
}

func processDirectory(dirPath string) ([]FileInfo, error) {
	var fileInfos []FileInfo

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".tex") {
			chineseCount, englishCount, err := countWords(path)
			if err != nil {
				return err
			}
			fileInfos = append(fileInfos, FileInfo{Path: path, ChineseCount: chineseCount, EnglishCount: englishCount})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(fileInfos, func(i, j int) bool {
		return fileInfos[i].ChineseCount+fileInfos[i].EnglishCount > fileInfos[j].ChineseCount+fileInfos[j].EnglishCount
	})

	return fileInfos, nil
}

func printTable(fileInfos []FileInfo) {
	totalChineseCount := 0
	totalEnglishCount := 0
	fmt.Println("+----------------------+--------------+--------------+--------------+")
	fmt.Println("| File                 | Chinese Count| English Count| Total Count  |")
	fmt.Println("+----------------------+--------------+--------------+--------------+")
	for _, fi := range fileInfos {
		filename := filepath.Base(fi.Path)
		if len(filename) > 20 {
			filename = filename[:17] + "..."
		}
		fmt.Printf("| %-20s | %12d | %12d | %12d |\n", filename, fi.ChineseCount, fi.EnglishCount, fi.ChineseCount+fi.EnglishCount)
		totalChineseCount += fi.ChineseCount
		totalEnglishCount += fi.EnglishCount
	}
	fmt.Println("+----------------------+--------------+--------------+--------------+")
	fmt.Printf("| Total                | %12d | %12d | %12d |\n", totalChineseCount, totalEnglishCount, totalChineseCount+totalEnglishCount)
	fmt.Println("+----------------------+--------------+--------------+--------------+")
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
