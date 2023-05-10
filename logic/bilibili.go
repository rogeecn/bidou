package logic

import (
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/sync/errgroup"
)

func GetQRUrl() (*QRCodeURLResponse, error) {
	url := "https://passport.bilibili.com/x/passport-login/web/qrcode/generate?source=main_web"

	var qrResponse QRCodeURLResponse
	resp, err := client.R().SetSuccessResult(&qrResponse).Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsErrorState() {
		return nil, resp.ErrorResult().(error)
	}
	return &qrResponse, nil
}

func GetScanResult(key string) (*ScanResult, error) {
	url := fmt.Sprintf("https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=%s&source=main_web", key)

	var scanResult ScanResult
	resp, err := client.R().SetSuccessResult(&scanResult).Get(url)
	if err != nil {
		return nil, err
	}

	if resp.IsErrorState() {
		return nil, resp.ErrorResult().(error)
	}

	return &scanResult, nil
}

func FetchSSOCookies(url string) error {
	var ssoCookies SsoCookies
	resp, err := client.R().SetSuccessResult(&ssoCookies).Get(url)
	if err != nil {
		return err
	}
	if resp.IsErrorState() {
		return resp.ErrorResult().(error)
	}

	var eg errgroup.Group
	for _, url := range ssoCookies.Data.Sso {
		eg.Go(func() error {
			_, err := client.R().Get(url)
			return err
		})
	}

	return eg.Wait()
}
func ParseVideoID(rawURL string) (string, error) {
	// https://www.bilibili.com/video/BV1NP411o7MB/?spm_id_from=333.1007.top_right_bar_window_custom_collection.content.click

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	if !strings.HasSuffix(u.Host, ".bilibili.com") {
		return "", fmt.Errorf("invalid bilibili url")
	}

	path := strings.TrimRight(u.Path, "/")
	if !strings.HasPrefix(path, "/video/") {
		return "", fmt.Errorf("invalid bilibili video url")
	}

	return strings.TrimPrefix(path, "/video/"), nil
}
