package main

import (
	"flag"
	"image"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/disintegration/imaging"
)

func main() {
	flag.Parse()
	args := flag.Args()
	log.Println(args[0])

	filename := "half.jpg"
	img, err := imaging.Open(filename)
	if err != nil {
		log.Fatalf("failed to open %v", err)
	}

	r := img.Bounds()
	log.Println(r.Size())
	size := r.Size()
	delta, err := strconv.ParseInt(args[0], 0, 10)
	if err != nil {
		log.Fatalf("failed to parse int %v", err)
	}

	imgLeft := imaging.Crop(img, image.Rect(0, 0, size.X/2-int(delta), size.Y))
	err = imaging.Save(imgLeft, "half-dist-left.jpg")
	if err != nil {
		log.Fatalf("failed to save %v", err)
	}
	copyExif("half.jpg", "half-dist-left.jpg")

	imgRight := imaging.Crop(img, image.Rect(size.X/2+int(delta), 0, size.X, size.Y))
	err = imaging.Save(imgRight, "half-dist-right.jpg")
	if err != nil {
		log.Fatalf("failed to save %v", err)
	}
	copyExif("half.jpg", "half-dist-right.jpg")
}

func copyExif(source, destination string) error {
	// Prepare the command to run exiftool
	cmd := exec.Command("exiftool", "-overwrite_original", "-tagsFromFile", source, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}
