package insect

type Location struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
}

type Classification struct {
	Id      string
	Kingdom string
	Phylum  string
	Class   string
	Order   string
	Family  string
	Genus   string
	Species string
}

type Insect struct {
	Id                 string `xlsx:"主键id"`
	ClassId            string `xlsx:"分类id"`
	ChineseName        string `xlsx:"中文名称"`
	EnglishName        string `xlsx:"英文名称"`
	LatinName          string `xlsx:"拉丁学名"`
	FeatureDescription string `xlsx:"特征描述文本"`
	Distribution       string `xlsx:"分布情况文本"`
	ProtectionLevel    string `xlsx:"保护级别文本"`
	Classification     Classification
}
