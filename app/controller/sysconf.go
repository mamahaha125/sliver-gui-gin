package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pear-admin-go/app/core/cache"
	"pear-admin-go/app/global/e"
	"pear-admin-go/app/global/request"
	"pear-admin-go/app/global/response"
	"pear-admin-go/app/model"
	"pear-admin-go/app/service"
	"pear-admin-go/app/util/gomail"
)

func SiteEdit(c *gin.Context) {
	if c.Request.Method == "GET" {
		site, sysID := service.GetSiteConf()
		c.HTML(http.StatusOK, "site_config.html", gin.H{"id": sysID, "site": site})
	} else {
		var f request.SiteConfForm
		if err := c.ShouldBind(&f); err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e.SiteEdit, c.Request.PostForm).WriteJsonExit()
			return
		}
		if err := service.SiteEditService(f); err != nil {
			response.ErrorResp(c).SetType(model.OperEdit).Log(e.SiteEdit, c.Request.PostForm).WriteJsonExit()
			return
		}
		response.SuccessResp(c).SetType(model.OperEdit).Log(e.SiteEdit, c.Request.PostForm).WriteJsonExit()
	}
}

func MailList(c *gin.Context) {
	mail, sysID := service.GetMailConf()
	testMail := service.GetMailTestConf()
	c.HTML(http.StatusOK, "sys_mail_list.html", gin.H{"id": sysID, "mail": mail, "test": testMail})
}

func MailEdit(c *gin.Context) {
	var f request.MailConfForm
	if err := c.ShouldBind(&f); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e.MailEdit, c.Request.PostForm).WriteJsonExit()
		return
	}
	status := c.PostForm("email_status")
	if status == "on" {
		f.EmailStatus = 1
	} else {
		f.EmailStatus = 0
	}
	if err := service.MailEditService(f); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e.MailEdit, c.Request.PostForm).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetType(model.OperEdit).Log(e.MailEdit, c.Request.PostForm).WriteJsonExit()
}

func MailTest(c *gin.Context) {
	var f gomail.MailConfForm
	if err := c.ShouldBind(&f); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e.MailEdit, c.Request.PostForm).WriteJsonExit()
		return
	}
	var testMail gomail.Config
	testMail.Config = f
	testMail.MailTo = append(testMail.MailTo, f.EmailTest)
	testMail.Subject = f.EmailTestTitle
	if err := gomail.SendMail(testMail); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e.MailEdit, c.Request.PostForm).WriteJsonExit()
		return
	}
	// 保存邮件测试配置到缓存
	cache.Instance().Set(e.TestMailConf, model.MailTest{
		EmailTest:      f.EmailTest,
		EmailTestTitle: f.EmailTestTitle,
		EmailTemplate:  f.EmailTemplate,
	}, e.TestMailEffTime)
	response.SuccessResp(c).SetType(model.OperEdit).Log(e.MailEdit, c.Request.PostForm).WriteJsonExit()
}
