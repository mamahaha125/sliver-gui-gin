package service

import (
	"errors"
	"fmt"
	"pear-admin-go/app/core/log"
	"pear-admin-go/app/dao"
	"pear-admin-go/app/util/str"

	"pear-admin-go/app/global/request"
	"pear-admin-go/app/model"
	"time"
)

func ServerAdd(f request.TaskServerForm) error {
	f.Status = 1
	var s model.TaskServer
	err := str.CopyFields(&s, f)
	if err != nil {
		return err
	}
	s.CreateTime = time.Now()
	s.UpdateTime = time.Now()
	err = dao.NewTaskServerDaoImpl().Insert(s)
	if err != nil {
		return err
	}
	return nil
}

func ServerEdit(f request.TaskServerForm) error {
	s, err := dao.NewTaskServerDaoImpl().FindOne(f.Id)
	if err != nil {
		return err
	}
	if s == nil {
		return errors.New(fmt.Sprintf("Id:%d not found", f.Id))
	}
	err = dao.NewTaskServerDaoImpl().Update(*s, map[string]interface{}{
		"server_name":     f.ServerName,
		"server_account":  f.ServerAccount,
		"server_password": f.ServerPassword,
		"server_ip":       f.ServerIp,
		"port":            f.Port,
		"private_key_src": f.PrivateKeySrc,
		"public_key_src":  f.PublicKeySrc,
		"conn_type":       f.ConnType,
		"detail":          f.Detail,
		"update_time":     time.Now(),
	})
	if err != nil {
		return err
	}
	return nil
}

func ServerJson(f request.TaskServerPage) ([]model.TaskServer, int, error) {
	filters := make([]interface{}, 0)
	if f.ServerName != "" {
		filters = append(filters, "server_name LIKE ?", "%"+f.ServerName+"%")
	}
	if f.ServerIp != "" {
		filters = append(filters, "server_ip LIKE ?", "%"+f.ServerIp+"%")
	}
	if f.Detail != "" {
		filters = append(filters, "detail LIKE ?", "%"+f.Detail+"%")
	}
	ts, count, err := dao.NewTaskServerDaoImpl().FindByPage(f.Page, f.Limit, filters...)
	if err != nil {
		return nil, 0, err
	}
	return ts, count, nil
}

func FindServerById(id int) (*model.TaskServer, error) {
	s, err := dao.NewTaskServerDaoImpl().FindOne(id)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func ServerDel(id int) error {
	s, err := dao.NewTaskServerDaoImpl().FindOne(id)
	if err != nil {
		return err
	}
	if s == nil {
		return errors.New(fmt.Sprintf("Id:%d not found", id))
	}
	err = dao.NewTaskServerDaoImpl().Delete(*s)
	if err != nil {
		return err
	}
	return nil
}

func ServerList() ([]model.TaskServer, int, error) {
	t, c, err := dao.NewTaskServerDaoImpl().FindByPage(1, 100)
	if err != nil {
		log.Instance().Error(err.Error())
		return nil, 0, err
	}
	return t, c, nil
}
