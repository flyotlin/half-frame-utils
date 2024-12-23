package utils

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/disintegration/imaging"
)

// crop the image in half
func CropInHalf(path string, destDir string) {
	img, err := imaging.Open(path)
	if err != nil {
		log.Fatalf("failed to open [%v]: %v", path, err)
	}
	originalSize := getImageSize(img)
	basename, _ := strings.CutSuffix(path, ".jpg")

	names := strings.Split(basename, "/")
	filename := names[len(names)-1]

	var name string
	rects := []image.Rectangle{
		image.Rect(0, 0, originalSize.X/2, originalSize.Y),
		image.Rect(originalSize.X/2, 0, originalSize.X, originalSize.Y),
	}

	for i, rect := range rects {
		name = fmt.Sprintf("%v/%v-%v.jpg", destDir, filename, i+1)
		crop(img, rect, name)
		copyExif(basename+".jpg", name)
	}
}

func getImageSize(img image.Image) image.Point {
	b := img.Bounds()
	return b.Size()
}

func crop(img image.Image, cropSize image.Rectangle, name string) {
	croppedImage := imaging.Crop(img, cropSize)
	err := imaging.Save(croppedImage, name)
	if err != nil {
		log.Fatalf("failed to save image: [%v]", err)
	}
}

func copyExif(source, destination string) error {
	// Prepare the command to run exiftool
	cmd := exec.Command("exiftool", "-overwrite_original", "-tagsFromFile", source, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}
