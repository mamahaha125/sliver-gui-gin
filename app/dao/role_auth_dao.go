package dao

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
	"strings"
)

type RoleAuthDao interface {
	GetByRoleIds(RolesIds string) (authIds string, err error)
	Insert(auth model.RoleAuth, db *gorm.DB) error
	InsertByTx(auths []model.RoleAuth) error
	Delete(k, v string) error
	FindRoleAuthByRoleID(roleId int) ([]model.RoleAuth, error)
}

func NewRoleAuthDaoImpl() RoleAuthDao {
	roleAuth := new(RoleAuthDaoImpl)
	return roleAuth
}

type RoleAuthDaoImpl struct {
}

func (r *RoleAuthDaoImpl) GetByRoleIds(RolesIds string) (authIds string, err error) {
	ids := strings.Split(RolesIds, ",")
	var roleAuths []model.RoleAuth
	err = db.Instance().Model(model.RoleAuth{}).Where("role_id IN (?)", ids).Find(&roleAuths).Error
	if err != nil {
		return "", err
	}
	b := bytes.Buffer{}
	for _, v := range roleAuths {
		if v.AuthID != 0 && v.AuthID != 1 {
			b.WriteString(gconv.String(v.AuthID))
			b.WriteString(",")
		}
	}
	authIds = strings.TrimRight(b.String(), ",")
	return authIds, nil
}

func (r *RoleAuthDaoImpl) InsertByTx(auths []model.RoleAuth) error {
	for _, v := range auths {
		err := db.Instance().Create(&v).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *RoleAuthDaoImpl) Insert(auth model.RoleAuth, conn *gorm.DB) error {
	if conn == nil {
		conn = db.Instance()
	}
	err := conn.Model(model.RoleAuth{}).Where("auth_id = ? and role_id = ?", auth.AuthID, auth.RoleID).FirstOrCreate(&auth).Error
	return err
}

func (r *RoleAuthDaoImpl) Delete(k, v string) error {
	tx := db.Instance().Begin()
	err := tx.Where(k, v).Delete(model.RoleAuth{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error
}

func (r *RoleAuthDaoImpl) FindRoleAuthByRoleID(roleId int) (roles []model.RoleAuth, err error) {
	db := db.Instance()
	err = db.Model(model.RoleAuth{}).Where("role_id = ?", roleId).Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
