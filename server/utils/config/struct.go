package config

import "time"

// 配置文件结构体
type Configs struct {
	Server struct {
		Name string `yaml:"name"`
		HTTP struct {
			Bind string `yaml:"bind"`
			Port int    `yaml:"port"`
		} `yaml:"http"`
		Log struct {
			SaveDir   string `yaml:"saveDir"`
			Leve      string `yaml:"leve"`
			CutSizeMb int    `yaml:"cutSizeMb"`
			SaveNum   int    `yaml:"saveNum"`
		} `yaml:"log"`
		Mysql struct {
			Host   string `yaml:"host"`
			Port   int    `yaml:"port"`
			User   string `yaml:"user"`
			Pass   string `yaml:"pass"`
			DbName string `yaml:"dbName"`
		} `yaml:"mysql"`
		Postgres struct {
			Host     string `yaml:"host"`
			Port     int    `yaml:"port"`
			User     string `yaml:"user"`
			Pass     string `yaml:"pass"`
			DbName   string `yaml:"dbName"`
			TimeZone string `yaml:"timeZone"`
		} `yaml:"postgres"`
		Redis struct {
			Host                string        `yaml:"host"`
			Port                int           `yaml:"port"`
			Pass                string        `yaml:"pass"`
			DbName              int           `yaml:"dbName"`
			SessionIdExpiration time.Duration `yaml:"sessionIdExpiration"`
		} `yaml:"redis"`
		Cors struct {
			AllowOrigins     string `yaml:"allowOrigins"`
			AllowHeaders     string `yaml:"allowHeaders"`
			AllowMethods     string `yaml:"allowMethods"`
			AllowCredentials bool   `yaml:"allowCredentials"`
		} `yaml:"cors"`
		Casdoor struct {
			Endpoint     string `yaml:"endpoint"`
			ClientID     string `yaml:"clientId"`
			ClientSecret string `yaml:"clientSecret"`
			Certificate  string `yaml:"certificate"`
			Organization string `yaml:"organization"`
			Application  string `yaml:"application"`
		} `yaml:"casdoor"`
	} `yaml:"server"`
}
