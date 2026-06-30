package router

import (
	"github.com/lihongsheng/pay-gateway/log"
	"github.com/lihongsheng/pay-gateway/middleware"
	"github.com/lihongsheng/pay-gateway/plugin"
	"github.com/lihongsheng/pay-gateway/utils/response"

	"github.com/gin-gonic/gin"
)

// PluginRouter ：/api/v1/plugin/list 与每个插件自身路由 /api/v1/plugin/<name>/*
func PluginRouter(g *gin.Engine) {
	p := g.Group("/api/v1/plugin", middleware.JWTAuth(), middleware.CasbinAuth())
	p.GET("/list", func(c *gin.Context) {
		type item struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}
		out := []item{}
		for _, pl := range plugin.All() {
			out = append(out, item{pl.Name(), pl.Version()})
		}
		response.OK(c, gin.H{"list": out, "total": len(out)})
	})
	for _, pl := range plugin.All() {
		pl.RegisterRoute(g)
		log.Global().Info("plugin route mounted:" + pl.Name())
	}
}
