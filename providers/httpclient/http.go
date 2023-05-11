package httpclient

import (
	"log"

	"github.com/imroc/req/v3"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"
	"github.com/rogeecn/bidou/pkg/cookiejar"
	"github.com/spf13/viper"
)

const (
	UA_EDGE = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1774.35"
)

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options:  []opt.Option{},
	}
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	return container.Container.Provide(func(jar *cookiejar.Jar) (*req.Client, error) {
		var client *req.Client

		if viper.GetBool("DEBUG") {
			log.Println("dev mode")
			client = req.C().DevMode()
		} else {
			client = req.C()
		}
		client.SetUserAgent(UA_EDGE)
		client.SetCommonHeader("Referer", "https://www.bilibili.com/")
		client.SetCookieJar(jar)

		return client, nil
	}, o.DiOptions()...)
}
