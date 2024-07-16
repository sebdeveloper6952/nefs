package main

import (
	"encoding/base64"
	"fmt"
	"os"
)

func readFileToBase64(filePath string) (string, error) {
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("readFileToBase64: %w", err)
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil
}
