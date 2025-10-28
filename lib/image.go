package lib

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log/slog"
	"os"

	"github.com/nfnt/resize"
)

func GetExtensionByMimeType(mimeType string) (string, error) {
	switch mimeType {
	case "image/png":
		return "png", nil
	case "image/jpeg":
		return "jpeg", nil
	case "image/gif":
		return "gif", nil
	case "image/webp":
		return "webp", nil
	default:
		return "", fmt.Errorf("unknown image mime type: %v", mimeType)
	}
}

func CreateThumbnail(originalFullPath string, destFullPath string, mimeType string, logger *slog.Logger) error {
	input, openErr := os.Open(originalFullPath)
	if openErr != nil {
		logger.Error(fmt.Sprintf("Error opening %v: %v", originalFullPath, openErr.Error()))
		return openErr
	}
	defer input.Close()

	output, createErr := os.Create(destFullPath)
	if createErr != nil {
		logger.Error(fmt.Sprintf("Error creating %v: %v", destFullPath, createErr.Error()))
		return createErr
	}
	defer output.Close()

	originalImage, _, decodeErr := image.Decode(input)
	if decodeErr != nil {
		logger.Error(fmt.Sprintf("Error decoding %v: %v", originalFullPath, decodeErr.Error()))
		return decodeErr
	}

	newImage := resize.Resize(160, 0, originalImage, resize.Lanczos3)

	var encodeErr error
	switch mimeType {
	case "image/png":
		encodeErr = png.Encode(output, newImage)
	case "image/jpeg":
		encodeErr = jpeg.Encode(output, newImage, nil)
	case "image/gif":
		encodeErr = gif.Encode(output, newImage, nil)
	}

	if encodeErr != nil {
		logger.Error(fmt.Sprintf("Encode error: %v", encodeErr.Error()))
		return encodeErr
	}

	return nil
}
