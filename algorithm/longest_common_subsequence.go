package algorithm

import (
	"fmt"
	"strings"
)

// LongestCommonSubsequence 计算两个字符串的最长公共子序列长度
// 使用动态规划解法，时间复杂度O(mn)，空间复杂度O(mn)
func LongestCommonSubsequence(text1, text2 string) int {
	m, n := len(text1), len(text2)

	// 创建DP表，dp[i][j]表示text1[0:i]和text2[0:j]的LCS长度
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	// 填充DP表
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				// 当前字符匹配，LCS长度+1
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				// 当前字符不匹配，取两种情况的最大值
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return dp[m][n]
}

// LongestCommonSubsequenceOptimized 空间优化版本
// 时间复杂度O(mn)，空间复杂度O(n)
func LongestCommonSubsequenceOptimized(text1, text2 string) int {
	m, n := len(text1), len(text2)

	// 只使用一维数组，节省空间
	dp := make([]int, n+1)

	for i := 1; i <= m; i++ {
		prev := 0 // 保存dp[i-1][j-1]的值
		for j := 1; j <= n; j++ {
			temp := dp[j] // 保存dp[i-1][j]的值
			if text1[i-1] == text2[j-1] {
				dp[j] = prev + 1
			} else {
				dp[j] = max(dp[j], dp[j-1])
			}
			prev = temp
		}
	}

	return dp[n]
}

// LongestCommonSubsequenceRecursive 递归解法（带记忆化）
// 时间复杂度O(mn)，空间复杂度O(mn)
func LongestCommonSubsequenceRecursive(text1, text2 string) int {
	m, n := len(text1), len(text2)

	// 记忆化数组
	memo := make([][]int, m)
	for i := range memo {
		memo[i] = make([]int, n)
		for j := range memo[i] {
			memo[i][j] = -1
		}
	}

	return lcsHelper(text1, text2, m-1, n-1, memo)
}

// lcsHelper 递归辅助函数
func lcsHelper(text1, text2 string, i, j int, memo [][]int) int {
	// 边界条件
	if i < 0 || j < 0 {
		return 0
	}

	// 检查记忆化
	if memo[i][j] != -1 {
		return memo[i][j]
	}

	var result int
	if text1[i] == text2[j] {
		// 当前字符匹配
		result = lcsHelper(text1, text2, i-1, j-1, memo) + 1
	} else {
		// 当前字符不匹配，取两种情况的最大值
		result = max(
			lcsHelper(text1, text2, i-1, j, memo),
			lcsHelper(text1, text2, i, j-1, memo),
		)
	}

	memo[i][j] = result
	return result
}

// GetLongestCommonSubsequence 获取具体的LCS序列
// 返回最长公共子序列的字符串
func GetLongestCommonSubsequence(text1, text2 string) string {
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

	// 回溯构建LCS序列
	lcs := make([]byte, dp[m][n])
	index := dp[m][n] - 1
	i, j := m, n

	for i > 0 && j > 0 {
		if text1[i-1] == text2[j-1] {
			lcs[index] = text1[i-1]
			index--
			i--
			j--
		} else if dp[i-1][j] > dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	return string(lcs)
}

// LongestCommonSubstring 最长公共子串
// 注意：子串必须是连续的，而子序列可以不连续
func LongestCommonSubstring(text1, text2 string) int {
	m, n := len(text1), len(text2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	maxLen := 0
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				maxLen = max(maxLen, dp[i][j])
			}
		}
	}

	return maxLen
}

// GetLongestCommonSubstring 获取最长公共子串
func GetLongestCommonSubstring(text1, text2 string) string {
	m, n := len(text1), len(text2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	maxLen := 0
	endPos := 0

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
					endPos = i - 1
				}
			}
		}
	}

	if maxLen == 0 {
		return ""
	}

	startPos := endPos - maxLen + 1
	return text1[startPos : endPos+1]
}

// LongestCommonSubsequenceMultiple 多个序列的LCS（简化版本）
func LongestCommonSubsequenceMultiple(sequences []string) int {
	if len(sequences) == 0 {
		return 0
	}
	if len(sequences) == 1 {
		return len(sequences[0])
	}

	// 先计算前两个序列的LCS
	result := LongestCommonSubsequence(sequences[0], sequences[1])

	// 逐步与其他序列计算LCS
	for i := 2; i < len(sequences); i++ {
		// 简化处理：直接计算与当前序列的LCS
		result = min(result, LongestCommonSubsequence(sequences[0], sequences[i]))
	}

	return result
}

// WeightedLongestCommonSubsequence 带权重的LCS
func WeightedLongestCommonSubsequence(text1, text2 string, weights map[byte]int) int {
	m, n := len(text1), len(text2)

	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if text1[i-1] == text2[j-1] {
				weight := weights[text1[i-1]]
				dp[i][j] = dp[i-1][j-1] + weight
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	return dp[m][n]
}

// CompareDNASequences DNA序列比对，计算相似度
func CompareDNASequences(seq1, seq2 string) float64 {
	lcsLength := LongestCommonSubsequence(seq1, seq2)
	maxLength := max(len(seq1), len(seq2))

	if maxLength == 0 {
		return 1.0
	}

	// 计算相似度
	similarity := float64(lcsLength) / float64(maxLength)
	return similarity
}

// CompareFiles 文件差异比较（简化版本）
func CompareFiles(content1, content2 string) (int, string) {
	// 按行分割
	lines1 := strings.Split(content1, "\n")
	lines2 := strings.Split(content2, "\n")

	// 计算LCS
	lcsLength := LongestCommonSubsequence(strings.Join(lines1, ""), strings.Join(lines2, ""))
	lcsString := GetLongestCommonSubsequence(strings.Join(lines1, ""), strings.Join(lines2, ""))

	return lcsLength, lcsString
}

// PrintDPTable 打印DP表（用于调试）
func PrintDPTable1(text1, text2 string) {
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

	// 打印表头
	print("    ")
	for j := 0; j <= n; j++ {
		if j == 0 {
			print("  ")
		} else {
			fmt.Printf(" %c ", text2[j-1])
		}
	}
	println()

	// 打印DP表
	for i := 0; i <= m; i++ {
		if i == 0 {
			print("  ")
		} else {
			fmt.Printf(" %c ", text1[i-1])
		}

		for j := 0; j <= n; j++ {
			fmt.Printf(" %d ", dp[i][j])
		}
		println()
	}
}
