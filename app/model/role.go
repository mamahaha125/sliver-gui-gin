package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Role struct {
	ID        uint       `json:"id"`
	RoleName  string     `json:"role_name"`
	Detail    string     `json:"detail"`
	Status    int        `json:"status"`
	CreateId  int        `json:"create_id"`
	UpdateId  int        `json:"update_id"`
	RoleAuths []RoleAuth `json:"-" gorm:"foreignkey:RoleID"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (r *Role) BeforeCreate(scope *gorm.Scope) {
	scope.SetColumn("Status", 1)
}

func (r *Role) BeforeUpdate(scope *gorm.Scope) {
}

// 管理员信息修改页面展示
type RoleEditShow struct {
	ID       int
	RoleName string
	Status   int
	Checked  int
}
