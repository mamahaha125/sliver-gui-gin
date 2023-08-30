package dao

import (
	"github.com/jinzhu/gorm"
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
	"strings"
)

type TaskServerDao interface {
	Insert(server model.TaskServer) error
	FindOne(serverID int) (server *model.TaskServer, err error)
	Update(server model.TaskServer, attr map[string]interface{}) error
	Delete(server model.TaskServer) error
	FindServer(k, v string) (server *model.TaskServer, err error)
	FindServers(k, v string) (servers []model.TaskServer, err error)
	FindByPage(pageNum, limit int, filters ...interface{}) (servers []model.TaskServer, count int, err error)
}

func NewTaskServerDaoImpl() TaskServerDao {
	s := new(TaskServerDaoImpl)
	return s
}

type TaskServerDaoImpl struct {
}

func (r *TaskServerDaoImpl) FindByPage(pageNum, limit int, filters ...interface{}) (servers []model.TaskServer, count int, err error) {
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
	query := db.Instance().Model(model.TaskServer{})
	query.Where(strings.Join(queryArr, " AND "), values...).Count(&count)
	err = query.Where(strings.Join(queryArr, " AND "), values...).Order("id desc").Limit(limit).Offset(offset).Find(&servers).Error
	return
}

func (r *TaskServerDaoImpl) Insert(server model.TaskServer) error {
	err := db.Instance().Create(&server).Error
	return err
}

func (r *TaskServerDaoImpl) FindOne(serverID int) (*model.TaskServer, error) {
	var server model.TaskServer
	err := db.Instance().Model(model.TaskServer{}).First(&server, serverID).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &server, err
}

func (r *TaskServerDaoImpl) Update(server model.TaskServer, attr map[string]interface{}) error {
	err := db.Instance().Model(&server).Omit("id").Updates(attr).Error
	return err
}

func (r *TaskServerDaoImpl) Delete(server model.TaskServer) error {
	err := db.Instance().Delete(&server).Error
	return err
}

func (r *TaskServerDaoImpl) FindServer(k, v string) (server *model.TaskServer, err error) {
	err = db.Instance().Model(model.TaskServer{}).Where(k, v).First(&server).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return server, nil
}

func (r *TaskServerDaoImpl) FindServers(k, v string) (servers []model.TaskServer, err error) {
	err = db.Instance().Model(model.TaskServer{}).Where(k, v).Find(&servers).Error
	return
}
