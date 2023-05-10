package logic

import (
	"log"
	"os"
	"path/filepath"

	"github.com/rogeecn/bilibili-downloader/pkg/cookiejar"

	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

const (
	UA_EDGE = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.35"
)

var client *req.Client
var jar *cookiejar.Jar

func Init() error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	cookieFilePath := filepath.Join(cacheDir, "bilibili-downloader", "cookies.json")
	log.Println("cookie file path:", cookieFilePath)

	// check dir exists, if not exists then create it recursively
	dir := filepath.Dir(cookieFilePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	jar, err = cookiejar.New(&cookiejar.Options{
		PublicSuffixList: nil,
		Filename:         cookieFilePath,
	})
	if err != nil {
		return err
	}
	if viper.GetBool("DEBUG") {
		log.Println("dev mode")
		client = req.C().DevMode()
	} else {
		client = req.C()
	}
	client.SetUserAgent(UA_EDGE)
	client.SetCommonHeader("Referer", "https://www.bilibili.com/")
	client.SetCookieJar(jar)
	return nil
}

func SaveCookies() {
	jar.Save()
}

func GetCookieValue(name string) string {
	all := jar.KVData()
	if v, ok := all[name]; ok {
		return v
	}
	return ""
}
