package handler

import (
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/HunDun0Ben/bs_server/app/internal/dto"
	"github.com/HunDun0Ben/bs_server/app/pkg/helper"
)

// BigLotteryRandom godoc
// @Summary      生成大乐透随机号码
// @Description  生成大乐透随机号码，分为前区和后区两部分
// @Tags         LotteryController
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.SwaggerResponse{data=dto.LotteryResult} "成功响应，返回随机生成的大乐透号码"
// @Failure      500  {object}  dto.SwaggerResponse "服务器内部错误"
// @Router       /lottery/bigLottery/random [get]
func BigLotteryRandom(cxt *gin.Context) {
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Generate random numbers for upper half (e.g., 6 unique numbers between 1-35)
	upperHalf := generateUniqueNumbers(6, 1, 35, rng)

	// Generate random numbers for lower half (e.g., 1 number between 1-12)
	lowerHalf := generateUniqueNumbers(1, 1, 12, rng)

	result := dto.LotteryResult{
		UpperHalf: upperHalf,
		LowerHalf: lowerHalf,
	}

	helper.Success(cxt, result)
}

// generateUniqueNumbers generates n unique random numbers between min and max (inclusive)
func generateUniqueNumbers(n, min, max int, rng *rand.Rand) []int {
	if n > (max - min + 1) {
		n = max - min + 1
	}

	result := make([]int, 0, n)
	used := make(map[int]bool)

	for len(result) < n {
		num := rng.Intn(max-min+1) + min
		if !used[num] {
			used[num] = true
			result = append(result, num)
		}
	}

	return result
}
