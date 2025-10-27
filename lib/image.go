package lib

import (
	"fmt"
	"image"
	"image/png"
	"log/slog"
	"os"

	"golang.org/x/image/draw"
)

func CreateThumbnail(originalFilename string, destinationFilename string, logger *slog.Logger) error {
	logger.Info(fmt.Sprintf("Creating thumbnail from %v with filename %v", originalFilename, destinationFilename))

	input, openErr := os.Open(originalFilename)
	if openErr != nil {
		logger.Error(fmt.Sprintf("Error opening %v: %v", originalFilename, openErr.Error()))
		return openErr
	}
	defer input.Close()

	output, createErr := os.Create(destinationFilename)
	if createErr != nil {
		logger.Error(fmt.Sprintf("Error creating %v: %v", destinationFilename, createErr.Error()))
		return createErr
	}
	defer output.Close()

	// Decode the image (from PNG to image.Image):
	src, decodeErr := png.Decode(input)
	if decodeErr != nil {
		logger.Error(fmt.Sprintf("Decode error: %v", decodeErr.Error()))
		return decodeErr
	}

	// Set the expected size that you want:
	dst := image.NewRGBA(image.Rect(0, 0, src.Bounds().Max.X/2, src.Bounds().Max.Y/2))

	// Resize:
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, src.Bounds(), draw.Over, nil)

	// Encode to `output`:
	encodeErr := png.Encode(output, dst)
	if encodeErr != nil {
		logger.Error(fmt.Sprintf("Encode error: %v", encodeErr.Error()))
		return encodeErr
	}

	return nil
}
