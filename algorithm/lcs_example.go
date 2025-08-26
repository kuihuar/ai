package algorithm

import (
	"fmt"
	"strings"
)

// LCSExample 演示LCS算法的详细计算过程
func LCSExample() {
	fmt.Println("=== 最长公共子序列 (LCS) 算法演示 ===\n")

	// 示例1：基本示例
	demonstrateLCS("abcde", "ace", "示例1")

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例2：DNA序列比对
	demonstrateLCS("ATCGATCG", "ATCGATCG", "示例2 - 相同DNA序列")

	fmt.Println("\n" + strings.Repeat("=", 50) + "\n")

	// 示例3：文本相似度
	demonstrateLCS("GeeksforGeeks", "GeeksQuiz", "示例3 - 文本相似度")
}

// demonstrateLCS 演示LCS算法的详细过程
func demonstrateLCS(text1, text2, title string) {
	fmt.Printf("%s\n", title)
	fmt.Printf("字符串1: %s (长度: %d)\n", text1, len(text1))
	fmt.Printf("字符串2: %s (长度: %d)\n", text2, len(text2))

	// 计算LCS长度
	lcsLength := LongestCommonSubsequence(text1, text2)
	fmt.Printf("\n最长公共子序列长度: %d\n", lcsLength)

	// 获取具体的LCS序列
	lcsString := GetLongestCommonSubsequence(text1, text2)
	fmt.Printf("最长公共子序列: %s\n", lcsString)

	// 打印DP表
	fmt.Println("\n动态规划表:")
	printDPTableDetailed(text1, text2)

	// 计算相似度
	similarity := float64(lcsLength) / float64(max(len(text1), len(text2)))
	fmt.Printf("\n相似度: %.2f (%.1f%%)\n", similarity, similarity*100)
}

// printDPTableDetailed 打印详细的DP表
func printDPTableDetailed(text1, text2 string) {
	m, n := len(text1), len(text2)

	// 创建DP表
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// 填充DP表
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// 打印表头
	fmt.Print("    ")
	for j := 0; j <= n; j++ {
		if j == 0 {
			fmt.Print("  ")
		} else {
			fmt.Printf(" %c ", text2[j-1])
		}
	}
	fmt.Println()

	// 打印DP表
	for i := 0; i <= m; i++ {
		if i == 0 {
			fmt.Print("  ")
		} else {
			fmt.Printf(" %c ", text1[i-1])
		}

		for j := 0; j <= n; j++ {
			fmt.Printf(" %d ", dp[i][j])
		}
		fmt.Println()
	}

	// 解释DP表
	fmt.Println("\nDP表解释:")
	fmt.Println("- dp[i][j] 表示 text1[0:i] 和 text2[0:j] 的最长公共子序列长度")
	fmt.Println("- 当 text1[i-1] == text2[j-1] 时，dp[i][j] = dp[i-1][j-1] + 1")
	fmt.Println("- 否则，dp[i][j] = max(dp[i-1][j], dp[i][j-1])")
}

// LCSBacktrackingExample 演示LCS回溯过程
func LCSBacktrackingExample() {
	fmt.Println("=== LCS回溯过程演示 ===\n")

	text1 := "abcde"
	text2 := "ace"

	fmt.Printf("字符串1: %s\n", text1)
	fmt.Printf("字符串2: %s\n", text2)

	// 构建DP表
	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// 填充DP表
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	// 打印DP表
	fmt.Println("\nDP表:")
	printDPTableDetailed(text1, text2)

	// 回溯过程
	fmt.Println("\n回溯过程:")
	lcs := make([]byte, dp[m][n])
	index := dp[m][n] - 1
	i, j := m, n

	step := 1
	for i > 0 && j > 0 {
		fmt.Printf("步骤%d: 当前位置 dp[%d][%d] = %d\n", step, i, j, dp[i][j])

		if text1[i-1] == text2[j-1] {
			fmt.Printf("  字符匹配: text1[%d] = text2[%d] = '%c'\n", i-1, j-1, text1[i-1])
			fmt.Printf("  选择字符 '%c'，移动到 dp[%d][%d]\n", text1[i-1], i-1, j-1)
			lcs[index] = text1[i-1]
			index--
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			fmt.Printf("  字符不匹配，dp[%d][%d] > dp[%d][%d] (%d > %d)\n",
				i-1, j, i, j-1, dp[i-1][j], dp[i][j-1])
			fmt.Printf("  移动到 dp[%d][%d]\n", i-1, j)
			i--
		} else {
			fmt.Printf("  字符不匹配，dp[%d][%d] <= dp[%d][%d] (%d <= %d)\n",
				i-1, j, i, j-1, dp[i-1][j], dp[i][j-1])
			fmt.Printf("  移动到 dp[%d][%d]\n", i, j-1)
			j--
		}
		step++
	}

	fmt.Printf("\n最终结果: %s\n", string(lcs))
}

// LCSComparisonExample 比较不同算法的性能
func LCSComparisonExample() {
	fmt.Println("=== LCS算法性能比较 ===\n")

	text1 := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	text2 := "ZYXWVUTSRQPONMLKJIHGFEDCBA"

	fmt.Printf("测试字符串1: %s (长度: %d)\n", text1, len(text1))
	fmt.Printf("测试字符串2: %s (长度: %d)\n", text2, len(text2))

	// 测试不同算法
	fmt.Println("\n算法比较:")

	// 动态规划
	result1 := LongestCommonSubsequence(text1, text2)
	fmt.Printf("1. 动态规划: LCS长度 = %d\n", result1)

	// 空间优化版本
	result2 := LongestCommonSubsequenceOptimized(text1, text2)
	fmt.Printf("2. 空间优化: LCS长度 = %d\n", result2)

	// 递归版本
	result3 := LongestCommonSubsequenceRecursive(text1, text2)
	fmt.Printf("3. 递归记忆化: LCS长度 = %d\n", result3)

	// 验证结果一致性
	if result1 == result2 && result2 == result3 {
		fmt.Println("\n✓ 所有算法结果一致")
	} else {
		fmt.Println("\n✗ 算法结果不一致")
	}
}

// LCSApplicationsExample 演示LCS的实际应用
func LCSApplicationsExample() {
	fmt.Println("=== LCS实际应用演示 ===\n")

	// 1. DNA序列比对
	fmt.Println("1. DNA序列比对:")
	dna1 := "ATCGATCG"
	dna2 := "ATCGATCG"
	dna3 := "GCTAGCTA"

	similarity1 := CompareDNASequences(dna1, dna2)
	similarity2 := CompareDNASequences(dna1, dna3)

	fmt.Printf("   DNA序列1: %s\n", dna1)
	fmt.Printf("   DNA序列2: %s\n", dna2)
	fmt.Printf("   相似度: %.2f (%.1f%%)\n", similarity1, similarity1*100)

	fmt.Printf("   DNA序列1: %s\n", dna1)
	fmt.Printf("   DNA序列3: %s\n", dna3)
	fmt.Printf("   相似度: %.2f (%.1f%%)\n", similarity2, similarity2*100)

	// 2. 文本相似度
	fmt.Println("\n2. 文本相似度:")
	text1 := "GeeksforGeeks"
	text2 := "GeeksQuiz"
	text3 := "HelloWorld"

	lcs1 := LongestCommonSubsequence(text1, text2)
	lcs2 := LongestCommonSubsequence(text1, text3)

	similarity3 := float64(lcs1) / float64(max(len(text1), len(text2)))
	similarity4 := float64(lcs2) / float64(max(len(text1), len(text3)))

	fmt.Printf("   文本1: %s\n", text1)
	fmt.Printf("   文本2: %s\n", text2)
	fmt.Printf("   相似度: %.2f (%.1f%%)\n", similarity3, similarity3*100)

	fmt.Printf("   文本1: %s\n", text1)
	fmt.Printf("   文本3: %s\n", text3)
	fmt.Printf("   相似度: %.2f (%.1f%%)\n", similarity4, similarity4*100)

	// 3. 最长公共子串 vs 最长公共子序列
	fmt.Println("\n3. 最长公共子串 vs 最长公共子序列:")

	lcsLength := LongestCommonSubsequence(text1, text2)
	lcsString := GetLongestCommonSubsequence(text1, text2)

	lcsSubstringLength := LongestCommonSubstring(text1, text2)
	lcsSubstring := GetLongestCommonSubstring(text1, text2)

	fmt.Printf("   文本1: %s\n", text1)
	fmt.Printf("   文本2: %s\n", text2)
	fmt.Printf("   最长公共子序列: %s (长度: %d)\n", lcsString, lcsLength)
	fmt.Printf("   最长公共子串: %s (长度: %d)\n", lcsSubstring, lcsSubstringLength)
}

// 运行所有示例
func RunLCSExamples() {
	LCSExample()
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	LCSBacktrackingExample()
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	LCSComparisonExample()
	fmt.Println("\n" + strings.Repeat("=", 80) + "\n")

	LCSApplicationsExample()
}

// DetailedLCSExplanation 详细解说LongestCommonSubsequence算法
func DetailedLCSExplanation() {
	fmt.Println("🔍 LongestCommonSubsequence 算法详细解说")
	fmt.Println("=" + strings.Repeat("=", 60))

	// 使用具体例子
	text1 := "abcde"
	text2 := "ace"

	fmt.Printf("示例: text1 = \"%s\", text2 = \"%s\"\n\n", text1, text2)

	// 算法概述
	fmt.Println("📋 算法概述:")
	fmt.Println("这个算法使用动态规划来解决最长公共子序列问题。")
	fmt.Println("LCS是指两个字符串中按原顺序出现的最长公共字符序列")
	fmt.Println("（字符可以不连续，但必须保持相对顺序）。")

	// 逐行代码解释
	fmt.Println("\n📝 逐行代码解释:")

	fmt.Println("\n```go")
	fmt.Println("func LongestCommonSubsequence(text1, text2 string) int {")
	fmt.Println("```")
	fmt.Println("第1行: 函数定义，接收两个字符串参数，返回LCS的长度")

	fmt.Println("\n```go")
	fmt.Println("m, n := len(text1), len(text2)")
	fmt.Println("```")
	fmt.Printf("第2行: 获取两个字符串的长度\n")
	fmt.Printf("   - m = %d (text1 \"%s\" 的长度)\n", len(text1), text1)
	fmt.Printf("   - n = %d (text2 \"%s\" 的长度)\n", len(text2), text2)

	fmt.Println("\n```go")
	fmt.Println("dp := make([][]int, m+1)")
	fmt.Println("for i := range dp {")
	fmt.Println("    dp[i] = make([]int, n+1)")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println("第3-6行: 创建动态规划表")
	fmt.Printf("   - 创建 (%d+1) × (%d+1) = %d×%d 的二维数组\n", len(text1), len(text2), len(text1)+1, len(text2)+1)
	fmt.Println("   - dp[i][j] 表示 text1[0:i] 和 text2[0:j] 的LCS长度")
	fmt.Println("   - 初始时所有值都是0")

	// 显示DP表结构
	fmt.Println("\n📊 DP表结构:")
	fmt.Println("    \"\"  \"a\"  \"c\"  \"e\"")
	fmt.Println("\"\"   0   0   0   0")
	fmt.Println("\"a\"  0   ?   ?   ?")
	fmt.Println("\"b\"  0   ?   ?   ?")
	fmt.Println("\"c\"  0   ?   ?   ?")
	fmt.Println("\"d\"  0   ?   ?   ?")
	fmt.Println("\"e\"  0   ?   ?   ?")

	fmt.Println("\n```go")
	fmt.Println("for i := 1; i <= m; i++ {")
	fmt.Println("    for j := 1; j <= n; j++ {")
	fmt.Println("```")
	fmt.Println("第7-8行: 双重循环遍历DP表")
	fmt.Println("   - 从 i=1, j=1 开始（跳过第一行和第一列，它们都是0）")

	fmt.Println("\n```go")
	fmt.Println("if text1[i-1] == text2[j-1] {")
	fmt.Println("```")
	fmt.Println("第9行: 检查当前字符是否匹配")
	fmt.Println("   - text1[i-1] 是text1的第i个字符（因为数组索引从0开始）")
	fmt.Println("   - text2[j-1] 是text2的第j个字符")

	// 举例说明字符匹配
	fmt.Println("\n举例说明:")
	fmt.Printf("   - 当 i=1, j=1: text1[0]='%c', text2[0]='%c' → 匹配 ✓\n", text1[0], text2[0])
	fmt.Printf("   - 当 i=2, j=1: text1[1]='%c', text2[0]='%c' → 不匹配 ✗\n", text1[1], text2[0])

	fmt.Println("\n```go")
	fmt.Println("dp[i][j] = dp[i-1][j-1] + 1")
	fmt.Println("```")
	fmt.Println("第10行: 字符匹配时的状态转移")
	fmt.Println("   - 如果当前字符匹配，LCS长度 = 左上角的值 + 1")
	fmt.Println("   - 这表示在之前LCS的基础上加上当前匹配的字符")

	fmt.Println("\n```go")
	fmt.Println("} else {")
	fmt.Println("    dp[i][j] = max(dp[i-1][j], dp[i][j-1])")
	fmt.Println("}")
	fmt.Println("```")
	fmt.Println("第11-13行: 字符不匹配时的状态转移")
	fmt.Println("   - 取上方和左方的最大值")
	fmt.Println("   - 这表示选择更优的子问题解")

	// 完整执行过程演示
	fmt.Println("\n🔄 完整执行过程演示:")
	fmt.Println("让我们逐步填充DP表:")

	// 构建DP表
	m, n := len(text1), len(text2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	step := 1
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			fmt.Printf("\n步骤%d: i=%d, j=%d (text1[%d]='%c', text2[%d]='%c')\n",
				step, i, j, i-1, text1[i-1], j-1, text2[j-1])

			if text1[i-1] == text2[j-1] {
				fmt.Printf("  字符匹配: '%c' == '%c'\n", text1[i-1], text2[j-1])
				fmt.Printf("  dp[%d][%d] = dp[%d][%d] + 1 = %d + 1 = %d\n",
					i, j, i-1, j-1, dp[i-1][j-1], dp[i-1][j-1]+1)
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				fmt.Printf("  字符不匹配: '%c' != '%c'\n", text1[i-1], text2[j-1])
				fmt.Printf("  dp[%d][%d] = max(dp[%d][%d], dp[%d][%d]) = max(%d, %d) = %d\n",
					i, j, i-1, j, i, j-1, dp[i-1][j], dp[i][j-1], max(dp[i-1][j], dp[i][j-1]))
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
			step++
		}
	}

	// 显示最终DP表
	fmt.Println("\n📊 最终DP表:")
	fmt.Print("    ")
	for j := 0; j <= n; j++ {
		if j == 0 {
			fmt.Print("  ")
		} else {
			fmt.Printf(" %c ", text2[j-1])
		}
	}
	fmt.Println()

	for i := 0; i <= m; i++ {
		if i == 0 {
			fmt.Print("  ")
		} else {
			fmt.Printf(" %c ", text1[i-1])
		}

		for j := 0; j <= n; j++ {
			fmt.Printf(" %d ", dp[i][j])
		}
		fmt.Println()
	}

	fmt.Println("\n```go")
	fmt.Println("return dp[m][n]")
	fmt.Println("```")
	fmt.Printf("第15行: 返回最终结果\n")
	fmt.Printf("   - dp[%d][%d] = %d 就是LCS的长度\n", m, n, dp[m][n])

	// 结果验证
	fmt.Println("\n✅ 结果验证:")
	fmt.Printf("   - LCS长度: %d\n", dp[m][n])
	lcsString := GetLongestCommonSubsequence(text1, text2)
	fmt.Printf("   - 实际LCS序列: \"%s\" (可以通过回溯获得)\n", lcsString)
	fmt.Printf("   - 验证: \"%s\" 是 \"%s\" 和 \"%s\" 的公共子序列，且长度为%d\n",
		lcsString, text1, text2, dp[m][n])

	// 算法复杂度
	fmt.Println("\n📈 算法复杂度:")
	fmt.Println("   - 时间复杂度: O(m×n) - 需要填充整个DP表")
	fmt.Println("   - 空间复杂度: O(m×n) - 需要存储整个DP表")

	// 核心思想
	fmt.Println("\n💡 核心思想:")
	fmt.Println("   将大问题分解为小问题，通过填表的方式自底向上解决。")
	fmt.Println("   每个dp[i][j]的值都依赖于其左上方、上方、左方的值，")
	fmt.Println("   体现了动态规划的最优子结构性质。")
}

// RunDetailedLCSExplanation 运行详细解说
func RunDetailedLCSExplanation() {
	DetailedLCSExplanation()
}
