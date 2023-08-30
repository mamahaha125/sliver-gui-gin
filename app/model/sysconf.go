package model

import "time"

type SysConf struct {
	ID        uint      `json:"id"`
	Type      SysType   `json:"type" gorm:"unique"`
	Info      string    `json:"info" gorm:"type:varchar(4000);"`
	Status    uint8     `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SysType uint

const (
	SysSiteConf SysType = 1
	SysMailConf SysType = 2
)

type SiteConf struct {
	WebName     string `json:"web_name"`
	WebUrl      string `json:"web_url"`
	LogoUrl     string `json:"logo_url"`
	KeyWords    string `json:"key_words"`
	Description string `json:"description"`
	Copyright   string `json:"copyright"`
	Icp         string `json:"icp"`
	SiteStatus  uint8  `json:"site_status"`
}

type MailConf struct {
	EmailName   string `json:"email_name"`
	EmailHost   string `json:"email_host"`
	EmailPort   string `json:"email_port"`
	EmailUser   string `json:"email_user"`
	EmailPwd    string `json:"email_pwd"`
	EmailStatus int    `json:"email_status"`
}

type MailTest struct {
	EmailTest      string `json:"email_test"`
	EmailTestTitle string `json:"email_test_title"`
	EmailTemplate  string `json:"email_template"`
}
