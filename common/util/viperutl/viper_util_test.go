package viperutl_test

import (
	"testing"

	"github.com/HunDun0Ben/bs_server/common/conf"
	"github.com/HunDun0Ben/bs_server/common/util/viperutl"
)

func TestViper(t *testing.T) {
	viperutl.PrintViperSetting(conf.GlobalViper)
}
