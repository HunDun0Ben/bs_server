package config

type SVMConfig struct {
	Type   int     `json:"type"`
	Kernel int     `json:"kernel"`
	Degree float64 `json:"degree"`
	Gamma  float64 `json:"gamma"`
	Coeff0 float64 `json:"coef0"`
	Nu     float64 `json:"nu"`
	P      float64 `json:"p"`
}

func NewSVMConfig() *SVMConfig {
	return &SVMConfig{
		Type:   100,
		Kernel: 0,
		Degree: 0,
		Gamma:  1,
		Coeff0: 1,
		Nu:     0,
		P:      0,
	}
}
