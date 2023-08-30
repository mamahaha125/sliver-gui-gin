package db

import (
	"fmt"
	"github.com/gchaincl/dotsql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"os"
	"pear-admin-go/app/core/config"
	"pear-admin-go/app/core/log"
	"pear-admin-go/app/global/initial"
	"pear-admin-go/app/model"
	"pear-admin-go/app/util/file"
)

var conn *gorm.DB

func Instance() *gorm.DB {
	if conn == nil {
		InitConn()
	}
	return conn
}

func InitConn() {
	switch config.Instance().DB.DBType {
	case "mysql":
		conn = GormMysql()
	case "sqlite":
		conn = GormSqlite()
	default:
		log.Instance().Fatal("No DBType")
	}
}

var (
	db  *gorm.DB
	err error
)

func GormMysql() *gorm.DB {
	m := config.Instance().DB
	if m.DBName == "" {
		return nil
	}
	dsn := m.DBUser + ":" + m.DBPwd + "@tcp(" + m.DBHost + ")/" + m.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		fmt.Println("MySQL启动异常", err.Error())
		os.Exit(0)
	}
	db.DB().SetMaxIdleConns(100)
	db.DB().SetMaxOpenConns(300)
	db.SingularTable(true)
	db.LogMode(true)
	initTables()
	initMyCallbacks()
	return db
}

func GormSqlite() *gorm.DB {
	dbFile := fmt.Sprintf("%s.db", config.Instance().DB.DBName)
	if file.CheckNotExist(dbFile) {
		if err := createDB(dbFile); err != nil {
			log.Instance().Fatal("创建数据库文件失败：" + err.Error())
		}
	}
	db, err = gorm.Open("sqlite3", dbFile)
	if err != nil {
		log.Instance().Fatal("连接数据库失败：" + err.Error())
	}
	db.SingularTable(true)
	db.LogMode(true)
	initTables()
	initMyCallbacks()
	return db
}

func createDB(path string) error {
	fp, err := os.Create(path) // 如果文件已存在，会将文件清空。
	if err != nil {
		return err
	}
	defer fp.Close() //关闭文件，释放资源。
	return nil
}

func initTables() {
	checkTableData(&model.Admin{})
	checkTableData(&model.AdminOnline{})
	checkTableData(&model.Auth{})
	checkTableData(&model.LoginInfo{})
	checkTableData(&model.OperLog{})
	checkTableData(&model.Role{})
	checkTableData(&model.RoleAuth{})
	checkTableData(&model.SysConf{})
	checkTableData(&model.PearConfig{})
	checkTableData(&model.Task{})
	checkTableData(&model.TaskLog{})
	checkTableData(&model.TaskServer{})
}

func checkTableData(tb interface{}) {
	if db.HasTable(tb) == false {
		if config.Instance().DB.DBType == "sqlite" {
			if err := db.Debug().CreateTable(tb).Error; err != nil {
				log.Instance().Fatal("创建数据表失败: " + err.Error())
			}
		} else {
			if err := db.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").CreateTable(tb).Error; err != nil {
				log.Instance().Fatal("创建数据表失败: " + err.Error())
			}
		}
		var sqlName string
		if _, ok := tb.(*model.Admin); ok {
			sqlName = "create-admin"
		} else if _, ok := tb.(*model.Auth); ok {
			sqlName = "create-auth"
		} else if _, ok := tb.(*model.Role); ok {
			sqlName = "create-role"
		} else if _, ok := tb.(*model.RoleAuth); ok {
			sqlName = "create-role-auth"
		} else if _, ok := tb.(*model.PearConfig); ok {
			sqlName = "create-pear-config"
		}
		if sqlName != "" {
			initData(sqlName)
		}
	} else {
		// 已存在的表校验一下是否有新增字段
		if config.Instance().DB.DBType == "sqlite" {
			if err := db.Debug().AutoMigrate(tb).Error; err != nil {
				log.Instance().Fatal("数据库初始化失败: " + err.Error())
			}
		} else {
			if err := db.Debug().Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8").AutoMigrate(tb).Error; err != nil {
				log.Instance().Fatal("数据库初始化失败: " + err.Error())
			}
		}
	}
}

func initData(sqlName string) {
	dot, err := dotsql.LoadFromString(initial.SqlInfo)
	//dot, err := dotsql.LoadFromFile("database/db.sql")
	if err != nil {
		log.Instance().Fatal("无法加载初始数据")
		return
	}
	_, err = dot.Exec(db.DB(), sqlName)
	if err != nil {
		log.Instance().Fatal("执行 " + sqlName + " 失败，" + err.Error())
		return
	}
}

func initMyCallbacks() {
	db.Callback().Create().Replace("gorm:update_time_stamp", model.ForBeforeCreate)
	db.Callback().Update().Replace("gorm:update_time_stamp", model.ForBeforeUpdate)
}
