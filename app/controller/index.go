package controller

import (
	"github.com/gin-gonic/gin"
	"html/template"
	"net/http"
	e2 "pear-admin-go/app/global/e"
	response2 "pear-admin-go/app/global/response"
	"pear-admin-go/app/model"
	"pear-admin-go/app/service"
	pkg "pear-admin-go/app/util/file"
)

//func Index(c *gin.Context) {
//	user := service.GetProfile(c)
//	if pkg.CheckNotExist(service.GetImgSavePath(user.Avatar)) {
//		user.Avatar = e2.DefaultAvatar
//	}
//	site, _ := service.GetSiteConf()
//	c.HTML(http.StatusOK, "index.html", gin.H{
//		"site":      site,
//		"user":      user,
//		"copyright": template.HTML(site.Copyright), // 防止转义
//	})
//}

func Index(c *gin.Context) {
	user := &model.Admin{}
	site := &model.SiteConf{}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"site":      site,
		"user":      user,
		"copyright": template.HTML(site.Copyright), // 防止转义
	})
}

func SliverIndex(c *gin.Context) {
	user := service.GetProfile(c)
	if pkg.CheckNotExist(service.GetImgSavePath(user.Avatar)) {
		user.Avatar = e2.DefaultAvatar
	}
	site, _ := service.GetSiteConf()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"site":      site,
		"user":      user,
		"copyright": template.HTML(site.Copyright), // 防止转义
	})
}

func FramePage(c *gin.Context) {
	response2.SuccessResp(c).WriteJsonExit()
}
