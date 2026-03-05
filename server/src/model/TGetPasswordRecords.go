package model

import (
	"time"
)

type GetPasswordRecordsStatus int

const (
	// 未知错误
	Unknown GetPasswordRecordsStatus = iota
	// 处理成功
	Success
	// 已拒绝
	Refused
	// 记录不存在
	NotFound
)

type TGetPasswordRecords struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	// 请求者id
	UserId string `gorm:"index;not null;comment:用户id"`
	// 密码记录id
	PasswordId   string        `gorm:"not null;comment:密码记录id"`
	PasswordInfo TPasswordInfo `gorm:"foreignKey:PasswordId;references:ID"`
	// 请求状态
	Status GetPasswordRecordsStatus `gorm:"not null;comment:请求状态"`
	// ip地址
	IP string `gorm:"not null;comment:ip地址"`
}

func (T *TGetPasswordRecords) TableName() string {
	return "t_get_password_records"
}
