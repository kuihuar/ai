package redpackage



import (
	// "fmt"
	"math/rand"
	"time"
)

// DivideRedEnvelope 实现二倍均值法拆分红包
func DivideRedEnvelope(totalAmount float64, totalCount int) []float64 {
	// 存储每个红包的金额
	result := make([]float64, totalCount)
	// 剩余金额初始化为总金额
	remainAmount := totalAmount
	// 剩余红包个数初始化为总个数
	remainCount := totalCount

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < totalCount-1; i++ {
		// 计算当前可分配的最大金额
		max := remainAmount / float64(remainCount) * 2
		// 生成 0 到 max 之间的随机金额，最小为 0.01
		amount := rand.Float64()*(max-0.01) + 0.01
		// 保留两位小数
		amount = float64(int(amount*100)) / 100
		result[i] = amount
		// 更新剩余金额
		remainAmount -= amount
		// 更新剩余红包个数
		remainCount--
	}
	// 最后一个红包为剩余金额
	// result[totalCount-1] = float64(int(remainAmount*100)) / 100

	result[totalCount-1] = remainAmount

	return result
}

