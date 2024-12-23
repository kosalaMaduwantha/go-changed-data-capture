package fileOps

import (
    "fmt"
    "io"
    "os"
	"bufio"
	"log/slog"
)

var logger = slog.Default()

func ReadFile(filePath string) *os.File {
	file_obj, err := os.Open(filePath)

	if err != nil {
		errorString := err.Error()
		logger.Error(fmt.Sprintf("File reading error: %s", errorString))
		return nil
	}
	return file_obj
}

func FileWriter(filePath string, data []byte) {
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while creating the file: %s", err))
		return
	}
}

func ReadFileContent(filePath string) []byte {
	file := ReadFile(filePath)
	if file == nil {
		logger.Error("Error while reading the file")
		return nil
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while reading the file content: %s", err))
		return nil
	}
	return fileContent
}

func ReadFileContentLineByLine(filePath string, startLine int) ([]string, int) {
	file := ReadFile(filePath)
	if file == nil {
		logger.Error("Error while reading the file")
		return nil, 0
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var fileContent []string
	lineNumber := 1
	for scanner.Scan() {
		if lineNumber > startLine {
			fileContent = append(fileContent, scanner.Text())
		}
		lineNumber++
	}
	currentLineNo := lineNumber -1
	return fileContent, currentLineNo
}
