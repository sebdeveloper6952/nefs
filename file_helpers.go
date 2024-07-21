package nefs

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

func ReadFileToBase64(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("readFileToBase64: %w", err)
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

func DetectFileTypeAndExtension(bytes []byte) (string, string) {
	mime := mimetype.Detect(bytes)

	return mime.String(), mime.Extension()
}
