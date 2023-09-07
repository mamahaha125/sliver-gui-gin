package router

import (
	"pear-admin-go/app/core/config"

	"github.com/gin-gonic/gin"
	"pear-admin-go/app/middleware"
	"pear-admin-go/app/util/session"
)

func InitRouter() *gin.Engine {
	gin.SetMode(config.Instance().App.RunMode)
	r := gin.New()

	r.Static(config.Instance().App.ImgUrlPath, config.Instance().App.ImgSavePath)
	r.Static("/runtime/file", "runtime/file")

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())
	r.Use(session.EnableCookieSession(config.Instance().App.JwtSecret))
	CommonRouter(r)
	SystemRouter(r)
	TaskRouter(r)
	SliverRouter(r)
	return r
}
