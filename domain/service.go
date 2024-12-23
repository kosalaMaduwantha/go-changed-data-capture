package domain

import (
	"cdc-file-processor/packages/fileOps"
	"cdc-file-processor/packages/hashOps"
	"cdc-file-processor/packages/jsonOps"
	"time"

	"github.com/sirupsen/logrus"
)


func Cdc_run() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})

	const inputFilePath string = "../files/input.txt"
	const configFilePath string = "../files/hashStore.json"

	for {
		newFileHash := hashOps.HashFile(inputFilePath)
		fileContent := fileOps.ReadFileContent(configFilePath)
		hashData := jsonOps.JsonDeserializeFileHashData(fileContent)

		if newFileHash != hashData.FileHash {
			logger.Info("File is changed")

			lineCount := hashData.LineCount
			if lineCount == 0 {
				lineCount = 1
			}

			fileContent, currentLineNo := fileOps.ReadFileContentLineByLine(
				inputFilePath, lineCount)

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
