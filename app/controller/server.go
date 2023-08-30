package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	rpcs "pear-admin-go/app/core/rpc"
	"pear-admin-go/app/global/request"
	"pear-admin-go/app/global/response"
	"pear-admin-go/app/service"
	"pear-admin-go/app/util/gconv"
	"pear-admin-go/app/util/validate"
	"strings"
)

func ServerList(c *gin.Context) {
	c.HTML(http.StatusOK, "server_list.html", nil)
}

func ServerJson(c *gin.Context) {
	var f request.TaskServerPage
	if err := c.ShouldBind(&f); err != nil {
		response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
		return
	}
	data, count, err := service.ServerJson(f)
	if err != nil {
		response.SuccessResp(c).SetCode(0).SetMsg(err.Error()).SetCount(count).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetCode(0).SetCount(count).SetData(data).WriteJsonExit()
}

func DelImplants(c *gin.Context) {
	name := c.PostForm("name")
	namelist := strings.Split(name, ",")
	_, err := service.DeleteImplant(namelist)
	if err != nil {
		response.ErrorResp(c).SetCode(0).SetData("err").WriteJsonExit()
	}
	response.SuccessResp(c).SetCode(0).WriteJsonExit()
}

func ImplantsList(c *gin.Context) {

	con := rpcs.GetInstance().GetCon()
	data, count := service.GetBuilds(con)

	response.SuccessResp(c).SetCode(0).SetCount(count).SetData(data).WriteJsonExit()
}

func ServerAdd(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "server_add.html", nil)
	} else {
		var f request.TaskServerForm
		if err := c.ShouldBind(&f); err != nil {
			response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
			return
		}
		err := service.ServerAdd(f)
		if err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
			return
		}
		response.SuccessResp(c).WriteJsonExit()
		return
	}
}

func ServerEdit(c *gin.Context) {
	if c.Request.Method == "GET" {
		id := c.Query("id")
		s, _ := service.FindServerById(gconv.Int(id))
		c.HTML(http.StatusOK, "server_edit.html", gin.H{"s": s})
	} else {
		var f request.TaskServerForm
		if err := c.ShouldBind(&f); err != nil {
			response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
			return
		}
		err := service.ServerEdit(f)
		if err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
			return
		}
		response.SuccessResp(c).WriteJsonExit()
		return
	}
}

func ServerDel(c *gin.Context) {
	id := c.PostForm("id")
	err := service.ServerDel(gconv.Int(id))
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).WriteJsonExit()
	return
}

func DownImplants(c *gin.Context) {
	filename := c.Query("implants_name")

	fileContents, err := service.DownImplants(filename)
	if err != nil {
		response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).WriteJsonExit()
	}

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", "attachment; filename="+filename) // 用来指定下载下来的文件名
	c.Header("Content-Transfer-Encoding", "binary")

	c.Data(http.StatusOK, "application/octet-stream", fileContents)
}
