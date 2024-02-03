package constant

type Flags uint

var FeatureTypeMap map[string]int

const (
	FEATURE_S int = 1 << (iota)
	FEATURE_HOG
)

func init() {
	FeatureTypeMap = make(map[string]int)
	FeatureTypeMap["SURF/SFIT"] = FEATURE_S
	FeatureTypeMap["HOG"] = FEATURE_HOG
}
