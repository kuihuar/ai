package algorithm

import (
	"container/list"
	"fmt"
	"strconv"
)

// BFS 广度优先搜索 - 基本版本
// 时间复杂度: O(V + E) - V为顶点数，E为边数
// 空间复杂度: O(V) - 队列大小
func BFS(graph map[int][]int, start int) []int {
	visited := make(map[int]bool)
	result := []int{}
	queue := list.New()

	// 将起始节点加入队列
	queue.PushBack(start)
	visited[start] = true

	for queue.Len() > 0 {
		// 从队列头部取出节点
		node := queue.Remove(queue.Front()).(int)
		result = append(result, node)

		// 访问所有邻居节点
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue.PushBack(neighbor)
			}
		}
	}

	return result
}

// BFSWithLevel 带层级的广度优先搜索
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func BFSWithLevel(graph map[int][]int, start int) [][]int {
	visited := make(map[int]bool)
	result := [][]int{}
	queue := list.New()

	// 将起始节点加入队列，记录层级
	queue.PushBack([]int{start, 0}) // [节点, 层级]
	visited[start] = true

	currentLevel := 0
	currentLevelNodes := []int{}

	for queue.Len() > 0 {
		item := queue.Remove(queue.Front()).([]int)
		node := item[0]
		level := item[1]

		// 如果进入新层级，保存当前层级节点
		if level > currentLevel {
			result = append(result, currentLevelNodes)
			currentLevelNodes = []int{}
			currentLevel = level
		}

		currentLevelNodes = append(currentLevelNodes, node)

		// 访问所有邻居节点
		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue.PushBack([]int{neighbor, level + 1})
			}
		}
	}

	// 添加最后一层
	if len(currentLevelNodes) > 0 {
		result = append(result, currentLevelNodes)
	}

	return result
}

// BFSShortestPath 使用BFS找最短路径
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func BFSShortestPath(graph map[int][]int, start, target int) []int {
	if start == target {
		return []int{start}
	}

	visited := make(map[int]bool)
	parent := make(map[int]int) // 记录父节点，用于重建路径
	queue := list.New()

	queue.PushBack(start)
	visited[start] = true

	for queue.Len() > 0 {
		node := queue.Remove(queue.Front()).(int)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				parent[neighbor] = node
				queue.PushBack(neighbor)

				// 找到目标节点
				if neighbor == target {
					// 重建路径
					path := []int{target}
					current := target
					for current != start {
						current = parent[current]
						path = append([]int{current}, path...)
					}
					return path
				}
			}
		}
	}

	return []int{} // 没有找到路径
}

// BFSAllShortestPaths 找到所有最短路径
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func BFSAllShortestPaths(graph map[int][]int, start, target int) [][]int {
	if start == target {
		return [][]int{{start}}
	}

	visited := make(map[int]bool)
	parent := make(map[int][]int) // 每个节点可能有多个父节点
	distance := make(map[int]int) // 记录到每个节点的距离
	queue := list.New()

	queue.PushBack(start)
	visited[start] = true
	distance[start] = 0

	for queue.Len() > 0 {
		node := queue.Remove(queue.Front()).(int)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				distance[neighbor] = distance[node] + 1
				parent[neighbor] = []int{node}
				queue.PushBack(neighbor)
			} else if distance[neighbor] == distance[node]+1 {
				// 如果距离相同，说明是另一条最短路径
				parent[neighbor] = append(parent[neighbor], node)
			}
		}
	}

	// 重建所有最短路径
	var result [][]int
	var buildPaths func(current int, path []int)
	buildPaths = func(current int, path []int) {
		if current == start {
			// 创建路径副本
			pathCopy := make([]int, len(path))
			copy(pathCopy, path)
			result = append(result, pathCopy)
			return
		}

		for _, p := range parent[current] {
			buildPaths(p, append([]int{current}, path...))
		}
	}

	buildPaths(target, []int{})
	return result
}

// BFSTree 二叉树广度优先搜索（层序遍历）
// 时间复杂度: O(n) - n为节点数
// 空间复杂度: O(w) - w为树的最大宽度
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

		// 处理当前层的所有节点
		for i := 0; i < levelSize; i++ {
			node := queue.Remove(queue.Front()).(*TreeNode)
			level = append(level, node.Val)

			// 添加子节点到队列
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

// BFSZigzag 二叉树之字形层序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(w)
func BFSZigzag(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	result := [][]int{}
	queue := list.New()
	queue.PushBack(root)
	leftToRight := true

	for queue.Len() > 0 {
		levelSize := queue.Len()
		level := make([]int, levelSize)

		// 处理当前层的所有节点
		for i := 0; i < levelSize; i++ {
			node := queue.Remove(queue.Front()).(*TreeNode)

			// 根据方向决定插入位置
			index := i
			if !leftToRight {
				index = levelSize - 1 - i
			}
			level[index] = node.Val

			// 添加子节点到队列
			if node.Left != nil {
				queue.PushBack(node.Left)
			}
			if node.Right != nil {
				queue.PushBack(node.Right)
			}
		}

		result = append(result, level)
		leftToRight = !leftToRight
	}

	return result
}

// BFSBottomUp 二叉树自底向上层序遍历
// 时间复杂度: O(n)
// 空间复杂度: O(w)
func BFSBottomUp(root *TreeNode) [][]int {
	if root == nil {
		return [][]int{}
	}

	result := [][]int{}
	queue := list.New()
	queue.PushBack(root)

	for queue.Len() > 0 {
		levelSize := queue.Len()
		level := []int{}

		// 处理当前层的所有节点
		for i := 0; i < levelSize; i++ {
			node := queue.Remove(queue.Front()).(*TreeNode)
			level = append(level, node.Val)

			// 添加子节点到队列
			if node.Left != nil {
				queue.PushBack(node.Left)
			}
			if node.Right != nil {
				queue.PushBack(node.Right)
			}
		}

		// 在开头插入，实现自底向上
		result = append([][]int{level}, result...)
	}

	return result
}

// BFSConnectedComponents 使用BFS找连通分量
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func BFSConnectedComponents(graph map[int][]int) [][]int {
	visited := make(map[int]bool)
	var components [][]int

	var bfsComponent func(start int) []int
	bfsComponent = func(start int) []int {
		component := []int{}
		queue := list.New()

		queue.PushBack(start)
		visited[start] = true

		for queue.Len() > 0 {
			node := queue.Remove(queue.Front()).(int)
			component = append(component, node)

			for _, neighbor := range graph[node] {
				if !visited[neighbor] {
					visited[neighbor] = true
					queue.PushBack(neighbor)
				}
			}
		}

		return component
	}

	for node := range graph {
		if !visited[node] {
			component := bfsComponent(node)
			components = append(components, component)
		}
	}

	return components
}

// BFSBipartite 检测二分图
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func BFSBipartite(graph map[int][]int) bool {
	if len(graph) == 0 {
		return true
	}

	color := make(map[int]int) // 0: 未染色, 1: 红色, -1: 蓝色
	queue := list.New()

	// 从每个未访问的节点开始BFS
	for start := range graph {
		if color[start] != 0 {
			continue
		}

		queue.PushBack(start)
		color[start] = 1

		for queue.Len() > 0 {
			node := queue.Remove(queue.Front()).(int)

			for _, neighbor := range graph[node] {
				if color[neighbor] == 0 {
					// 给邻居染上相反的颜色
					color[neighbor] = -color[node]
					queue.PushBack(neighbor)
				} else if color[neighbor] == color[node] {
					// 相邻节点颜色相同，不是二分图
					return false
				}
			}
		}
	}

	return true
}

// BFSWordLadder 单词接龙（最短转换序列）
// 时间复杂度: O(26 * wordLength * wordListSize)
// 空间复杂度: O(wordListSize)
func BFSWordLadder(beginWord, endWord string, wordList []string) int {
	// 将单词列表转换为集合，便于查找
	wordSet := make(map[string]bool)
	for _, word := range wordList {
		wordSet[word] = true
	}

	if !wordSet[endWord] {
		return 0
	}

	queue := list.New()
	queue.PushBack([]string{beginWord, "1"}) // [单词, 步数]
	visited := make(map[string]bool)
	visited[beginWord] = true

	for queue.Len() > 0 {
		item := queue.Remove(queue.Front()).([]string)
		word := item[0]
		steps, _ := strconv.Atoi(item[1])

		if word == endWord {
			return steps
		}

		// 尝试改变单词的每个字符
		for i := 0; i < len(word); i++ {
			for c := 'a'; c <= 'z'; c++ {
				newWord := word[:i] + string(c) + word[i+1:]

				if wordSet[newWord] && !visited[newWord] {
					visited[newWord] = true
					queue.PushBack([]string{newWord, strconv.Itoa(steps + 1)})
				}
			}
		}
	}

	return 0
}

// BFSMatrix 二维矩阵的BFS
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
func BFSMatrix(grid [][]int, start []int) [][]int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return [][]int{}
	}

	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	result := [][]int{}
	queue := list.New()
	queue.PushBack(start)
	visited[start[0]][start[1]] = true

	// 四个方向：上、下、左、右
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for queue.Len() > 0 {
		pos := queue.Remove(queue.Front()).([]int)
		result = append(result, pos)

		// 检查四个方向的邻居
		for _, dir := range directions {
			newRow := pos[0] + dir[0]
			newCol := pos[1] + dir[1]

			// 检查边界和是否已访问
			if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols &&
				!visited[newRow][newCol] && grid[newRow][newCol] != 0 {
				visited[newRow][newCol] = true
				queue.PushBack([]int{newRow, newCol})
			}
		}
	}

	return result
}

// BFSIslands 计算岛屿数量
// 时间复杂度: O(m*n)
// 空间复杂度: O(m*n)
func BFSIslands(grid [][]int) int {
	if len(grid) == 0 || len(grid[0]) == 0 {
		return 0
	}

	rows, cols := len(grid), len(grid[0])
	visited := make([][]bool, rows)
	for i := range visited {
		visited[i] = make([]bool, cols)
	}

	count := 0
	directions := [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	var bfsIsland func(row, col int)
	bfsIsland = func(row, col int) {
		queue := list.New()
		queue.PushBack([]int{row, col})
		visited[row][col] = true

		for queue.Len() > 0 {
			pos := queue.Remove(queue.Front()).([]int)

			for _, dir := range directions {
				newRow := pos[0] + dir[0]
				newCol := pos[1] + dir[1]

				if newRow >= 0 && newRow < rows && newCol >= 0 && newCol < cols &&
					!visited[newRow][newCol] && grid[newRow][newCol] == 1 {
					visited[newRow][newCol] = true
					queue.PushBack([]int{newRow, newCol})
				}
			}
		}
	}

	// 遍历整个网格
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if grid[i][j] == 1 && !visited[i][j] {
				count++
				bfsIsland(i, j)
			}
		}
	}

	return count
}

// BFSWithDistance 带距离的BFS
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func BFSWithDistance(graph map[int][]int, start int) map[int]int {
	visited := make(map[int]bool)
	distance := make(map[int]int)
	queue := list.New()

	queue.PushBack(start)
	visited[start] = true
	distance[start] = 0

	for queue.Len() > 0 {
		node := queue.Remove(queue.Front()).(int)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				distance[neighbor] = distance[node] + 1
				queue.PushBack(neighbor)
			}
		}
	}

	return distance
}

// BFSMultiSource 多源BFS
// 时间复杂度: O(V + E)
// 空间复杂度: O(V)
func BFSMultiSource(graph map[int][]int, sources []int) map[int]int {
	visited := make(map[int]bool)
	distance := make(map[int]int)
	queue := list.New()

	// 将所有源点加入队列
	for _, source := range sources {
		queue.PushBack(source)
		visited[source] = true
		distance[source] = 0
	}

	for queue.Len() > 0 {
		node := queue.Remove(queue.Front()).(int)

		for _, neighbor := range graph[node] {
			if !visited[neighbor] {
				visited[neighbor] = true
				distance[neighbor] = distance[node] + 1
				queue.PushBack(neighbor)
			}
		}
	}

	return distance
}

// 测试函数
func TestBFS() {
	fmt.Println("=== 广度优先搜索算法测试 ===")

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

	// 基本BFS
	bfsResult := BFS(graph, 0)
	fmt.Printf("基本BFS结果: %v\n", bfsResult)

	// 带层级的BFS
	bfsLevelResult := BFSWithLevel(graph, 0)
	fmt.Printf("带层级BFS结果: %v\n", bfsLevelResult)

	// 最短路径
	shortestPath := BFSShortestPath(graph, 0, 5)
	fmt.Printf("从0到5的最短路径: %v\n", shortestPath)

	// 所有最短路径
	allShortestPaths := BFSAllShortestPaths(graph, 0, 4)
	fmt.Printf("从0到4的所有最短路径: %v\n", allShortestPaths)

	// 连通分量
	components := BFSConnectedComponents(graph)
	fmt.Printf("连通分量: %v\n", components)

	// 二分图检测
	isBipartite := BFSBipartite(graph)
	fmt.Printf("是否为二分图: %t\n", isBipartite)

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

	fmt.Println("\n=== 二叉树BFS测试 ===")

	// 层序遍历
	levelOrder := BFSTree(root)
	fmt.Printf("层序遍历: %v\n", levelOrder)

	// 之字形遍历
	zigzag := BFSZigzag(root)
	fmt.Printf("之字形遍历: %v\n", zigzag)

	// 自底向上遍历
	bottomUp := BFSBottomUp(root)
	fmt.Printf("自底向上遍历: %v\n", bottomUp)

	// 二维矩阵测试
	matrix := [][]int{
		{1, 1, 0, 0},
		{1, 1, 0, 0},
		{0, 0, 1, 0},
		{0, 0, 0, 1},
	}

	fmt.Println("\n=== 二维矩阵BFS测试 ===")

	// 矩阵BFS
	matrixBFS := BFSMatrix(matrix, []int{0, 0})
	fmt.Printf("矩阵BFS结果: %v\n", matrixBFS)

	// 岛屿数量
	islands := BFSIslands(matrix)
	fmt.Printf("岛屿数量: %d\n", islands)

	// 距离计算
	distance := BFSWithDistance(graph, 0)
	fmt.Printf("从0到各节点的距离: %v\n", distance)

	// 多源BFS
	multiSourceDistance := BFSMultiSource(graph, []int{0, 2})
	fmt.Printf("多源BFS距离: %v\n", multiSourceDistance)
}
