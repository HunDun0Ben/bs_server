package dto

// LotteryResult represents the structure for lottery results with upper and lower halves
type LotteryResult struct {
	UpperHalf []int `json:"upperHalf"`
	LowerHalf []int `json:"lowerHalf"`
}