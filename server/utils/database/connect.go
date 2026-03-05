package database

import (
	"fmt"

	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ConnectDB connect to db
func ConnectDB() {
	c := config.Config.Server.Postgres

	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s", c.Host, c.User, c.Pass, c.DbName, c.Port, c.TimeZone)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	DB.AutoMigrate(
		&model.TPasswordInfo{},
		&model.TSm2Info{},
		&model.TGetPasswordRecords{},
	)
	fmt.Println("Database Migrated")

	fmt.Println("init insert sql")
	// 执行一些初始化sql插入
	if err := model.InsertInitial(DB); err != nil {
		panic(err.Error())
	}
	fmt.Println("init insert sql ok")
}
