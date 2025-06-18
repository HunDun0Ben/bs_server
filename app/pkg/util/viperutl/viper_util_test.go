package viperutl_test

import (
	"testing"

	"github.com/HunDun0Ben/bs_server/app/pkg/conf"
	"github.com/HunDun0Ben/bs_server/app/pkg/util/viperutl"
)

func TestViper(t *testing.T) {
	viperutl.PrintViperSetting(conf.GlobalViper)
}
