package dao

import (
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"

	"strings"
)

type RoleDao interface {
	Insert(role model.Role) (uint, error)
	FindOne(roleID string) (role model.Role, err error)
	Update(role model.Role, attr map[string]interface{}) error
	Delete(role model.Role) error
	FindRole(k, v string) (role model.Role, err error)
	FindRoles(k, v string) (roles []model.Role, err error)
	FindByPage(pageNum, limit int, filters ...interface{}) (admins []model.Role, count int, err error)
}

func NewRoleDaoImpl() RoleDao {
	role := new(RoleDaoImpl)
	return role
}

type RoleDaoImpl struct {
}

func (r *RoleDaoImpl) FindByPage(pageNum, limit int, filters ...interface{}) (roles []model.Role, count int, err error) {
	offset := (pageNum - 1) * limit

	var queryArr []string
	var values []interface{}
	if len(filters) > 0 {
		l := len(filters)
		for k := 0; k < l; k += 2 {
			queryArr = append(queryArr, gconv.String(filters[k]))
			values = append(values, filters[k+1])
		}
	}

	//query := db.Instance().Model(model.Role{}).Where("status = 1")
	query := db.Instance().Model(model.Role{})
	query.Where(strings.Join(queryArr, " AND "), values...).Count(&count)
	err = query.Where(strings.Join(queryArr, " AND "), values...).Order("id desc").Limit(limit).Offset(offset).Find(&roles).Error
	return
}

func (r *RoleDaoImpl) Insert(role model.Role) (uint, error) {
	err := db.Instance().Create(&role).Error
	return role.ID, err
}

func (r *RoleDaoImpl) FindOne(roleID string) (role model.Role, err error) {
	err = db.Instance().Preload("RoleAuths").First(&role, roleID).Error // 预加载关联查询
	//client.First(&role,roleID)
	return role, err
}

func (r *RoleDaoImpl) Update(role model.Role, attr map[string]interface{}) error {
	err := db.Instance().Model(&role).Omit("id").Updates(attr).Error
	return err
}

func (r *RoleDaoImpl) Delete(role model.Role) error {
	err := db.Instance().Delete(&role).Error
	return err
}

func (r *RoleDaoImpl) FindRole(k, v string) (role model.Role, err error) {
	err = db.Instance().Model(model.Role{}).Where(k, v).First(&role).Error
	return
}

func (r *RoleDaoImpl) FindRoles(k, v string) (roles []model.Role, err error) {
	err = db.Instance().Model(model.Role{}).Where(k, v).Find(&roles).Error
	return
}
