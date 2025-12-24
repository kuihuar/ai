package suanfa

import "fmt"

// Adjacency List
func BFS(root *TreeNode, start, end int) {
	queue := []*TreeNode{root}
	visited := make(map[int]bool)
	visited[root.Val] = true
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node.Val == end {
			fmt.Printf("BFS found %d \n", end)
			return
		}
		if node.Left != nil && !visited[node.Left.Val] {
			queue = append(queue, node.Left)
			visited[node.Left.Val] = true
		}
		if node.Right != nil && !visited[node.Right.Val] {
			queue = append(queue, node.Right)
			visited[node.Right.Val] = true
		}
	}
}

func DFS(root *TreeNode, start, end int) {
	stack := []*TreeNode{root}
	visited := make(map[int]bool)
	visited[root.Val] = true
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if node.Val == end {
			fmt.Printf("DFS found %d \n", end)
		}
		if node.Right != nil && !visited[node.Right.Val] {
			stack = append(stack, node.Right)
			visited[node.Right.Val] = true
		}
		if node.Left != nil && !visited[node.Left.Val] {
			stack = append(stack, node.Left)
			visited[node.Left.Val] = true
		}
	}
}

func DFS1(root *TreeNode, start, end int) {
	visited := make(map[int]bool)

	var dfs func(node *TreeNode) bool
	dfs = func(node *TreeNode) bool {
		if node == nil {
			return false
		}
		if node.Val == end {
			fmt.Printf("DFS found %d \n", end)
			return true
		}
		if visited[node.Val] {
			return false
		}
		visited[node.Val] = true
		return dfs(node.Left) || dfs(node.Right)
	}
}

// 102 二叉树的层序遍历，分层遍历，时间复杂度为O(n)，空间复杂度为O(n)
// BFS 广度优先搜索
func LevelOrderBFS(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}
	queue := []*TreeNode{root}
	result := [][]int{}

	for len(queue) > 0 {
		level := []int{}
		size := len(queue)
		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]
			level = append(level, node.Val)
			if node.Left != nil {
				queue = append(queue, node.Left)
			}
			if node.Right != nil {
				queue = append(queue, node.Right)
			}
		}
		result = append(result, level)

	}
	return result
}

// 深度优先搜索, 时间复杂度为O(n)，空间复杂度为O(n)
func LevelOrderDFS(root *TreeNode) [][]int {

	var result [][]int
	var dfs func(node *TreeNode, level int)
	dfs = func(node *TreeNode, level int) {
		// terminator
		if node == nil {
			return
		}
		// process current logic
		if len(result) == level {
			result = append(result, []int{})
		}
		result[level] = append(result[level], node.Val)
		// drill down
		dfs(node.Left, level+1)
		dfs(node.Right, level+1)

	}
	dfs(root, 0)
	return result
}

// 104 二叉树的最大深度 111 二叉树的最小深度

func MaxDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	leftDepth := MaxDepth(root.Left)
	rightDepth := MaxDepth(root.Right)
	return max(leftDepth, rightDepth) + 1
}
func MinDepth(root *TreeNode) int {
	if root == nil {
		return 0
	}
	// if root.Left == nil && root.Right == nil {
	// 	return 1
	// }
	// divide and conquer
	leftDepth := MinDepth(root.Left)
	rightDepth := MinDepth(root.Right)
	if leftDepth == 0 || rightDepth == 0 {
		// or return leftDepth + rightDepth + 1;
		return max(leftDepth, rightDepth) + 1

	}
	return min(leftDepth, rightDepth) + 1
}

func MinDepth1(root *TreeNode) int {
	if root == nil {
		return 0
	}
	if root.Left == nil {
		return 1 + MinDepth1(root.Right)
	}
	if root.Right == nil {
		return 1 + MinDepth1(root.Left)
	}
	leftDepth := MinDepth1(root.Left)
	rightDepth := MinDepth1(root.Right)
	// if leftDepth == 0 || rightDepth == 0 {
	// 	return max(leftDepth, rightDepth) + 1
	// }
	// process  subproblems' results
	result := min(leftDepth, rightDepth) + 1
	return result
	// return min(leftDepth, rightDepth) + 1
}
