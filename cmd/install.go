/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	_ "embed"
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

//go:embed downloader.py
var Downloader []byte

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use: "install",
	Run: func(cmd *cobra.Command, args []string) {
		pkgs := []string{"asyncio", "bilibili-api-python", "httpx", "argparse"}
		for _, pkg := range pkgs {
			log.Println("installing", pkg)
			exeCmd := exec.Command("pip3", "install", pkg)
			if err := exeCmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

		// save downloader.py to /usr/local/bin/bider
		if err := os.WriteFile("/usr/local/bin/bider", Downloader, 0755); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
