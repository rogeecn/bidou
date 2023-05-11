package main

import (
	"os"

	"github.com/rogeecn/bidou/database/migrations"
	"github.com/rogeecn/bidou/database/seeders"
	"github.com/rogeecn/bidou/modules/bidou"
	"github.com/rogeecn/bidou/providers/httpclient"
	"github.com/rogeecn/bidou/providers/jar"

	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/providers/database/redis"
	"github.com/rogeecn/atom/services"
	"github.com/spf13/cobra"
)

func main() {
	providers := atom.DefaultHTTP(redis.DefaultProvider())
	providers = append(providers, jar.DefaultProvider(), httpclient.DefaultProvider())
	providers = append(providers, bidou.Providers()...)

	opts := []atom.Option{
		atom.Name("http"),
		atom.RunE(func(cmd *cobra.Command, args []string) error {
			return services.ServeHttp()
		}),
		atom.Seeders(seeders.Seeders...),
		atom.Migrations(migrations.Migrations...),
	}

	if err := atom.Serve(providers, opts...); err != nil {
		os.Exit(1)
	}
}
