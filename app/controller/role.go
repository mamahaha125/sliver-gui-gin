package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pear-admin-go/app/core/cache"
	dao2 "pear-admin-go/app/dao"
	e2 "pear-admin-go/app/global/e"
	"pear-admin-go/app/global/request"
	"pear-admin-go/app/global/response"
	"pear-admin-go/app/model"
	"pear-admin-go/app/service"
	"pear-admin-go/app/util/gconv"
	pkg "pear-admin-go/app/util/str"
	"pear-admin-go/app/util/validate"
	"strconv"
)

func IconShow(c *gin.Context) {
	c.HTML(http.StatusOK, "icon.html", nil)
}

func AdminAdd(c *gin.Context) {
	if c.Request.Method == "GET" {
		var rolesShow []model.RoleEditShow
		roles, err := dao2.NewRoleDaoImpl().FindRoles("status = ?", "1") // 查找全部的分组
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		for _, v := range roles {
			rolesShow = append(rolesShow, model.RoleEditShow{
				ID:       gconv.Int(v.ID),
				RoleName: v.RoleName,
				Status:   v.Status,
			})
		}
		c.HTML(http.StatusOK, "admin_add.html", gin.H{"role": rolesShow})
	} else {
		roles := c.PostFormArray("role_ids")
		roleIds := pkg.Array2Str(roles)
		status := c.PostForm("status")
		if status == "" {
			status = "0"
		} else if status == "on" {
			status = "1"
		}
		var f request.AdminAddForm
		if err := c.ShouldBind(&f); err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperAdd).Log(e2.AdminAddHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
			return
		}
		if err := service.AdminAddHandlerService(roleIds, status, f, c); err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperAdd).Log(e2.AdminAddHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
			return
		}
		response.SuccessResp(c).SetMsg("创建成功!").SetType(model.OperAdd).Log(e2.AdminAddHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
		return
	}
}

func AdminList(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_list.html", gin.H{})
}

func AdminEdit(c *gin.Context) {
	if c.Request.Method == "GET" {
		uid := c.Query("id")
		if uid == "" {
			c.String(http.StatusOK, "请检查参数")
			return
		}
		show, rolesShow, err := service.AdminEditService(uid)
		if err != nil {
			c.String(http.StatusOK, err.Error())
			return
		}
		c.HTML(http.StatusOK, "admin_edit.html", gin.H{"show": show, "role": rolesShow})
	} else {
		roles := c.PostFormArray("role_ids")
		roleIds := pkg.Array2Str(roles)
		status := c.PostForm("status")
		if status == "" {
			status = "0"
		} else if status == "on" {
			status = "1"
		}
		var f request.AdminEditForm
		if err := c.ShouldBind(&f); err != nil {
			response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).SetType(model.OperEdit).Log(e2.AdminEditHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
			return
		}
		f.RoleIds = roleIds
		f.Status = gconv.Int(status)
		if err := service.UpdateAdminAttrService(f); err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e2.AdminEditHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
			return
		}
		response.SuccessResp(c).SetMsg("更新成功").SetType(model.OperEdit).Log(e2.AdminEditHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
		return
	}
}

func AdminStatus(c *gin.Context) {
	user := service.GetProfile(c)
	if service.IsAdmin(user) == false {
		response.ErrorResp(c).SetMsg("权限不足，无法修改").SetType(model.OperEdit).Log(e2.AdminEditHandler, gin.H{"form": c.Request.Form}).WriteJsonExit()
		return
	}
	id := c.PostForm("id")
	status := c.PostForm("status")
	if id == "" || status == "" {
		response.ErrorResp(c).SetMsg("参数不足").SetType(model.OperEdit).Log(e2.AdminEditHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
		return
	}
	err := service.UpdateAdminStatus(id, status)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e2.AdminEditHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetMsg("更新成功").SetType(model.OperEdit).Log(e2.AdminEditHandler, gin.H{"form": c.Request.PostForm}).WriteJsonExit()
	return
}

func AdminJson(c *gin.Context) {
	var f request.AdminForm
	if err := c.ShouldBind(&f); err != nil {
		response.SuccessResp(c).SetCode(0).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	count, data, err := service.AdminListJsonService(f)
	if err != nil {
		response.SuccessResp(c).SetCode(0).SetMsg(err.Error()).SetCount(count).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetCode(0).SetCount(count).SetData(data).WriteJsonExit()
}

func AdminDelete(c *gin.Context) {
	uid := c.Query("id")
	if uid == "" {
		response.ErrorResp(c).SetMsg("id不能为空").SetType(model.OperOther).Log(e2.AdminDelete, gin.H{"uid": uid, "header": c.Request.Header.Get("User-Agent")}).WriteJsonExit()
		return
	}
	if err := service.AdminDeleteService(uid, c); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e2.AdminDelete, gin.H{"uid": uid, "header": c.Request.Header.Get("User-Agent")}).WriteJsonExit()
	}
	response.SuccessResp(c).SetType(model.OperOther).Log(e2.AdminDelete, gin.H{"uid": uid, "header": c.Request.Header.Get("User-Agent")}).WriteJsonExit()
}

func RoleList(c *gin.Context) {
	c.HTML(http.StatusOK, "role_list.html", gin.H{})
}

func RoleJson(c *gin.Context) {
	var f request.RoleForm
	if err := c.ShouldBind(&f); err != nil {
		response.SuccessResp(c).SetCode(0).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	count, data, err := service.RoleListJsonService(f)
	if err != nil {
		response.SuccessResp(c).SetCode(0).SetMsg(err.Error()).SetCount(count).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetCode(0).SetCount(count).SetData(data).WriteJsonExit()
}

func RoleAdd(c *gin.Context) {
	if c.Request.Method == "GET" {
		c.HTML(http.StatusOK, "role_add.html", gin.H{})
	} else {
		var f request.RoleAddForm
		if err := c.ShouldBind(&f); err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperAdd).Log(e2.RoleAddHandler, c.Request.PostForm).WriteJsonExit()
			return
		}
		if err := service.RoleAddHandlerService(f, c); err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperAdd).Log(e2.RoleAddHandler, c.Request.PostForm).WriteJsonExit()
			return
		}
		response.SuccessResp(c).SetType(model.OperAdd).Log(e2.RoleAddHandler, c.Request.PostForm).WriteJsonExit()
	}
}

func RolePower(c *gin.Context) {
	id := c.Query("id")
	c.HTML(http.StatusOK, "role_power.html", gin.H{"id": id})
}

func GetRolePower(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
	data := service.FindAuthPower(id)
	c.JSON(http.StatusOK, data)
}

func SaveRolePower(c *gin.Context) {
	roleId := c.PostForm("roleId")
	powerIds := c.PostForm("powerIds")
	err := service.SaveRoleAuth(roleId, powerIds)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).Log(e2.RoleSave, c.Request.PostForm).WriteJsonExit()
		return
	}
	response.SuccessResp(c).Log(e2.RoleSave, c.Request.PostForm).WriteJsonExit()
}

func RoleEdit(c *gin.Context) {
	if c.Request.Method == "GET" {
		id := c.Query("id")
		role, err := service.RoleEditService(id)

		if err != nil {
			c.String(http.StatusOK, err.Error())
		}
		c.HTML(http.StatusOK, "role_edit.html", gin.H{"role": role})
	} else {
		var f request.RoleEditForm
		if err := c.ShouldBindJSON(&f); err != nil {
			response.ErrorResp(c).SetMsg(validate.GetValidateError(err)).SetType(model.OperEdit).Log(e2.RoleEditHandler, c.Request.Form).WriteJsonExit()
			return
		}
		if err := service.RoleEditHandlerService(f); err != nil {
			response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperEdit).Log(e2.RoleEditHandler, f).WriteJsonExit()
			return
		}
		response.SuccessResp(c).SetType(model.OperEdit).Log(e2.RoleEditHandler, f).WriteJsonExit()
	}
}

func RoleDeleteHandler(c *gin.Context) {
	id := c.PostForm("id")
	if err := service.RoleDeleteHandlerService(id); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e2.RoleDeleteHandler, c.Request.PostForm).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetType(model.OperOther).Log(e2.RoleDeleteHandler, c.Request.PostForm).WriteJsonExit()
}

func AuthList(c *gin.Context) {
	c.HTML(http.StatusOK, "auth_list.html", gin.H{})
}

func AuthEdit(c *gin.Context) {
	if c.Request.Method == "GET" {
		authID := c.Query("id")
		resp, err := service.FindAuthByID(authID)
		if err != nil {
			c.String(http.StatusOK, err.Error())
		}
		firstAuths := service.FindAuthName(0)
		secondAuths := service.FindAuthName(1)
		c.HTML(http.StatusOK, "auth_edit.html", gin.H{"parents": firstAuths, "seconds": secondAuths, "auth": resp})
	} else {
		var req request.AuthNodeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			response.ErrorResp(c).SetType(model.OperOther).SetMsg(err.Error()).Log(e2.AuthNode, nil).WriteJsonExit()
			return
		}
		if req.ID == "" {
			if err := service.AuthInsert(req); err != nil {
				response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e2.AuthNodeAdd, req).WriteJsonExit()
				return
			}
			cache.Instance().Delete(e2.MenuCache + gconv.String(service.GetUid(c))) // 删除栏目列表缓存，重新进行设置
			response.SuccessResp(c).SetType(model.OperEdit).Log(e2.AuthNodeAdd, req).WriteJsonExit()
		} else {
			if err := service.AuthUpdate(req); err != nil {
				response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e2.AuthNodeEdit, req).WriteJsonExit()
				return
			}
			cache.Instance().Delete(e2.MenuCache + gconv.String(service.GetUid(c))) // 删除栏目列表缓存，重新进行设置
			response.SuccessResp(c).SetType(model.OperEdit).Log(e2.AuthNodeEdit, req).WriteJsonExit()
		}
	}
}

func AuthDelete(c *gin.Context) {
	authID := c.PostForm("id")
	if err := service.AuthDelete(authID); err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).SetType(model.OperOther).Log(e2.AuthDelete, c.Request.PostForm).WriteJsonExit()
		return
	}
	cache.Instance().Delete(e2.MenuCache + gconv.String(service.GetUid(c))) // 删除栏目列表缓存，重新进行设置
	response.SuccessResp(c).SetType(model.OperOther).Log(e2.AuthDelete, c.Request.PostForm).WriteJsonExit()
}

func AuthNode(c *gin.Context) {
	authID := c.PostForm("id")
	resp, err := service.FindAuthByID(authID)
	if err != nil {
		response.ErrorResp(c).SetMsg(err.Error()).WriteJsonExit()
		return
	}
	response.SuccessResp(c).SetData(resp).WriteJsonExit()
}

func AuthNodes(c *gin.Context) {
	resp, count := service.FindAuths()
	response.SuccessResp(c).SetCount(gconv.Int(count)).SetData(resp).WriteJsonExit()
}

func AddNode(c *gin.Context) {
	firstAuths := service.FindAuthName(0)
	secondAuths := service.FindAuthName(1)
	c.HTML(http.StatusOK, "auth_add.html", gin.H{"parents": firstAuths, "seconds": secondAuths})
}

func Parent(c *gin.Context) {
	data := service.FindAllPower()
	c.JSON(http.StatusOK, data)
}
