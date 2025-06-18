package viperutl

import (
	"log/slog"

	"github.com/spf13/viper"
)

func PrintViperSetting(viper *viper.Viper) {
	settingMap := viper.AllSettings()
	for key, value := range settingMap {
		slog.Info("config", key, value)
	}
}
