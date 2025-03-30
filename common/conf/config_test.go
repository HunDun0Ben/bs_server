package conf_test

import (
	"encoding/json"
	"log/slog"
	"testing"

	"github.com/HunDun0Ben/bs_server/common/conf"
)

func TestLoadAllConfig(t *testing.T) {
	slog.Info("", "Viper Config Map", conf.GlobalViper.AllSettings())
	json, _ := json.Marshal(conf.GlobalViper.AllSettings())
	slog.Info("Viper Config JSON = ", "", string(json))
}
