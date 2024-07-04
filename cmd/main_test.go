package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCountChineseCharacters(t *testing.T) {
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

	// Test countChineseCharacters
	expectedCount := 6
	count, err := countChineseCharacters(tempFile.Name())
	if err != nil {
		t.Errorf("countChineseCharacters returned an error: %v", err)
	}
	if count != expectedCount {
		t.Errorf("Expected %d Chinese characters, got %d", expectedCount, count)
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
	fileContents := map[string]int{
		"file1.tex": 5, // This file should contain 5 Chinese characters
		"file2.tex": 0, // This file should contain 0 Chinese characters
	}
	for filename, charCount := range fileContents {
		content := strings.Repeat("寶", charCount) + " This is a test."
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
		expectedCount, exists := fileContents[filepath.Base(fileInfo.Path)]
		if !exists {
			t.Errorf("Unexpected file: %s", fileInfo.Path)
			continue
		}
		if fileInfo.WordCount != expectedCount {
			t.Errorf("Expected %d Chinese characters in %s, got %d", expectedCount, fileInfo.Path, fileInfo.WordCount)
		}
	}
}
