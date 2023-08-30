package config

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"pear-admin-go/app/global/initial"
	"pear-admin-go/app/util/file"
)

var c *conf

func Instance() *conf {
	if c == nil {
		InitConfig("./config.toml")
	}
	return c
}

type conf struct {
	DB     DBConf
	Redis  RedisConf
	App    SettingConf
	Zaplog ZapLogConf
}

type DBConf struct {
	DBType string
	DBUser string
	DBPwd  string
	DBHost string
	DBName string
}

type RedisConf struct {
	RedisAddr string
	RedisPWD  string
	RedisDB   int
}

type SettingConf struct {
	HttpPort    int `json:"http-port"`
	RunMode     string
	PageSize    int
	JwtSecret   string
	ImgSavePath string
	ImgUrlPath  string
}

type ZapLogConf struct {
	Director string ` json:"director"  yaml:"director"`
}

func InitConfig(tomlPath ...string) {
	if len(tomlPath) > 1 {
		log.Fatal("配置路径数量不正确")
	}
	if file.CheckNotExist(tomlPath[0]) {
		err := ioutil.WriteFile(tomlPath[0], []byte(initial.ConfigToml), 0777)
		if err != nil {
			log.Fatal("无法写入配置模板", err.Error())
		}
		log.Fatal("配置文件不存在，已在根目录下生成示例模板文件【config.toml】，请修改后重新启动程序！")
	}
	v := viper.New()
	v.SetConfigFile(tomlPath[0])
	err := v.ReadInConfig()
	if err != nil {
		log.Fatal("配置文件读取失败: ", err.Error())
	}
	err = v.Unmarshal(&c)
	if err != nil {
		log.Fatal("配置解析失败:", err.Error())
	}
}
