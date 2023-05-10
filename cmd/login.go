package cmd

import (
	"log"
	"time"

	"github.com/rogeecn/bilibili-downloader/logic"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use: "login",
	Run: func(cmd *cobra.Command, args []string) {
		qr, err := logic.GetQRUrl()
		if err != nil {
			log.Fatal(err)
		}

		logic.ShowQRCode(qr.Data.URL)
		ticker := time.NewTicker(time.Second * 3)
		defer ticker.Stop()
		defer logic.SaveCookies()

		for range ticker.C {
			scanResult, err := logic.GetScanResult(qr.Data.QRcodeKey)
			if err != nil {
				log.Fatal(err)
			}

			if scanResult.Data.Message == "二维码已失效" {
				log.Fatal("二维码已失效")
			}

			if scanResult.Data.RefreshToken != "" {

				log.Println("login: fetch sso cookies")
				logic.FetchSSOCookies(scanResult.Data.URL)
				log.Println("login: success")
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
}
