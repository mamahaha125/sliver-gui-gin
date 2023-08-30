package request

import "pear-admin-go/app/model"

type LoginForm struct {
	UserName string `json:"username" form:"username" binding:"required,min=3,max=30" zh:"用户名"`
	Password string `json:"password" form:"password" binding:"required,min=3,max=30" zh:"密码"`
}

type LayerListForm struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type RoleForm struct {
	LayerListForm
	ID       string `json:"id" form:"id"`
	RoleName string `json:"role_name" form:"role_name"`
	Detail   string `json:"detail" form:"detail"`
}

type AdminForm struct {
	LayerListForm
	ID        string `json:"id" form:"id"`
	LoginName string `json:"login_name" form:"login_name"`
	RealName  string `json:"real_name" form:"real_name"`
	Phone     string `json:"phone" form:"phone"`
	Email     string `json:"email" form:"email"`
}

type AdminEditForm struct {
	ID        string `json:"id" form:"id" binding:"required" zh:"用户ID"`
	LoginName string `json:"login_name" form:"login_name"`
	RealName  string `json:"real_name" form:"real_name"`
	Phone     string `json:"phone" form:"phone"`
	Email     string `json:"email" form:"email"`
	Status    int
	RoleIds   string
}

type AdminAddForm struct {
	LoginName string `json:"login_name" form:"login_name"`
	Password  string `json:"password" form:"password"`
	RealName  string `json:"real_name" form:"real_name"`
	Phone     string `json:"phone" form:"phone"`
	Email     string `json:"email" form:"email"`
	Status    int
	RoleIds   string
}

type AuthNodeReq struct {
	ID        string `json:"id" form:"id" column:"id"`
	Pid       string `json:"parentId" form:"parentId"  column:"pid"`
	AuthName  string `json:"auth_name" form:"auth_name"  column:"auth_name"`
	AuthUrl   string `json:"auth_url" form:"auth_url"  column:"auth_url"`
	PowerType string `json:"power_type" form:"power_type"  column:"power_type"`
	Sort      string `json:"sort" form:"sort"  column:"sort"`
	IsShow    string `json:"is_show" form:"is_show"  column:"is_show"`
	Icon      string `json:"icon" form:"icon"  column:"icon"`
}

type RoleEditForm struct {
	ID       string `json:"id" form:"id" binding:"required" zh:"权限ID"`
	RoleName string `json:"role_name" form:"role_name" binding:"required" zh:"权限名称"`
	Detail   string `json:"detail" form:"detail"`
	Status   string `json:"status" form:"status"`
	//NodesData string `json:"nodes_data" form:"nodes_data"`
}

type RoleAddForm struct {
	RoleName  string `json:"role_name" form:"role_name" binding:"required" zh:"权限名称"`
	Detail    string `json:"detail" form:"detail"`
	NodesData string `json:"nodes_data" form:"nodes_data"`
}

type SiteConfForm struct {
	ID          uint   `json:"-" form:"id"`
	WebName     string `json:"web_name" form:"web_name"`
	WebUrl      string `json:"web_url" form:"web_url"`
	LogoUrl     string `json:"logo_url" form:"logo_url"`
	KeyWords    string `json:"key_words" form:"key_words"`
	Description string `json:"description" form:"description"`
	Copyright   string `json:"copyright" form:"copyright"`
	Icp         string `json:"icp" form:"icp"`
	SiteStatus  uint8  `json:"site_status" form:"site_status"`
}

type MailConfForm struct {
	ID             uint   `json:"-" form:"id"`
	EmailName      string `json:"email_name" form:"email_name"`
	EmailHost      string `json:"email_host" form:"email_host"`
	EmailPort      string `json:"email_port" form:"email_port"`
	EmailUser      string `json:"email_user" form:"email_user"`
	EmailPwd       string `json:"email_pwd" form:"email_pwd"`
	EmailTest      string `json:"email_test" form:"email_test"`
	EmailTestTitle string `json:"email_test_title" form:"email_test_title"`
	EmailTemplate  string `json:"email_template" form:"email_template"`
	EmailStatus    int    `json:"email_status"`
}

type AvatarForm struct {
	Avatar string `json:"avatar" form:"avatar" binding:"required" zh:"头像"`
}

type ProfileForm struct {
	RealName string `json:"real_name" form:"real_name"`
	Email    string `json:"email" form:"email"`
	Phone    string `json:"phone" form:"phone"`
	Remark   string `json:"remark" form:"remark"`
}

type PasswordForm struct {
	OldPwd     string `json:"old_pwd" form:"old_pwd" binding:"required" zh:"原密码"`
	NewPwd     string `json:"new_pwd" form:"new_pwd" binding:"required" zh:"新密码"`
	ConfirmPwd string `json:"confirm_pwd" form:"confirm_pwd" binding:"required" zh:"重复密码"`
}

type BookMarkForm struct {
	LayerListForm
	SiteName string `json:"site_name"`
	UrlName  string `json:"url_name"`
	Menu     string `json:"menu"`
}

type BookMarkAddForm struct {
	Menu     string `json:"menu" form:"menu" `           // 收藏夹栏目名
	UrlName  string `json:"url_name" form:"url_name" `   // 收藏网址名称
	Url      string `json:"url" form:"url" `             // 收藏网址
	SiteName string `json:"site_name" form:"site_name" ` // 网站名称
	SiteUrl  string `json:"site_url" form:"site_url" `   // 网站地址
	AddDate  string `json:"add_date" form:"add_date" `   // 收藏日期
}

type BookMarkEditForm struct {
	ID uint `json:"id" form:"id" binding:"required"`
	BookMarkAddForm
}

type BookMarkSiteForm struct {
	LayerListForm
	Url         string `json:"url"`
	Title       string `json:"title"`
	Keywords    string `json:"keywords"`
	Description string `json:"description"`
}

type BookMarkSiteEditForm struct {
	ID          uint   `json:"id" form:"id" binding:"required"`
	Url         string `json:"url" form:"url"`
	Title       string `json:"title" form:"title"`
	Keywords    string `json:"keywords" form:"keywords"`
	Description string `json:"description" form:"description"`
}

type WebScreenForm struct {
	LayerListForm
	Url   string `json:"url"`
	Title string `json:"title"`
}

type SpiderAddForm struct {
	Ulr    string `json:"ulr" form:"url"`
	ByYear int    `json:"by_year" form:"by_year"`
	Down   int    `json:"down" form:"down"`
}

type SpiderListForm struct {
	LayerListForm
	UserName string `json:"user_name" form:"user_name"`
	Info     string `json:"info" form:"info"`
	Url      string `json:"url" form:"url"`
}

type TaskForm struct {
	model.TaskCommon
	CreateTime string
}

type TaskServerForm struct {
	model.ServerCommon
	CreateTime string `json:"create_time,omitempty"`
	UpdateTime string `json:"update_time,omitempty"`
}

type TaskServerPage struct {
	LayerListForm
	ServerName string `json:"server_name,omitempty" form:"server_name" zh:"服务器名称"`
	ServerIp   string `json:"server_ip,omitempty" form:"server_ip" zh:"服务器IP"`
	Detail     string `json:"detail,omitempty" form:"detail"`
}

type TaskPage struct {
	LayerListForm
}
