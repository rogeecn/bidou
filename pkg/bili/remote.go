package bili

import (
	"encoding/json"
	"fmt"

	"github.com/imroc/req/v3"
)

type RemoteResponse struct {
	client *req.Request `json:"-"`

	Code    int             `json:"code"`
	Message string          `json:"message"`
	TTL     int             `json:"ttl"`
	Data    json.RawMessage `json:"data"`
}

func New(client *req.Client) *RemoteResponse {
	return &RemoteResponse{client: client.R()}
}

func (r *RemoteResponse) GetError() error {
	if r.Code != 0 {
		return fmt.Errorf("response error: %s", r.Message)
	}
	return nil
}

func (r *RemoteResponse) With(f func(*req.Request)) *RemoteResponse {
	f(r.client)
	return r
}

func (r *RemoteResponse) Get(url string) error {
	resp, err := r.client.SetSuccessResult(&r).Get(url)
	if err != nil {
		return err
	}
	if resp.IsErrorState() {
		return resp.ErrorResult().(error)
	}
	if err := r.GetError(); err != nil {
		return err
	}

	return nil
}

func (r *RemoteResponse) Decode(v interface{}) error {
	if err := json.Unmarshal(r.Data, v); err != nil {
		return err
	}
	return nil
}

type QRCodeInfo struct {
	URL       string `json:"url"`
	QRcodeKey string `json:"qrcode_key"`
}

func (d QRCodeInfo) From(data json.RawMessage) (*QRCodeInfo, error) {
	err := json.Unmarshal(data, &d)
	return &d, err
}

type ScanResult struct {
	URL          string `json:"url"`
	RefreshToken string `json:"refresh_token"`
	Timestamp    int    `json:"timestamp"`
	Code         int    `json:"code"`
	Message      string `json:"message"`
}

type SsoCookies struct {
	Sso []string `json:"sso"`
}

type VideoInfoDetail struct {
	Bvid      string `json:"bvid"`
	Aid       int    `json:"aid"`
	Videos    int    `json:"videos"`
	Tid       int    `json:"tid"`
	Tname     string `json:"tname"`
	Copyright int    `json:"copyright"`
	Pic       string `json:"pic"`
	Title     string `json:"title"`
	Pubdate   int    `json:"pubdate"`
	Ctime     int    `json:"ctime"`
	Desc      string `json:"desc"`
	DescV2    []struct {
		RawText string `json:"raw_text"`
		Type    int    `json:"type"`
		BizID   int    `json:"biz_id"`
	} `json:"desc_v2"`
	State    int `json:"state"`
	Duration int `json:"duration"`
	Rights   struct {
		Bp            int `json:"bp"`
		Elec          int `json:"elec"`
		Download      int `json:"download"`
		Movie         int `json:"movie"`
		Pay           int `json:"pay"`
		Hd5           int `json:"hd5"`
		NoReprint     int `json:"no_reprint"`
		Autoplay      int `json:"autoplay"`
		UgcPay        int `json:"ugc_pay"`
		IsCooperation int `json:"is_cooperation"`
		UgcPayPreview int `json:"ugc_pay_preview"`
		NoBackground  int `json:"no_background"`
		CleanMode     int `json:"clean_mode"`
		IsSteinGate   int `json:"is_stein_gate"`
		Is360         int `json:"is_360"`
		NoShare       int `json:"no_share"`
		ArcPay        int `json:"arc_pay"`
		FreeWatch     int `json:"free_watch"`
	} `json:"rights"`
	Owner struct {
		Mid  int    `json:"mid"`
		Name string `json:"name"`
		Face string `json:"face"`
	} `json:"owner"`
	Stat struct {
		Aid        int    `json:"aid"`
		View       int    `json:"view"`
		Danmaku    int    `json:"danmaku"`
		Reply      int    `json:"reply"`
		Favorite   int    `json:"favorite"`
		Coin       int    `json:"coin"`
		Share      int    `json:"share"`
		NowRank    int    `json:"now_rank"`
		HisRank    int    `json:"his_rank"`
		Like       int    `json:"like"`
		Dislike    int    `json:"dislike"`
		Evaluation string `json:"evaluation"`
		ArgueMsg   string `json:"argue_msg"`
	} `json:"stat"`
	Dynamic   string `json:"dynamic"`
	Cid       int    `json:"cid"`
	Dimension struct {
		Width  int `json:"width"`
		Height int `json:"height"`
		Rotate int `json:"rotate"`
	} `json:"dimension"`
	Premiere           interface{} `json:"premiere"`
	TeenageMode        int         `json:"teenage_mode"`
	IsChargeableSeason bool        `json:"is_chargeable_season"`
	IsStory            bool        `json:"is_story"`
	IsUpowerExclusive  bool        `json:"is_upower_exclusive"`
	IsUpowerPlay       bool        `json:"is_upower_play"`
	NoCache            bool        `json:"no_cache"`
	Pages              []struct {
		Cid       int    `json:"cid"`
		Page      int    `json:"page"`
		From      string `json:"from"`
		Part      string `json:"part"`
		Duration  int    `json:"duration"`
		Vid       string `json:"vid"`
		Weblink   string `json:"weblink"`
		Dimension struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			Rotate int `json:"rotate"`
		} `json:"dimension"`
		FirstFrame string `json:"first_frame"`
	} `json:"pages"`
	Subtitle struct {
		AllowSubmit bool `json:"allow_submit"`
		List        []struct {
			ID          int64  `json:"id"`
			Lan         string `json:"lan"`
			LanDoc      string `json:"lan_doc"`
			IsLock      bool   `json:"is_lock"`
			SubtitleURL string `json:"subtitle_url"`
			Type        int    `json:"type"`
			IDStr       string `json:"id_str"`
			AiType      int    `json:"ai_type"`
			AiStatus    int    `json:"ai_status"`
			Author      struct {
				Mid            int    `json:"mid"`
				Name           string `json:"name"`
				Sex            string `json:"sex"`
				Face           string `json:"face"`
				Sign           string `json:"sign"`
				Rank           int    `json:"rank"`
				Birthday       int    `json:"birthday"`
				IsFakeAccount  int    `json:"is_fake_account"`
				IsDeleted      int    `json:"is_deleted"`
				InRegAudit     int    `json:"in_reg_audit"`
				IsSeniorMember int    `json:"is_senior_member"`
			} `json:"author"`
		} `json:"list"`
	} `json:"subtitle"`
	IsSeasonDisplay bool `json:"is_season_display"`
	UserGarb        struct {
		URLImageAniCut string `json:"url_image_ani_cut"`
	} `json:"user_garb"`
	HonorReply struct {
	} `json:"honor_reply"`
	LikeIcon   string `json:"like_icon"`
	NeedJumpBv bool   `json:"need_jump_bv"`
}
