package algorithm

import (
	"fmt"
	"sort"
)

// LengthOfLIS 计算最长递增子序列的长度
// 使用动态规划解法，时间复杂度O(n²)，空间复杂度O(n)
func LengthOfLIS(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	// dp[i] 表示以 nums[i] 结尾的最长递增子序列长度
	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1 // 初始化为1，因为单个元素本身就是长度为1的递增子序列
	}

	maxLen := 1
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
		maxLen = max(maxLen, dp[i])
	}

	return maxLen
}

// LengthOfLISOptimized 优化版本，使用二分查找
// 时间复杂度O(n log n)，空间复杂度O(n)
func LengthOfLISOptimized(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	// tails[i] 表示长度为i+1的递增子序列的最小末尾值
	tails := make([]int, 0)

	for _, num := range nums {
		// 二分查找插入位置
		pos := binarySearch(tails, num)
		if pos == len(tails) {
			// 如果num大于所有tails中的值，则添加到末尾
			tails = append(tails, num)
		} else {
			// 否则替换tails[pos]为num
			tails[pos] = num
		}
	}

	return len(tails)
}

// binarySearch 二分查找，找到第一个大于等于target的位置
func binarySearch(nums []int, target int) int {
	left, right := 0, len(nums)
	for left < right {
		mid := left + (right-left)/2
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid
		}
	}
	return left
}

// GetLongestIncreasingSubsequence 获取具体的最长递增子序列
func GetLongestIncreasingSubsequence(nums []int) []int {
	if len(nums) == 0 {
		return []int{}
	}

	// dp[i] 表示以 nums[i] 结尾的最长递增子序列长度
	dp := make([]int, len(nums))
	// prev[i] 记录前驱节点，用于回溯
	prev := make([]int, len(nums))

	for i := range dp {
		dp[i] = 1
		prev[i] = -1
	}

	maxLen := 1
	maxIndex := 0

	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				prev[i] = j
			}
		}
		if dp[i] > maxLen {
			maxLen = dp[i]
			maxIndex = i
		}
	}

	// 回溯构建LIS序列
	lis := make([]int, maxLen)
	index := maxLen - 1
	current := maxIndex

	for current != -1 {
		lis[index] = nums[current]
		index--
		current = prev[current]
	}

	return lis
}

// LengthOfLISRecursive 递归解法（带记忆化）
func LengthOfLISRecursive(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	// 记忆化数组
	memo := make([]int, len(nums))
	for i := range memo {
		memo[i] = -1
	}

	maxLen := 0
	for i := 0; i < len(nums); i++ {
		maxLen = max(maxLen, lisHelper(nums, i, memo))
	}

	return maxLen
}

// lisHelper 递归辅助函数
func lisHelper(nums []int, index int, memo []int) int {
	if memo[index] != -1 {
		return memo[index]
	}

	result := 1
	for i := 0; i < index; i++ {
		if nums[index] > nums[i] {
			result = max(result, lisHelper(nums, i, memo)+1)
		}
	}

	memo[index] = result
	return result
}

// LengthOfLISWithPath 带路径的LIS算法
func LengthOfLISWithPath(nums []int) (int, []int) {
	if len(nums) == 0 {
		return 0, []int{}
	}

	// dp[i] 表示以 nums[i] 结尾的最长递增子序列长度
	dp := make([]int, len(nums))
	// path[i] 记录以 nums[i] 结尾的LIS路径
	path := make([][]int, len(nums))

	for i := range dp {
		dp[i] = 1
		path[i] = []int{nums[i]}
	}

	maxLen := 1
	maxIndex := 0

	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				// 复制路径并添加当前元素
				path[i] = make([]int, len(path[j]))
				copy(path[i], path[j])
				path[i] = append(path[i], nums[i])
			}
		}
		if dp[i] > maxLen {
			maxLen = dp[i]
			maxIndex = i
		}
	}

	return maxLen, path[maxIndex]
}

// CountLIS 统计所有可能的最长递增子序列数量
func CountLIS(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	// dp[i] 表示以 nums[i] 结尾的最长递增子序列长度
	dp := make([]int, len(nums))
	// count[i] 表示以 nums[i] 结尾的最长递增子序列数量
	count := make([]int, len(nums))

	for i := range dp {
		dp[i] = 1
		count[i] = 1
	}

	maxLen := 1

	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
					count[i] = count[j]
				} else if dp[j]+1 == dp[i] {
					count[i] += count[j]
				}
			}
		}
		maxLen = max(maxLen, dp[i])
	}

	// 统计所有长度为maxLen的LIS数量
	totalCount := 0
	for i := 0; i < len(nums); i++ {
		if dp[i] == maxLen {
			totalCount += count[i]
		}
	}

	return totalCount
}

// LengthOfLISWithConstraints 带约束条件的LIS
// 例如：相邻元素差值不能超过k
func LengthOfLISWithConstraints(nums []int, k int) int {
	if len(nums) == 0 {
		return 0
	}

	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1
	}

	maxLen := 1
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] && nums[i]-nums[j] <= k {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
		maxLen = max(maxLen, dp[i])
	}

	return maxLen
}

// LengthOfLIS2D 二维LIS问题
// 给定二维点集，找到最长的递增序列
type Point struct {
	X, Y int
}

func LengthOfLIS2D(points []Point) int {
	if len(points) == 0 {
		return 0
	}

	// 按x坐标排序，如果x相同则按y排序
	sort.Slice(points, func(i, j int) bool {
		if points[i].X != points[j].X {
			return points[i].X < points[j].X
		}
		return points[i].Y < points[j].Y
	})

	// 转换为y坐标的LIS问题
	yCoords := make([]int, len(points))
	for i, p := range points {
		yCoords[i] = p.Y
	}

	return LengthOfLISOptimized(yCoords)
}

// PrintDPTable 打印DP表（用于调试）
func PrintDPTable(nums []int) {
	if len(nums) == 0 {
		return
	}

	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1
	}

	// 填充DP表
	for i := 1; i < len(nums); i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				dp[i] = max(dp[i], dp[j]+1)
			}
		}
	}

	// 打印表头
	fmt.Print("索引: ")
	for i := 0; i < len(nums); i++ {
		fmt.Printf(" %2d ", i)
	}
	fmt.Println()

	fmt.Print("数值: ")
	for _, num := range nums {
		fmt.Printf(" %2d ", num)
	}
	fmt.Println()

	fmt.Print("DP值: ")
	for _, val := range dp {
		fmt.Printf(" %2d ", val)
	}
	fmt.Println()
}

// 辅助函数
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
