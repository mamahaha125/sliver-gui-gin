package service

import (
	"pear-admin-go/app/core/config"
	"pear-admin-go/app/core/db"
	"testing"
)

func init() {
	config.InitConfig("../../config.toml")
	db.InitConn()
}

func TestRunTask(t *testing.T) {
	//task, _ := dao.NewTaskDaoImpl().FindOne(17)
	select {}
}
