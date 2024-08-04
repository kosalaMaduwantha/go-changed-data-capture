package hashOps

import(
	"fmt"
	"io"
	"crypto/sha256"
	"encoding/hex"
	"cdc-file-processor/packages/fileOps"
	"log/slog"
)

var logger = slog.Default()

func HashFile(filePath string) string {
	fileObj := fileOps.ReadFile(filePath)
	hasher := sha256.New()

	if fileObj == nil {
		logger.Error("Error while reading the file")
		logger.Error("File is empty")
		return ""
	}

	if _, err := io.Copy(hasher, fileObj); err != nil {
		logger.Error(fmt.Sprintf("Error while copying the file content to the hasher: %s", err))
		return ""
	}
	defer fileObj.Close()
	hashInHex := hex.EncodeToString(hasher.Sum(nil))
	return hashInHex
}