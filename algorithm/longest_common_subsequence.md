# 最长公共子序列 (Longest Common Subsequence, LCS)

## 问题描述

给定两个字符串 `text1` 和 `text2`，返回这两个字符串的最长公共子序列的长度。

**子序列**：在不改变字符相对顺序的情况下，删除某些字符（也可以不删除）后得到的新序列。

### 示例
```
输入：text1 = "abcde", text2 = "ace"
输出：3
解释：最长公共子序列是 "ace"，长度为 3。

输入：text1 = "abc", text2 = "abc"
输出：3
解释：最长公共子序列是 "abc"，长度为 3。

输入：text1 = "abc", text2 = "def"
输出：0
解释：两个字符串没有公共子序列，返回 0。
```

## 算法实现

### 1. 动态规划解法

#### 核心思想
- 使用二维DP表，`dp[i][j]` 表示 `text1[0:i]` 和 `text2[0:j]` 的最长公共子序列长度
- 状态转移方程：
  - 如果 `text1[i-1] == text2[j-1]`：`dp[i][j] = dp[i-1][j-1] + 1`
  - 否则：`dp[i][j] = max(dp[i-1][j], dp[i][j-1])`

#### 实现代码
```go
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

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}
```

#### 空间优化版本
```go
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
```

### 2. 递归解法（带记忆化）

```go
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
```

### 3. 获取具体的LCS序列

```go
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
```

## 算法复杂度分析

### 时间复杂度
- **动态规划**：O(m × n)
- **递归记忆化**：O(m × n)
- **空间优化版本**：O(m × n)

### 空间复杂度
- **标准动态规划**：O(m × n)
- **空间优化版本**：O(n)
- **递归记忆化**：O(m × n)

## 应用场景

### 1. 生物信息学
- **DNA序列比对**：比较两个DNA序列的相似性
- **蛋白质序列分析**：分析蛋白质序列的保守区域

### 2. 文本处理
- **文件差异比较**：如Git的diff算法
- **拼写检查**：计算编辑距离
- **文本相似度**：评估两个文本的相似程度

### 3. 版本控制
- **代码差异检测**：比较不同版本的代码
- **文档版本管理**：跟踪文档的变更

### 4. 数据挖掘
- **模式识别**：发现序列中的共同模式
- **时间序列分析**：比较时间序列的相似性

## 变种问题

### 1. 最长公共子串 (Longest Common Substring)
```go
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
```

### 2. 多个序列的LCS
```go
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
        // 这里需要修改为处理中间结果
        // 简化处理：直接计算与当前序列的LCS
        result = min(result, LongestCommonSubsequence(sequences[0], sequences[i]))
    }
    
    return result
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
```

### 3. 带权重的LCS
```go
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
```

## 性能优化技巧

### 1. 空间优化
- 使用滚动数组，只保存必要的状态
- 对于大字符串，考虑分块处理

### 2. 剪枝优化
- 如果两个字符串长度差异很大，可以先处理较短的字符串
- 使用哈希表预处理字符位置

### 3. 并行化
- 对于大规模数据，可以使用并行算法
- 分块并行计算DP表

## 实际应用示例

### 1. 文件差异比较
```go
func CompareFiles(file1, file2 string) []string {
    // 读取文件内容
    content1 := readFile(file1)
    content2 := readFile(file2)
    
    // 按行分割
    lines1 := strings.Split(content1, "\n")
    lines2 := strings.Split(content2, "\n")
    
    // 计算LCS
    lcs := GetLongestCommonSubsequence(strings.Join(lines1, ""), strings.Join(lines2, ""))
    
    // 返回差异
    return findDifferences(lines1, lines2, lcs)
}
```

### 2. DNA序列比对
```go
func CompareDNASequences(seq1, seq2 string) float64 {
    lcsLength := LongestCommonSubsequence(seq1, seq2)
    maxLength := max(len(seq1), len(seq2))
    
    // 计算相似度
    similarity := float64(lcsLength) / float64(maxLength)
    return similarity
}
```

## 总结

最长公共子序列问题是动态规划的经典应用，具有重要的理论和实践价值：

- **算法特点**：时间复杂度O(mn)，空间复杂度可优化到O(n)
- **应用广泛**：生物信息学、文本处理、版本控制等领域
- **变种丰富**：子串、多序列、带权重等变种问题
- **优化空间**：空间优化、并行化、剪枝等优化技巧

选择合适的算法需要根据具体问题的规模、精度要求和实时性需求来决定。 