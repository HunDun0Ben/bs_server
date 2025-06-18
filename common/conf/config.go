package conf

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/HunDun0Ben/bs_server/app/conf/confmodel"
)

var (
	GlobalViper *viper.Viper            = viper.New()
	AppViperMap map[string]*viper.Viper = make(map[string]*viper.Viper)
	AppConfig   confmodel.AppConfig

	fs afero.Fs = afero.NewOsFs()
)

// 初始化全局配置.
func init() {
	// 启用环境变量支持覆盖配置文件
	GlobalViper.AutomaticEnv()
	loadAllConfig()
	GlobalViper.Unmarshal(&AppConfig)
}

func loadAllConfig() {
	var app string
	app, ok := os.LookupEnv("GOAPP")
	if !ok {
		app = "./"
	}
	if err := loadConfigFiles(app); err != nil {
		panic(fmt.Errorf("加载配置文件失败: %v", err))
	}
}

// 加载目录下的所有 YAML 文件到对应的 viper 中.
// 并且聚合生成全局的 GlobalViper.
func loadConfigFiles(dir string) error {
	// 遍历目录下的所有文件
	err := afero.Walk(fs, dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			// 只处理 YAML 文件
			if filepath.Ext(path) == ".yaml" || filepath.Ext(path) == ".yml" {
				// 获取文件名作为 viper 的命名空间
				configName := strings.Split(info.Name(), ".")[0]
				// 根据每一个 yaml 文件, 创建一个对应的 viper 实例
				v := viper.New()
				// 设置文件路径并读取配置
				v.SetConfigFile(configName)
				v.SetConfigName(configName)
				// 默认加载 conf 目录下的配置文件
				v.AddConfigPath(filepath.Join(dir, "conf"))
				v.AddConfigPath(dir)
				v.SetConfigType("yaml")
				if err := v.ReadInConfig(); err != nil {
					return fmt.Errorf("读取配置文件 %s 错误: %v", path, err)
				}
				AppViperMap[configName] = v
				// 将此 viper 实例中的配置合并到全局配置中
				// 这里我们使用 Set() 方法将不同的配置内容按照文件名作为键值存储
				GlobalViper.MergeConfigMap(v.AllSettings())
			}
			return nil
		})
	return err
}

// 返回对应配置的 viper 实例
//
// Parameters:
//   - name: 配置文件名
//
// Returns:
//   - *viper.Viper: 配置实例
//   - bool: 对应的配置存在返回 ture, 不存在返回false
func GetConfig(name string) (*viper.Viper, bool) {
	if AppViperMap[name] != nil {
		return AppViperMap[name], true
	}
	return nil, false
}
