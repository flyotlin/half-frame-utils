/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/flyotlin/half-frame-utils/internal"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload [file or directory]",
	Short: "Upload file or directory to immich",
	Long: `Upload file or directory to immich.

Config precedence:
	env (not implemented yet) > cli argument > config file`,
	Args: cobra.ExactArgs(1),
	Run:  uploadRun,
}

var uploadSrc string

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().StringVar(&config.ImmichUrl, "immich-url", "", "immich url")
	uploadCmd.Flags().StringVar(&config.ImmichApiKey, "immich-api-key", "", "immich api key")
}

var config internal.HFUtilsConfig

func uploadRun(cmd *cobra.Command, args []string) {
	uploadSrc = args[0]
	readHFUtilsConfig()

	execute("immich", "login", config.ImmichUrl, config.ImmichApiKey)
	execute("immich", "upload", uploadSrc)
}

func readHFUtilsConfig() {
	var c internal.HFUtilsConfig

	file, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("failed to read config file: [%v]", err)
	}
	err = toml.Unmarshal(file, &c)
	if err != nil {
		log.Fatalf("failed to unmarshal config file to go-struct: [%v]", err)
	}

	if config.ImmichUrl == "" {
		config.ImmichUrl = c.ImmichUrl
	}

	if config.ImmichApiKey == "" {
		config.ImmichApiKey = c.ImmichApiKey
	}
}

func execute(name string, arg ...string) {
	c := exec.Command(name, arg...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	if err != nil {
		log.Fatalf("failed to execute: [%v]", err)
	}
}
