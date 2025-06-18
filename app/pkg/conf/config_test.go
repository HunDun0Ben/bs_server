package conf_test

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/HunDun0Ben/bs_server/app/conf/confmodel"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
)

func TestLoadAllConfig(t *testing.T) {
	slog.Info("", "Viper Config Map", conf.GlobalViper.AllSettings())
	json, _ := json.Marshal(conf.GlobalViper.AllSettings())
	slog.Info("Viper Config JSON = ", "", string(json))
}

func TestGetDuration(t *testing.T) {
	slog.Info("get viper Duration configuration.", "expiration", conf.GlobalViper.GetDuration("jwt.expire"))
}

func TestConfig(t *testing.T) {
	var cfg confmodel.AppConfig
	conf.GlobalViper.Unmarshal(&cfg)
	slog.Info("Config", "ServerConfig", cfg)
}
