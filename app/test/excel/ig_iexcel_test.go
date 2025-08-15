package test

import (
	"log/slog"
	"testing"

	"github.com/HunDun0Ben/bs_server/app/internal/model/insect"
	"github.com/HunDun0Ben/bs_server/app/pkg/util/iexcel"
)

func TestWriteExcel(t *testing.T) {
	insects := []insect.Insect{
		{ID: "1", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "2", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "3", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "4", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "5", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
	}
	builder := iexcel.New().AddSheet("蝴蝶信息", insects)
	if builder.Error() != nil {
		slog.Error("生成 excel 的时候有问题", "error", builder.Error())
	}
	builder.Save("butterfly.xlsx")
}
