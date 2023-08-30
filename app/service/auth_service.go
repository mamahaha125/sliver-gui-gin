package service

import (
	"errors"
	"fmt"
	"github.com/bishopfox/sliver/client/assets"
	"github.com/bishopfox/sliver/client/transport"
	"github.com/gin-gonic/gin"
	"path/filepath"
	"pear-admin-go/app/core/log"
	dao2 "pear-admin-go/app/dao"
	"pear-admin-go/app/util/gconv"
	pkg "pear-admin-go/app/util/str"
	"sort"

	"pear-admin-go/app/global/request"
	"pear-admin-go/app/model"
	"strings"
)

func AuthInsert(req request.AuthNodeReq) (err error) {
	var auth model.Auth
	auth.Pid = gconv.Int(req.Pid)
	auth.AuthName = req.AuthName
	auth.AuthUrl = req.AuthUrl
	auth.Icon = fmt.Sprintf("layui-icon %s", req.Icon)
	auth.Sort = gconv.Int(req.Sort)
	auth.IsShow = gconv.Int(req.IsShow)
	auth.PowerType = gconv.Int(req.PowerType)
	auth.Status = 1
	_, err = dao2.NewAuthDaoImpl().Insert(auth)
	return err
}

func AuthUpdate(req request.AuthNodeReq) (err error) {
	auth, err := dao2.NewAuthDaoImpl().FindOne(gconv.Int(req.ID))
	if err != nil {
		return err
	}
	if req.IsShow == "" {
		req.IsShow = "0"
	}
	if strings.HasPrefix(req.Icon, "layui-icon ") == false {
		req.Icon = fmt.Sprintf("layui-icon %s", req.Icon)
	}
	err = dao2.NewAuthDaoImpl().Update(auth, pkg.Struct2MapByTag(req, "column"))
	return err
}

func AuthDelete(authID string) error {
	count, err := dao2.NewAuthDaoImpl().FindChildNode(gconv.Int(authID))
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("请先删除子节点")
	}
	err = dao2.NewRoleAuthDaoImpl().Delete("auth_id = ?", authID)
	if err != nil {
		return err
	}
	err = dao2.NewAuthDaoImpl().Delete(gconv.Int(authID))
	return err
}

func FindAuthByID(authID string) (*model.NodeResp, error) {
	auth, err := dao2.NewAuthDaoImpl().FindOne(gconv.Int(authID))
	if err != nil {
		return nil, err
	}
	return &model.NodeResp{
		ID:        gconv.Int(auth.ID),
		Pid:       auth.Pid,
		AuthName:  auth.AuthName,
		AuthUrl:   auth.AuthUrl,
		Sort:      auth.Sort,
		IsShow:    auth.IsShow,
		Icon:      auth.Icon,
		PowerType: auth.PowerType}, nil
}

func FindAuths() ([]model.FrontAuthResp, int64) {
	filters := make([]interface{}, 0)
	filters = append(filters, "status = ?", 1)
	auths, count := dao2.NewAuthDaoImpl().Find(1, 5000, filters...)
	var resp []model.FrontAuthResp
	for _, v := range auths {
		var powerCode string
		if v.IsShow == 0 && v.Pid != 0 {
			powerCode = v.AuthUrl
		}
		resp = append(resp, model.FrontAuthResp{
			PowerID:   gconv.String(v.ID),
			PowerName: v.AuthName,
			PowerType: gconv.String(v.PowerType),
			PowerCode: powerCode,
			PowerURL:  v.AuthUrl,
			OpenType:  "",
			ParentID:  gconv.String(v.Pid),
			Icon:      v.Icon,
			Sort:      v.Sort,
			CheckArr:  "0",
		})
	}
	return resp, count
}

func FindAuthName(powerType int) []model.AuthResp {
	var authNames []model.AuthResp
	filters := make([]interface{}, 0)
	filters = append(filters, "status = ?", 1)
	filters = append(filters, "power_type = ?", powerType)
	auths, _ := dao2.NewAuthDaoImpl().Find(1, 5000, filters...)
	if len(auths) < 1 {
		return authNames
	}
	for _, v := range auths {
		var resp model.AuthResp
		resp.ID = int(v.ID)
		resp.Pid = v.Pid
		resp.Name = v.AuthName
		resp.PowerType = v.PowerType
		authNames = append(authNames, resp)
	}
	return authNames
}

func GetAuthList() []model.AuthListResp {
	authList, err := dao2.NewAuthDaoImpl().AuthList()
	if err != nil {
		log.Instance().Warn(err.Error())
		return nil
	}
	return authList
}

func FindAuthPower(id int) (power model.RolePower) {
	power.Status.Code = 200
	power.Status.Message = "success"

	filters := make([]interface{}, 0)
	filters = append(filters, "status = ?", 1)
	auths, _ := dao2.NewAuthDaoImpl().Find(1, 5000, filters...)
	var (
		powerData []model.RolePowerData
		roleAuth  []string
	)
	ras, err := dao2.NewRoleAuthDaoImpl().FindRoleAuthByRoleID(id)
	if err != nil {
		log.Instance().Warn(err.Error())
		return power
	}
	for _, ra := range ras {
		roleAuth = append(roleAuth, gconv.String(ra.AuthID))
	}
	for _, auth := range auths {
		checkArr := "0"
		if pkg.IsContain(roleAuth, gconv.String(auth.ID)) {
			checkArr = "1"
		}
		openType := ""
		if auth.PowerType == 1 {
			openType = "_iframe"
		}
		powerData = append(powerData, model.RolePowerData{
			CheckArr:   checkArr,
			Enable:     1,
			Icon:       auth.Icon,
			OpenType:   openType,
			ParentID:   gconv.String(auth.Pid),
			PowerID:    gconv.String(auth.ID),
			PowerName:  auth.AuthName,
			PowerType:  gconv.String(auth.PowerType),
			PowerURL:   auth.AuthUrl,
			Sort:       auth.Sort,
			UpdateTime: auth.UpdatedAt,
		})
	}
	power.Data = powerData
	return power
}

func FindAllPower() (power model.RolePower) {
	power.Status.Code = 200
	power.Status.Message = "success"

	filters := make([]interface{}, 0)
	filters = append(filters, "status = ?", 1)
	auths, _ := dao2.NewAuthDaoImpl().Find(1, 5000, filters...)
	var powerData []model.RolePowerData
	powerData = append(powerData, model.RolePowerData{
		ParentID:  "-1",
		PowerID:   "0",
		PowerName: "顶级权限",
	})
	for _, auth := range auths {
		openType := ""
		if auth.PowerType == 1 {
			openType = "_iframe"
		}
		powerData = append(powerData, model.RolePowerData{
			Enable:     1,
			Icon:       auth.Icon,
			OpenType:   openType,
			ParentID:   gconv.String(auth.Pid),
			PowerID:    gconv.String(auth.ID),
			PowerName:  auth.AuthName,
			PowerType:  gconv.String(auth.PowerType),
			PowerURL:   auth.AuthUrl,
			Sort:       auth.Sort,
			CreateTime: auth.CreatedAt,
			UpdateTime: auth.UpdatedAt,
		})
	}
	power.Data = powerData
	return power
}

func ImportCfg(c *gin.Context) bool {
	FileHandler, err := c.FormFile("file")
	if err != nil {
		return false
	}
	fmt.Println(FileHandler.Filename)
	FileDir := assets.GetConfigDir()
	fliepath := filepath.Join(FileDir, FileHandler.Filename)
	fmt.Println(fliepath)
	err = c.SaveUploadedFile(FileHandler, fliepath)
	if err != nil {
		return false
	}
	return true
}

func ChooseConf(rpc string) bool {
	//var rpc rpcpb.SliverRPCClient
	//var ln *grpc.ClientConn
	var err error

	configs := assets.GetConfigs()
	ClientKey := configs[rpc]
	_, _, err = transport.MTLSConnect(ClientKey)
	if err != nil {
		return false
	}
	return true
}

func GetConf(c *gin.Context) []string {
	var keys []string
	configs := assets.GetConfigs()
	for k := range configs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
