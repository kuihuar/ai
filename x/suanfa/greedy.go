package suanfa

//
// 贪心算法思路
// 核心思想：只要后一天的价格比前一天高，就累加这两天的差价作为利润。所有局部最优操作（上涨阶段买卖）的累计即为全局最优解。

// 数学证明：多次买卖的累计利润等于所有独立上涨阶段的利润之和
func MaxProfit(prices []int) int {
	maxProfix := 0

	for i := 1; i < len(prices); i++ {
		if prices[i] > prices[i-1] {
			maxProfix += prices[i] - prices[i-1]
		}
	}
	return maxProfix
}

//
// 贪心算法思路
// 核心思想：每次选择当前可获得的最大利润，直到无法继续交易。
//
// 数学证明：局部最优解是全局最优解的一部分
// 每天的状态只与前一天的状态有关

// 动态规划通常用于解决有重叠子问题和最优子结构的问题
// 动态规划的时间复杂度通常为O(n)，空间复杂度为O(n)。
func MaxProfit1(prices []int) int {
	// n := len(prices)

	// dp0, dp1 := 0, -prices[0]
	// for i := 1; i < n; i++ {
	// 	dp0, dp1 = max(dp0, dp1+prices[i]), max(dp1, dp0-prices[i])
	// }
	// return dp0

	n := len(prices)

	dp := make([][2]int, n)

	dp[0][0] = 0          // 第0天未持有股票的利润
	dp[0][1] = -prices[0] // 第0天持有股票的利润（买入成本）

	for i := 1; i < n; i++ {
		dp[i][0] = max(dp[i-1][0], dp[i-1][1]+prices[i])
		dp[i][1] = max(dp[i-1][1], dp[i-1][0]-prices[i])
	}
	return dp[n-1][0]
}

func MaxProfit2(prices []int) int {
	// n := len(prices)
	// dp0, dp1 := 0, -prices[0]
	// for i := 1; i < n; i++ {
	// 	dp0, dp1 = max(dp0, dp1+prices[i]), max(dp1, dp0-prices[i])
	// }
	// return dp0
	n := len(prices)

	preNotHold := 0
	preHold := -prices[0]
	for i := 1; i < n; i++ {
		currentNotHold := max(preNotHold, preHold+prices[i])
		currentHold := max(preHold, preNotHold-prices[i])
		// 将当天的计算结果传递给“前一天的”状态变量，实现状态的递推更新
		preNotHold = currentNotHold
		preHold = currentHold
	}
	return preNotHold
}

// 动态规划时间复杂度为O(n)，空间复杂度为O(n)。
func MaxProfit3(prices []int) int {
	n := len(prices)
	preNotHold := 0 // 第0天未持有股票的利润
	preHold := -prices[0]
	for i := 1; i < n; i++ {
		preNotHold, preHold = max(preNotHold, preHold+prices[i]), max(preHold, preNotHold-prices[i])
	}
	return preNotHold
}

// 122 买卖股票的最佳时机 II
// 最佳时机：尽可能地完成更多的交易（多次买卖一支股票）。
// 数学证明：局部最优解是全局最优解的一部分
// 每天的状态只与前一天的状态有关
// 1. 多少股: 最多持有一股股票，或者0股。
// 2. 多少次: 当天最多可以完成一次交易。要么买进/卖出，要么不操作。
// 3. 何时: 可以在任何时候进行买入/卖出操作。
// 4. 如何: 买入后必须卖出才能再次买入。
// 5. 利润: 卖出时获得的利润。
// 6. 成本: 买入时的成本。
// 7. 手续费: 每次交易需要支付一定的手续费，本题没有。

// 贪心算法时间复杂度为O(n)，空间复杂度为O(1)。
func MaxProfit4(prices []int) int {
	n := len(prices)
	profit := 0
	for i := 1; i < n; i++ {
		if prices[i] > prices[i-1] {
			profit += prices[i] - prices[i-1]
		}
	}
	return profit
}

// DFS 代码实现
// 时间复杂度为O(n)，空间复杂度为O(n)。
func MaxProfit5(prices []int) int {
	n := len(prices)

	var dfs func(i int, hold bool) int
	dfs = func(i int, hold bool) int {
		if i >= n {
			return 0 // 边界条件, 没有股票可交易
		}
		profit := dfs(i+1, hold) // 不操作
		if hold {
			profit = max(profit, dfs(i+1, false)+prices[i]) // 卖出

			// 当前持有：可以选择卖出
			sell := dfs(i+1, false) + prices[i]
			profit = max(profit, sell)
		} else {
			//profit = max(profit, dfs(i+1, true)-prices[i]) // 买入
			// 当前未持有：可以选择买入
			buy := dfs(i+1, true) - prices[i]
			profit = max(profit, buy)
		}
		return profit
	}
	return dfs(0, false)
}

//		dfs + 记忆化搜索 代码实现 记忆化搜索是一种优化技术，用于减少递归函数的重复计算。
//		它通过将已经计算过的结果存储在一个数组或哈希表中，
//		从而避免重复计算相同的子问题。
//	 memo[i][0] 表示第 i 天未持有股票的最大利润，
//		memo[i][1] 表示第 i 天持有股票的最大利润。
//
// 时间复杂度为O(n)，空间复杂度为O(n)。
func MaxProfit7(prices []int) int {
	n := len(prices)
	memo := make([][2]int, n)
	for i := range memo {
		memo[i] = [2]int{-1, -1}
	}
	var dfs func(i int, hold bool) int
	dfs = func(i int, hold bool) int {
		if i >= n {
			return 0 // 边界条件，没有股票可交易
		}
		holdIdx := 0
		if hold {
			holdIdx = 1
		}
		if memo[i][holdIdx] != -1 {
			return memo[i][holdIdx]
		}
		profit := dfs(i+1, hold) // 不操作
		// 利润初始化为0，因为如果不操作，利润为0
		// 这里计算在当前天 i 不进行任何交易的情况下的利润，递归调用下一天
		// 每一天都是一个子问题
		if hold {
			profit = max(profit, dfs(i+1, false)+prices[i]) // 卖出
		} else {
			profit = max(profit, dfs(i+1, true)-prices[i]) // 买入
		}
		memo[i][holdIdx] = profit
		return profit
	}
	return dfs(0, false)
}

// 递归性质：

// 递归的本质是将问题分解为子问题。每个子问题的解必须独立计算，因此在每次调用时都需要初始化相关变量

// [7, 1, 5, 3, 6, 4]
func MaxProfit8(prices []int) int {
	n := len(prices)

	var dfs func(i int, hold bool) int
	dfs = func(i int, hold bool) int {
		if i >= n {
			return 0 // 边界条件，没有股票可交易
		}
		profit := dfs(i+1, hold) // 不操作
		// 利润初始化为0，因为如果不操作，利润为0
		// 这里计算在当前天 i 不进行任何交易的情况下的利润，递归调用下一天
		// 每一天都是一个子问题
		if hold {
			profit = max(profit, dfs(i+1, false)+prices[i]) // 卖出
		} else {
			profit = max(profit, dfs(i+1, true)-prices[i]) // 买入
		}
		return profit
	}
	return dfs(0, false)
}
