package conf_test

import (
	"demo/common/conf"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestLoadAllConfig(t *testing.T) {
	conf.InitConfig()
	slog.Info("", "Viper Config Map", conf.GlobalViper.AllSettings())
	json, _ := json.Marshal(conf.GlobalViper.AllSettings())
	slog.Info("Viper Config JSON = ", "", string(json))
}
