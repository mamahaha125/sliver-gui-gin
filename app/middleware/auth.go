package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pear-admin-go/app/core/cache"
	e2 "pear-admin-go/app/global/e"
	"pear-admin-go/app/service"
	"pear-admin-go/app/util/gconv"
	pkg "pear-admin-go/app/util/str"
	"strings"
)

func AuthMiddleware(c *gin.Context) {
	if service.IsSignedIn(c) == false {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	} else {
		if AuthByPage(c) == false {
			c.Redirect(http.StatusFound, "/not_found")
			c.Abort()
		} else {
			if c.Request.URL.Path == "/system" {
				c.Redirect(http.StatusFound, "/system/")
			}
			c.Next()
		}
	}
}

func AuthByPage(c *gin.Context) bool {
	user := service.GetProfile(c)
	if service.IsAdmin(user) { // 管理员不限制
		return true
	}
	url := c.Request.URL.Path
	if strings.HasSuffix(url, "/") {
		url = strings.TrimRight(url, "/")
	}

	allowAuthArr := strings.Split(e2.AllowAuth, ",") // 校验公共路径
	if pkg.IsContain(allowAuthArr, url) {
		return true
	}

	var allowUrlArr []string
	menuCache, found := cache.Instance().Get(e2.MenuCache + gconv.String(user.ID))
	if found && menuCache != nil { //从缓存取菜单
		menu := menuCache.(service.CacheMenuV2)
		allowUrlArr = strings.Split(menu.AllowUrl, ",")
	} else {
		result := service.GetAuth(user)
		for _, v := range result {
			if pkg.IsContain([]string{"", "/"}, v.AuthUrl) == false {
				allowUrlArr = append(allowUrlArr, v.AuthUrl)
			}
		}
	}
	if pkg.IsContain(allowUrlArr, url) == false { // 校验用户路径
		return false
	}
	return true
}

func CheckLoginPage(c *gin.Context) {
	if service.IsSignedIn(c) == true {
		c.Redirect(http.StatusFound, "/system/index")
	}
}

func CheckDefaultPage(c *gin.Context) {
	if service.IsSignedIn(c) == false {
		c.Redirect(http.StatusFound, "/login")
	} else {
		c.Redirect(http.StatusFound, "/system/index")
	}
}

func CheckDefaultConf(c *gin.Context) {
	if service.ImportCfg(c) == true {
		c.Redirect(http.StatusFound, "/")
	}
}
