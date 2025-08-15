package conf_test

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
	"github.com/HunDun0Ben/bs_server/app/pkg/conf/confmodel"
)

func TestLoadAllConfig(t *testing.T) {
	settings := conf.GlobalViper.AllSettings()
	slog.Info("", "Viper Config Map", settings)
	json, _ := json.Marshal(conf.GlobalViper.AllSettings())
	slog.Info("settings", slog.Any("json", json))
}

func TestGetDuration(t *testing.T) {
	slog.Info("get viper Duration configuration.", "expiration", conf.GlobalViper.GetDuration("jwt.expire"))
}

func TestConfig(t *testing.T) {
	var cfg confmodel.AppConfig
	conf.GlobalViper.Unmarshal(&cfg)
	slog.Info("Config", "ServerConfig", cfg)
}
