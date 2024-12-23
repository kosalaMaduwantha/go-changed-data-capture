package domain

import (
	"cdc-file-processor/packages/fileOps"
	"cdc-file-processor/packages/hashOps"
	"cdc-file-processor/packages/jsonOps"
	"cdc-file-processor/adapters/rabitMqAdapter"
	"os"
	"time"
	"github.com/gofrs/flock"
	"github.com/sirupsen/logrus"

)

func Cdc_run() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.JSONFormatter{})
	queueService := rabitMqAdapter.NewRabbitMqAdapter()

	// Close the queue connection when the function ends
	defer func() {
        if err := queueService.(*rabitMqAdapter.RabbitMqAdapter).Close(); err != nil {
            logger.Error("Failed to close queue connection:", err)
        }
    }()

	var inputFilePath string = os.Getenv("INPUT_FILE_PATH")
	var configFilePath string = os.Getenv("CONFIG_FILE_PATH")

	fileLock := flock.New(inputFilePath)

	for {
		// Try to lock the file before reading it
		locked, err := fileLock.TryLock()
		if err != nil {
			logger.Error("Error while trying to lock the file")
			continue
		}
		if !locked {
			logger.Info("File is currently being written to, retrying...")
			time.Sleep(1 * time.Second)
			continue
		}

		newFileHash := hashOps.HashFile(inputFilePath)
		configFileContent := fileOps.ReadFileContent(configFilePath)
		configFileData := jsonOps.FileS{}

		if len(configFileContent) == 0 {
			jsonData := jsonOps.JsonSerializeFileHashData(
				newFileHash, configFilePath, 0, 0)
			fileOps.FileWriter(configFilePath, jsonData)
			fileLock.Unlock()
			continue
		}

		configFileData = jsonOps.JsonDeserializeFileHashData(configFileContent)
		if newFileHash != configFileData.FileHash {
			logger.Info("File is changed")

			lineCount := configFileData.LineCount
			fileContent, currentLineNo, lenPrevLine := fileOps.ReadFileContentLineByLine(
				inputFilePath,
				lineCount,
				configFileData.LenPrevLine)

			jsonData := jsonOps.JsonSerializeFileHashData(
				newFileHash,
				configFilePath,
				currentLineNo,
				lenPrevLine)

			fileOps.FileWriter(configFilePath, jsonData)
			if len(fileContent) == 0 {
				logger.Info("No new lines added")
				fileLock.Unlock()
				continue
			}
			for i := 0; i < len(fileContent); i++ {
				logger.Info(fileContent[i])
				queueService.SendMessage(fileContent[i], "cdc")
			}

		} else {
			logger.Info("File is not changed")
		}
		fileLock.Unlock()
		time.Sleep(5 * time.Second)
	}

	
}
