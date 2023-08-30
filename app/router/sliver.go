package router

import (
	"github.com/gin-gonic/gin"
	controller "pear-admin-go/app/controller"
)

func SliverRouter(route *gin.Engine) {

	sliver := route.Group("conf")
	{
		//配置相关
		//sliver.POST("/client/create", controller.SliverCfgimport)
		sliver.POST("/client/import", controller.ImportConf)
		//sliver.GET("/client/connect", controller.RpcConnect)
		sliver.POST("/client/user", controller.ChooseConf)

		//sliver.POST("/client/rpc", controller.ConnectRPC)
	}
	slivers := route.Group("sliver")
	{
		slivers.GET("index", controller.Index)
		slivers.GET("jobs", controller.GetJob)
		slivers.GET("sessions", controller.GetSession)
	}
	implant := route.Group("implant")
	{
		implant.POST("build", controller.Build)
		implant.GET("getimplants", controller.GetImplants)
		implant.GET("testb", controller.ImplantsList)
		implant.POST("del", controller.DelImplants)
		implant.GET("download", controller.DownImplants)
		implant.GET("session/list", controller.GetSessionsPage)
		implant.GET("session/lists", controller.GetSessionList)
		implant.GET("session/choose", controller.ChooseSession)
		implant.GET("session/file", controller.GetFilePage)
	}
	shell := route.Group("shell")
	{
		shell.GET("choose", controller.ChooseSession)
		shell.POST("file", controller.GetFile)
		shell.POST("file/back", controller.BackFile)
		shell.POST("ps", controller.GetPS)
		shell.POST("exec", controller.Exec)
		shell.GET("websocket", controller.WebSokcetListen)
	}

	mtls := route.Group("mtls")
	{
		mtls.POST("listening", controller.MtlsListening)
	}
}
