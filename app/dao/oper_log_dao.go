package dao

import (
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
	"strings"
)

type OperLogDao interface {
	Insert(oper model.OperLog) error
	Delete() error
	FindByPage(pageNum, limit int, filters ...interface{}) (opers []model.OperLog, count int, err error)
}

func NewOperLogDaoImpl() OperLogDao {
	admin := new(OperLogDaoImpl)
	return admin
}

type OperLogDaoImpl struct {
}

func (op *OperLogDaoImpl) Insert(oper model.OperLog) error {
	err := db.Instance().Create(&oper).Error
	return err
}

func (op *OperLogDaoImpl) Delete() error {
	return nil
}

func (op *OperLogDaoImpl) FindByPage(pageNum, limit int, filters ...interface{}) (opers []model.OperLog, count int, err error) {
	client := db.Instance()
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

	query := client.Model(model.OperLog{})
	query.Where(strings.Join(queryArr, " AND "), values...).Count(&count)
	err = query.Where(strings.Join(queryArr, " AND "), values...).Order("id desc").Limit(limit).Offset(offset).Find(&opers).Error
	return
}
