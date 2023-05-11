package jar

import (
	"os"
	"path/filepath"

	"github.com/rogeecn/bidou/pkg/cookiejar"

	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/providers/log"
	"github.com/rogeecn/atom/utils/opt"
)

func DefaultProvider() container.ProviderContainer {
	return container.ProviderContainer{
		Provider: Provide,
		Options:  []opt.Option{},
	}
}

func Provide(opts ...opt.Option) error {
	o := opt.New(opts...)
	return container.Container.Provide(func() (*cookiejar.Jar, error) {
		cacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, err
		}
		cookieFilePath := filepath.Join(cacheDir, "bidou", "cookies.json")
		log.Info("cookie file path:", cookieFilePath)

		// check dir exists, if not exists then create it recursively
		dir := filepath.Dir(cookieFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				return nil, err
			}
		}

		jar, err := cookiejar.New(&cookiejar.Options{
			PublicSuffixList: nil,
			Filename:         cookieFilePath,
		})
		if err != nil {
			return nil, err
		}
		return jar, nil
	}, o.DiOptions()...)
}
