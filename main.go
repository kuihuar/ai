package main

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	//"./algorithm"
	"github.com/kuihuar/ai/algorithm"
)

func modifyArr(arr *[3]int) {
	for i := 0; i < len(arr); i++ {
		arr[i] *= 2
	}
}
func modifyArrayWithPointer(arr *[3]int) {
	for i := 0; i < len(arr); i++ {
		(*arr)[i] *= 2
	}
}

func main() {
	fmt.Println("=== Go GMP 模型详解 ===")

	explainGMPModel()

	// 演示 GMP 实际运行
	fmt.Println("\n\n=== GMP 实际运行演示 ===")
	demonstrateGMP()
}

// 解释三色标记算法
func explainThreeColorMarking() {
	fmt.Println("\n📚 三色标记算法（Tri-color Marking）")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n🎨 三种颜色：")
	fmt.Println("   1. ⚪ 白色（White）：未被访问的对象")
	fmt.Println("      → 表示不可达对象")
	fmt.Println("      → ✅ 将被删除（回收）")
	fmt.Println()
	fmt.Println("   2. ⚫ 灰色（Gray）：已被访问，但引用的对象还未扫描")
	fmt.Println("      → 表示正在处理的对象")
	fmt.Println("      → 需要继续扫描其引用的对象")
	fmt.Println()
	fmt.Println("   3. ⚫ 黑色（Black）：已被访问，且所有引用的对象都已扫描")
	fmt.Println("      → 表示可达对象（不会被回收）")
	fmt.Println("      → 所有引用都已处理完成")

	fmt.Println("\n🔄 GC 执行流程：")
	fmt.Println("   阶段1: 初始标记（Stop The World，短暂）")
	fmt.Println("      → 所有对象标记为白色")
	fmt.Println("      → 从根对象（全局变量、栈变量等）开始，标记为灰色")
	fmt.Println()
	fmt.Println("   阶段2: 并发标记（与程序并发执行）")
	fmt.Println("      → 从灰色对象队列中取出对象")
	fmt.Println("      → 扫描该对象引用的所有对象")
	fmt.Println("      → 将引用的对象标记为灰色（如果还是白色）")
	fmt.Println("      → 将当前对象标记为黑色")
	fmt.Println("      → 重复直到灰色队列为空")
	fmt.Println()
	fmt.Println("   阶段3: 标记完成（Stop The World，短暂）")
	fmt.Println("      → 处理在并发标记期间新分配的对象")
	fmt.Println("      → 重新扫描可能被修改的栈")
	fmt.Println()
	fmt.Println("   阶段4: 清除（与程序并发执行）")
	fmt.Println("      → ✅ 删除所有白色对象（不可达对象）")
	fmt.Println("      → 保留黑色对象（可达对象）")

	fmt.Println("\n💡 关键点：")
	fmt.Println("   ✅ 删除的是：⚪ 白色对象（不可达对象）")
	fmt.Println("   ✅ 保留的是：⚫ 黑色对象（可达对象）")
	fmt.Println("   ⚠️  灰色对象：正在处理中，最终会变成黑色")
}

// 演示实际 GC
func demonstrateGC() {
	// 创建一些对象来演示
	fmt.Println("\n1. 创建对象...")

	// 可达对象（会被保留）
	reachable := make([]byte, 1024*1024) // 1MB
	_ = reachable

	// 不可达对象（会被回收）
	func() {
		unreachable := make([]byte, 10*1024*1024) // 10MB
		_ = unreachable
		// 函数返回后，unreachable 变为不可达
	}()

	fmt.Println("   创建了可达对象（1MB）和不可达对象（10MB）")

	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	fmt.Printf("\n2. GC 前堆内存: %d KB\n", m1.HeapAlloc/1024)

	fmt.Println("\n3. 触发 GC...")
	runtime.GC()

	runtime.ReadMemStats(&m2)
	fmt.Printf("4. GC 后堆内存: %d KB\n", m2.HeapAlloc/1024)
	fmt.Printf("   ✅ 回收了约 %d KB 内存（白色对象被删除）\n",
		(m1.HeapAlloc-m2.HeapAlloc)/1024)
	fmt.Printf("   ✅ GC 次数: %d\n", m2.NumGC)
}

// 分配内存的函数
func allocateMemory(size int) {
	data := make([]byte, size)
	// 使用数据，避免被优化掉
	for i := range data {
		data[i] = byte(i % 256)
	}
	_ = data // 数据变为不可达，等待 GC
}

// 监控 GC 事件
func monitorGC() {
	var lastGC uint32
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		if m.NumGC > lastGC {
			fmt.Printf("   [GC] 触发！次数: %d, 堆内存: %d KB\n",
				m.NumGC, m.HeapAlloc/1024)
			lastGC = m.NumGC
		}
	}
}

// 解释 GMP 模型
func explainGMPModel() {
	fmt.Println("\n📚 Go GMP 并发模型")
	fmt.Println(strings.Repeat("=", 60))

	fmt.Println("\n🔤 GMP 三个核心组件：")
	fmt.Println("   G - Goroutine（协程）")
	fmt.Println("      → 轻量级线程，Go 程序的执行单元")
	fmt.Println("      → 由 Go 运行时管理，不是 OS 线程")
	fmt.Println("      → 初始栈大小：2KB（可动态增长）")
	fmt.Println()
	fmt.Println("   M - Machine（机器/OS 线程）")
	fmt.Println("      → 真正的操作系统线程")
	fmt.Println("      → 由操作系统调度")
	fmt.Println("      → 执行 G 的代码")
	fmt.Println()
	fmt.Println("   P - Processor（处理器/上下文）")
	fmt.Println("      → 逻辑处理器，管理 G 的执行")
	fmt.Println("      → 数量 = GOMAXPROCS（默认 = CPU 核心数）")
	fmt.Println("      → 包含本地 G 队列、运行队列等")

	fmt.Println("\n🔗 GMP 关系：")
	fmt.Println("   M 必须绑定 P 才能执行 G")
	fmt.Println("   P 管理一组 G（本地队列）")
	fmt.Println("   多个 M 可以绑定同一个 P（但同一时刻只有一个 M 在工作）")

	fmt.Println("\n📊 默认配置：")
	fmt.Printf("   CPU 核心数: %d\n", runtime.NumCPU())
	fmt.Printf("   GOMAXPROCS (P 的数量): %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("   当前 Goroutine 数: %d\n", runtime.NumGoroutine())

	fmt.Println("\n🔄 GMP 调度流程：")
	fmt.Println("   1. 创建 G（Goroutine）")
	fmt.Println("      → G 被放入某个 P 的本地队列")
	fmt.Println("      → 或放入全局队列（如果本地队列满）")
	fmt.Println()
	fmt.Println("   2. M 获取 G")
	fmt.Println("      → M 绑定 P 后，从 P 的本地队列获取 G")
	fmt.Println("      → 如果本地队列为空，从全局队列获取")
	fmt.Println("      → 如果全局队列也为空，从其他 P 偷取（work-stealing）")
	fmt.Println()
	fmt.Println("   3. M 执行 G")
	fmt.Println("      → M 执行 G 的代码")
	fmt.Println("      → G 可能阻塞（系统调用、channel 操作等）")
	fmt.Println()
	fmt.Println("   4. G 执行完成或阻塞")
	fmt.Println("      → 如果完成：G 结束，M 继续获取下一个 G")
	fmt.Println("      → 如果阻塞：M 和 G 解绑，M 可以执行其他 G")
	fmt.Println("      → 阻塞的 G 在条件满足后重新调度")

	fmt.Println("\n⚡ 关键特性：")
	fmt.Println("   ✅ M:N 模型：M 个 Goroutine 映射到 N 个 OS 线程")
	fmt.Println("   ✅ 工作窃取（Work Stealing）：空闲 P 从其他 P 偷取 G")
	fmt.Println("   ✅ 抢占式调度：长时间运行的 G 会被抢占")
	fmt.Println("   ✅ 系统调用优化：阻塞时 M 和 G 解绑，不阻塞其他 G")
}

// 演示 GMP 实际运行
func demonstrateGMP() {
	fmt.Println("\n1. 查看当前 GMP 状态：")
	fmt.Printf("   CPU 核心数: %d\n", runtime.NumCPU())
	fmt.Printf("   P 数量 (GOMAXPROCS): %d\n", runtime.GOMAXPROCS(0))
	fmt.Printf("   当前 Goroutine 数: %d\n", runtime.NumGoroutine())

	fmt.Println("\n2. 创建多个 Goroutine 观察 GMP：")

	// 创建多个 goroutine
	for i := 0; i < 10; i++ {
		go func(id int) {
			// 模拟一些工作
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("   [G-%d] 执行中，当前 Goroutine 数: %d\n",
				id, runtime.NumGoroutine())
		}(i)
	}

	fmt.Printf("   创建后 Goroutine 数: %d\n", runtime.NumGoroutine())

	fmt.Println("\n3. 等待所有 Goroutine 完成...")
	time.Sleep(200 * time.Millisecond)

	fmt.Printf("   完成后 Goroutine 数: %d\n", runtime.NumGoroutine())

	fmt.Println("\n4. 演示系统调用（阻塞场景）：")
	fmt.Println("   当 G 执行系统调用时，M 和 G 会解绑")
	fmt.Println("   M 可以继续执行其他 G，提高并发效率")
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
