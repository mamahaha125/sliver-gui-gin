package router

import (
	"github.com/gin-gonic/gin"
	controller "pear-admin-go/app/controller"
)

func SystemRouter(r *gin.Engine) {
	sr := r.Group("system")

	// default
	sr.GET("pear", controller.Pear)
	sr.GET("menu", controller.Menu)
	sr.GET("server_err", controller.ServerErr)
	sr.GET("ui/icon", controller.IconShow)
	sr.GET("file", controller.ShowFile)

	//sr.Use(middleware.AuthMiddleware) // 下面的使用中间件验证  上面的路由为基础路由 不需要验证权限

	// default
	sr.POST("upload", controller.Upload)

	// index
	sr.GET("/", controller.Index)
	//sr.GET("index", controller.Index)
	sr.GET("main", controller.FramePage)

	// log
	sr.GET("log/list", controller.LogList)
	sr.GET("log/operate", controller.LogOperate)
	sr.GET("log/login", controller.LogLogin)

	// role 管理员列表页
	sr.GET("admin/list", controller.AdminList)
	sr.GET("admin/json", controller.AdminJson)
	sr.GET("admin/add", controller.AdminAdd)
	sr.POST("admin/add", controller.AdminAdd)
	sr.GET("admin/edit", controller.AdminEdit)
	sr.POST("admin/edit", controller.AdminEdit)
	sr.POST("admin/status", controller.AdminStatus)
	r.DELETE("admin/delete", controller.AdminDelete)

	// role 角色管理列表页
	sr.GET("role/list", controller.RoleList)
	sr.GET("role/json", controller.RoleJson)
	sr.GET("role/add", controller.RoleAdd)
	sr.POST("role/add", controller.RoleAdd)
	sr.GET("role/power", controller.RolePower)
	sr.GET("role/getRolePower", controller.GetRolePower)
	sr.POST("role/saveRolePower", controller.SaveRolePower)
	sr.GET("role/edit", controller.RoleEdit)
	sr.POST("role/edit", controller.RoleEdit)
	sr.POST("role/delete", controller.RoleDeleteHandler)

	// role 权限因子列表页
	sr.GET("auth/list", controller.AuthList)
	sr.GET("auth/nodes", controller.AuthNodes)    // 权限配置页面
	sr.GET("auth/add", controller.AddNode)        // 新增权限
	sr.GET("auth/edit", controller.AuthEdit)      // 修改权限
	sr.POST("auth/edit", controller.AuthEdit)     // 新增、修改权限
	sr.POST("auth/node", controller.AuthNode)     // 权限因子列表页
	sr.POST("auth/delete", controller.AuthDelete) // 权限因子列表页
	sr.GET("auth/parent", controller.Parent)      // 权限列表

	// 站点设置
	sr.GET("site/edit", controller.SiteEdit)
	sr.POST("site/edit", controller.SiteEdit)

	// mail 邮件配置暂未启用
	sr.GET("mail/list", controller.MailList)
	sr.POST("mail/edit", controller.MailEdit)
	sr.POST("mail/test", controller.MailTest)

	// user
	sr.POST("user/avatar", controller.AvatarEdit)
	sr.GET("user/uploadPage", controller.UploadPage)
	sr.GET("user/edit", controller.UserEdit)
	sr.POST("user/edit", controller.UserEdit)
	sr.GET("user/pwd", controller.PwdEdit)
	sr.POST("user/pwd", controller.PwdEdit)
}
