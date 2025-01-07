package images

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
)

type LocalImageLoader struct{}

func (l *LocalImageLoader) LoadImage(filePath string, sessionId string, imageHash string) (string, error) {
	indexJson := fmt.Sprintf("%s/index.json", filePath)
	indexJsonContent, err := os.ReadFile(indexJson)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	// Parse the JSON content into a map
	var index map[string]int
	err = json.Unmarshal(indexJsonContent, &index)
	if err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}
	imageIdx, found := index[imageHash]
	if !found {
		return "", fmt.Errorf("image not found")
	}
	// Load the image file
	imageFile := fmt.Sprintf("%s/%s.png", filePath, fmt.Sprint(imageIdx))
	imageContent, err := os.ReadFile(imageFile)
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}

	// Convert image content to base64
	encodedImage := base64.StdEncoding.EncodeToString(imageContent)
	return encodedImage, nil
}
