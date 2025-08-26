package algorithm

import (
	"fmt"
	"strings"
)

// LISExample 演示LIS算法的详细计算过程
func LISExample() {
	fmt.Println("=== 最长递增子序列 (LIS) 算法演示 ===\n")

	// 示例1：基本示例
	demonstrateLIS([]int{10, 9, 2, 5, 3, 7, 101, 18}, "示例1")

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例2：简单示例
	demonstrateLIS([]int{1, 3, 6, 7, 9, 4, 10, 5, 6}, "示例2")

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例3：递减序列
	demonstrateLIS([]int{5, 4, 3, 2, 1}, "示例3 - 递减序列")

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例4：递增序列
	demonstrateLIS([]int{1, 2, 3, 4, 5}, "示例4 - 递增序列")
}

// demonstrateLIS 演示LIS算法的详细过程
func demonstrateLIS(nums []int, title string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("数组: %v (长度: %d)\n", nums, len(nums))

	// 计算LIS长度
	lisLength := LengthOfLIS(nums)
	fmt.Printf("\n最长递增子序列长度: %d\n", lisLength)

	// 获取具体的LIS序列
	lisSequence := GetLongestIncreasingSubsequence(nums)
	fmt.Printf("最长递增子序列: %v\n", lisSequence)

	// 打印DP表
	fmt.Println("\n动态规划表:")
	PrintDPTable(nums)

	// 测试优化版本
	optimizedLength := LengthOfLISOptimized(nums)
	fmt.Printf("\n优化版本结果: %d\n", optimizedLength)

	// 验证结果一致性
	if lisLength == optimizedLength {
		fmt.Println("✓ 标准版本和优化版本结果一致")
	} else {
		fmt.Println("✗ 结果不一致")
	}
}

// LISBacktrackingExample 演示LIS回溯过程
func LISBacktrackingExample() {
	fmt.Println("=== LIS回溯过程演示 ===\n")

	nums := []int{10, 9, 2, 5, 3, 7, 101, 18}

	fmt.Printf("数组: %v\n", nums)

	// 构建DP表和前驱数组
	dp := make([]int, len(nums))
	prev := make([]int, len(nums))

	for i := range dp {
		dp[i] = 1
		prev[i] = -1
	}

	maxLen := 1
	maxIndex := 0

	// 填充DP表
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

	// 打印DP表
	fmt.Println("\nDP表:")
	PrintDPTable(nums)

	// 回溯过程
	fmt.Println("\n回溯过程:")
	lis := make([]int, maxLen)
	index := maxLen - 1
	current := maxIndex

	step := 1
	for current != -1 {
		fmt.Printf("步骤%d: 当前位置 = %d, 数值 = %d\n", step, current, nums[current])
		fmt.Printf("  选择元素 nums[%d] = %d\n", current, nums[current])
		lis[index] = nums[current]
		index--
		current = prev[current]
		step++
	}

	fmt.Printf("\n最终结果: %v\n", lis)
}

// LISComparisonExample 比较不同算法的性能
func LISComparisonExample() {
	fmt.Println("=== LIS算法性能比较 ===\n")

	nums := []int{10, 9, 2, 5, 3, 7, 101, 18, 19, 20, 21, 22, 23, 24, 25}

	fmt.Printf("测试数组: %v (长度: %d)\n", nums, len(nums))

	// 测试不同算法
	fmt.Println("\n算法比较:")

	// 动态规划
	result1 := LengthOfLIS(nums)
	fmt.Printf("1. 动态规划: LIS长度 = %d\n", result1)

	// 优化版本
	result2 := LengthOfLISOptimized(nums)
	fmt.Printf("2. 二分查找优化: LIS长度 = %d\n", result2)

	// 递归版本
	result3 := LengthOfLISRecursive(nums)
	fmt.Printf("3. 递归记忆化: LIS长度 = %d\n", result3)

	// 验证结果一致性
	if result1 == result2 && result2 == result3 {
		fmt.Println("\n✓ 所有算法结果一致")
	} else {
		fmt.Println("\n✗ 算法结果不一致")
	}
}

// LISApplicationsExample 演示LIS的实际应用
func LISApplicationsExample() {
	fmt.Println("=== LIS实际应用演示 ===\n")

	// 1. 股票价格分析
	fmt.Println("1. 股票价格分析:")
	prices := []int{100, 80, 120, 90, 130, 110, 140, 95, 150}
	fmt.Printf("   股票价格序列: %v\n", prices)
	lisLength := LengthOfLIS(prices)
	lisSequence := GetLongestIncreasingSubsequence(prices)
	fmt.Printf("   最长上涨序列长度: %d\n", lisLength)
	fmt.Printf("   最长上涨序列: %v\n", lisSequence)

	// 2. 身高排序问题
	fmt.Println("\n2. 身高排序问题:")
	heights := []int{160, 165, 170, 155, 175, 180, 150, 185}
	fmt.Printf("   身高序列: %v\n", heights)
	heightLIS := LengthOfLIS(heights)
	fmt.Printf("   最长递增身高序列长度: %d\n", heightLIS)

	// 3. 带约束条件的LIS
	fmt.Println("\n3. 带约束条件的LIS:")
	constraintNums := []int{1, 3, 6, 7, 9, 4, 10, 5, 6}
	k := 3 // 相邻元素差值不能超过3
	constraintLIS := LengthOfLISWithConstraints(constraintNums, k)
	fmt.Printf("   原数组: %v\n", constraintNums)
	fmt.Printf("   约束条件: 相邻元素差值 ≤ %d\n", k)
	fmt.Printf("   满足约束的最长递增子序列长度: %d\n", constraintLIS)

	// 4. 二维LIS问题
	fmt.Println("\n4. 二维LIS问题:")
	points := []Point{
		{1, 1}, {2, 3}, {3, 2}, {4, 4}, {5, 1},
		{6, 5}, {7, 3}, {8, 6}, {9, 2}, {10, 7},
	}
	fmt.Printf("   二维点集: %v\n", points)
	lis2D := LengthOfLIS2D(points)
	fmt.Printf("   二维LIS长度: %d\n", lis2D)

	// 5. 统计LIS数量
	fmt.Println("\n5. 统计LIS数量:")
	countNums := []int{1, 3, 5, 4, 7}
	lisCount := CountLIS(countNums)
	fmt.Printf("   数组: %v\n", countNums)
	fmt.Printf("   最长递增子序列的数量: %d\n", lisCount)
}

// DetailedLISExplanation 详细解说LengthOfLIS算法
func DetailedLISExplanation() {
	fmt.Println("🔍 LengthOfLIS 算法详细解说")
	fmt.Println("=" + strings.Repeat("=", 60))

	// 使用具体例子
	nums := []int{10, 9, 2, 5, 3, 7, 101, 18}

	fmt.Printf("示例: nums = %v\n\n", nums)

	// 算法概述
	fmt.Println("📋 算法概述:")
	fmt.Println("这个算法使用动态规划来解决最长递增子序列问题。")
	fmt.Println("LIS是指在一个序列中找到一个最长的子序列，使得这个子序列中的数字严格递增。")

	// 逐行代码解释
	fmt.Println("\n📝 逐行代码解释:")

	fmt.Println("\n```go")
	fmt.Println("func LengthOfLIS(nums []int) int {")
	fmt.Println("```")
	fmt.Println("第1行: 函数定义，接收一个整数数组，返回LIS的长度")

	fmt.Println("\n```go")
	fmt.Println("if len(nums) == 0 {")
	fmt.Println("    return 0")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println("第2-4行: 边界条件检查")
	fmt.Println("   - 如果数组为空，返回0")

	fmt.Println("\n```go")
	fmt.Println("dp := make([]int, len(nums))")
	fmt.Println("```")
	fmt.Printf("第5行: 创建动态规划数组\n")
	fmt.Printf("   - 创建长度为 %d 的数组\n", len(nums))
	fmt.Println("   - dp[i] 表示以 nums[i] 结尾的最长递增子序列长度")

	fmt.Println("\n```go")
	fmt.Println("for i := range dp {")
	fmt.Println("    dp[i] = 1")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println("第6-8行: 初始化DP数组")
	fmt.Println("   - 每个位置初始化为1")
	fmt.Println("   - 因为单个元素本身就是长度为1的递增子序列")

	fmt.Println("\n```go")
	fmt.Println("maxLen := 1")
	fmt.Println("```")
	fmt.Println("第9行: 记录全局最大长度")

	fmt.Println("\n```go")
	fmt.Println("for i := 1; i < len(nums); i++ {")
	fmt.Println("```")
	fmt.Println("第10行: 从第二个元素开始遍历")
	fmt.Println("   - 因为第一个元素的LIS长度已经确定为1")

	fmt.Println("\n```go")
	fmt.Println("for j := 0; j < i; j++ {")
	fmt.Println("```")
	fmt.Println("第11行: 遍历当前元素之前的所有元素")
	fmt.Println("   - 寻找可以接在当前元素前面的递增子序列")

	fmt.Println("\n```go")
	fmt.Println("if nums[i] > nums[j] {")
	fmt.Println("```")
	fmt.Println("第12行: 检查是否可以接在nums[j]后面")
	fmt.Println("   - 只有当nums[i] > nums[j]时，nums[i]才能接在nums[j]后面")

	fmt.Println("\n```go")
	fmt.Println("dp[i] = max(dp[i], dp[j]+1)")
	fmt.Println("```")
	fmt.Println("第13行: 状态转移方程")
	fmt.Println("   - 如果nums[i]可以接在nums[j]后面，则dp[i] = max(dp[i], dp[j]+1)")
	fmt.Println("   - 这表示选择更长的递增子序列")

	// 完整执行过程演示
	fmt.Println("\n🔄 完整执行过程演示:")
	fmt.Println("让我们逐步填充DP数组:")

	// 构建DP数组
	dp := make([]int, len(nums))
	for i := range dp {
		dp[i] = 1
	}

	fmt.Printf("\n初始状态: dp = %v\n", dp)

	step := 1
	for i := 1; i < len(nums); i++ {
		fmt.Printf("\n步骤%d: i=%d, nums[%d]=%d\n", step, i, i, nums[i])
		fmt.Printf("  检查 nums[%d]=%d 是否可以接在之前的元素后面:\n", i, nums[i])

		for j := 0; j < i; j++ {
			fmt.Printf("    j=%d, nums[%d]=%d: ", j, j, nums[j])
			if nums[i] > nums[j] {
				oldVal := dp[i]
				dp[i] = max(dp[i], dp[j]+1)
				fmt.Printf("nums[%d]=%d > nums[%d]=%d ✓\n", i, nums[i], j, nums[j])
				fmt.Printf("      dp[%d] = max(dp[%d], dp[%d]+1) = max(%d, %d+1) = %d\n",
					i, i, j, oldVal, dp[j], dp[i])
			} else {
				fmt.Printf("nums[%d]=%d <= nums[%d]=%d ✗\n", i, nums[i], j, nums[j])
			}
		}
		fmt.Printf("  最终 dp[%d] = %d\n", i, dp[i])
		step++
	}

	fmt.Printf("\n最终DP数组: dp = %v\n", dp)

	// 找到最大值
	maxLen := 1
	for _, val := range dp {
		maxLen = max(maxLen, val)
	}

	fmt.Println("\n```go")
	fmt.Println("return maxLen")
	fmt.Println("```")
	fmt.Printf("第17行: 返回最终结果\n")
	fmt.Printf("   - maxLen = %d 就是LIS的长度\n", maxLen)

	// 结果验证
	fmt.Println("\n✅ 结果验证:")
	fmt.Printf("   - LIS长度: %d\n", maxLen)
	lisSequence := GetLongestIncreasingSubsequence(nums)
	fmt.Printf("   - 实际LIS序列: %v (可以通过回溯获得)\n", lisSequence)
	fmt.Printf("   - 验证: %v 是严格递增的，且长度为%d\n", lisSequence, maxLen)

	// 算法复杂度
	fmt.Println("\n📈 算法复杂度:")
	fmt.Println("   - 时间复杂度: O(n²) - 需要两层嵌套循环")
	fmt.Println("   - 空间复杂度: O(n) - 需要存储DP数组")

	// 优化版本说明
	fmt.Println("\n🚀 优化版本 (二分查找):")
	fmt.Println("   - 时间复杂度: O(n log n)")
	fmt.Println("   - 使用tails数组维护递增子序列的最小末尾值")
	fmt.Println("   - 通过二分查找优化插入过程")

	// 核心思想
	fmt.Println("\n💡 核心思想:")
	fmt.Println("   将大问题分解为小问题，通过填表的方式自底向上解决。")
	fmt.Println("   每个dp[i]的值都依赖于其前面所有满足条件的dp[j]值，")
	fmt.Println("   体现了动态规划的最优子结构性质。")
}

// RunLISExamples 运行所有示例
func RunLISExamples() {
	LISExample()
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	LISBacktrackingExample()
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	LISComparisonExample()
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	LISApplicationsExample()
}

// RunDetailedLISExplanation 运行详细解说
func RunDetailedLISExplanation() {
	DetailedLISExplanation()
}
