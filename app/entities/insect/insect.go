package insect

type Location struct {
	Latitude  float64 `bson:"latitude"`
	Longitude float64 `bson:"longitude"`
	Altitude  float64 `bson:"altitude"`
}

type Classification struct {
	Kingdom string `xlsx:"界" bson:"kingdom,omitempty"`
	Phylum  string `xlsx:"门" bson:"phylum,omitempty"`
	Class   string `xlsx:"纲" bson:"class,omitempty"`
	Order   string `xlsx:"目" bson:"order,omitempty"`
	Family  string `xlsx:"科" bson:"family,omitempty"`
	Genus   string `xlsx:"属" bson:"genus,omitempty"`
	Species string `xlsx:"种" bson:"species,omitempty"`
}

type Insect struct {
	Id                 string          `xlsx:"主键id" bson:"_id"`
	ChineseName        string          `xlsx:"中文名称" bson:"chinese_name,omitempty"`
	EnglishName        string          `xlsx:"英文名称" bson:"english_name,omitempty"`
	LatinName          string          `xlsx:"拉丁学名" bson:"latin_name,omitempty"`
	FeatureDescription string          `xlsx:"特征描述文本" bson:"feature_description,omitempty"`
	Distribution       string          `xlsx:"分布情况文本" bson:"distribution,omitempty"`
	ProtectionLevel    string          `xlsx:"保护级别文本" bson:"protection_level,omitempty"`
	Classification     *Classification `xlsx:"分类类型" bson:"classification,omitempty"`
}
