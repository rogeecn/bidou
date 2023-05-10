/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/rogeecn/bilibili-downloader/logic"
	"github.com/spf13/cobra"
)

type videoGetParams struct {
	BVID     string
	SavePath string
}

var vParams = &videoGetParams{}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use: "get",
	Run: runGet,
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&vParams.BVID, "bvid", "B", "", "bvid")
	getCmd.Flags().StringVarP(&vParams.SavePath, "save_path", "S", "/opt/downloads/bilibili", "save to dir")
}

func runGet(cmd *cobra.Command, args []string) {
	if vParams.BVID == "" && len(args) == 0 {
		fmt.Println("Please video url or bvid")
		return
	}

	var err error
	bvID := vParams.BVID
	if bvID == "" {
		bvID, err = logic.ParseVideoID(args[0])
		if err != nil {
			log.Fatal(err)
		}
	}
	log.Println("BVID: ", bvID)

	argsItems := []string{
		"--bvid", bvID,
		"--buvid3", logic.GetCookieValue("buvid3"),
		"--bili_jct", logic.GetCookieValue("bili_jct"),
		"--sessdata", logic.GetCookieValue("sessdata"),
		"--save_path", vParams.SavePath,
	}

	log.Println("cmd: python", strings.Join(argsItems, " "))
	execCmd := exec.CommandContext(context.Background(), "/usr/local/bin/bider", argsItems...)
	output, err := execCmd.CombinedOutput()
	if err != nil {
		log.Println("ERR: ", err)
		return
	}
	log.Printf("\n\nOUTPUT: \n\n %s\n", output)
}
