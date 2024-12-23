/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/flyotlin/half-frame-utils/internal"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// cropCmd represents the crop command
var cropCmd = &cobra.Command{
	Use:   "crop [file or directory]",
	Short: "Crop images",
	Long: `Two modes are supported in image cropping.

1. Auto: hf-utils calculates width based on color distance for you
2. Manual: set width on your own`,
	Args: cobra.ExactArgs(1),
	Run:  cropRun,
}

var (
	src      string
	dest     string
	cropConf internal.CropConfig
	pb       *progressbar.ProgressBar
)

func init() {
	rootCmd.AddCommand(cropCmd)

	cropCmd.Flags().StringVarP(&dest, "dest-dir", "d", "./", "destination directory path")
	cropCmd.Flags().IntVarP(&cropConf.Width, "width", "w", -1, "width for cropping (80 is recommended)")
	cropCmd.Flags().StringVar(&cropConf.Prefix, "prefix", "", "prefix for cropped images")
	cropCmd.Flags().StringVar(&cropConf.Suffix, "suffix", "", "suffix for cropped images")
}

func cropRun(cmd *cobra.Command, args []string) {
	src = args[0]
	stat := statSrc(src)
	statDestDir(dest)
	if stat.IsDir() {
		log.Printf("start to crop images in directory [%v]...\n", src)
	} else {
		log.Printf("start to crop image [%v]...\n", src)
	}
	if cropConf.Width == -1 {
		log.Println("crop in auto mode...")
	} else {
		log.Println("crop in manual mode...")
	}

	if stat.IsDir() { // Directory
		count := countImages(src)
		pb = progressbar.Default(int64(count), "Cropping images")
		filepath.WalkDir(src, visitDir)
	} else { // File
		internal.CropInHalf(src, dest, cropConf)
	}
}

func statSrc(path string) fs.FileInfo {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Fatalf("[%v] not exist: [%v]", path, err)
	} else if os.IsPermission(err) {
		log.Fatalf("[%v] permission denied: [%v]", path, err)
	} else if err != nil {
		log.Fatalf("failed to stat [%v]: [%v]", path, err)
	}

	return stat
}

func statDestDir(path string) fs.FileInfo {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Printf("[%v] not exist, create a new one", path)
		err = os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Fatalf("failed to create a new dir [%v]: %v", path, err)
		}
	} else if os.IsPermission(err) {
		log.Fatalf("[%v] permission denied: [%v]", path, err)
	} else if err != nil {
		log.Fatalf("failed to stat [%v]: [%v]", path, err)
	}

	return stat
}

func countImages(path string) int {
	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatalf("failed to list directory %v: [%v]", path, err)
	}
	count := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".jpg") {
			count += 1
		}
	}
	return count
}

func visitDir(path string, d os.DirEntry, err error) error {
	if !strings.HasSuffix(path, ".jpg") {
		return nil
	}
	internal.CropInHalf(path, dest, cropConf)
	pb.Add(1)
	return nil
}
