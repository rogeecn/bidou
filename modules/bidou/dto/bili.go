package dto

import (
	"encoding/json"

	"github.com/rogeecn/bidou/pkg/bili"
)

type DownloadVideoItem struct {
	BVID    string `json:"bvid"`
	AID     string `json:"aid"`
	CID     string `json:"cid"`
	Album   string `json:"album"`
	Title   string `json:"title"`
	Retries int    `json:"retries"`
}

func (v DownloadVideoItem) String() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v DownloadVideoItem) Path() string {
	if v.Album == "" {
		return bili.CleanName(v.Title)
	}
	return bili.CleanName(v.Album) + "/" + bili.CleanName(v.Title)
}
