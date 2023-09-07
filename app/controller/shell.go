package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	rpcs "pear-admin-go/app/core/rpc"
	"pear-admin-go/app/global/response"
	"pear-admin-go/app/service"
	"strings"
)

// @Title GetFilePage
// @Description TODO
// @Author mamahaha125 2023-08-24 11:22:58
// @Update mamahaha125 2023-08-24 11:22:58
// @Param c
func GetFilePage(c *gin.Context) {
	c.HTML(http.StatusOK, "sessions_file.html", nil)
}

// @Title ChooseSession
// @Description TODO
// @Author mamahaha125 2023-08-24 11:22:56
// @Update mamahaha125 2023-08-24 11:22:56
// @Param c
func ChooseSession(c *gin.Context) {
	sessions := c.Query("name")
	err := service.ChooseSessionByID(sessions)
	if err != nil {
		response.ErrorResp(c).WriteJsonExit()
	}
	response.SuccessResp(c).WriteJsonExit()
	//c.HTML(http.StatusOK, "sessions_file.html", nil)
}

func GetSessionsPage(c *gin.Context) {

	c.HTML(http.StatusOK, "sessions_list.html", nil)

}

// @Title GetSessionList
// @Description SessionsList
// @Author mamahaha125 2023-08-24 11:22:51
// @Update mamahaha125 2023-08-24 11:22:51
// @Param c
func GetSessionList(c *gin.Context) {
	con := rpcs.GetInstance().GetCon()
	data, err := service.GetSessions(con)
	datas, count := service.ResSessions(data)
	if err != nil {
		response.ErrorResp(c).WriteJsonExit()
	}

	response.SuccessResp(c).SetCode(0).SetCount(count).SetData(datas).WriteJsonExit()

}

func GetFile(c *gin.Context) {
	path := c.PostForm("path")
	data, pwd := service.GetFiles(path)
	response.SuccessResp(c).SetData(data).SetMsg(pwd).WriteJsonExit()
}

func BackFile(c *gin.Context) {
	path := c.PostForm("path")
	data, pwd := service.GetFiles(filepath.Dir(path))
	response.SuccessResp(c).SetData(data).SetMsg(pwd).WriteJsonExit()
}

// @Title GetPS
// @Description TODO
// @Author mamahaha125 2023-08-29 12:20:54
// @Update mamahaha125 2023-08-29 12:20:54
// @Param c
func GetPS(c *gin.Context) {
	//path := c.PostForm("path")
	data := service.PSList()
	response.SuccessResp(c).SetData(data).WriteJsonExit()
}

func Exec(c *gin.Context) {
	args := strings.SplitAfterN(c.PostForm("args"), " ", 2)
	data, err := service.ExecCommand(args)

	if err != nil {
		response.ErrorResp(c).WriteJsonExit()
	}
	response.SuccessResp(c).SetCode(0).SetData(data).WriteJsonExit()
}
