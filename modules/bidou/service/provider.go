package service

import (
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/utils/opt"
)

func Provide(opts ...opt.Option) error {
	container.Container.Provide(NewBidouService)
	container.Container.Provide(NewDownloaderService)
	return nil
}
