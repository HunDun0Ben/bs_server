package gexcel_test

import (
	"testing"

	"github.com/HunDun0Ben/bs_server/app/internal/model/insect"
	"github.com/HunDun0Ben/bs_server/app/pkg/util/gexcel"
)

func Test(t *testing.T) {
	insects := []insect.Insect{
		{ID: "1", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "2", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "3", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "4", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{ID: "5", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
	}
	gexcel.WriteData(insects)
}
