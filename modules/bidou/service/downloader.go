package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rogeecn/atom/providers/log"
	"github.com/rogeecn/bidou/modules/bidou/dto"
	"github.com/rogeecn/bidou/pkg/consts"
	"github.com/rogeecn/bidou/pkg/cookiejar"
)

type DownloaderService struct {
	redis       *redis.Client
	jar         *cookiejar.Jar
	rateLimiter chan struct{}
}

func NewDownloaderService(
	redis *redis.Client,
	jar *cookiejar.Jar,
) *DownloaderService {
	svc := &DownloaderService{redis: redis, jar: jar, rateLimiter: make(chan struct{}, 3)}
	go svc.Start()
	return svc
}

func (svc *DownloaderService) Start() {
	log.Info("downloader service start")
	for {
		svc.getToken()
		video, err := svc.getVideo()
		if err != nil {
			svc.releaseToken()
			time.Sleep(time.Second * 5)
			continue
		}
		go svc.RunDownloader(video)
	}
}

func (svc *DownloaderService) getVideo() (*dto.DownloadVideoItem, error) {
	log.Info("wait for task")
	items, err := svc.redis.BLPop(context.Background(), time.Minute, consts.TASK_KEY).Result()
	if err != nil {
		return nil, err
	}

	if len(items) != 2 {
		return nil, errors.New("invalid task item")
	}

	var video dto.DownloadVideoItem
	if err := json.Unmarshal([]byte(items[1]), &video); err != nil {
		log.Errorf("json unmarshal error: %v,raw: %s", err, items[1])
		return nil, err
	}
	return &video, nil
}

func (svc *DownloaderService) RunDownloader(video *dto.DownloadVideoItem) {
	defer svc.releaseToken()

	argsItems := []string{
		"--buvid3", svc.GetCookieValue("buvid3"),
		"--bili_jct", svc.GetCookieValue("bili_jct"),
		"--sessdata", svc.GetCookieValue("sessdata"),
		"--save_path", filepath.Join(consts.SAVE_PATH, video.Path()),
		"--bvid", video.BVID,
		"--cid", video.CID,
	}

	log.Debug("cmd: /usr/local/bin/bider ", strings.Join(argsItems, " "))
	cmd := exec.CommandContext(context.Background(), "/usr/local/bin/bider", argsItems...)

	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout

	if err != nil {
		log.Error(err)
		return
	}

	if err = cmd.Start(); err != nil {
		log.Error(err)
		return
	}

	// 从管道中实时获取输出并打印到终端
	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}

	if err = cmd.Wait(); err != nil {
		video.Retries = video.Retries + 1
		if video.Retries > 10 {
			log.Error("drop task: ", video.BVID, " ", video.CID, " ", video.Title, ", case: ", err)
			return
		}
		svc.redis.RPush(context.Background(), consts.TASK_KEY, video.String())
		log.Error(err)
	}
}

func (svc *DownloaderService) GetCookieValue(name string) string {
	all := svc.jar.KVData()
	if v, ok := all[name]; ok {
		return v
	}
	return ""
}

func (svc *DownloaderService) releaseToken() {
	<-svc.rateLimiter
}

func (svc *DownloaderService) getToken() {
	svc.rateLimiter <- struct{}{}
}
