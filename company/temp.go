package company

import (
	"container/list"
)

// bfs 图的遍历
// graph:= 邻接矩阵

// graph := [][]int{
// // 下标   邻居列表
//   0:  {1, 2},
//   1:  {0, 3, 4},
//   2:  {0, 5},
//   3:  {1},
//   4:  {1, 5},
//   5:  {2, 4},
// }

//       0
//     /   \
//    1     2
//   / \     \
//  3   4 --- 5

func bfs(graph [][]int, start int) []int {

	visited := make(map[int]bool)

	result := []int{}
	queue := list.New()
	queue.PushBack(start)
	visited[start] = true

	for queue.Len() > 0 {
		node := queue.Front()
		queue.Remove(node)
		result = append(result, node.Value.(int))
		for _, neighbor := range graph[node.Value.(int)] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue.PushBack(neighbor)
			}
		}
	}
	return result
}

func bfs1(graph [][]int, start int) []int {

	queue := []int{start}

	visited := make(map[int]bool)

	visited[start] = true

	result := []int{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		result = append(result, node)
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {

				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
	return result
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

func LevelOrder(root *TreeNode) [][]int {

	var res [][]int
	if root == nil {
		return res
	}

	queue := list.New()

	queue.PushBack(root)

	for queue.Len() > 0 {
		level := []int{}
		levelSize := queue.Len()

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
		res = append(res, level)
	}
	return res
}

type ListNode struct {
	Next *ListNode
}

func reverseList(head *ListNode) *ListNode {
	var prev *ListNode

	for head != nil {
		next := head.Next

		head.Next = prev
		prev = head
		head = next
	}

	return prev
}

func dfs(grapth [][]int, start int) []int {

	visited := make(map[int]bool)

	result := []int{}

	var dfsHelper func(node int)
	dfsHelper = func(node int) {
		if visited[node] {
			return
		}

		visited[node] = true
		result = append(result, node)
		for _, neighbor := range grapth[node] {
			dfsHelper(neighbor)
		}
	}
	dfsHelper(start)
	return result
}

func dfsInterative(grapth [][]int, start int) []int {

	visited := make(map[int]bool)

	result := []int{}

	stack := list.New()
	stack.PushBack(start)

	for stack.Len() > 0 {
		node := stack.Remove(stack.Back()).(int)
		if visited[node] {
			continue
		}
		visited[node] = true
		result = append(result, node)

		for _, neighbor := range grapth[node] {

			stack.PushBack(neighbor)
		}

		for i := len(grapth[node]) - 1; i >= 0; i-- {
			neighbor := grapth[node][i]
			stack.PushBack(neighbor)
		}
	}

	return result
}

func QuickSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}

	pivot := arr[0]

	var left, right []int
	for _, num := range arr[1:] {
		if num < pivot {
			left = append(left, num)
		} else {
			right = append(right, num)
		}
	}
	left = QuickSort(left)
	right = QuickSort(right)

	result := append(left, pivot)
	result = append(result, right...)
	return result
}

func QuickSortInPlaxce(arr []int, low, hight int) {
	var partition func(int, int) int

	partition = func(low, hight int) int {

		pivot := arr[low]
		left := low
		right := hight
		for left < right {
			for left < right && arr[right] > pivot {
				right--
			}
			for left < right && arr[left] <= pivot {
				left++
			}
			arr[left], arr[right] = arr[right], arr[left]
		}
		arr[low], arr[left] = arr[left], arr[low]
		return left
	}

	if low < hight {
		pivotIndex := partition(low, hight)
		QuickSortInPlaxce(arr, low, pivotIndex-1)
		QuickSortInPlaxce(arr, pivotIndex+1, hight)
	}
}

// QuickSortInPlaxce(arr, 0, len(arr)-1)
