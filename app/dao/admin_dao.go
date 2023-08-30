package dao

import (
	"pear-admin-go/app/core/db"
	"strings"

	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
)

type AdminDao interface {
	Insert(admin model.Admin) (adminID uint, err error)
	FindAdmin(k, v string) (admin model.Admin, err error)
	FindByName(name string) (admin model.Admin, err error)
	Update(admin model.Admin, updateAttrMap map[string]interface{}) error
	Delete(id int) error
	FindByPage(pageNum, limit int, filters ...interface{}) (admins []model.Admin, count int, err error)
}

func NewAdminDaoImpl() AdminDao {
	admin := new(AdminDaoImpl)
	return admin
}

type AdminDaoImpl struct {
}

func (a *AdminDaoImpl) Insert(admin model.Admin) (adminID uint, err error) {
	err = db.Instance().Create(&admin).Error
	return admin.ID, err
}

func (a *AdminDaoImpl) FindAdmin(k, v string) (admin model.Admin, err error) {
	db.Instance().Model(model.Admin{}).Where(k, v).First(&admin)
	return admin, nil
}

func (a *AdminDaoImpl) FindByName(name string) (admin model.Admin, err error) {
	db.Instance().Model(model.Admin{}).Where("login_name = ?", name).First(&admin)
	return admin, nil
}

// 有多个条件筛选时 满足其中一个即可
func (a *AdminDaoImpl) FindByPage(pageNum, limit int, filters ...interface{}) (admins []model.Admin, count int, err error) {
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
	query := db.Instance().Model(model.Admin{})
	query.Where(strings.Join(queryArr, " OR "), values...).Count(&count)
	err = query.Where(strings.Join(queryArr, " OR "), values...).Order("id desc").Limit(limit).Offset(offset).Find(&admins).Error
	return
}

func (a *AdminDaoImpl) Update(admin model.Admin, updateAttrMap map[string]interface{}) error {
	err := db.Instance().Model(&admin).Where("id = ?", admin.ID).Updates(updateAttrMap).Error
	return err
}

func (a *AdminDaoImpl) Delete(id int) error {
	err := db.Instance().Where("id = ?", id).Delete(model.Admin{}).Error
	return err
}
