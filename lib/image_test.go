package lib

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/gabriel-vasile/mimetype"
)

const ImagePath = "../testdata"

func createThumbnail(t *testing.T, originalFilename string) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	destinationFilename := GetThumbnailFilename(originalFilename)

	// If the thumbnail exists already, remove it
	_, existsErr := os.Stat(destinationFilename)
	if existsErr == nil {
		removeErr := os.Remove(destinationFilename)
		if removeErr != nil {
			t.Fatal(removeErr)
		}
	}

	mimeType, mimeTypeErr := mimetype.DetectFile(originalFilename)
	if mimeTypeErr != nil {
		t.Fatal(mimeTypeErr)
	}

	thumbnailErr := CreateThumbnail(originalFilename, destinationFilename, mimeType.String(), logger)
	if thumbnailErr != nil {
		t.Fatal(thumbnailErr)
	}

	_, err := os.Stat(destinationFilename)
	if err != nil {
		t.Fatal("failed to create destination thumbnail")
	}

	imageWidthHeight, _ := GetImageWidthAndHeight(destinationFilename, logger)
	if imageWidthHeight.Width != ThumbnailMaxWidth {
		t.Fatal("destination thumbnail width should be ", ThumbnailMaxWidth)
	}

	deleteErr := os.Remove(destinationFilename)
	if deleteErr != nil {
		t.Fatal(deleteErr)
	}
}

func TestCreateThumbnailJpg(t *testing.T) {
	originalFilename := fmt.Sprintf("%s/purple.jpg", ImagePath)
	createThumbnail(t, originalFilename)
}

func TestCreateThumbnailPng(t *testing.T) {
	originalFilename := fmt.Sprintf("%s/red.png", ImagePath)
	createThumbnail(t, originalFilename)
}

func TestCreateThumbnailGif(t *testing.T) {
	originalFilename := fmt.Sprintf("%s/samurai-doge.gif", ImagePath)
	createThumbnail(t, originalFilename)
}
