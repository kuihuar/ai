// 贪心算法
// 在对问题求解时，总是做出在当前看来是最好的选择，不追求最优解，快速找到满意解。
// 贪心算法的时间复杂度为O(n)，空间复杂度为O(1)


### 适用贪心算法的场景
1. 问题能够分成子问题，子问题的最优解能够递推到最终问题的最优解。这种子问题最优解成为最优子结构
2. 贪心算法与动态规划的不同在于它对每个子问题的解决方案都做出选择，不能回退。动态规划则会保存以前的运算结果，并根据以前的结果对当前进行选择，有回退功能。最后得到全局最优解。

### 122 买卖股票的最佳时机 II

1. DFSO(2^n)
2. 贪心算法，只要后天价格比前一天高就卖出O(n)
3. 动态规划，每天的最大利润为前一天的最大利润加上今天的利润，或者前一天的最大利润减去今天的价格，取最大值 O(n)

#### 贪心算法思路
- 核心思想：只要后一天的价格比前一天高，就累加这两天的差价作为利润。所有局部最优操作（上涨阶段买卖）的累计即为全局最优解。

- 数学证明：多次买卖的累计利润等于所有独立上涨阶段的利润之和


### DFS方法
```go
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
```
#### 调用过程
输入：prices = [7, 1, 5, 3, 6, 4]
调用dfs(0, false)
dfs(0, false)调用dfs(1, false)  // 不操作
dfs(1, false)调用dfs(2, false)  // 不操作
dfs(2, false)调用dfs(3, false)  // 不操作
dfs(3, false)调用dfs(4, false)  // 不操作
dfs(4, false)调用dfs(5, false)  // 不操作
dfs(5, false)调用dfs(6, false)  // 不操作
dfs(6, false)返回0，因为i >= n，没有股票可交易
dfs(5, false)返回dfs(5, false) 
操作结果：
- hold取反，prices[5]，max(0, dfs(i+1, true) - 4)，
- 继续调用dfs(6, true)
- profit = max(0, 0 - 4) 
- profit = max(0, -4) = 0
dfs(4, false)返回dfs(4, false) 回溯到这里了
回溯到状态：dfs(4, false)
1. 初始状态：

i = 4，hold = false。
先调用 dfs(5, false)，不进行任何操作。
2. 调用 dfs(5, false)：

进入 dfs(5, false)：
调用 dfs(6, false)，得到返回值 0（基线条件）。
3. 计算利润：

在 dfs(5, false) 中：
profit = max(0, dfs(5, true) - prices[5]) // prices[5] = 4
4. 调用 dfs(5, true)：

进入 dfs(5, true)：
先调用 dfs(6, true)，得到返回值 0（基线条件）。
5. 计算利润：

在 dfs(5, true) 中：

profit = max(0, dfs(6, false) + prices[5]) // prices[5] = 4
这里 dfs(6, false) 返回 0：

profit = max(0, 0 + 4) = 4
返回 4 到 dfs(5, false)。
回到 dfs(5, false)：

计算：

profit = max(0, 4 - 4) = max(0, 0) = 0
返回 0 到 dfs(4, false)。


#### 返回到 dfs(4, false)
现在继续在 dfs(4, false) 中的其他选择：

1. 选择买入：
计算：
profit = max(0, dfs(5, true) - prices[4]) // prices[4] = 6

2. 已经计算过 dfs(5, true)，返回值为 4，所以：
计算：
profit = max(0, 4 - 6) = max(0, -2) = 0

3. 返回 dfs(4, false)：
此时，profit 为 0。
```
---
dfs(3, false)返回dfs(3, false) + prices[3] = 0 + 3 = 3
dfs(2, false)返回dfs(2, false) + prices[2] = 0 + 5 = 5
dfs(1, false)返回dfs(1, false) + prices[1] = 0 + 1 = 1
dfs(0, false)返回dfs(0, false) + prices[0] = 0 + 7 = 7