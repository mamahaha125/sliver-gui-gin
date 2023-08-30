package service

import (
	"errors"
	"pear-admin-go/app/dao"
	"pear-admin-go/app/global/e"
	"pear-admin-go/app/global/request"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/runtask"
	"pear-admin-go/app/util/str"
	"runtime"
	"strings"
	"time"
)

func TaskJson(f request.TaskPage) ([]model.TaskResp, int, error) {
	ts, count, err := dao.NewTaskDaoImpl().FindByPage(f.Page, f.Limit)
	if err != nil {
		return nil, 0, err
	}
	var resp []model.TaskResp
	for _, v := range ts {
		var r model.TaskResp
		err = str.CopyFields(&r, v)
		if err != nil {
			return nil, 0, err
		}
		r.CreateTime = v.CreateTime.Format(e.TimeFormat)
		count, err := dao.NewTaskLogDaoImpl().CountByTaskId(v.Id)
		if err != nil {
			return nil, 0, err
		}
		r.TaskLogNum = count
		resp = append(resp, r)
	}
	return resp, count, nil
}

func TaskAdd(f request.TaskForm) error {
	if f.SourceType == e.Local && f.DstType == e.Local {
		return errors.New("亲，此模式尚未开发，本地复制请使用复制粘贴功能~")
	}
	var s model.Task
	err := str.CopyFields(&s, f)
	if err != nil {
		return err
	}
	if runtime.GOOS == "windows" {
		f.SourcePath = strings.ReplaceAll(f.SourcePath, "\\", "/")
		f.DstPath = strings.ReplaceAll(f.DstPath, "\\", "/")
	}
	s.CreateTime = time.Now()
	tid, err := dao.NewTaskDaoImpl().Insert(s)
	if err != nil {
		return err
	}
	s.Id = tid
	go runtask.NewRunTask(s).SetSourceClient().SetDstClient().Run()
	return nil
}
