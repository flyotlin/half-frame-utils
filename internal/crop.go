package internal

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"

	"github.com/disintegration/imaging"
)

// crop the image in half
func CropInHalf(path string, destDir string, config CropConfig) {
	img, err := imaging.Open(path)
	if err != nil {
		log.Fatalf("failed to open [%v]: %v", path, err)
	}
	originalSize := getImageSize(img)
	basename, _ := strings.CutSuffix(path, ".jpg")

	names := strings.Split(basename, "/")
	filename := fmt.Sprintf("%v%v%v", config.Prefix, names[len(names)-1], config.Suffix)

	var w int
	if config.Width == -1 {
		left, right := CalculateIntervalWidth(img)
		w = int(math.Min(float64(left), float64(right)))
	} else {
		w = config.Width
	}
	var name string
	rects := []image.Rectangle{
		image.Rect(0, 0, originalSize.X/2-w, originalSize.Y),
		image.Rect(originalSize.X/2+w, 0, originalSize.X, originalSize.Y),
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

func CalculateIntervalWidth(img image.Image) (int, int) {
	size := getImageSize(img)
	base := getBase(img)
	threshold := 30

	// left
	left := size.X / 2
	for left >= 0 {
		// calculate diff for each column (sample 20 points)
		// if any diff for each column greater than threshold, break
		interval := size.Y / 50
		j := 0
		br := false
		for j < size.Y {
			c := img.At(left, j).(color.YCbCr)
			x := math.Pow(float64(int(c.Y)-int(base.Y)), 2)
			y := math.Pow(float64(int(c.Cb)-int(base.Cb)), 2)
			z := math.Pow(float64(int(c.Cr)-int(base.Cr)), 2)
			diff := math.Sqrt(x + y + z)
			if diff >= float64(threshold) {
				br = true
				break
			}
			j += interval
		}
		if br {
			break
		}
		left--
	}

	// right
	right := size.X / 2
	for right < size.Y {
		interval := size.Y / 50
		j := 0
		br := false
		for j < size.Y {
			c := img.At(right, 0).(color.YCbCr)
			x := math.Pow(float64(int(c.Y)-int(base.Y)), 2)
			y := math.Pow(float64(int(c.Cb)-int(base.Cb)), 2)
			z := math.Pow(float64(int(c.Cr)-int(base.Cr)), 2)
			diff := math.Sqrt(x + y + z)
			if diff >= float64(threshold) {
				br = true
				break
			}
			j += interval
		}
		if br {
			break
		}
		right++
	}

	// single sample point
	// right := size.X / 2
	// for right < size.Y {
	// 	c := img.At(right, 0).(color.YCbCr)
	// 	x := math.Pow(float64(int(c.Y)-int(base.Y)), 2)
	// 	y := math.Pow(float64(int(c.Cb)-int(base.Cb)), 2)
	// 	z := math.Pow(float64(int(c.Cr)-int(base.Cr)), 2)
	// 	diff := math.Sqrt(x + y + z)
	// 	if diff >= float64(threshold) {
	// 		break
	// 	}
	// 	right++
	// }

	left = size.X/2 - left
	right = right - size.X/2
	return left, right
}

func getBase(img image.Image) color.YCbCr {
	size := getImageSize(img)

	var (
		y  int
		cb int
		cr int
	)
	y = 0
	cb = 0
	cr = 0
	for i := range size.Y {
		cc := img.At(size.X/2, i).(color.YCbCr)
		y += int(cc.Y)
		cb += int(cc.Cb)
		cr += int(cc.Cr)
	}

	return color.YCbCr{uint8(y / size.Y), uint8(cb / size.Y), uint8(cr / size.Y)}
}
