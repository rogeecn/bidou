package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rogeecn/atom"
	"github.com/rogeecn/atom/container"
	"github.com/rogeecn/atom/providers/http"
	"github.com/rogeecn/atom/utils/opt"
	"github.com/rogeecn/bidou/modules/bidou/controller"
	"github.com/rogeecn/gen"
)

func Provide(opts ...opt.Option) error {
	return container.Container.Provide(NewRoute, atom.GroupRoutes)
}

type Route struct {
	engine     *gin.Engine
	controller *controller.BidouController
}

func NewRoute(svc http.Service, ctrl *controller.BidouController) http.Route {
	engine := svc.GetEngine().(*gin.Engine)
	return &Route{engine: engine, controller: ctrl}
}

func (r *Route) Register() {
	// r.engine.StaticFile("/", "./resources/statics/index.html") // gen.DataFunc(r.controller.Index))
	r.engine.GET("/login", gen.DataFunc(r.controller.Login))
	r.engine.GET("/crawl", gen.DataFunc(r.controller.Crawl))
}
