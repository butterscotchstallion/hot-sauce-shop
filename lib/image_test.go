package lib

import (
	"fmt"
	"image"
	"log"
	"os"
	"testing"

	"github.com/gabriel-vasile/mimetype"
)

const ImagePath = "../testdata"

func getImageWidthAndHeight(imagePath string) (width int, height int) {
	reader, openErr := os.Open(imagePath)
	if openErr != nil {
		log.Fatal(openErr)
	}
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}
	bounds := m.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	return w, h
}

func createThumbnail(t *testing.T, originalFilename string) {
	destinationFilename := GetThumbnailFilename(originalFilename)

	// If thumbnail exists already, remove it
	_, existsErr := os.Stat(destinationFilename)
	if existsErr == nil {
		removeErr := os.Remove(destinationFilename)
		if removeErr != nil {
			t.Error(removeErr)
		}
	}

	mimeType, mimeTypeErr := mimetype.DetectFile(originalFilename)
	if mimeTypeErr != nil {
		t.Error(mimeTypeErr)
	}

	thumbnailErr := CreateThumbnail(originalFilename, destinationFilename, mimeType.String())
	if thumbnailErr != nil {
		t.Error(thumbnailErr)
	}

	_, err := os.Stat(destinationFilename)
	if err != nil {
		t.Error("failed to create destination thumbnail")
	}

	w, _ := getImageWidthAndHeight(destinationFilename)

	if w != ThumbnailMaxWidth {
		t.Error("destination thumbnail width should be ", ThumbnailMaxWidth)
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
