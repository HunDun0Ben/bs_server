package conf_test

import (
	"demo/common/conf"
	"log/slog"
	"testing"
)

func TestLoadAllConfig(t *testing.T) {
	conf.InitConfig()
	slog.Info("Viper map = ", "", conf.GlobalViper.AllSettings())
}
