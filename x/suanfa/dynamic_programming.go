package suanfa

import (
	"fmt"
	"math"
)

// import "github.com/mailru/easyjson/opt"

func fib(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fib(n-1) + fib(n-2)
}
func fib1(n int) int {
	if n == 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	memo := make([]int, n+1)
	if memo[n] != 0 {
		return memo[n]
	}
	memo[n] = fib1(n-1) + fib1(n-2)
	return memo[n]
}

// 动态规划
// 展开fib调用函数过程
// 去掉记忆化的重复步骤
// 自底向上，从最小的问题开始，逐步向上解决问题
// 发现最小的问题是f(0)和f(1)，所以从f(0)和f(1)开始，逐步向上解决问题
// f(2) = f(1) + f(0)， 循环i从2到n，依次计算f(i)，直到f(n)
// 结果放在一个数组里，f(0)和f(1)已经知道了，所以从f(2)开始，依次计算f(i)，直到f(n)
// 最后返回f(n)

// f(3) = f(2) + f(1)
func Fib(n int) int {
	db := make([]int, n+1)
	db[0] = 0
	db[1] = 1
	for i := 2; i <= n; i++ {
		db[i] = db[i-1] + db[i-2]
	}
	return db[n]
}

func CountPath(grid [][]int) int {
	// 障碍物
	// if validSquare(grid, row, col) {
	// 	return 0
	// }
	// // 到达终点
	// if row == len(grid)-1 && col == len(grid[0])-1 {
	// 	return 1
	// }
	// return CountPath(grid, row+1, col) + CountPath(grid, row, col+1)

	// 动态规划
	// 递推，管相邻状态就可以
	//opt[i,j] = opt[i-1,j] + opt[i,j-1]
	// 空地
	// if grid[i][j] == 0 {
	// 	opt[i,j] = opt[i-1,j] + opt[i,j-1]
	// }
	// return opt[m-1,n-1]
	m := len(grid)
	n := len(grid[0])
	fmt.Println("m: ", m, "n: ", n)
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	fmt.Println("grid: ", grid)
	fmt.Printf("dp   : %v \n", dp)
	fmt.Println("------------------------")
	dp[m-1][n-1] = 1

	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if i == m-1 && j == n-1 {
				continue // 跳过终点
			}
			if grid[i][j] == 1 {
				dp[i][j] = 0 // 障碍物，不可达
			} else {
				// 向右和向下的路径数之和
				// right := 0
				// if j+1 < n {
				// 	right = dp[i][j+1]
				// }
				// down := 0
				// if i+1 < m {
				// 	down = dp[i+1][j]
				// }
				// dp[i][j] = right + down
				dp[i][j] = dp[i][j+1] + dp[i+1][j]
			}
		}
	}
	fmt.Printf("dp   : %v \n", dp)
	return dp[0][0]
}

// 从右上角往左上角递推, (i--,j--)
func CountPath1(grid [][]int) int {

	m := len(grid)
	n := len(grid[0])
	fmt.Println("m: ", m, "n: ", n)
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	dp[m-1][n-1] = 1

	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if i == m-1 && j == n-1 {
				//fmt.Printf("跳过终点dp[%d][%d] = %d \n", i, j, dp[i][j])
				continue // 跳过终点
			}
			if grid[i][j] == 1 {
				dp[i][j] = 0 // 障碍物，不可达
			} else {
				// 向右和向下的路径数之和
				right := 0
				if j+1 < n {
					right = dp[i][j+1]
				}
				down := 0
				if i+1 < m {
					down = dp[i+1][j]
				}
				dp[i][j] = right + down
			}
		}
	}

	return dp[0][0]

}

func CountPath2(grid [][]int) int {

	m := len(grid)
	n := len(grid[0])
	fmt.Println("m: ", m, "n: ", n)
	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}
	dp[0][0] = 1
	// 初始化第一行
	for j := 1; j < n; j++ {
		if grid[0][j] == 1 {
			dp[0][j] = 0
		} else {
			dp[0][j] = dp[0][j-1]
		}
	}
	//初始化第一列
	for i := 1; i < m; i++ {
		if grid[i][0] == 1 {
			dp[i][0] = 0
		} else {
			dp[i][0] = dp[i-1][0]
		}
	}
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if grid[i][j] == 1 {
				dp[i][j] = 0 // 障碍物，不可达
			} else {
				dp[i][j] = dp[i-1][j] + dp[i][j-1]
			}
		}
	}

	return dp[m-1][n-1]

}

// 70 爬楼梯
// 1. 定义子问题
// 2. 写出子问题的递推关系
// 3. 确定 DP 数组的计算顺序
// 4. 空间优化（可选）
func ClimbStairs(n int) int {
	if n == 1 {
		return 1
	}
	if n == 2 {
		return 2
	}
	dp := make([]int, n+1)
	dp[1] = 1
	dp[2] = 2
	// dp[i] 表示到达第 i 级楼梯的不同方法数
	for i := 3; i <= n; i++ {
		dp[i] = dp[i-1] + dp[i-2]
	}
	return dp[n]
}

// 走到第 4 阶的方法数 = 走到第 3 阶的方法数（再迈 1 阶） + 走到第 2 阶的方法数（再迈 2 阶）。

// 即 dp[4] = dp[3] + dp[2]。
func ClimbStairs1(n int) int {
	if n <= 2 {
		return n
	}
	//first 表示走到第 1 阶的方法数，second 表示走到第 2 阶的方法数。
	first := 1  // 走到第 1 阶的方法数,只有一种方法，一步到达
	second := 2 // 走到第 2 阶的方法数，有两种方法，一步到达，两步到达
	var current int
	// 从第 3 级开始，使用 current 计算当前级的方法数。
	for i := 3; i <= n; i++ {
		// 这表示到达第 i 级的方法数
		current = first + second
		// 将 first 更新为当前的 second,
		// 因为在下一次迭代中，first 需要表示第 i-2 级的方法数。
		first = second
		// 将 second 更新为当前的 current，
		// 因为在下一次迭代中，second 需要表示第 i-1 级的方法数
		second = current
	}
	fmt.Printf("first: %v, second: %v, current: %v \n", first, second, current)
	return second
	// 定义子问题
}

func ClimbStairs2(n int) int {
	if n <= 2 {
		return n
	}
	//first 表示走到第 1 阶的方法数，second 表示走到第 2 阶的方法数。
	first := 1  // 走到第 1 阶的方法数,只有一种方法，一步到达
	second := 1 // 走到第 2 阶的方法数，有两种方法，一步到达，两步到达
	for i := 2; i <= n; i++ {
		first, second = second, first+second
	}
	return second
}

// 120 三角形最小路径和

// 1. dp实现
// 2. 至顶向下，用回溯或者递归实现，将所有结果算出来，然后取最小值

func MinimumTotal(triangle [][]int) int {
	// dp 初始化为三角形的最后一行 triangle[n-1]，
	// 表示从底层各个位置出发的最小路径和（即它们自身的值）
	mini := triangle[len(triangle)-1]
	// 外层循环：从倒数第二行（i = n-2）向上遍历到第 0 行。
	for i := len(triangle) - 2; i >= 0; i-- {
		// 内层循环：对当前行的每个位置 j，计算 dp[j] = 当前值 + min(下一层的dp[j], dp[j+1])
		for j := 0; j < len(triangle[i]); j++ {
			// 这里的 dp[j] 和 dp[j+1] 是下一层已经计算好的最小路径和。
			// 更新后的 dp[j] 表示从当前位置 (i,j) 到底层的最小路径和
			// 每次计算完当前行的 dp[j] 后，直接覆盖原来的值。由于计算顺序是从左到右，且 dp[j+1] 在下一轮计算中才会被用到，因此不会影响正确性。

			mini[j] = triangle[i][j] + min(mini[j], mini[j+1])
		}
	}

	return mini[0]
}

func MinimumTotal1(triangle [][]int) int {
	n := len(triangle)
	dp := make([][]int, n)

	for i := range dp {
		dp[i] = make([]int, len(triangle[i]))
		copy(dp[i], triangle[i])
	}
	fmt.Printf("dp: %v \n", dp)

	for i := n - 2; i >= 0; i-- {
		for j := 0; j < len(triangle[i]); j++ {
			dp[i][j] = triangle[i][j] + min(dp[i+1][j], dp[i+1][j+1])
		}
	}

	return dp[0][0]
}

// 152 乘积最大子数组
// 1. 暴力求解， recursion,如何解
// 2. dynamic programming
// 3. 状态转移方程
// 4. 状态定义

func Maxproduct(nums []int) int {

	dp := make([][2]int, len(nums))

	dp[0][0] = nums[0]
	dp[0][1] = nums[0]
	res := nums[0]
	for i := 1; i < len(nums); i++ {

		if nums[i] > 0 {
			dp[i][0] = max(dp[i-1][0]*nums[i], nums[i])
			dp[i][1] = min(dp[i-1][1]*nums[i], nums[i])
		} else {
			dp[i][0] = max(dp[i-1][1]*nums[i], nums[i])
			dp[i][1] = min(dp[i-1][0]*nums[i], nums[i])
		}
		res = max(res, dp[i][0])
	}
	return res
}
func MaxProduct(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	maxProd, minProd, res := nums[0], nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] < 0 {
			maxProd, minProd = minProd, maxProd // 交换最大值和最小值
		}
		maxProd = max(nums[i], maxProd*nums[i])
		minProd = min(nums[i], minProd*nums[i])
		res = max(res, maxProd)
	}
	return res
}

func MaxProduct1(nums []int) int {
	res, maxVal, minVal := nums[0], nums[0], nums[0]

	for i := 1; i < len(nums); i++ {
		num := nums[i]
		maxVal, minVal = max(maxVal*num, minVal*num, num), min(maxVal*num, minVal*num, num)
		res = max(res, maxVal)
	}
	return res
}

// 121(每天只能买卖1次),122（无数次）,123（2次）,309（cool down）188（k次）,714（fee）

// dp[i] 表示到i天的最大利润maxProfit
// 结果就是到最后一天的最大利润 dp[len(prices)-1]
// 状态转移方程：
// dp[i] = dp[i-1] + （prices[i]卖） +  （-prices[i-1]买）
// 手上还有没有股票
// i 范围，从0到len(prices)-1表示天
// j 范围，0表示没有股票，1表示有股票
// 状态转移方程：
// dp[i-1][0]表示不交易，
// dp[i-1][1] + prices[i] 表示前天天有一股，卖出去
// dp[i][0] = max(dp[i-1][0], dp[i-1][1] + prices[i])
// dp[i-1][1]表示前一天有股票，今天不交易
// dp[i-1][0] - prices[i] 表示前一天没有股票，今天买股票
// dp[i][1] = max(dp[i-1][1], dp[i-1][0] - prices[i])
// k 表示之前交易了多少次

// dp[i][k][0] = max(dp[i-1][k][0], dp[i-1][k-1][1] + prices[i])
// dp[i][k][1] = max(dp[i-1][k][1], dp[i-1][k][0] - prices[i])

// 结果： max(dp[len(prices)-1][0~k][0]) // 从0～k

// cool down  k{0~1}

//最多拥有X股票，用j表示。每次只能买一股，卖一股，j有边界
// dp[i][k][j] = max(
// dp[i-1][k][j], //不动
// dp[i-1][k-1][j+1] + prices[i] //卖
// dp[i-1][k-1][j-1] - prices[i] //买
// )
// 算法复杂度 O（N*K）

// 最多一次交易
func maxProfitOne(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	dp := make([][3]int, len(prices))
	dp[0][0] = 0          // 没有买入
	dp[0][1] = -prices[0] // 买入
	dp[0][2] = 0          // 卖出
	res := 0
	for i := 1; i < len(prices); i++ {
		dp[i][0] = dp[i-1][0]
		dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
		dp[i][2] = dp[i-1][1] + prices[i]
		res = max(res, dp[i][0], dp[i][2], dp[i][1]) //dp[i][1]可以不考虑
		// dp[i][0] = max(dp[i-1][0], dp[i-1][1]+prices[i])
		// dp[i][1] = max(dp[i-1][1], -prices[i])
	}
	return res
}

// 最多两次交易
func maxProfitTwo(prices []int) int {
	if len(prices) == 0 {
		return 0
	}
	// 3表示最多交易2次，0表示没有买入，1表示第一次买入，2表示第一次卖出，
	// 最后一维0或者1表示是否持有股票
	// 中间一唯表示交易了多少次
	//第一维表示天数
	dp := make([][3][2]int, len(prices))
	dp[0][0][0], dp[0][0][1] = 0, -prices[0]
	dp[0][1][0], dp[0][1][1] = math.MinInt, math.MinInt
	dp[0][2][0], dp[0][2][1] = math.MinInt, math.MinInt

	for i := 1; i < len(prices); i++ {
		dp[i][0][0] = dp[i-1][0][0]
		dp[i][0][1] = max(dp[i-1][0][1], dp[i-1][0][0]-prices[i])

		dp[i][1][0] = max(dp[i-1][1][0], dp[i-1][0][1]+prices[i])
		dp[i][1][1] = max(dp[i-1][1][1], dp[i-1][1][0]-prices[i])

		dp[i][2][0] = max(dp[i-1][2][0], dp[i-1][1][1]+prices[i])
		//最后这个没啥用
		// dp[i][2][1] = max(dp[i-1][2][1], dp[i-1][2][0]-prices[i])
	}
	end := len(prices) - 1
	//不需要统计最后一维为1的
	return max(dp[end][0][0], dp[end][1][0], dp[end][2][0])
}

// 最多k次交易
func maxProfitK(prices []int, k int) int {
	if len(prices) == 0 {
		return 0
	}
	dp := make([][][2]int, len(prices))
	for i := range dp {
		dp[i] = make([][2]int, k+1)
	}
	dp[0][0][0], dp[0][0][1] = 0, -prices[0]
	for kk := 1; kk <= k; kk++ {
		dp[0][kk][0], dp[0][kk][1] = math.MinInt, math.MinInt
	}
	for i := 1; i < len(prices); i++ {
		for kk := 0; kk <= k; kk++ {
			if kk == 0 {
				dp[i][kk][0] = dp[i-1][kk][0]
				dp[i][kk][1] = max(dp[i-1][kk][1], dp[i-1][kk][0]-prices[i])
			} else {
				dp[i][kk][0] = max(dp[i-1][kk][0], dp[i-1][kk-1][1]+prices[i])
				dp[i][kk][1] = max(dp[i-1][kk][1], dp[i-1][kk][0]-prices[i])
			}
		}
	}
	res := 0
	end := len(prices) - 1
	for i := 0; i <= k; i++ {
		res = max(res, dp[end][i][0])
	}

	return res
}
func UseMP() {
	prices := []int{7, 1, 5, 3, 6, 4}
	res := maxProfitOne(prices)
	fmt.Printf("maxProfit: %v \n", res)
}
