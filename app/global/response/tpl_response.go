package response

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	e2 "pear-admin-go/app/global/e"
)

// 通用tpl响应
type TplResp struct {
	c   *gin.Context
	tpl string
}

// 返回一个tpl响应
func BuildTpl(c *gin.Context, tpl string) *TplResp {
	var t = TplResp{
		c:   c,
		tpl: tpl,
	}
	return &t
}

// 返回一个错误的tpl响应
func ErrorTpl(c *gin.Context) *TplResp {
	var t = TplResp{
		c:   c,
		tpl: "error/error.html",
	}
	return &t
}

// 返回一个无操作权限tpl响应
func ForbiddenTpl(c *gin.Context) *TplResp {
	var t = TplResp{
		c:   c,
		tpl: "error/unauth.html",
	}
	return &t
}

// 输出页面模板
func (resp *TplResp) WriteTpl(params ...gin.H) {
	session := sessions.Default(resp.c)
	uid := session.Get(e2.Auth)
	if params == nil || len(params) == 0 {
		resp.c.HTML(http.StatusOK, resp.tpl, gin.H{"uid": uid})
	} else {
		params[0]["uid"] = uid
		resp.c.HTML(http.StatusOK, resp.tpl, params[0])
	}
}
