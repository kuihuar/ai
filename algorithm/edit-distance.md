# 编辑距离 (Edit Distance)

## 📖 概述

编辑距离是衡量两个字符串相似度的重要指标，表示将一个字符串转换为另一个字符串所需的最少操作次数。常见的操作包括插入、删除和替换字符。

## 🎯 应用场景

- **拼写检查** - 自动纠错和拼写建议
- **DNA序列比对** - 生物信息学中的序列分析
- **自然语言处理** - 文本相似度计算
- **版本控制** - 文件差异比较
- **语音识别** - 语音到文本的纠错

## 🔍 编辑距离类型

### 1. Levenshtein 距离
最常用的编辑距离，允许三种操作：插入、删除、替换。

### 2. Damerau-Levenshtein 距离
在 Levenshtein 基础上增加交换相邻字符的操作。

### 3. Hamming 距离
只允许替换操作，要求两个字符串长度相等。

### 4. Longest Common Subsequence (LCS)
计算最长公共子序列的长度。

## 🛠️ Go 语言实现

### 1. 动态规划解法

#### 基本思路
使用二维DP数组，`dp[i][j]` 表示将 `word1[0...i-1]` 转换为 `word2[0...j-1]` 所需的最少操作次数。

#### 状态转移方程
```go
if word1[i-1] == word2[j-1] {
    dp[i][j] = dp[i-1][j-1]  // 无需操作
} else {
    dp[i][j] = min(
        dp[i-1][j] + 1,      // 删除 word1[i-1]
        dp[i][j-1] + 1,      // 插入 word2[j-1]
        dp[i-1][j-1] + 1     // 替换 word1[i-1] 为 word2[j-1]
    )
}
```

#### Go 实现
```go
package main

import (
	"fmt"
)

// LevenshteinDistance 计算两个字符串的 Levenshtein 距离
func LevenshteinDistance(word1, word2 string) int {
	m, n := len(word1), len(word2)
	
	// 创建 DP 数组
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	
	// 初始化第一行和第一列
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	
	// 填充 DP 数组
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min(
					dp[i-1][j] + 1,      // 删除
					dp[i][j-1] + 1,      // 插入
					dp[i-1][j-1] + 1,    // 替换
				)
			}
		}
	}
	
	return dp[m][n]
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// min3 返回三个整数中的最小值
func min3(a, b, c int) int {
	return min(min(a, b), c)
}

func main() {
	word1 := "horse"
	word2 := "ros"
	distance := LevenshteinDistance(word1, word2)
	fmt.Printf("'%s' 到 '%s' 的编辑距离: %d\n", word1, word2, distance)
	// 输出: 'horse' 到 'ros' 的编辑距离: 3
}
```

### 2. 空间优化版本

#### 思路
由于每次计算只依赖上一行的数据，可以使用一维数组优化空间复杂度。

#### Go 实现
```go
// LevenshteinDistanceOptimized 空间优化的 Levenshtein 距离计算
func LevenshteinDistanceOptimized(word1, word2 string) int {
	m, n := len(word1), len(word2)
	
	// 使用一维数组
	dp := make([]int, n+1)
	
	// 初始化第一行
	for j := 0; j <= n; j++ {
		dp[j] = j
	}
	
	// 逐行计算
	for i := 1; i <= m; i++ {
		prev := dp[0] // 保存左上角的值
		dp[0] = i
		
		for j := 1; j <= n; j++ {
			temp := dp[j] // 保存当前值作为下一次的左上角
			if word1[i-1] == word2[j-1] {
				dp[j] = prev
			} else {
				dp[j] = min3(dp[j] + 1, dp[j-1] + 1, prev + 1)
			}
			prev = temp
		}
	}
	
	return dp[n]
}
```

### 3. Damerau-Levenshtein 距离

#### Go 实现
```go
// DamerauLevenshteinDistance 计算 Damerau-Levenshtein 距离（包含交换操作）
func DamerauLevenshteinDistance(word1, word2 string) int {
	m, n := len(word1), len(word2)
	
	// 创建 DP 数组
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	
	// 初始化
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	
	// 填充 DP 数组
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min3(
					dp[i-1][j] + 1,      // 删除
					dp[i][j-1] + 1,      // 插入
					dp[i-1][j-1] + 1,    // 替换
				)
				
				// 检查交换操作
				if i > 1 && j > 1 && 
				   word1[i-1] == word2[j-2] && 
				   word1[i-2] == word2[j-1] {
					dp[i][j] = min(dp[i][j], dp[i-2][j-2] + 1)
				}
			}
		}
	}
	
	return dp[m][n]
}
```

## 📊 算法分析

### 时间复杂度
- **时间复杂度**: O(m × n)
- **空间复杂度**: O(m × n) (优化后为 O(min(m, n)))

### 算法特点
- **优点**: 准确计算最小编辑距离
- **缺点**: 对于长字符串计算较慢

## 🎯 实际应用

### 1. 拼写检查器
```go
// SpellChecker 简单的拼写检查器
func SpellChecker(word string, dictionary []string, threshold int) []Suggestion {
	var suggestions []Suggestion
	
	for _, dictWord := range dictionary {
		distance := LevenshteinDistance(word, dictWord)
		if distance <= threshold {
			suggestions = append(suggestions, Suggestion{
				Word:     dictWord,
				Distance: distance,
			})
		}
	}
	
	// 按距离排序
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Distance < suggestions[j].Distance
	})
	
	return suggestions
}

type Suggestion struct {
	Word     string
	Distance int
}

func main() {
	dictionary := []string{"hello", "world", "python", "algorithm", "computer"}
	word := "helo"
	suggestions := SpellChecker(word, dictionary, 2)
	
	fmt.Printf("'%s' 的建议:\n", word)
	for _, s := range suggestions {
		fmt.Printf("  %s (距离: %d)\n", s.Word, s.Distance)
	}
}
```

### 2. 文本相似度计算
```go
// TextSimilarity 计算两个文本的相似度
func TextSimilarity(text1, text2 string) float64 {
	distance := LevenshteinDistance(text1, text2)
	maxLen := max(len(text1), len(text2))
	
	if maxLen == 0 {
		return 1.0
	}
	
	similarity := 1.0 - float64(distance)/float64(maxLen)
	return similarity
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	text1 := "hello world"
	text2 := "hello python"
	similarity := TextSimilarity(text1, text2)
	fmt.Printf("相似度: %.2f\n", similarity)
}
```

### 3. DNA序列比对
```go
// DNASequenceCompare DNA序列比对
func DNASequenceCompare(seq1, seq2 string) int {
	distance := LevenshteinDistance(seq1, seq2)
	maxLen := max(len(seq1), len(seq2))
	
	fmt.Printf("序列1: %s\n", seq1)
	fmt.Printf("序列2: %s\n", seq2)
	fmt.Printf("编辑距离: %d\n", distance)
	fmt.Printf("相似度: %.2f\n", 1.0-float64(distance)/float64(maxLen))
	
	return distance
}

func main() {
	seq1 := "ATCGATCG"
	seq2 := "ATCGATCC"
	DNASequenceCompare(seq1, seq2)
}
```

## 🚀 优化技巧

### 1. 早期终止
```go
// LevenshteinDistanceWithEarlyStop 带早期终止的编辑距离计算
func LevenshteinDistanceWithEarlyStop(word1, word2 string, maxDistance int) int {
	m, n := len(word1), len(word2)
	
	// 如果长度差超过最大距离，直接返回
	if abs(m-n) > maxDistance {
		return maxDistance + 1
	}
	
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	
	for i := 0; i <= m; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= n; j++ {
		dp[0][j] = j
	}
	
	for i := 1; i <= m; i++ {
		minInRow := int(^uint(0) >> 1) // 最大整数值
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = min3(
					dp[i-1][j] + 1,
					dp[i][j-1] + 1,
					dp[i-1][j-1] + 1,
				)
			}
			if dp[i][j] < minInRow {
				minInRow = dp[i][j]
			}
		}
		
		// 如果当前行的最小值超过阈值，提前终止
		if minInRow > maxDistance {
			return maxDistance + 1
		}
	}
	
	return dp[m][n]
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
```

### 2. 并行计算
```go
import (
	"sync"
)

// ParallelEditDistance 并行计算多个单词的编辑距离
func ParallelEditDistance(words []string, target string) []Suggestion {
	var wg sync.WaitGroup
	suggestions := make([]Suggestion, len(words))
	
	for i, word := range words {
		wg.Add(1)
		go func(index int, w string) {
			defer wg.Done()
			distance := LevenshteinDistance(w, target)
			suggestions[index] = Suggestion{
				Word:     w,
				Distance: distance,
			}
		}(i, word)
	}
	
	wg.Wait()
	return suggestions
}

func main() {
	words := []string{"hello", "world", "python", "algorithm"}
	target := "helo"
	results := ParallelEditDistance(words, target)
	
	fmt.Printf("与 '%s' 的距离:\n", target)
	for _, r := range results {
		fmt.Printf("  %s: %d\n", r.Word, r.Distance)
	}
}
```

## 🧪 测试用例

### 1. 基本测试
```go
import "testing"

func TestEditDistance(t *testing.T) {
	testCases := []struct {
		word1    string
		word2    string
		expected int
	}{
		{"horse", "ros", 3},
		{"intention", "execution", 5},
		{"", "hello", 5},
		{"hello", "", 5},
		{"", "", 0},
		{"same", "same", 0},
	}
	
	for _, tc := range testCases {
		result := LevenshteinDistance(tc.word1, tc.word2)
		if result != tc.expected {
			t.Errorf("错误: %s -> %s, 期望 %d, 得到 %d", 
				tc.word1, tc.word2, tc.expected, result)
		}
		t.Logf("✓ %s -> %s: %d", tc.word1, tc.word2, result)
	}
}
```

### 2. 性能测试
```go
import (
	"math/rand"
	"testing"
	"time"
)

func BenchmarkEditDistance(b *testing.B) {
	// 生成测试数据
	rand.Seed(time.Now().UnixNano())
	word1 := generateRandomString(100)
	word2 := generateRandomString(100)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LevenshteinDistance(word1, word2)
	}
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
```

## 📚 扩展阅读

### 1. 相关算法
- **Wagner-Fischer 算法** - 经典的动态规划解法
- **Myers 算法** - 线性空间复杂度的解法
- **Bit-parallel 算法** - 使用位运算优化

### 2. 应用领域
- **生物信息学** - 序列比对和进化分析
- **自然语言处理** - 文本相似度和纠错
- **信息检索** - 模糊搜索和推荐系统

### 3. 进阶主题
- **加权编辑距离** - 不同操作有不同的代价
- **近似算法** - 快速近似计算
- **并行算法** - 大规模数据的并行处理

## 🎯 练习题目

### 1. 基础练习
1. 实现 Hamming 距离算法
2. 实现最长公共子序列算法
3. 实现加权编辑距离

### 2. 进阶练习
1. 实现 Myers 算法
2. 实现 Bit-parallel 算法
3. 设计并行编辑距离算法

### 3. 应用练习
1. 构建完整的拼写检查器
2. 实现文本相似度搜索引擎
3. 设计 DNA 序列比对工具

---

**编辑距离是字符串处理和算法设计中的基础概念，掌握它对于理解更复杂的文本处理算法非常重要！** 🎉
 