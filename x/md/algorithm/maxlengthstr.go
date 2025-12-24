// 1. 无重复字符的最长子串
package algorithm

// abcabcbb
func lengthOfLongestSubstring(s string) int {
	n := len(s)
	if n == 0 {
		return 0
	}
	left := 0
	maxLength := 0
	charIndexMap := make(map[byte]int)

	for right := 0; right < n; right++ {
		if index, exists := charIndexMap[s[right]]; exists {
			left = max(left, index+1)
		}
		charIndexMap[s[right]] = right
		maxLength = max(maxLength, right-left+1)
	}
	return maxLength
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// func main() {
//     s := "abcabcbb"
//     fmt.Println(lengthOfLongestSubstring(s))
// }
// 在上述代码中，使用 left 和 right 指针构建滑动窗口，charIndexMap 用于记录字符的最新位置。right 指针右移，若遇到重复字符，更新 left 指针，同时更新最大长度。

// 2. 最长公共子串
// package main

// import (
//     "fmt"
// )

func longestCommonSubstring(s1, s2 string) int {
	m := len(s1)
	n := len(s2)
	dp := make([][]int, m+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	maxLength := 0

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
				maxLength = max(maxLength, dp[i][j])
			} else {
				dp[i][j] = 0
			}
		}
	}
	return maxLength
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func main() {
// 	s1 := "abcde"
// 	s2 := "abfde"
// 	fmt.Println(longestCommonSubstring(s1, s2))
// }

// 此代码使用动态规划，dp[i][j] 表示以 s1 的第 i 个字符和 s2 的第 j 个字符结尾的最长公共子串的长度。若字符相等则更新 dp 值并记录最大长度，不相等则置为 0。

// 3. 最长回文子串
// package main

// import (
//     "fmt"
// )

func longestPalindrome(s string) string {
	if len(s) < 2 {
		return s
	}
	start, maxLength := 0, 0

	expandAroundCenter := func(left, right int) {
		for left >= 0 && right < len(s) && s[left] == s[right] {
			if right-left+1 > maxLength {
				start = left
				maxLength = right - left + 1
			}
			left--
			right++
		}
	}

	for i := 0; i < len(s); i++ {
		expandAroundCenter(i, i)
		expandAroundCenter(i, i+1)
	}
	return s[start : start+maxLength]
}

// func main() {
// 	s := "babad"
// 	fmt.Println(longestPalindrome(s))
// }

// 这里运用中心扩展法，对每个字符分别以其为中心（奇数长度回文串）和以其与下一个字符为中心（偶数长度回文串）向两边扩展，记录最长回文串的起始位置和长度。
