package gexcel_test

import (
	"demo/app/entities/insect"
	"demo/common/util/gexcel"
	"testing"
)

func Test(t *testing.T) {
	insects := []insect.Insect{
		{Id: "1", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{Id: "2", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{Id: "3", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{Id: "4", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
		{Id: "5", ChineseName: "黑脉金斑蝶", LatinName: "Danaus plexippus", EnglishName: "Monarch butterfly "},
	}
	gexcel.WriteData(insects)
}
