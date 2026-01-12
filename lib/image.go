package lib

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log/slog"
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

type ImageWidthHeight struct {
	Height int
	Width  int
}

func GetImageWidthAndHeight(imagePath string, logger *slog.Logger) (ImageWidthHeight, error) {
	reader, openErr := os.Open(imagePath)
	if openErr != nil {
		return ImageWidthHeight{}, openErr
	}
	defer func() {
		if closeErr := reader.Close(); closeErr != nil {
			// Log the close error if needed, but don't override the main error
			logger.Error(fmt.Sprintf("Error closing image file: %v", closeErr.Error()))
		}
	}()
	m, _, err := image.Decode(reader)
	if err != nil {
		return ImageWidthHeight{}, openErr
	}
	bounds := m.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	return ImageWidthHeight{
		Width:  w,
		Height: h,
	}, nil
}

func GetThumbnailFilename(originalFilename string) string {
	extension := filepath.Ext(originalFilename)
	return fmt.Sprintf(
		"%s_thumbnail%s",
		strings.TrimSuffix(originalFilename, extension),
		extension,
	)
}

func CreateThumbnail(originalFullPath string, destFullPath string, mimeType string, logger *slog.Logger) error {
	input, openErr := os.Open(originalFullPath)
	if openErr != nil {
		return fmt.Errorf("error opening %v: %v", originalFullPath, openErr.Error())
	}
	defer func() {
		if closeErr := input.Close(); closeErr != nil {
			logger.Error(fmt.Sprintf("Error closing image file: %v", closeErr.Error()))
		}
	}()

	output, createErr := os.Create(destFullPath)
	if createErr != nil {
		return fmt.Errorf("error creating %v: %v", destFullPath, createErr.Error())
	}
	defer func() {
		if closeErr := output.Close(); closeErr != nil {
			logger.Error(fmt.Sprintf("Error closing thumbnail file: %v", closeErr.Error()))
		}
	}()

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
