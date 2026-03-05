package casdoor

import (
	"fmt"
	"time"

	"gitee.com/ouhaoqiang/passwordserver/server/utils/config"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type builtinApplicationRow struct {
	Owner        string `gorm:"column:owner"`
	Name         string `gorm:"column:name"`
	Organization string `gorm:"column:organization"`
	ClientID     string `gorm:"column:client_id"`
	ClientSecret string `gorm:"column:client_secret"`
	CertName     string `gorm:"column:cert"`
}

type builtinCertRow struct {
	Certificate string `gorm:"column:certificate"`
}

func shouldLoadBuiltinConfig() bool {
	// 源码仓库默认模板没有填写 organization / application，
	// 在 docker-compose 一键部署场景下自动读取 Casdoor 内置 app-built-in 的配置。
	c := config.Config.Server.Casdoor
	return c.Organization == "" || c.Application == ""
}

func loadBuiltinCasdoorConfig() error {
	pg := config.Config.Server.Postgres
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
		pg.Host,
		pg.User,
		pg.Pass,
		"casdoor",
		pg.Port,
		pg.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	var app builtinApplicationRow
	if err := db.Table("application").
		Select("owner", "name", "organization", "client_id", "client_secret", "cert").
		Where("owner = ? AND name = ?", "admin", "app-built-in").
		First(&app).Error; err != nil {
		return err
	}

	var cert builtinCertRow
	if err := db.Table("cert").
		Select("certificate").
		Where("owner = ? AND name = ?", app.Owner, app.CertName).
		First(&cert).Error; err != nil {
		return err
	}

	c := &config.Config.Server.Casdoor
	c.ClientID = app.ClientID
	c.ClientSecret = app.ClientSecret
	c.Certificate = cert.Certificate
	c.Organization = app.Organization
	c.Application = app.Name

	fmt.Printf("casdoor built-in config loaded: organization=%s application=%s clientId=%s\n", c.Organization, c.Application, c.ClientID)
	return nil
}

func InitCasdoorConfig() {
	if shouldLoadBuiltinConfig() {
		var lastErr error
		for i := 0; i < 60; i++ {
			if err := loadBuiltinCasdoorConfig(); err == nil {
				lastErr = nil
				break
			} else {
				lastErr = err
				time.Sleep(time.Second)
			}
		}
		if lastErr != nil {
			fmt.Printf("load casdoor built-in config failed, fallback to yaml config: %v\n", lastErr)
		}
	}

	c := config.Config.Server.Casdoor
	casdoorsdk.InitConfig(
		c.Endpoint,
		c.ClientID,
		c.ClientSecret,
		c.Certificate,
		c.Organization,
		c.Application,
	)

}
