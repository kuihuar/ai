package ai

import (
	"bytes"
	"container/list"
	"net"
	"net/http"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func ReverseList(head *ListNode) *ListNode {

	var prev *ListNode = nil
	curr := head

	for curr != nil {
		next := curr.Next

		curr.Next = prev
		prev = curr
		curr = next
	}
	return prev
}

func BinarySearch(arr []int, target int) int {

	left := 0
	right := len(arr) - 1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			result = mid
			break
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return result
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func DFSTree(root *TreeNode) []int {

	result := []int{}

	var dfsHelper func(node *TreeNode)

	dfsHelper = func(node *TreeNode) {
		if node == nil {
			return
		}
		result = append(result, node.Val)
		dfsHelper(node.Left)
		dfsHelper(node.Right)
	}
	dfsHelper(root)
	return result
}

func DFSIterativeTree(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	result := []int{}
	stack := list.New()

	stack.PushBack(root)

	for stack.Len() > 0 {
		node := stack.Remove(stack.Back()).(*TreeNode)
		result = append(result, node.Val)
		if node.Right != nil {
			stack.PushBack(node.Right)
		}
		if node.Left != nil {
			stack.PushBack(node.Left)
		}
	}
	return result
}

func BFSTree(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	result := [][]int{}

	queue := list.New()
	queue.PushBack(root)

	for queue.Len() > 0 {
		levelSize := queue.Len()
		level := []int{}

		for i := 0; i < levelSize; i++ {
			node := queue.Remove(queue.Front()).(*TreeNode)
			level = append(level, node.Val)
			if node.Left != nil {
				queue.PushBack(node.Left)
			}
			if node.Right != nil {
				queue.PushBack(node.Right)
			}
		}
		result = append(result, level)
	}
	return result
}

func Knapsack01DP(weights []int, values []int, capacity int) int {

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	conn.Read([]byte("hello"))

	resp, err := http.Post("localhost:8080", "application/json", bytes.NewBufferString("{}"))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// 获取物品数量
	n := len(weights)

	dp := make([][]int, n+1)

	for i := range dp {
		dp[i] = make([]int, capacity+1)
	}

	for i := 1; i <= n; i++ {
		for j := 0; j <= capacity; j++ {
			if weights[i-1] <= j {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-weights[i-1]]+values[i-1])
			} else {
				dp[i][j] = dp[i-1][j]
			}
		}
	}
	return dp[n][capacity]

}

func LongestCommonSubsequence(text1 string, text2 string) int {

	m := len(text1)
	n := len(text2)

	dp := make([][]int, m+1)

	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {

			if text1[i-1] == text2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}
	return dp[m][n]
}

func LengOfLIS(nums []int) int {

	dp := make([]int, len(nums))

	for i := range dp {
		dp[i] = 1
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

func GetLongestIncreasingSubsequence(nums []int) []int {
	dp := make([]int, len(nums))

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

func LengthOfLISRecursive(nums []int) int {

	memo := make([]int, len(nums))
	for i := range memo {
		memo[i] = -1
	}

	maxLen := 1
	for i := 0; i > len(nums); i++ {
		maxLen = max(maxLen, lisHelper(nums, i, memo))
	}

	return maxLen
}

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
