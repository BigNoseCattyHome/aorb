package conf

import (
	"log"

	"github.com/spf13/viper"
)

var AppConfig *viper.Viper

const (
	defaultConfigPath = "conf/"
	defaultConfigName = "config"
	defaultConfigType = "toml"
)

// LoadConfig 在main.go中调用，用于加载配置文件
func LoadConfig() error {
	viper.AddConfigPath(defaultConfigPath)
	viper.SetConfigName(defaultConfigName)
	viper.SetConfigType(defaultConfigType)
	err := viper.ReadInConfig() // 查找并读取配置文件
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
		return err
	}
	AppConfig = viper.GetViper()
	return nil
}

// AppConfig的使用方法
// conf.AppConfig.GetString("db.mongodb_url")
