package e

import "time"

const (
	DefaultAvatar  = "/static/admin/images/avatar.jpg"
	DefUploadSize  = 2 * 1024 * 1024 // 默认最大上传
	DefaultSaltLen = 10

	AllowAuth = "/system,/system/index,/system/main,/system/pear,/system/menu" // 不需要验证的地址放在这里

	// 日期格式化
	TimeFormatDay = "20060102"
	TimeFormat    = "2006-01-02 15:04:05"

	// 缓存KEY
	Menu            = "menu_list"
	AdminInfo       = "admin_info"
	Auth            = "auth"
	UserLoginErr    = "user_pwd_err_"
	UserLock        = "user_lock_"
	MaxErrNum       = 5
	MenuCache       = "menu_cache"
	AuthList        = "auth_list"
	TestMailConf    = "test_mail_conf"
	TestMailEffTime = time.Hour * 48
	PearConfigCache = "pear_cache"

	// log title
	AdminAddHandler  = "新增管理员"
	AdminEditHandler = "修改管理员信息"
	LoginHandler     = "登陆"
	AdminDelete      = "删除管理员"

	AuthDelete   = "删除节点"
	AuthNodeEdit = "修改节点"
	AuthNodeAdd  = "新增节点"
	AuthNode     = "权限配置"

	RoleEditHandler   = "修改角色权限"
	RoleAddHandler    = "新增角色权限"
	RoleDeleteHandler = "删除角色权限"
	RoleSave          = "权限分配"

	DefaultUpload = "上传图片"

	SiteEdit = "更新站点设置"
	MailEdit = "更新邮件配置"

	AvatarEdit     = "更新头像"
	ProfileEdit    = "更新资料"
	PwdEditHandler = "修改密码"

	ImgHeight    = 80
	ImgWidth     = 240
	ImgKeyLength = 4

	// task
	Local    = 1
	Remote   = 2
	MaxPool  = 10
	ToLocal  = 1
	ToRemote = 2
)
