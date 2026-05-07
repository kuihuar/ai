package main

import (
	"container/heap"
	"fmt"
	"sort"
)

func main() {
	data := []int{3, 11, 5, 8, 1, 10, 2, 7, 15, 4}

	fmt.Println(topK(data, 3))
	fmt.Println(quickSort(data))
}

type intHeap []int

func (h intHeap) Len() int {
	return len(h)
}
func (h intHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h intHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *intHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *intHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
func topK(nums []int, k int) []int {

	h := &intHeap{}
	heap.Init(h)

	for _, num := range nums {
		if h.Len() < k {
			heap.Push(h, num)
		} else if num > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, num)
		}
	}

	return *h
}

func quickSort(nums []int) []int {

	if len(nums) <= 1 {
		return nums
	}

	pivot := nums[0]

	var left, right []int

	for i := 1; i < len(nums); i++ {
		if nums[i] <= pivot {
			left = append(left, nums[i])
		} else {
			right = append(right, nums[i])
		}
	}

	left = quickSort(left)
	right = quickSort(right)

	result := append(left, pivot)
	result = append(result, right...)

	return result
}

func lengthOfLIS(nums []int) int {

	n := len(nums)

	if n == 0 {
		return 0
	}
	dp := make([]int, n)

	for i := range dp {
		dp[i] = 1
	}
	maxLen := 1

	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[i] > nums[j] {
				dp[i] = dp[j] + 1
			}
		}
		if dp[i] > maxLen {
			maxLen = dp[i]
		}
	}
	return maxLen
}

func lcs(s1, s2 string) (int, string) {
	m, n := len(s1), len(s2)
	dp := make([][]int, m+1)

	for i := range dp {
		dp[i] = make([]int, n+1)
	}
	maxLen := 0
	for i := 1; i < m; i++ {
		for j := 1; j < n; j++ {
			if s1[i-1] == s2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1

				if dp[i][j] > maxLen {
					maxLen = dp[i][j]
				}
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i][j-1])
			}
		}
	}

	res := []byte{}

	i, j := m, n

	for i > 0 && j > 0 {
		if s1[i-1] == s2[j-1] {
			res = append(res, s1[i-1])
			i--
			j--
		} else if dp[i-1][j] >= dp[i][j-1] {
			i--
		} else {
			j--
		}
	}

	for a, b := 0, len(res)-1; a < b; a, b = a+1, b-1 {
		res[a], res[b] = res[b], res[a]
	}
	return dp[m][n], string(res)
}

func threeSum(nums []int) [][]int {

	res := [][]int{}
	n := len(nums)
	if n < 3 {
		return res
	}
	sort.Ints(nums)

	for i := 0; i < n-2; i++ {
		if nums[i] > 0 {
			break
		}

		left, right := i+1, n-1

		for left < right {
			sum := nums[i] + nums[left] + nums[right]

			if sum == 0 {
				res = append(res, []int{nums[i], nums[left], nums[right]})
			} else if sum < 0 {
				left++
			} else {
				right--
			}
		}
	}
	return res
}

func trap(height []int) int {
	if len(height) == 0 {
		return 0
	}
	var res int

	for i := 1; i < len(height)-1; i++ {
		leftMax := 0

		for left := 0; left < i; left++ {
			leftMax = max(leftMax, height[left])
		}
		rightMax := 0

		for right := i + 1; right < len(height); right++ {
			rightMax = max(rightMax, height[right])
		}

		trap := min(leftMax, rightMax) - height[i]

		if trap > 0 {
			res += trap
		}
	}
	return res

}

type Listnode struct {
	Next *Listnode
}

func reverseList(head *Listnode) *Listnode {

	var prev *Listnode = nil

	curr := head

	for curr != nil {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}
	return prev
}
