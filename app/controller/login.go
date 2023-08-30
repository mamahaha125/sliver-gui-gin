package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/mssola/user_agent"
	"net/http"
	rpcs "pear-admin-go/app/core/rpc"
	"pear-admin-go/app/dao"
	"pear-admin-go/app/global/e"
	"pear-admin-go/app/global/request"

	response2 "pear-admin-go/app/global/response"
	"pear-admin-go/app/model"
	"pear-admin-go/app/service"
	"pear-admin-go/app/util/captcha"
	"pear-admin-go/app/util/clientIP"
	"pear-admin-go/app/util/gconv"
	"pear-admin-go/app/util/ip"
	pkg "pear-admin-go/app/util/str"
	"pear-admin-go/app/util/validate"
	"strings"
	"time"
)

func Login(c *gin.Context) {
	confs := service.GetConf(c)
	c.HTML(http.StatusOK, "login3.html", gin.H{"confs": confs})
}

func ConnectServer(c *gin.Context) {

	var CFG struct {
		File []byte `json:"file"`
		Conf string `json:"conf"`
	}

	if err := c.ShouldBindBodyWith(&CFG, binding.JSON); err != nil {
		fmt.Println(err)
	}
	fmt.Println(CFG)
	rpc, con, ln := service.ConnectRpc(CFG.Conf)
	fmt.Println(rpc, ln)
	service.SetCon(con, rpc)
	response2.SuccessResp(c).SetMsg("登陆成功").SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()

}

func Test(c *gin.Context) {
	a := rpcs.GetInstance()
	fmt.Println(a)
}

func LoginHandler(c *gin.Context) {
	var req request.LoginForm
	if err := c.ShouldBind(&req); err != nil {
		response2.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	isLock := service.CheckLock(req.UserName)
	if isLock {
		response2.ErrorResp(c).SetMsg("密码错误次数超限，账号已锁定,请稍后再试").SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
		return
	}
	userAgent := c.Request.Header.Get("User-Agent")
	ua := user_agent.New(userAgent)
	ub, _ := ua.Browser()

	var info model.LoginInfo
	info.LoginName = req.UserName
	info.IpAddr = clientIP.GetIp(c.Request)
	info.Os = ua.OS()
	info.Browser = ub
	info.LoginTime = time.Now()
	info.LoginLocation = ip.GetCityByIp(clientIP.GetIp(c.Request))

	if sid, err := service.SignIn(req.UserName, req.Password, c); err != nil {
		errNums := service.SetPwdErrNum(req.UserName)
		having := e.MaxErrNum - errNums
		info.Msg = "账号或密码错误"
		info.Status = "0"
		err = dao.NewLoginInfoImpl().Insert(info)
		if err != nil {
			response2.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		}
		if having > 0 {
			response2.ErrorResp(c).SetMsg("账号或密码不正确,还有"+gconv.String(having)+"次之后账号将锁定").SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		} else {
			response2.ErrorResp(c).SetMsg("密码错误次数超限，账号已锁定,请稍后再试").SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		}
	} else {
		var online model.AdminOnline
		err = pkg.CopyFields(&online, info)
		if err != nil {
			response2.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		}
		online.SessionID = sid
		online.Status = "on_line"
		online.ExpireTime = 1440
		online.StartTimestamp = time.Now().Unix()
		online.LastAccessTime = time.Now().Format(e.TimeFormat)
		dao.NewAdminOnlineDaoImpl().Delete(sid)
		dao.NewAdminOnlineDaoImpl().Insert(online)
		service.RemovePwdErrNum(req.UserName)

		info.Msg = "登陆成功"
		info.Status = "1"
		err := dao.NewLoginInfoImpl().Insert(info)
		if err != nil {
			response2.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
		}
		response2.SuccessResp(c).SetMsg("登陆成功").SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
	}
}

// sliver登录
func RpcConnectHandle(c *gin.Context) {
	var req request.LoginForm
	if err := c.ShouldBind(&req); err != nil {
		response2.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	isLock := service.CheckLock(req.UserName)
	if isLock {
		response2.ErrorResp(c).SetMsg("密码错误次数超限，账号已锁定,请稍后再试").SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
		return
	}
	userAgent := c.Request.Header.Get("User-Agent")
	ua := user_agent.New(userAgent)
	ub, _ := ua.Browser()

	var info model.LoginInfo
	info.LoginName = req.UserName
	info.IpAddr = clientIP.GetIp(c.Request)
	info.Os = ua.OS()
	info.Browser = ub
	info.LoginTime = time.Now()
	info.LoginLocation = ip.GetCityByIp(clientIP.GetIp(c.Request))

	if sid, err := service.SignIn(req.UserName, req.Password, c); err != nil {
		errNums := service.SetPwdErrNum(req.UserName)
		having := e.MaxErrNum - errNums
		info.Msg = "账号或密码错误"
		info.Status = "0"
		err = dao.NewLoginInfoImpl().Insert(info)
		if err != nil {
			response2.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		}
		if having > 0 {
			response2.ErrorResp(c).SetMsg("账号或密码不正确,还有"+gconv.String(having)+"次之后账号将锁定").SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		} else {
			response2.ErrorResp(c).SetMsg("密码错误次数超限，账号已锁定,请稍后再试").SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		}
	} else {
		var online model.AdminOnline
		err = pkg.CopyFields(&online, info)
		if err != nil {
			response2.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e.LoginHandler, req).WriteJsonExit()
			return
		}
		online.SessionID = sid
		online.Status = "on_line"
		online.ExpireTime = 1440
		online.StartTimestamp = time.Now().Unix()
		online.LastAccessTime = time.Now().Format(e.TimeFormat)
		dao.NewAdminOnlineDaoImpl().Delete(sid)
		dao.NewAdminOnlineDaoImpl().Insert(online)
		service.RemovePwdErrNum(req.UserName)

		info.Msg = "登陆成功"
		info.Status = "1"
		err := dao.NewLoginInfoImpl().Insert(info)
		if err != nil {
			response2.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
		}
		response2.SuccessResp(c).SetMsg("登陆成功").SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
	}
}

// 注销
func Logout(c *gin.Context) {
	if service.IsSignedIn(c) {
		service.SignOut(c)
	}
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
}

func NotFound(c *gin.Context) {
	c.HTML(http.StatusOK, "not_found.html", gin.H{})
}

func GetCaptcha(c *gin.Context) {
	id, b64s, err := captcha.CaptMake()
	if err != nil {
		response2.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response2.SuccessResp(c).SetData(response2.CaptchaResponse{CaptchaId: id, PicPath: b64s}).WriteJsonExit()
}

func CaptchaVerify(c *gin.Context) {
	id := c.PostForm("id")
	capt := c.PostForm("capt")
	if id == "" || capt == "" {
		response2.ErrorResp(c).SetMsg("请填写完整信息").WriteJsonExit()
		return
	}

	if captcha.CaptVerify(id, strings.ToLower(capt)) == true {
		response2.SuccessResp(c).WriteJsonExit()
	} else {
		response2.ErrorResp(c).SetMsg("验证码有误，请重新输入").WriteJsonExit()
	}
	return
}
