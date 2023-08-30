package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"pear-admin-go/app/core/cache"
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/core/log"
	dao2 "pear-admin-go/app/dao"
	"pear-admin-go/app/util/gconv"
	pkg "pear-admin-go/app/util/str"

	e2 "pear-admin-go/app/global/e"
	"pear-admin-go/app/global/request"
	"pear-admin-go/app/model"
	"strings"
)

func BuildRoleFilter(r request.RoleForm) []interface{} {
	filters := make([]interface{}, 0)
	if r.ID != "" {
		filters = append(filters, "id = ?", r.ID)
	}
	if r.RoleName != "" {
		filters = append(filters, "role_name LIKE ?", "%"+r.RoleName+"%")
	}
	if r.Detail != "" {
		filters = append(filters, "detail LIKE ?", "%"+r.Detail+"%")
	}
	return filters
}

func RoleListJsonService(f request.RoleForm) (count int, data []map[string]interface{}, err error) {
	if f.Page == 0 {
		f.Page = 1
	}
	if f.Limit == 0 {
		f.Limit = 10
	}
	filters := BuildRoleFilter(f)
	list, count, err := dao2.NewRoleDaoImpl().FindByPage(f.Page, f.Limit, filters...)
	if err != nil {
		return count, data, err
	}
	for _, v := range list {
		data = append(data, pkg.Struct2MapByTag(model.RoleShow{
			ID:        gconv.Int(v.ID),
			RoleName:  v.RoleName,
			Detail:    v.Detail,
			CreatedAt: v.CreatedAt.Format(e2.TimeFormat),
			UpdatedAt: v.UpdatedAt.Format(e2.TimeFormat),
		}, "json"))
	}
	return count, data, nil
}

func RoleEditService(roleID string) (model.Role, error) {
	role, err := dao2.NewRoleDaoImpl().FindOne(roleID)
	return role, err
}

func RoleAddHandlerService(f request.RoleAddForm, c *gin.Context) error {
	hasRole, err := dao2.NewRoleDaoImpl().FindRole("role_name = ?", f.RoleName)
	if err == nil || hasRole.ID > 0 {
		return errors.New("角色名已存在")
	}
	roleID, err := dao2.NewRoleDaoImpl().Insert(model.Role{
		RoleName: f.RoleName,
		Detail:   f.Detail,
		CreateId: GetUid(c),
	})
	if err != nil {
		return err
	}
	if strings.HasSuffix(f.NodesData, ",") {
		f.NodesData = string([]rune(f.NodesData)[:len(f.NodesData)-1])
	}
	nodesArr := strings.Split(f.NodesData, ",")
	for _, v := range nodesArr {
		err := dao2.NewRoleAuthDaoImpl().Insert(model.RoleAuth{
			RoleID: gconv.Uint64(roleID),
			AuthID: gconv.Uint(v),
		}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func RoleEditHandlerService(f request.RoleEditForm) error {
	role, err := dao2.NewRoleDaoImpl().FindOne(f.ID)
	if err != nil {
		return err
	}
	attr := make(map[string]interface{})
	if f.Status == "" {
		attr["status"] = 0
	} else {
		attr["status"] = 1
	}
	attr["role_name"] = f.RoleName
	attr["detail"] = f.Detail
	if err := dao2.NewRoleDaoImpl().Update(role, attr); err != nil {
		return err
	}
	return nil
}

func RoleDeleteHandlerService(id string) error {
	role, err := dao2.NewRoleDaoImpl().FindOne(id)
	if err != nil {
		return err
	}
	err = dao2.NewRoleDaoImpl().Delete(role)
	if err != nil {
		return err
	}
	err = dao2.NewRoleAuthDaoImpl().Delete("role_id = ?", gconv.String(role.ID))
	return err
}

func SaveRoleAuth(roleId, authIds string) (err error) {
	authIdMap := strings.Split(authIds, ",")
	if len(authIdMap) < 1 {
		return errors.New("权限分配出错")
	}
	if err := dao2.NewRoleAuthDaoImpl().Delete("role_id = ?", gconv.String(roleId)); err != nil {
		return err
	}
	var roleAuth model.RoleAuth
	tx := db.Instance().Begin()
	for _, v := range authIdMap {
		roleAuth.AuthID = gconv.Uint(v)
		roleAuth.RoleID = gconv.Uint64(roleId)
		err = dao2.NewRoleAuthDaoImpl().Insert(roleAuth, tx)
		if err != nil {
			log.Instance().Warn("InsertRoleAuth.Insert error:" + err.Error())
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	deleteMenuCache()
	return nil
}

func deleteMenuCache() {
	items := cache.Instance().Items()
	for k, _ := range items {
		if strings.HasPrefix(k, e2.MenuCache) {
			cache.Instance().Delete(k)
		}
	}
}
