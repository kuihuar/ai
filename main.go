package main

import (
	"fmt"
	"strings"
	"time"

	//"./algorithm"
	"github.com/kuihuar/ai/algorithm"
)

func main() {
	mainSubsequence()
	// closure.DemoClosure()
	// return
}

func readData(ch <-chan string) {

	val, ok := <-ch
	if !ok {
		return
	}
	fmt.Println(val)
}
func mainSubsequence() {

	fmt.Println("🚀 最长公共子序列 (LCS) 算法演示程序")
	fmt.Println("=" + strings.Repeat("=", 50))

	for {
		fmt.Println("\n请选择操作:")
		fmt.Println("1. 运行预设示例")
		fmt.Println("2. 交互式LCS计算")
		fmt.Println("3. 性能测试")
		fmt.Println("4. 详细回溯演示")
		fmt.Println("5. 实际应用演示")
		fmt.Println("0. 退出程序")

		var choice int
		fmt.Print("\n请输入选择 (0-5): ")
		fmt.Scanf("%d", &choice)

		switch choice {
		case 1:
			runPresetExamples()
		// case 2:
		// 	interactiveLCS()
		case 3:
			performanceTest()
		// case 4:
		// 	backtrackingDemo()
		case 5:
			applicationsDemo()
		case 0:
			fmt.Println("👋 感谢使用！")
			return
		default:
			fmt.Println("❌ 无效选择，请重新输入")
		}
	}
}

func runPresetExamples() {
	fmt.Println("\n📋 运行预设示例...")

	examples := []struct {
		text1, text2 string
		description  string
	}{
		{"abcde", "ace", "基本示例"},
		{"abc", "abc", "完全相同字符串"},
		{"abc", "def", "无公共子序列"},
		{"GeeksforGeeks", "GeeksQuiz", "文本相似度"},
		{"ATCGATCG", "ATCGATCG", "DNA序列"},
		{"ABCDGH", "AEDFHR", "经典示例"},
	}

	for i, example := range examples {
		fmt.Printf("\n--- 示例 %d: %s ---\n", i+1, example.description)
		fmt.Printf("字符串1: %s\n", example.text1)
		fmt.Printf("字符串2: %s\n", example.text2)

		// 计算LCS
		lcsLength := algorithm.LongestCommonSubsequence(example.text1, example.text2)
		lcsString := algorithm.GetLongestCommonSubsequence(example.text1, example.text2)

		fmt.Printf("最长公共子序列长度: %d\n", lcsLength)
		fmt.Printf("最长公共子序列: %s\n", lcsString)

		// 计算相似度
		similarity := float64(lcsLength) / float64(max(len(example.text1), len(example.text2)))
		fmt.Printf("相似度: %.2f (%.1f%%)\n", similarity, similarity*100)
	}
}

// func interactiveLCS() {
// 	fmt.Println("\n🎯 交互式LCS计算")

// 	reader := bufio.NewReader(os.Stdin)

// 	fmt.Print("请输入第一个字符串: ")
// 	text1, _ := reader.ReadString('\n')
// 	text1 = strings.TrimSpace(text1)

// 	fmt.Print("请输入第二个字符串: ")
// 	text2, _ := reader.ReadString('\n')
// 	text2 = strings.TrimSpace(text2)

// 	if text1 == "" || text2 == "" {
// 		fmt.Println("❌ 字符串不能为空")
// 		return
// 	}

// 	fmt.Printf("\n计算结果:\n")
// 	fmt.Printf("字符串1: %s (长度: %d)\n", text1, len(text1))
// 	fmt.Printf("字符串2: %s (长度: %d)\n", text2, len(text2))

// 	// 计算LCS
// 	lcsLength := algorithm.LongestCommonSubsequence(text1, text2)
// 	lcsString := algorithm.GetLongestCommonSubsequence(text1, text2)

// 	fmt.Printf("最长公共子序列长度: %d\n", lcsLength)
// 	fmt.Printf("最长公共子序列: %s\n", lcsString)

// 	// 计算相似度
// 	similarity := float64(lcsLength) / float64(max(len(text1), len(text2)))
// 	fmt.Printf("相似度: %.2f (%.1f%%)\n", similarity, similarity*100)

// 	// 显示DP表
// 	fmt.Print("\n是否显示动态规划表? (y/n): ")
// 	showDP, _ := reader.ReadString('\n')
// 	if strings.TrimSpace(strings.ToLower(showDP)) == "y" {
// 		fmt.Println("\n动态规划表:")
// 		algorithm.PrintDPTable(text1, text2)
// 	}
// }

func performanceTest() {
	fmt.Println("\n⚡ 性能测试")

	// 生成测试数据
	testCases := []struct {
		name  string
		text1 string
		text2 string
	}{
		{"小规模", "ABCDEF", "DEFGHI"},
		{"中等规模", "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "ZYXWVUTSRQPONMLKJIHGFEDCBA"},
		{"大规模", strings.Repeat("ABCDEFGHIJKLMNOPQRSTUVWXYZ", 10),
			strings.Repeat("ZYXWVUTSRQPONMLKJIHGFEDCBA", 10)},
	}

	for _, testCase := range testCases {
		fmt.Printf("\n--- %s测试 ---\n", testCase.name)
		fmt.Printf("字符串1长度: %d\n", len(testCase.text1))
		fmt.Printf("字符串2长度: %d\n", len(testCase.text2))

		// 测试动态规划版本
		start := time.Now()
		result1 := algorithm.LongestCommonSubsequence(testCase.text1, testCase.text2)
		duration1 := time.Since(start)

		// 测试空间优化版本
		start = time.Now()
		result2 := algorithm.LongestCommonSubsequenceOptimized(testCase.text1, testCase.text2)
		duration2 := time.Since(start)

		// 测试递归版本（仅对小规模数据）
		var duration3 time.Duration
		var result3 int
		if len(testCase.text1) <= 20 && len(testCase.text2) <= 20 {
			start = time.Now()
			result3 = algorithm.LongestCommonSubsequenceRecursive(testCase.text1, testCase.text2)
			duration3 = time.Since(start)
		}

		fmt.Printf("LCS长度: %d\n", result1)
		fmt.Printf("动态规划版本: %v\n", duration1)
		fmt.Printf("空间优化版本: %v\n", duration2)

		if len(testCase.text1) <= 20 && len(testCase.text2) <= 20 {
			fmt.Printf("递归记忆化版本: %v\n", duration3)
		} else {
			fmt.Printf("递归记忆化版本: 跳过（数据规模过大）\n")
		}

		// 验证结果一致性
		if result1 == result2 && (len(testCase.text1) > 20 || result1 == result3) {
			fmt.Println("✓ 结果验证通过")
		} else {
			fmt.Println("✗ 结果验证失败")
		}
	}
}

// func backtrackingDemo() {
// 	fmt.Println("\n🔍 详细回溯演示")

// 	text1 := "abcde"
// 	text2 := "ace"

// 	fmt.Printf("字符串1: %s\n", text1)
// 	fmt.Printf("字符串2: %s\n", text2)

// 	// 显示DP表
// 	fmt.Println("\n动态规划表:")
// 	algorithm.PrintDPTable(text1, text2)

// 	// 获取LCS
// 	lcsString := algorithm.GetLongestCommonSubsequence(text1, text2)
// 	fmt.Printf("\n最长公共子序列: %s\n", lcsString)

// 	// 手动演示回溯过程
// 	fmt.Println("\n回溯过程演示:")
// 	m, n := len(text1), len(text2)

// 	// 构建DP表
// 	dp := make([][]int, m+1)
// 	for i := range dp {
// 		dp[i] = make([]int, n+1)
// 	}

// 	for i := 1; i <= m; i++ {
// 		for j := 1; j <= n; j++ {
// 			if text1[i-1] == text2[j-1] {
// 				dp[i][j] = dp[i-1][j-1] + 1
// 			} else {
// 				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
// 			}
// 		}
// 	}

// 	// 回溯
// 	i, j := m, n
// 	step := 1
// 	for i > 0 && j > 0 {
// 		fmt.Printf("步骤%d: 当前位置 dp[%d][%d] = %d\n", step, i, j, dp[i][j])

// 		if text1[i-1] == text2[j-1] {
// 			fmt.Printf("  ✓ 字符匹配: '%c' == '%c'\n", text1[i-1], text2[j-1])
// 			fmt.Printf("  → 选择字符 '%c'，移动到 dp[%d][%d]\n", text1[i-1], i-1, j-1)
// 			i--
// 			j--
// 		} else if dp[i-1][j] > dp[i][j-1] {
// 			fmt.Printf("  ✗ 字符不匹配: '%c' != '%c'\n", text1[i-1], text2[j-1])
// 			fmt.Printf("  → dp[%d][%d] > dp[%d][%d]，移动到 dp[%d][%d]\n", i-1, j, i, j-1, i-1, j)
// 			i--
// 		} else {
// 			fmt.Printf("  ✗ 字符不匹配: '%c' != '%c'\n", text1[i-1], text2[j-1])
// 			fmt.Printf("  → dp[%d][%d] <= dp[%d][%d]，移动到 dp[%d][%d]\n", i-1, j, i, j-1, i, j-1)
// 			j--
// 		}
// 		step++
// 	}

// 	fmt.Printf("\n最终LCS: %s\n", lcsString)
// }

func applicationsDemo() {
	fmt.Println("\n🌍 实际应用演示")

	// 1. DNA序列比对
	fmt.Println("1. DNA序列比对:")
	dnaSequences := []struct {
		name, seq1, seq2 string
	}{
		{"相同序列", "ATCGATCG", "ATCGATCG"},
		{"相似序列", "ATCGATCG", "ATCGATCC"},
		{"不同序列", "ATCGATCG", "GCTAGCTA"},
	}

	for _, dna := range dnaSequences {
		similarity := algorithm.CompareDNASequences(dna.seq1, dna.seq2)
		fmt.Printf("   %s: %.2f (%.1f%%)\n", dna.name, similarity, similarity*100)
	}

	// 2. 文本相似度
	fmt.Println("\n2. 文本相似度:")
	texts := []struct {
		name, text1, text2 string
	}{
		{"相似文本", "GeeksforGeeks", "GeeksQuiz"},
		{"部分相似", "Hello World", "Hello Go"},
		{"不同文本", "Python Programming", "Java Development"},
	}

	for _, text := range texts {
		lcsLength := algorithm.LongestCommonSubsequence(text.text1, text.text2)
		similarity := float64(lcsLength) / float64(max(len(text.text1), len(text.text2)))
		fmt.Printf("   %s: %.2f (%.1f%%)\n", text.name, similarity, similarity*100)
	}

	// 3. 最长公共子串 vs 子序列
	fmt.Println("\n3. 最长公共子串 vs 最长公共子序列:")
	text1 := "GeeksforGeeks"
	text2 := "GeeksQuiz"

	lcsLength := algorithm.LongestCommonSubsequence(text1, text2)
	lcsString := algorithm.GetLongestCommonSubsequence(text1, text2)

	lcsSubstringLength := algorithm.LongestCommonSubstring(text1, text2)
	lcsSubstring := algorithm.GetLongestCommonSubstring(text1, text2)

	fmt.Printf("   文本: %s vs %s\n", text1, text2)
	fmt.Printf("   最长公共子序列: %s (长度: %d)\n", lcsString, lcsLength)
	fmt.Printf("   最长公共子串: %s (长度: %d)\n", lcsSubstring, lcsSubstringLength)

	// 4. 带权重的LCS
	fmt.Println("\n4. 带权重的LCS:")
	weights := map[byte]int{
		'A': 1, 'B': 2, 'C': 3, 'D': 4,
		'E': 5, 'F': 6, 'G': 7, 'H': 8,
	}

	weightedText1 := "ABC"
	weightedText2 := "ABD"

	weightedResult := algorithm.WeightedLongestCommonSubsequence(weightedText1, weightedText2, weights)
	fmt.Printf("   文本: %s vs %s\n", weightedText1, weightedText2)
	fmt.Printf("   带权重LCS值: %d\n", weightedResult)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
