package model

import "time"

type TaskServer struct {
	ServerCommon
	CreateTime time.Time
	UpdateTime time.Time
}

type ServerCommon struct {
	Id             int    `json:"id" form:"id" zh:"ID"`
	ServerName     string `json:"server_name,omitempty" form:"server_name" binding:"required,min=1,max=50" zh:"服务器名称"`
	ServerAccount  string `json:"server_account,omitempty" form:"server_account" binding:"required,min=1,max=50" zh:"用户名"`
	ServerPassword string `json:"server_password" form:"server_password" zh:"密码"`
	ServerIp       string `json:"server_ip,omitempty" form:"server_ip" binding:"required" zh:"服务器IP"`
	Port           int    `json:"port,omitempty" form:"port" binding:"required,min=1" zh:"端口号"`
	PrivateKeySrc  string `json:"private_key_src,omitempty" form:"private_key_src" zh:"私钥地址"` // ssh-keygen -t rsa -f pp_rsa
	PublicKeySrc   string `json:"public_key_src,omitempty" form:"public_key_src" zh:"公钥地址"`
	ConnType       int    `json:"conn_type,omitempty" form:"conn_type" zh:"登陆类型"` // 登陆类型 1=密码登陆 2=密钥登陆
	Detail         string `json:"detail,omitempty" form:"detail"`
	Status         int    `json:"status,omitempty" form:"status"` // 0=初始 1=正常 2=停止
}
