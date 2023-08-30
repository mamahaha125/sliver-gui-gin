package controller

import (
	"fmt"
	"github.com/bishopfox/sliver/client/assets"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	rpcs "pear-admin-go/app/core/rpc"
	"pear-admin-go/app/global/e"
	response2 "pear-admin-go/app/global/response"
	"strconv"

	"path/filepath"
	"pear-admin-go/app/model"
	GuiModels "pear-admin-go/app/model"
	"pear-admin-go/app/service"
	"sort"
)

// @Title ImportConf
// @Description TODO
// @Author mamahaha125 2023-08-24 06:08:30
// @Update mamahaha125 2023-08-24 06:08:30
// @Param c
func ImportConf(c *gin.Context) {
	FileDir := assets.GetConfigDir()
	FileHandler, _ := c.FormFile("file")
	fmt.Println(FileHandler)
	fliepath := filepath.Join(FileDir, FileHandler.Filename)
	fmt.Println(fliepath)
	if err := c.SaveUploadedFile(FileHandler, fliepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully!",
	})
}

// 查询本地客户端cfg文件
func SearchConf(c *gin.Context) {
	conf := GuiModels.UserBasic{}
	ClientCfg := assets.GetConfigs()
	conf.ConfigDir = assets.GetConfigDir()

	for k := range ClientCfg {
		conf.ConnectRpc = append(conf.ConnectRpc, k)
	}
	sort.Strings(conf.ConnectRpc)
	c.JSON(http.StatusOK, gin.H{
		"message": "Client's config!",
		"client":  conf.ConnectRpc,
	})
}

// @Title ChooseConf
// @Description TODO
// @Author mamahaha125 2023-08-24 06:05:13
// @Update mamahaha125 2023-08-24 06:05:13
// @Param c
func ChooseConf(c *gin.Context) {
	clientName := c.PostForm("client")
	service.ConnectRpc(clientName)
}

// @Title GetSession
// @Description TODO
// @Author mamahaha125 2023-08-24 06:05:05
// @Update mamahaha125 2023-08-24 06:05:05
// @Param c
func GetSession(c *gin.Context) {
	con := rpcs.GetInstance().GetCon()
	nums, err := service.GetSessions(con)

	if err != nil {
		response2.ErrorResp(c).SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
	}
	if len(nums) != 0 {
		response2.ErrorResp(c).SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
	}
	response2.SuccessResp(c).SetMsg("ok").SetType(model.OperOther).SetData(len(nums)).Log(e.LoginHandler, nil).WriteJsonExit()
}

// @Title GetJob
// @Description 获取监听器
// @Author mamahaha125 2023-08-24 06:05:01
// @Update mamahaha125 2023-08-24 06:05:01
// @Param c
func GetJob(c *gin.Context) {
	con := rpcs.GetInstance().GetCon()
	jobs := service.GetJobs(con)
	if jobs.Active == nil {
		response2.ErrorResp(c).SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
	}
	response2.SuccessResp(c).SetMsg("ok").SetType(model.OperOther).SetData(len(jobs.Active)).Log(e.LoginHandler, nil).WriteJsonExit()

}

// @Title Build
// @Description 编译生成木马
// @Author mamahaha125 2023-08-24 06:04:57
// @Update mamahaha125 2023-08-24 06:04:57
// @Param c
func Build(c *gin.Context) {
	Conf := model.BuildConf{
		Name:   "",
		Os:     "",
		Arch:   "",
		Dns:    "",
		Http:   "",
		Mtls:   "",
		Format: "",
	}
	if err := c.ShouldBindBodyWith(&Conf, binding.JSON); err != nil {
		response2.ErrorResp(c).SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
	}
	_, err := service.Build(Conf)
	if err != nil {
		response2.ErrorResp(c).SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()
	}
	response2.SuccessResp(c).SetMsg("ok").SetType(model.OperOther).Log(e.LoginHandler, nil).WriteJsonExit()

}

// @Title GetImplants
// @Description 获取服务器已生成木马
// @Author mamahaha125 2023-08-24 06:04:53
// @Update mamahaha125 2023-08-24 06:04:53
// @Param c
func GetImplants(c *gin.Context) {
	c.HTML(http.StatusOK, "server_list.html", nil)
	//Con := rpcs.GetInstance()
	//service.GetBuilds(Con.GetCon())

}

// @Title MtlsListening
// @Description 启动mtls监听
// @Author mamahaha125 2023-08-24 06:08:36
// @Update mamahaha125 2023-08-24 06:08:36
// @Param c
func MtlsListening(c *gin.Context) {
	var conf struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}

	if err := c.ShouldBindBodyWith(&conf, binding.JSON); err != nil {
		response2.ErrorResp(c).SetMsg("启动失败")
	}
	ports, err := strconv.ParseUint(conf.Port, 10, 32)
	if err != nil {
		response2.ErrorResp(c).SetMsg("正输入正确端口")
	}
	msg, e := service.MtlsListen(conf.Host, uint32(ports), false)
	if e != nil {
		response2.ErrorResp(c).SetData(msg).SetMsg("启动失败")
	}
	response2.SuccessResp(c).SetMsg("启动成功")
}

func WebSokcetListen(c *gin.Context) {
	service.WebSocket(c.Writer, c.Request)
	response2.SuccessResp(c).SetMsg("启动成功")
}
