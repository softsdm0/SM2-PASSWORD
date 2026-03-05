package config

import (
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var (
	// 配置文件路径
	ConfigPath string
)

// 配置实例
var Config *Configs = new(Configs)

// 初始化配置
func InitConfig() error {
	// 格式配置路径
	fp := filepath.ToSlash(ConfigPath)
	// 配置文件名称
	f := filepath.Base(fp)
	// 配置所在目录
	dir := filepath.Dir(fp)
	// 配置扩展名
	ext := strings.Split(filepath.Ext(fp), ".")[1]

	viper.SetConfigName(f)
	viper.SetConfigType(ext)
	viper.AddConfigPath(dir)

	// 读取配置
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// 将配置解析到结构体
	err = viper.Unmarshal(Config)
	if err != nil {
		return err
	}
	// viper.Debug()

	return nil
}
