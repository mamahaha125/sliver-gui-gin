package dao

import (
	"pear-admin-go/app/core/db"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/gconv"
	"strings"
)

type TaskLogDao interface {
	Insert(model.TaskLog) error
	Update(model.TaskLog, map[string]interface{}) error
	Delete(int) error
	FindByPage(pageNum, limit int, filters ...interface{}) ([]model.TaskLog, int, error)
	CountByTaskId(int) (int, error)
}

func NewTaskLogDaoImpl() TaskLogDao {
	return &TaskLogDaoImpl{}
}

type TaskLogDaoImpl struct {
}

func (this *TaskLogDaoImpl) Insert(data model.TaskLog) error {
	err := db.Instance().Create(&data).Error
	return err
}

func (this *TaskLogDaoImpl) Delete(id int) error {
	err := db.Instance().Where("id = ?", id).Delete(model.TaskLog{}).Error
	return err
}

func (this *TaskLogDaoImpl) Update(info model.TaskLog, updateAttrMap map[string]interface{}) error {
	err := db.Instance().Model(&info).Where("id = ?", info.Id).Updates(updateAttrMap).Error
	return err
}

func (this *TaskLogDaoImpl) FindByPage(pageNum, limit int, filters ...interface{}) (infos []model.TaskLog, count int, err error) {
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

	query := client.Model(model.TaskLog{})
	query.Where(strings.Join(queryArr, " AND "), values...).Count(&count)
	err = query.Where(strings.Join(queryArr, " AND "), values...).Order("id desc").Limit(limit).Offset(offset).Find(&infos).Error
	return
}

func (this *TaskLogDaoImpl) CountByTaskId(taskId int) (int, error) {
	var count int
	err := db.Instance().Model(model.TaskLog{}).Where("task_id = ?", taskId).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
