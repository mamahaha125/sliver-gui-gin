package dao

import (
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
	"strings"
	"sync"
)

type LoginInfoDao interface {
	Insert(info model.LoginInfo) error
	FindByPage(pageNum, limit int, filters ...interface{}) (info []model.LoginInfo, count int, err error)
}

func NewLoginInfoImpl() LoginInfoDao {
	info := new(LoginInfoDaoImpl)
	return info
}

type LoginInfoDaoImpl struct {
	rw *sync.RWMutex
}

func (l *LoginInfoDaoImpl) FindByPage(pageNum, limit int, filters ...interface{}) (info []model.LoginInfo, count int, err error) {
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
	query := db.Instance().Model(model.LoginInfo{})
	query.Where(strings.Join(queryArr, " AND "), values...).Count(&count)
	err = query.Where(strings.Join(queryArr, " AND "), values...).Order("info_id desc").Limit(limit).Offset(offset).Find(&info).Error
	return
}

func (l *LoginInfoDaoImpl) Insert(info model.LoginInfo) error {
	err := db.Instance().Create(&info).Error
	return err
}
