package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/imroc/req/v3"
	"github.com/rogeecn/atom/providers/log"
	"github.com/rogeecn/bidou/modules/bidou/dto"
	"github.com/rogeecn/bidou/pkg/bili"
	"github.com/rogeecn/bidou/pkg/consts"
	"github.com/rogeecn/bidou/pkg/cookiejar"
	"golang.org/x/sync/errgroup"
)

type BidouService struct {
	mutex sync.Mutex

	http       *req.Client
	jar        *cookiejar.Jar
	checkQRKey string
	redis      *redis.Client
}

func NewBidouService(redis *redis.Client, http *req.Client, jar *cookiejar.Jar) *BidouService {
	svc := &BidouService{http: http, jar: jar, redis: redis}
	go svc.CheckQRLoginStatus()
	return svc
}

func (svc *BidouService) GetName(ctx context.Context) (string, error) {
	return "Bidou.GetName", nil
}

// 获取登录二维码链接地址, 并检查扫码状态
func (svc *BidouService) GetQrCodeInfo(ctx context.Context) (*bili.QRCodeInfo, error) {
	url := "https://passport.bilibili.com/x/passport-login/web/qrcode/generate?source=main_web"

	client := bili.New(svc.http)
	if err := client.Get(url); err != nil {
		return nil, err
	}

	var resp bili.QRCodeInfo
	if err := client.Decode(&resp); err != nil {
		return nil, err
	}

	// 写入KEY来启动检查扫码状态
	svc.SetQRKey(resp.QRcodeKey)

	return &resp, nil
}

// 检查扫码状态, 这儿是个死循环, 一直检查扫码状态
func (svc *BidouService) CheckQRLoginStatus() {
	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for range ticker.C {
		svc.mutex.Lock()
		key := svc.checkQRKey
		svc.mutex.Unlock()

		if key == "" {
			continue
		}
		log.Infof("check qr login status by key: %s", key)

		scanResult, err := svc.GetScanResult(key)
		if err != nil {
			log.Error("CheckQRLoginStatus->GetScanResult err: ", err)
			continue
		}

		if scanResult.Message == "二维码已失效" {
			svc.RemoveQRKey()
			continue
		}

		if scanResult.RefreshToken != "" {
			log.Info("login success: fetch sso cookies")
			svc.RemoveQRKey()
			svc.FetchSSOCookies(scanResult.URL)
			svc.jar.Save()
			log.Info("login done")
		}
	}
}

// 检查扫码状态
func (svc *BidouService) GetScanResult(key string) (*bili.ScanResult, error) {
	url := fmt.Sprintf("https://passport.bilibili.com/x/passport-login/web/qrcode/poll?qrcode_key=%s&source=main_web", key)

	client := bili.New(svc.http)
	if err := client.Get(url); err != nil {
		return nil, err
	}

	var scanResult bili.ScanResult
	if err := client.Decode(&scanResult); err != nil {
		return nil, err
	}

	return &scanResult, nil
}

func (svc *BidouService) SetQRKey(key string) {
	svc.mutex.Lock()
	svc.checkQRKey = key
	svc.mutex.Unlock()
}
func (svc *BidouService) RemoveQRKey() {
	svc.mutex.Lock()
	svc.checkQRKey = ""
	svc.mutex.Unlock()
}

func (svc *BidouService) FetchSSOCookies(url string) error {
	client := bili.New(svc.http)
	if err := client.Get(url); err != nil {
		return err
	}

	var ssoCookies bili.SsoCookies
	if err := client.Decode(&ssoCookies); err != nil {
		return err
	}

	var eg errgroup.Group
	for _, url := range ssoCookies.Sso {
		url := url
		eg.Go(func() error {
			log.Infof("request url: %s", url)
			_, err := svc.http.R().Get(url)
			return err
		})
	}

	return eg.Wait()
}

func (svc *BidouService) PushPendingDownloadTasks(bvid string) error {
	url := fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvid)
	client := bili.New(svc.http)
	if err := client.Get(url); err != nil {
		return err
	}

	var vInfo bili.VideoInfoDetail
	if err := client.Decode(&vInfo); err != nil {
		return err
	}

	isAlbum := len(vInfo.Pages) > 1
	for _, v := range vInfo.Pages {
		video := dto.DownloadVideoItem{
			AID:   fmt.Sprintf("%d", vInfo.Aid),
			BVID:  vInfo.Bvid,
			CID:   fmt.Sprintf("%d", v.Cid),
			Album: "",
			Title: v.Part,
		}

		if isAlbum {
			video.Album = vInfo.Title
		}

		svc.redis.RPush(context.Background(), consts.TASK_KEY, video.String())
		log.Info("push video to redis: ", video.String())
	}

	return nil
}
