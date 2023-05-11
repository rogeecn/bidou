package bidou

import (
	"github.com/rogeecn/bidou/modules/bidou/controller"
	"github.com/rogeecn/bidou/modules/bidou/routes"
	"github.com/rogeecn/bidou/modules/bidou/service"

	"github.com/rogeecn/atom/container"
)

func Providers() container.Providers {
	return container.Providers{
		{Provider: service.Provide},
		{Provider: controller.Provide},
		{Provider: routes.Provide},
	}
}
