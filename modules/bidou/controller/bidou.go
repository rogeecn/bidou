package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/rogeecn/atom/providers/log"
	"github.com/rogeecn/bidou/modules/bidou/service"
	"github.com/rogeecn/bidou/pkg/bili"
)

type BidouController struct {
	svc *service.BidouService
}

func NewBidouController(bidouSvc *service.BidouService, _ *service.DownloaderService) *BidouController {
	return &BidouController{
		svc: bidouSvc,
	}
}

func (c *BidouController) Index(ctx *gin.Context) (string, error) {
	return c.svc.GetName(ctx)
}

func (c *BidouController) Login(ctx *gin.Context) (*bili.QRCodeInfo, error) {
	return c.svc.GetQrCodeInfo(ctx)
}

func (c *BidouController) Crawl(ctx *gin.Context) ([]string, error) {
	rawUrl := ctx.Query("url")
	if rawUrl == "" {
		return nil,errors.New("Invalid Url")
	}
	bvid, err := bili.GetBVIDFromURL(rawUrl)
	if err != nil {
		return nil,err
	}

	// add bvid to crawl list
	log.Infof("BVID: %s", bvid)

	return c.svc.PushPendingDownloadTasks(bvid)
}
