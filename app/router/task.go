package router

import (
	"github.com/gin-gonic/gin"
	"pear-admin-go/app/controller"
)

func TaskRouter(r *gin.Engine) {
	tr := r.Group("system")
	tr.GET("server/list", controller.ServerList)
	tr.GET("server/json", controller.ServerJson)
	tr.GET("server/add", controller.ServerAdd)
	tr.POST("server/add", controller.ServerAdd)
	tr.GET("server/edit", controller.ServerEdit)
	tr.POST("server/edit", controller.ServerEdit)
	tr.POST("server/del", controller.ServerDel)

	tr.GET("task/list", controller.TaskList)
	tr.GET("task/json", controller.TaskJson)
	tr.GET("task/add", controller.TaskAdd)
	tr.POST("task/add", controller.TaskAdd)
	tr.GET("task/edit", controller.TaskEdit)
	tr.POST("task/edit", controller.TaskEdit)
	tr.POST("task/del", controller.TaskDel)
}
