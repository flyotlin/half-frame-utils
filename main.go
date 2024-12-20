package main

import (
	"fmt"
	"image"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

const ROOT_DIR = "00008800 kodak gold200"
const DIST_DIR = "dist"

func main() {
	// flag.Parse()
	// args := flag.Args()
	// log.Println(args[0])

	filepath.WalkDir(ROOT_DIR, visit)

	// filename := "half.jpg"
	// img, err := imaging.Open(filename)
	// if err != nil {
	// 	log.Fatalf("failed to open %v", err)
	// }

	// r := img.Bounds()
	// log.Println(r.Size())
	// size := r.Size()
	// delta, err := strconv.ParseInt(args[0], 0, 10)
	// if err != nil {
	// 	log.Fatalf("failed to parse int %v", err)
	// }

	// imgLeft := imaging.Crop(img, image.Rect(0, 0, size.X/2-int(delta), size.Y))
	// err = imaging.Save(imgLeft, "half-dist-left.jpg")
	// if err != nil {
	// 	log.Fatalf("failed to save %v", err)
	// }
	// copyExif("half.jpg", "half-dist-left.jpg")

	// imgRight := imaging.Crop(img, image.Rect(size.X/2+int(delta), 0, size.X, size.Y))
	// err = imaging.Save(imgRight, "half-dist-right.jpg")
	// if err != nil {
	// 	log.Fatalf("failed to save %v", err)
	// }
	// copyExif("half.jpg", "half-dist-right.jpg")
}

func copyExif(source, destination string) error {
	// Prepare the command to run exiftool
	cmd := exec.Command("exiftool", "-overwrite_original", "-tagsFromFile", source, destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run the command
	return cmd.Run()
}

func visit(path string, d os.DirEntry, err error) error {
	if !strings.HasSuffix(path, ".jpg") {
		return nil
	}
	log.Println(path)
	log.Println(d.Name())

	img, er := imaging.Open(path)
	if er != nil {
		log.Fatalf("failed to open %v", er)
	}

	r := img.Bounds()
	// log.Println(r.Size())
	size := r.Size()
	delta := 90
	// delta, er := strconv.ParseInt(args[0], 0, 10)
	// if er != nil {
	// 	log.Fatalf("failed to parse int %v", er)
	// }

	name, _ := strings.CutSuffix(d.Name(), ".jpg")

	imgLeftName := fmt.Sprintf("%v/%v-1.jpg", DIST_DIR, name)
	log.Println(imgLeftName)
	imgLeft := imaging.Crop(img, image.Rect(0, 0, size.X/2-int(delta), size.Y))
	checkDistDir()
	er = imaging.Save(imgLeft, imgLeftName)
	if er != nil {
		log.Fatalf("failed to save %v", er)
	}
	copyExif(path, imgLeftName)

	imgRightName := fmt.Sprintf("%v/%v-2.jpg", DIST_DIR, name)
	log.Println(imgRightName)
	imgRight := imaging.Crop(img, image.Rect(size.X/2+int(delta), 0, size.X, size.Y))
	checkDistDir()
	er = imaging.Save(imgRight, imgRightName)
	if er != nil {
		log.Fatalf("failed to save %v", er)
	}
	copyExif(path, imgRightName)
	return nil
}

func checkDistDir() {
	_, err := os.Stat(DIST_DIR)
	if os.IsNotExist(err) {
		os.Mkdir(DIST_DIR, os.ModePerm)
	} else if err != nil {
		log.Fatalf("failed to stat: %v", err)
	}
}
