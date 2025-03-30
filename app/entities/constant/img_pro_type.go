package constant

type Flags uint

var FeatureTypeMap map[string]int

const (
	FeatureS int = 1 << (iota)
	FeatureHog
)

func init() {
	FeatureTypeMap = make(map[string]int)
	FeatureTypeMap["SURF/SFIT"] = FeatureS
	FeatureTypeMap["HOG"] = FeatureHog
}
