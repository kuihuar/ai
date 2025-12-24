package suanfa






ch1:= make(chan int)

// 300. 最长递增子序列
// 不能位置调换
// 10, 9, 2, 5, 3, 7, 101, 18,20
// result : 2,3,7,18,20
// 暴力求解：2^n

func LengthOfLISForce(nums []int) int {
	n := len(nums)
	if n == 0 {
		return 0
	}
	maxLen := 0
	var dfs func(int, int, int)
	dfs = func(i, prev int, length int) {
		if i == n {
			if length > maxLen {
				maxLen = length
			}
			return
		}
		// 不选择当前元素
		dfs(i+1, prev, length)
		// 如果当前元素大于前一个元素，选择当前元素
		if nums[i] > prev {
			dfs(i+1, nums[i], length+1)
		}
	}
	dfs(0, -1, 0)
	return maxLen
}

// 动态规划
// 1. 定义子问题
// 2. 写出子问题的递推关系
// dp[i]表示头元素到第i个元素（且把i个元素要选上）的最长递增子序列的长度
// 也就是以第i个元素为结尾的最长递增子序列的长度
// 结果： 从dp0,dp1,dp2...dpn中找出最大的值
// 递推关系：
// for i:=0 ~n-1
// 	for j:=0 ~i-1
// 1. 如果nums[i] > nums[j] 那么dp[i] = max(dp[i], dp[j]+1)
// 2. 如果nums[i] <= nums[j] 那么dp[i] = max(dp[i], 1)
// 3. 初始化dp[i] = 1

// O(N^2)

// 定义状态：dp[i] 表示以第 i 个元素结尾的最长上升子序列的长度
// 初始状态：对于每个元素，最短的上升子序列就是它自己，所以初始时所有 dp[i] = 1。
// 状态转移：对于每个 i，我们检查所有 j < i 的元素：
// 如果 nums[i] > nums[j]，说明 nums[i] 可以接在 nums[j] 结尾的子序列后面
// 因此 dp[i] 可以更新为 dp[j] + 1（如果这个值比当前 dp[i] 大）
func LengthOfLISDP(nums []int) int {
	n := len(nums)
	res := 1 // 初始化结果为1，因为最短的子序列长度是1

	dp := make([]int, n+1)
	// 初始化dp数组，每个元素至少可以单独构成一个子序列
	for i := 0; i < n; i++ {
		dp[i] = 1
	}
	for i := 1; i < n; i++ {
		// 遍历前面的所有元素
		// fmt.Printf("i: %v", i)
		for j := 0; j < i; j++ {
			// 如果当前数字比前面某个数字大，
			// 那么以当前数字结尾的LIS长度至少可以是以那个数字结尾的LIS长度+1。
			if nums[i] > nums[j] {
				dp[i] = max(dp[i], dp[j]+1)
			}
			// fmt.Printf("j: %v, i: %v \t", j, i)
		}
		// fmt.Println()
		res = max(res, dp[i])
	}
	return res
}

func lowerBoundManual(arr []int, target int) int {
	left, right := 0, len(arr)
	for left < right {
		mid := left + (right-left)/2
		if arr[mid] >= target {
			right = mid
		} else {
			left = mid + 1
		}
	}
	return left
}

// O(nlogn)
func LengthOfLISDP2(nums []int) int {
	res := make([]int, 0)
	for _, num := range nums {
		// 二分查找，找到第一个大于等于nums[i]的位置
		// 如果找到，更新这个位置的值为nums[i]
		// 如果找不到，说明nums[i]比res中所有元素都大，将nums[i]添加到res的末尾
		// 这样可以保证res中的元素是递增的，并且长度最长
		left, right := 0, len(res)
		for left < right {
			mid := left + (right-left)/2
			if res[mid] >= num {
				right = mid
			} else {
				left = mid + 1
			}
		}
		// 如果left等于len(res)，说明nums[i]比res中所有元素都大，将nums[i]添加到res的末尾
		if left == len(res) {
			res = append(res, num)
			// 否则，更新res[left]的值为nums[i]
		} else {
			res[left] = num
		}

	}
	return len(res)
}

func LengthOfLISDP3(nums []int) int {
	lis := make([]int, 0)
	for _, num := range nums {
		i := lowerBoundManual(lis, num)
		if i == len(lis) {
			lis = append(lis, num)
		} else {
			lis[i] = num
		}
	}
	return len(lis)
}

// 322 零钱兑换

func CoinChange(coins []int, amount int) int {

	if amount == 0 {
		return 0
	}
	if len(coins) == 0 {
		return -1
	}
	dp := make([]int, amount+1)

	dp[0] = 0 // 当金额为0时，不需要任何硬币
	// 初始化DP数组，每个金额初始设为amount+1（表示不可达）
	for i := 1; i <= amount; i++ {
		// 初始化每个金额所需的最少硬币数为amount+1
		// 这里假设了最最大值
		dp[i] = amount + 1
	}

	// dp[i] 表示凑出金额 i 所需的最少硬币数
	for i := 1; i <= amount; i++ {
		//dp[i] = amount + 1
		// 状态转移：
		// 对于每个金额 i，尝试所有可能的硬币 coin
		// 如果 i >= coin，则 dp[i] 可以取 dp[i-coin]+1 的最小值
		for _, coin := range coins {
			// 只有当硬币面值不大于目标金额时才考虑
			// 存在一种情况，就是硬币面值大于目标金额，
			if i >= coin {
				// dp[i]：凑出金额 i 所需的最少硬币数量
				// dp[i-coin]：凑出金额 i-coin 所需的最少硬币数
				// coin：当前考虑的硬币面值
				// dp[i-coin]+1：表示如果使用coin这枋硬币，那么剩下的金额是i-coin
				// 凑出i-coin的最少硬币数dp[i-coin]，
				// 再加上当前硬币coin，就是+1
				// 取所有可能的硬币面值中，保留最小的硬币数量到dp[i]
				dp[i] = min(dp[i], dp[i-coin]+1)
			}
		}
	}
	if dp[amount] > amount {
		return -1
	}
	return dp[amount]
}

// 看看可视化推理，还有点有不明白
// 5. 可视化过程
// 以 coins = [1, 2, 5], amount = 6 为例：

// 金额i	dp[i] 计算过程								结果		公式
// 0		0（基准情况）								  0		dp[i] = min(dp[i], dp[i-coin]+1)
// 1		min(dp[0]+1) = 1							1     dp[1] = min(dp[1], dp[0]+1) //dp[1]已经初始化最大值
// 2		min(dp[1]+1, dp[0]+1) = min(2,1)			1
// 3		min(dp[2]+1, dp[1]+1) = min(2,2)			2
// 4		min(dp[3]+1, dp[2]+1) = min(3,2)			2
// 5		min(dp[4]+1, dp[3]+1, dp[0]+1)=min(3,3,1)	1
// 6		min(dp[5]+1, dp[4]+1, dp[1]+1)=min(2,3,2)	2

// 72 编辑距离
// 1. 暴力解法， 时间复杂度O(3^(m+n))
func minDistance(word1, word2 string) int {
	var editDistance func(word1, word2 string, i, j int) int
	editDistance = func(word1, word2 string, i, j int) int {
		if i == 0 {
			return j
		}
		if j == 0 {
			return i
		}
		if word1[i-1] == word2[j-1] {
			return editDistance(word1, word2, i-1, j-1)
		} else {
			return 1 + min(
				editDistance(word1, word2, i-1, j-1), //替找
				editDistance(word1, word2, i-1, j),   //删除
				editDistance(word1, word2, i, j-1))   //插入
		}
	}

	return editDistance(word1, word2, len(word1), len(word2))
}

func minDistanceMem(word1, word2 string) int {
	memo := make(map[[2]int]int)
	var editDistance func(i, j int) int
	editDistance = func(i, j int) int {
		if v, ok := memo[[2]int{i, j}]; ok {
			return v
		}
		if i == 0 {
			return j
		}
		if j == 0 {
			return i
		}
		if word1[i-1] == word2[j-1] {
			memo[[2]int{i, j}] = editDistance(i-1, j-1)
		} else {
			memo[[2]int{i, j}] = 1 + min(
				editDistance(i-1, j-1), //替找
				editDistance(i-1, j),   //删除
				editDistance(i, j-1))   //插入
		}
		return memo[[2]int{i, j}]
	}

	return editDistance(len(word1), len(word2))
}

// 2. 动态规划
// 定义状态：dp[i][j] 表示将 word1 的前 i 个字符转换为 word2 的前 j 个字符所需的最少操作次数。
// 初始状态：dp[i][0] = i，dp[0][j] = j，因为当一个字符串为空时，需要进行插入或删除操作。
// 状态转移：
// 如果 word1[i-1] == word2[j-1]，则 dp[i][j] = dp[i-1][j-1]，因为不需要进行任何操作。
// 如果 word1[i-1] != word2[j-1]，则 dp[i][j] = 1 + min(dp[i-1][j-1], dp[i-1][j], dp[i][j-1])，表示进行替换、删除或插入操作。
// 最终结果：dp[m][n]，其中 m 和 n 分别是 word1 和 word2 的长度。
// 时间复杂度：O(m * n)，空间复杂度：O(m * n)。
func minDistanceDP(word1, word2 string) int {
	m, n := len(word1), len(word2)
	// 创建DP表，大小为 (m+1) x (n+1)
	// 多出来的一行和一列用于表示空字符串（即 i=0 或 j=0 的情况）
	// word1[0..m-1] 表示 word1 的所有字符（长度为 m）
	// 字符串下标从0开始，取值范围从0到m, 共m+1种情况
	// dp[0][j]：
	// 将空字符串 "" 转换为 word2 的前 j 个字符，需要 j 次插入操作。
	// 例如："" → "abc" 需要插入 3 次
	// dp[i][0]：
	// 将 word1 的前 i 个字符转换为空字符串 ""，需要 i 次删除操作。
	// 例如："abc" → "" 需要删除 3 次
	// (2) 为什么必须包括空字符串？
	// 动态规划的基准条件：
	// 所有子问题的解最终都递归到 i=0 或 j=0 的情况（即一个字符串为空）
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	//第一个单词前i个字符变为空串需要i次删除
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	//第一个单词word1为空串变为word2前j个字符需要j次插入
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j-1], dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[m][n]
}
