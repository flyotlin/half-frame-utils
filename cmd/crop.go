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

	"github.com/spf13/cobra"

	"github.com/flyotlin/half-frame-utils/utils"
)

// cropCmd represents the crop command
var cropCmd = &cobra.Command{
	Use:   "crop",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: cropRun,
}

var (
	src   string
	dest  string
	width int
)

func init() {
	rootCmd.AddCommand(cropCmd)

	// TODO: src becomes a required positional argument
	cropCmd.Flags().StringVarP(&src, "src", "s", "", "source file/directory path")
	cropCmd.Flags().StringVarP(&dest, "dest", "d", "./", "destination directory path")
	cropCmd.Flags().IntVarP(&width, "width", "w", 0, "width for cropping")
	// TODO: suffix, prefix

	cropCmd.MarkFlagRequired("src")
	// cropCmd.MarkFlagRequired("dest")
}

func cropRun(cmd *cobra.Command, args []string) {
	stat := statPath(src)
	if stat.IsDir() {
		log.Printf("start to crop images in directory [%v]...\n", src)
		statPath(dest)
	} else {
		log.Printf("start to crop image [%v]...\n", src)
	}

	if stat.IsDir() { // Directory
		filepath.WalkDir(src, visitDir)
	} else { // File
		utils.CropInHalf(src, dest)
	}
}

func statPath(path string) fs.FileInfo {
	if path == "" {
		log.Fatalln("--dest is not given")
	}

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

func visitDir(path string, d os.DirEntry, err error) error {
	if !strings.HasSuffix(path, ".jpg") {
		return nil
	}
	utils.CropInHalf(path, dest)
	return nil
}
