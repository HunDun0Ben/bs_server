package conf

import (
	"log"

	"github.com/spf13/viper"
)

var GlobalViper *viper.Viper

// 初始化全局配置
func InitConfig() {
	GlobalViper = viper.New()

	// 设置配置文件名、路径及格式
	GlobalViper.SetConfigName("application") // 配置文件名，不带扩展名
	GlobalViper.AddConfigPath("./conf")      // 配置文件所在路径
	GlobalViper.SetConfigType("yaml")        // 配置文件类型

	// 尝试读取配置文件
	err := GlobalViper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// 也可以启用环境变量支持
	GlobalViper.AutomaticEnv()
}
