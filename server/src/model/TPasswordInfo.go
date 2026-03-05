package model

import (
	"gitee.com/ouhaoqiang/passwordserver/server/utils/checkpassword"
	"gorm.io/gorm"
)

type TPasswordInfo struct {
	gorm.Model

	// 密码记录id
	ID string `gorm:"uniqueIndex;not null"`
	// 创建者id
	UserId string `gorm:"primaryKey;not null;comment:用户id"`

	// 应用名称
	AppName string `gorm:"primaryKey;comment:应用名称"`

	// 账户类型,一个账户可以有多种类型
	AccountType string `gorm:"primaryKey;comment:账户类型,一个账户可以有多种类型"`

	// 账户
	Account string `gorm:"primaryKey;comment:账户"`

	// 密码
	Password string `gorm:"comment:密码"`

	// url
	Url string `gorm:"default:'';comment:url"`

	// 密码强度
	PasswordStrength checkpassword.PasswordStrength `gorm:"comment:密码强度"`

	// 备注
	Notes string `gorm:"comment:备注"`
}

func (T *TPasswordInfo) TableName() string {
	return "t_password_info"
}
