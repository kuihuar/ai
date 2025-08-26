package algorithm

import (
	"container/list"
	"fmt"
)

// TreeNode 二叉树节点
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// GraphNode 图节点
type GraphNode struct {
	Val       int
	Neighbors []*GraphNode
}

// DFS 深度优先搜索 - 递归版本
// 时间复杂度: O(V + E) - V为顶点数，E为边数
// 空间复杂度: O(V) - 递归调用栈深度
func DFS(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	result := []int{}

	var dfsHelper func(node int)
	dfsHelper = func(node int) {
		if visited[node] {
			return
		}

		visited[node] = true
		result = append(result, node)

		for _, neighbor := range graph[node] {
			dfsHelper(neighbor)
		}
	}

	dfsHelper(start)
	return result
}

// DFSIterative 深度优先搜索 - 迭代版本（使用栈）
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func DFSIterative(graph map[int][]int, start int) []int {
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

		// 将邻居节点按逆序压入栈中（保持正确的访问顺序）
		for i := len(graph[node]) - 1; i >= 0; i-- {
			neighbor := graph[node][i]
			if !visited[neighbor] {
				stack.PushBack(neighbor)
			}
		}
	}

	return result
}

// DFSTree 二叉树深度优先搜索
// 时间复杂度: O(n) - n为节点数
// 空间复杂度: O(h) - h为树的高度
func DFSTree(root *TreeNode) []int {
	result := []int{}

	var dfsHelper func(node *TreeNode)
	dfsHelper = func(node *TreeNode) {
		if node == nil {
			return
		}

		// 前序遍历：根 -> 左 -> 右
		result = append(result, node.Val)
		dfsHelper(node.Left)
		dfsHelper(node.Right)
	}

	dfsHelper(root)
	return result
}

// DFSInorder 二叉树中序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(h)
func DFSInorder(root *TreeNode) []int {
	result := []int{}

	var inorderHelper func(node *TreeNode)
	inorderHelper = func(node *TreeNode) {
		if node == nil {
			return
		}

		// 中序遍历：左 -> 根 -> 右
		inorderHelper(node.Left)
		result = append(result, node.Val)
		inorderHelper(node.Right)
	}

	inorderHelper(root)
	return result
}

// DFSPostorder 二叉树后序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(h)
func DFSPostorder(root *TreeNode) []int {
	result := []int{}

	var postorderHelper func(node *TreeNode)
	postorderHelper = func(node *TreeNode) {
		if node == nil {
			return
		}

		// 后序遍历：左 -> 右 -> 根
		postorderHelper(node.Left)
		postorderHelper(node.Right)
		result = append(result, node.Val)
	}

	postorderHelper(root)
	return result
}

// DFSIterativeTree 二叉树迭代深度优先搜索
// 时间复杂度: O(n)
// 空间复杂度: O(n)
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

		// 先压入右子节点，再压入左子节点（栈的特性）
		if node.Right != nil {
			stack.PushBack(node.Right)
		}
		if node.Left != nil {
			stack.PushBack(node.Left)
		}
	}

	return result
}

// DFSInorderIterative 二叉树迭代中序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func DFSInorderIterative(root *TreeNode) []int {
	result := []int{}
	stack := list.New()
	current := root

	for current != nil || stack.Len() > 0 {
		// 一直向左走到底
		for current != nil {
			stack.PushBack(current)
			current = current.Left
		}

		// 弹出栈顶元素
		current = stack.Remove(stack.Back()).(*TreeNode)
		result = append(result, current.Val)

		// 转向右子树
		current = current.Right
	}

	return result
}

// DFSPostorderIterative 二叉树迭代后序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(n)
func DFSPostorderIterative(root *TreeNode) []int {
	if root == nil {
		return []int{}
	}

	result := []int{}
	stack := list.New()
	stack.PushBack(root)

	for stack.Len() > 0 {
		node := stack.Remove(stack.Back()).(*TreeNode)
		result = append([]int{node.Val}, result...) // 在开头插入

		// 先压入左子节点，再压入右子节点
		if node.Left != nil {
			stack.PushBack(node.Left)
		}
		if node.Right != nil {
			stack.PushBack(node.Right)
		}
	}

	return result
}

// DFSWithPath 带路径的深度优先搜索
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func DFSWithPath(graph map[int][]int, start, target int) []int {
	visited := make(map[int]bool)
	path := []int{}

	var dfsHelper func(node int) bool
	dfsHelper = func(node int) bool {
		if visited[node] {
			return false
		}

		visited[node] = true
		path = append(path, node)

		if node == target {
			return true
		}

		for _, neighbor := range graph[node] {
			if dfsHelper(neighbor) {
				return true
			}
		}

		// 回溯
		path = path[:len(path)-1]
		return false
	}

	if dfsHelper(start) {
		return path
	}

	return []int{}
}

// DFSAllPaths 查找所有路径
// 时间复杂度: O(V^V) - 最坏情况
// 空间复杂度: O(V)
func DFSAllPaths(graph map[int][]int, start, target int) [][]int {
	var result [][]int
	visited := make(map[int]bool)

	var dfsHelper func(node int, path []int)
	dfsHelper = func(node int, path []int) {
		if visited[node] {
			return
		}

		visited[node] = true
		path = append(path, node)

		if node == target {
			// 创建路径副本
			pathCopy := make([]int, len(path))
			copy(pathCopy, path)
			result = append(result, pathCopy)
		} else {
			for _, neighbor := range graph[node] {
				dfsHelper(neighbor, path)
			}
		}

		// 回溯
		visited[node] = false
	}

	dfsHelper(start, []int{})
	return result
}

// DFSConnectedComponents 查找连通分量
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func DFSConnectedComponents(graph map[int][]int) [][]int {
	visited := make(map[int]bool)
	var components [][]int

	var dfsHelper func(node int, component *[]int)
	dfsHelper = func(node int, component *[]int) {
		if visited[node] {
			return
		}

		visited[node] = true
		*component = append(*component, node)

		for _, neighbor := range graph[node] {
			dfsHelper(neighbor, component)
		}
	}

	for node := range graph {
		if !visited[node] {
			var component []int
			dfsHelper(node, &component)
			components = append(components, component)
		}
	}

	return components
}

// DFSCycleDetection 检测图中是否有环
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func DFSCycleDetection(graph map[int][]int) bool {
	visited := make(map[int]bool)
	recStack := make(map[int]bool) // 递归栈，用于检测有向图的环

	var hasCycle func(node int) bool
	hasCycle = func(node int) bool {
		if recStack[node] {
			return true // 发现环
		}

		if visited[node] {
			return false
		}

		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if hasCycle(neighbor) {
				return true
			}
		}

		recStack[node] = false
		return false
	}

	for node := range graph {
		if !visited[node] {
			if hasCycle(node) {
				return true
			}
		}
	}

	return false
}

// DFSTopologicalSort 拓扑排序
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func DFSTopologicalSort(graph map[int][]int) []int {
	visited := make(map[int]bool)
	recStack := make(map[int]bool)
	result := []int{}

	var dfsHelper func(node int) bool
	dfsHelper = func(node int) bool {
		if recStack[node] {
			return false // 有环
		}

		if visited[node] {
			return true
		}

		visited[node] = true
		recStack[node] = true

		for _, neighbor := range graph[node] {
			if !dfsHelper(neighbor) {
				return false
			}
		}

		recStack[node] = false
		result = append([]int{node}, result...) // 在开头插入
		return true
	}

	for node := range graph {
		if !visited[node] {
			if !dfsHelper(node) {
				return []int{} // 有环，无法拓扑排序
			}
		}
	}

	return result
}

// DFSBacktracking 回溯算法框架
// 时间复杂度: 取决于具体问题
// 空间复杂度: O(n) - n为递归深度
func DFSBacktracking(n int) [][]int {
	var result [][]int

	var backtrack func(path []int, used []bool)
	backtrack = func(path []int, used []bool) {
		if len(path) == n {
			// 创建路径副本
			pathCopy := make([]int, len(path))
			copy(pathCopy, path)
			result = append(result, pathCopy)
			return
		}

		for i := 0; i < n; i++ {
			if !used[i] {
				used[i] = true
				path = append(path, i)

				backtrack(path, used)

				// 回溯
				path = path[:len(path)-1]
				used[i] = false
			}
		}
	}

	backtrack([]int{}, make([]bool, n))
	return result
}

// DFSMatrix 二维矩阵深度优先搜索
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
func DFSMatrix(matrix [][]int, startRow, startCol int) [][]int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return [][]int{}
	}

	rows, cols := len(matrix), len(matrix[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	result := [][]int{}
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} // 上下左右

	var dfsHelper func(row, col int)
	dfsHelper = func(row, col int) {
		if row < 0 || row >= rows || col < 0 || col >= cols || visited[row][col] {
			return
		}

		visited[row][col] = true
		result = append(result, []int{row, col})

		for _, dir := range directions {
			newRow, newCol := row+dir[0], col+dir[1]
			dfsHelper(newRow, newCol)
		}
	}

	dfsHelper(startRow, startCol)
	return result
}

// DFSIslandCount 计算岛屿数量
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
func DFSIslandCount(grid [][]byte) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	rows, cols := len(grid), len(grid[0])
	count := 0

	var dfsHelper func(row, col int)
	dfsHelper = func(row, col int) {
		if row < 0 || row >= rows || col < 0 || col >= cols || grid[row][col] == '0' {
			return
		}

		grid[row][col] = '0' // 标记为已访问

		// 四个方向
		dfsHelper(row-1, col)
		dfsHelper(row+1, col)
		dfsHelper(row, col-1)
		dfsHelper(row, col+1)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == '1' {
				count++
				dfsHelper(i, j)
			}
		}
	}

	return count
}

// 测试函数
func TestDFS() {
	fmt.Println("=== 深度优先搜索算法测试 ===")

	// 图测试
	graph := map[int][]int{
		0: {1, 2},
		1: {3, 4},
		2: {5},
		3: {},
		4: {},
		5: {},
	}

	fmt.Println("图结构:", graph)

	// 基本DFS
	dfsResult := DFS(graph, 0)
	fmt.Printf("递归DFS结果: %v\n", dfsResult)

	// 迭代DFS
	dfsIterResult := DFSIterative(graph, 0)
	fmt.Printf("迭代DFS结果: %v\n", dfsIterResult)

	// 路径查找
	path := DFSWithPath(graph, 0, 5)
	fmt.Printf("从0到5的路径: %v\n", path)

	// 所有路径
	allPaths := DFSAllPaths(graph, 0, 4)
	fmt.Printf("从0到4的所有路径: %v\n", allPaths)

	// 连通分量
	components := DFSConnectedComponents(graph)
	fmt.Printf("连通分量: %v\n", components)

	// 环检测
	hasCycle := DFSCycleDetection(graph)
	fmt.Printf("图中是否有环: %t\n", hasCycle)

	// 二叉树测试
	root := &TreeNode{
		Val: 1,
		Left: &TreeNode{
			Val:   2,
			Left:  &TreeNode{Val: 4},
			Right: &TreeNode{Val: 5},
		},
		Right: &TreeNode{
			Val:   3,
			Left:  &TreeNode{Val: 6},
			Right: &TreeNode{Val: 7},
		},
	}

	fmt.Println("\n=== 二叉树DFS测试 ===")

	// 前序遍历
	preorder := DFSTree(root)
	fmt.Printf("前序遍历: %v\n", preorder)

	// 中序遍历
	inorder := DFSInorder(root)
	fmt.Printf("中序遍历: %v\n", inorder)

	// 后序遍历
	postorder := DFSPostorder(root)
	fmt.Printf("后序遍历: %v\n", postorder)

	// 迭代版本
	preorderIter := DFSIterativeTree(root)
	fmt.Printf("迭代前序遍历: %v\n", preorderIter)

	inorderIter := DFSInorderIterative(root)
	fmt.Printf("迭代中序遍历: %v\n", inorderIter)

	postorderIter := DFSPostorderIterative(root)
	fmt.Printf("迭代后序遍历: %v\n", postorderIter)

	// 回溯测试
	fmt.Println("\n=== 回溯算法测试 ===")
	backtrackResult := DFSBacktracking(3)
	fmt.Printf("3个元素的全排列: %v\n", backtrackResult)

	// 二维矩阵测试
	fmt.Println("\n=== 二维矩阵DFS测试 ===")
	matrix := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	matrixResult := DFSMatrix(matrix, 0, 0)
	fmt.Printf("从(0,0)开始的DFS: %v\n", matrixResult)

	// 岛屿数量测试
	grid := [][]byte{
		{'1', '1', '0', '0', '0'},
		{'1', '1', '0', '0', '0'},
		{'0', '0', '1', '0', '0'},
		{'0', '0', '0', '1', '1'},
	}
	islandCount := DFSIslandCount(grid)
	fmt.Printf("岛屿数量: %d\n", islandCount)
}
