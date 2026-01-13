package lib

import (
	"fmt"
	"log/slog"
	"os"
	"testing"

	"github.com/gabriel-vasile/mimetype"
)

const ImagePath = "../testdata"

func createThumbnailAndVerify(t *testing.T, originalFilename string) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	destinationFilename := GetThumbnailFilename(originalFilename)

	// If the thumbnail exists already, remove it
	_, destinationExistsErr := os.Stat(destinationFilename)
	if destinationExistsErr == nil {
		removeErr := os.Remove(destinationFilename)
		if removeErr != nil {
			t.Fatal(removeErr)
		}
	}

	mimeType, mimeTypeErr := mimetype.DetectFile(originalFilename)
	if mimeTypeErr != nil {
		t.Fatal(mimeTypeErr)
	}

	// Create the thumbnail and check if there was an error
	thumbnailErr := CreateThumbnail(originalFilename, destinationFilename, mimeType.String(), logger)
	if thumbnailErr != nil {
		t.Fatal(thumbnailErr)
	}

	// Check that the thumbnail exists now
	_, createDestinationThumbnailErr := os.Stat(destinationFilename)
	if createDestinationThumbnailErr != nil {
		t.Fatal("failed to create destination thumbnail")
	}

	imageWidthHeight, _ := GetImageWidthAndHeight(destinationFilename, logger)
	if imageWidthHeight.Width != ThumbnailMaxWidth {
		t.Fatal("destination thumbnail width should be ", ThumbnailMaxWidth)
	}

	// Clean up
	deleteErr := os.Remove(destinationFilename)
	if deleteErr != nil {
		t.Fatal(deleteErr)
	}
}

func TestCreateThumbnailJpg(t *testing.T) {
	testImages := []string{"purple.jpg", "cuttlefish.jpg", "diatom.jpg"}
	for _, image := range testImages {
		originalFilename := fmt.Sprintf("%s/%s", ImagePath, image)
		createThumbnailAndVerify(t, originalFilename)
	}
}

func TestCreateThumbnailPng(t *testing.T) {
	testImages := []string{"red.png", "pikachu.png", "shadow-alakazam.png"}
	for _, image := range testImages {
		originalFilename := fmt.Sprintf("%s/%s", ImagePath, image)
		createThumbnailAndVerify(t, originalFilename)
	}
}

func TestCreateThumbnailGif(t *testing.T) {
	testImages := []string{"samurai-doge.gif", "blastoise-pokemon.gif", "snorlax.gif"}
	for _, image := range testImages {
		originalFilename := fmt.Sprintf("%s/%s", ImagePath, image)
		createThumbnailAndVerify(t, originalFilename)
	}
}
