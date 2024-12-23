package main

import (
	"cdc-file-processor/packages/fileOps"
	"cdc-file-processor/packages/hashOps"
	"cdc-file-processor/packages/jsonOps"
	"time"

	"github.com/sirupsen/logrus"
)

func isFileChanged(inputFilePath string, configFilePath string) bool {
	newHashValue := hashOps.HashFile(inputFilePath)
	fileContent := fileOps.ReadFileContent(configFilePath)
	fileHash := jsonOps.JsonDeserializeFileHashData(fileContent).FileHash
	return newHashValue != fileHash
}

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	const inputFilePath string = "files/input.txt"
	const configFilePath string = "files/hashStore.json"
	for {
		if isFileChanged(inputFilePath, configFilePath) {
			logger.Info("File is changed")

			hashDataStore := jsonOps.JsonDeserializeFileHashData(
				fileOps.ReadFileContent(configFilePath))
			lineCount := hashDataStore.LineCount
			if lineCount == 0 {
				lineCount = 1
			}

			fileContent, currentLineNo := fileOps.ReadFileContentLineByLine(
				inputFilePath, lineCount)

			newFileHash := hashOps.HashFile(inputFilePath)
			jsonData := jsonOps.JsonSerializeFileHashData(
				newFileHash, configFilePath, currentLineNo)
			fileOps.FileWriter(configFilePath, jsonData)

			logger.WithFields(
				logrus.Fields{
					"content: ": fileContent,
				}).Info("File content is")
		} else {
			logger.Info("File is not changed")
		}
		time.Sleep(5 * time.Second)
	}
}
