package model

import "time"

type Auth struct {
	ID        uint
	AuthName  string
	AuthUrl   string
	UserId    int
	Pid       int
	Sort      int
	Icon      string
	IsShow    int
	Status    int
	PowerType int
	CreateId  int
	UpdateId  int
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NodesResp struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Open bool   `json:"open"`
	Pid  int    `json:"pId"`
}

type NodeResp struct {
	ID        int    `json:"id"`
	Pid       int    `json:"pid"`
	AuthName  string `json:"auth_name"`
	AuthUrl   string `json:"auth_url"`
	Sort      int    `json:"sort"`
	IsShow    int    `json:"is_show"`
	Icon      string `json:"icon"`
	PowerType int    `json:"power_type"`
}

type AuthResp struct {
	NodesResp
	PowerType int `json:"power_type"`
}

type FrontAuthResp struct {
	PowerID   string `json:"powerId"`
	PowerName string `json:"powerName"`
	PowerType string `json:"powerType"`
	PowerCode string `json:"powerCode"`
	PowerURL  string `json:"powerUrl"`
	OpenType  string `json:"openType"`
	ParentID  string `json:"parentId"`
	Icon      string `json:"icon"`
	Sort      int    `json:"sort"`
	CheckArr  string `json:"checkArr"`
}

type AuthListResp struct {
	CateId   int    `json:"cate_id"`
	CateName string `json:"cate_name"`
	MenuId   int    `json:"menu_id"`
	MenuName string `json:"menu_name"`
	MenuUrl  string `json:"menu_url"`
}

type RolePower struct {
	Status struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"status"`
	Data []RolePowerData `json:"data"`
}

type RolePowerData struct {
	CheckArr   string      `json:"checkArr,omitempty"`
	CreateTime interface{} `json:"create_time,omitempty"`
	Enable     int         `json:"enable,omitempty"`
	Icon       string      `json:"icon,omitempty"`
	OpenType   interface{} `json:"openType,omitempty"`
	ParentID   string      `json:"parentId"`
	PowerID    string      `json:"powerId"`
	PowerName  string      `json:"powerName"`
	PowerType  string      `json:"powerType,omitempty"`
	PowerURL   interface{} `json:"powerUrl,omitempty"`
	Sort       int         `json:"sort,omitempty"`
	UpdateTime interface{} `json:"update_time,omitempty"`
}
