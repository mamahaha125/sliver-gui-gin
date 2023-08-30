package model

import "time"

type LoginInfo struct {
	InfoId        int64     `json:"info_id" gorm:"primary_key"`
	LoginName     string    `json:"login_name"`
	IpAddr        string    `json:"ip_addr"`
	LoginLocation string    `json:"login_location"`
	Browser       string    `json:"browser"`
	Os            string    `json:"os"`
	Status        string    `json:"status"`
	Msg           string    `json:"msg"`
	LoginTime     time.Time `json:"login_time"`
}
