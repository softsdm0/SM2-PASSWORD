package model

import (
	"gorm.io/gorm"
)

type TSm2Info struct {
	gorm.Model
	PublicKey  string `gorm:"not null;comment:公钥"`
	PrivateKey string `gorm:"not null;comment:私钥"`
}

func (T *TSm2Info) TableName() string {
	return "t_sm2_info"
}
