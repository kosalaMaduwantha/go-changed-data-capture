package jsonOps

import (
	"encoding/json"
	"fmt"
	"strings"
	"log/slog"
)

var logger = slog.Default()

type FileS struct {
	FileName string `json:"file_name,omitempty"`
	FileHash string `json:"file_hash,omitempty"`
	LineCount int `json:"line_count,omitempty"`
	LenPrevLine int `json:"len_prev_line,omitempty"`
}

func JsonSerializeFileHashData(fileHash string, fileName string, lineCount int, lenPrevLine int) []byte {
	fileStruct:= FileS{
		FileName: strings.TrimSpace(fileName),
		FileHash: strings.TrimSpace(fileHash),
		LineCount: lineCount,
		LenPrevLine: lenPrevLine,
	}
	// encoding to json data
	jsonData, err := json.Marshal(fileStruct)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while encoding the json data: %s", err))
		return nil
	}
	return jsonData
}

func JsonDeserializeFileHashData(jsonData []byte) FileS {
	var fileStruct FileS
	err := json.Unmarshal(jsonData, &fileStruct)
	if err != nil {
		logger.Error(fmt.Sprintf("Error while decoding the json data: %s", err))
		return FileS{}
	}
	return fileStruct
}