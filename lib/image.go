package lib

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
)

const ThumbnailMaxWidth = 160

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

func GetThumbnailFilename(originalFilename string) string {
	extension := filepath.Ext(originalFilename)
	return fmt.Sprintf(
		"%s_thumbnail%s",
		strings.TrimSuffix(originalFilename, extension),
		extension,
	)
}

func CreateThumbnail(originalFullPath string, destFullPath string, mimeType string) error {
	input, openErr := os.Open(originalFullPath)
	if openErr != nil {
		return fmt.Errorf("error opening %v: %v", originalFullPath, openErr.Error())
	}
	defer input.Close()

	output, createErr := os.Create(destFullPath)
	if createErr != nil {
		return fmt.Errorf("error creating %v: %v", destFullPath, createErr.Error())
	}
	defer output.Close()

	originalImage, _, decodeErr := image.Decode(input)
	if decodeErr != nil {
		return fmt.Errorf("error decoding %v: %v", originalFullPath, decodeErr.Error())
	}

	newImage := resize.Resize(ThumbnailMaxWidth, 0, originalImage, resize.Lanczos3)

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
		return fmt.Errorf("encode error: %v", encodeErr.Error())
	}

	return nil
}
