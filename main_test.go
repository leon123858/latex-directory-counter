package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountWords(t *testing.T) {
	// Create a temporary file
	tempFile, err := ioutil.TempFile("", "*.tex")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write sample LaTeX content to the file
	latexContent := `% This is a comment
    Hello, 世界! This is a test. \textbf{加油！}
    \begin{equation}
    1 + 2 = 3
    \item 測試
    \end{equation}`
	if _, err := tempFile.WriteString(latexContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	// Test countWords
	expectedChineseCount := 6
	expectedEnglishCount := 5
	chineseCount, englishCount, err := countWords(tempFile.Name())
	if err != nil {
		t.Errorf("countWords returned an error: %v", err)
	}
	if chineseCount != expectedChineseCount {
		t.Errorf("Expected %d Chinese characters, got %d", expectedChineseCount, chineseCount)
	}
	if englishCount != expectedEnglishCount {
		t.Errorf("Expected %d English words, got %d", expectedEnglishCount, englishCount)
	}
}

func TestProcessDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create sample .tex files
	fileContents := map[string]struct {
		chineseCount int
		englishCount int
	}{
		"file1.tex": {5, 3}, // This file should contain 5 Chinese characters and 3 English words
		"file2.tex": {0, 5}, // This file should contain 0 Chinese characters and 5 English words
	}
	for filename, counts := range fileContents {
		content := strings.Repeat("寶", counts.chineseCount) + " " + strings.Repeat("test ", counts.englishCount)
		fullPath := filepath.Join(tempDir, filename)
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write to %s: %v", fullPath, err)
		}
	}

	// Test processDirectory
	fileInfos, err := processDirectory(tempDir)
	if err != nil {
		t.Errorf("processDirectory returned an error: %v", err)
	}
	if len(fileInfos) != len(fileContents) {
		t.Errorf("Expected %d FileInfo structs, got %d", len(fileContents), len(fileInfos))
	}

	// Verify each FileInfo struct
	for _, fileInfo := range fileInfos {
		expectedCounts, exists := fileContents[filepath.Base(fileInfo.Path)]
		if !exists {
			t.Errorf("Unexpected file: %s", fileInfo.Path)
			continue
		}
		if fileInfo.ChineseCount != expectedCounts.chineseCount {
			t.Errorf("Expected %d Chinese characters in %s, got %d", expectedCounts.chineseCount, fileInfo.Path, fileInfo.ChineseCount)
		}
		if fileInfo.EnglishCount != expectedCounts.englishCount {
			t.Errorf("Expected %d English words in %s, got %d", expectedCounts.englishCount, fileInfo.Path, fileInfo.EnglishCount)
		}
	}
}
