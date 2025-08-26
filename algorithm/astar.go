package algorithm

import (
	"container/heap"
	"fmt"
	"math"
)

// Node A*搜索节点
type Node struct {
	X, Y    int     // 坐标
	F, G, H float64 // F = G + H (总代价 = 实际代价 + 启发式代价)
	Parent  *Node   // 父节点，用于重建路径
	index   int     // 优先级队列索引
}

// PriorityQueue 优先级队列实现
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F // 按F值排序，F值越小优先级越高
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	old[n-1] = nil
	node.index = -1
	*pq = old[0 : n-1]
	return node
}

// Grid 网格地图
type Grid struct {
	Width, Height int
	Obstacles     map[string]bool // 障碍物位置
}

// NewGrid 创建新网格
func NewGrid(width, height int) *Grid {
	return &Grid{
		Width:     width,
		Height:    height,
		Obstacles: make(map[string]bool),
	}
}

// AddObstacle 添加障碍物
func (g *Grid) AddObstacle(x, y int) {
	g.Obstacles[fmt.Sprintf("%d,%d", x, y)] = true
}

// IsValid 检查位置是否有效
func (g *Grid) IsValid(x, y int) bool {
	return x >= 0 && x < g.Width && y >= 0 && y < g.Height
}

// IsWalkable 检查位置是否可通行
func (g *Grid) IsWalkable(x, y int) bool {
	return g.IsValid(x, y) && !g.Obstacles[fmt.Sprintf("%d,%d", x, y)]
}

// Heuristic 启发式函数
type Heuristic func(x1, y1, x2, y2 int) float64

// ManhattanDistance 曼哈顿距离启发式函数
// 适用于只能上下左右移动的网格
func ManhattanDistance(x1, y1, x2, y2 int) float64 {
	return float64(absWithAstar(x1-x2) + absWithAstar(y1-y2))
}

// EuclideanDistance 欧几里得距离启发式函数
// 适用于可以斜向移动的网格
func EuclideanDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

// ChebyshevDistance 切比雪夫距离启发式函数
// 适用于可以斜向移动且代价相同的网格
func ChebyshevDistance(x1, y1, x2, y2 int) float64 {
	return float64(maxWithAstar(absWithAstar(x1-x2), absWithAstar(y1-y2)))
}

// OctileDistance 八方向距离启发式函数
// 适用于可以斜向移动但斜向代价更高的网格
func OctileDistance(x1, y1, x2, y2 int) float64 {
	dx := absWithAstar(x1 - x2)
	dy := absWithAstar(y1 - y2)
	return float64(maxWithAstar(dx, dy)) + (math.Sqrt(2)-1)*float64(minWithAstar(dx, dy))
}

// AStar A*搜索算法主函数
// 时间复杂度: O(V log V) - V为节点数
// 空间复杂度: O(V)
func AStar(grid *Grid, startX, startY, goalX, goalY int, heuristic Heuristic) []Node {
	if !grid.IsWalkable(startX, startY) || !grid.IsWalkable(goalX, goalY) {
		return nil
	}

	// 初始化开放列表和关闭列表
	openList := &PriorityQueue{}
	heap.Init(openList)
	closedSet := make(map[string]bool)
	cameFrom := make(map[string]*Node)
	gScore := make(map[string]float64)
	fScore := make(map[string]float64)

	// 创建起始节点
	startNode := &Node{
		X: startX,
		Y: startY,
		G: 0,
		H: heuristic(startX, startY, goalX, goalY),
	}
	startNode.F = startNode.G + startNode.H

	// 将起始节点加入开放列表
	heap.Push(openList, startNode)
	gScore[fmt.Sprintf("%d,%d", startX, startY)] = 0
	fScore[fmt.Sprintf("%d,%d", startX, startY)] = startNode.F

	// 定义移动方向（8方向）
	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for openList.Len() > 0 {
		// 获取F值最小的节点
		current := heap.Pop(openList).(*Node)
		currentKey := fmt.Sprintf("%d,%d", current.X, current.Y)

		// 如果到达目标
		if current.X == goalX && current.Y == goalY {
			return reconstructPath(cameFrom, current)
		}

		// 将当前节点加入关闭列表
		closedSet[currentKey] = true

		// 检查所有邻居
		for _, dir := range directions {
			neighborX := current.X + dir[0]
			neighborY := current.Y + dir[1]
			neighborKey := fmt.Sprintf("%d,%d", neighborX, neighborY)

			// 检查邻居是否有效且可通行
			if !grid.IsWalkable(neighborX, neighborY) {
				continue
			}

			// 检查邻居是否已在关闭列表中
			if closedSet[neighborKey] {
				continue
			}

			// 计算移动代价
			moveCost := 1.0
			if dir[0] != 0 && dir[1] != 0 {
				moveCost = math.Sqrt(2) // 斜向移动代价更高
			}

			// 计算从起始点到邻居的代价
			tentativeGScore := gScore[currentKey] + moveCost

			// 检查邻居是否已在开放列表中
			if _, exists := gScore[neighborKey]; !exists {
				// 新节点，加入开放列表
				neighbor := &Node{
					X:      neighborX,
					Y:      neighborY,
					G:      tentativeGScore,
					H:      heuristic(neighborX, neighborY, goalX, goalY),
					Parent: current,
				}
				neighbor.F = neighbor.G + neighbor.H

				heap.Push(openList, neighbor)
				gScore[neighborKey] = tentativeGScore
				fScore[neighborKey] = neighbor.F
				cameFrom[neighborKey] = current
			} else if tentativeGScore < gScore[neighborKey] {
				// 找到更好的路径，更新节点
				gScore[neighborKey] = tentativeGScore
				fScore[neighborKey] = tentativeGScore + heuristic(neighborX, neighborY, goalX, goalY)
				cameFrom[neighborKey] = current

				// 更新开放列表中的节点
				for i, node := range *openList {
					if node.X == neighborX && node.Y == neighborY {
						node.G = tentativeGScore
						node.F = fScore[neighborKey]
						node.Parent = current
						heap.Fix(openList, i)
						break
					}
				}
			}
		}
	}

	// 没有找到路径
	return nil
}

// AStarWithWeights 带权重的A*搜索
// 可以设置不同类型地形的移动代价
func AStarWithWeights(grid *Grid, startX, startY, goalX, goalY int,
	heuristic Heuristic, terrainCosts map[string]float64) []Node {

	if !grid.IsWalkable(startX, startY) || !grid.IsWalkable(goalX, goalY) {
		return nil
	}

	openList := &PriorityQueue{}
	heap.Init(openList)
	closedSet := make(map[string]bool)
	cameFrom := make(map[string]*Node)
	gScore := make(map[string]float64)
	fScore := make(map[string]float64)

	startNode := &Node{
		X: startX,
		Y: startY,
		G: 0,
		H: heuristic(startX, startY, goalX, goalY),
	}
	startNode.F = startNode.G + startNode.H

	heap.Push(openList, startNode)
	gScore[fmt.Sprintf("%d,%d", startX, startY)] = 0
	fScore[fmt.Sprintf("%d,%d", startX, startY)] = startNode.F

	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*Node)
		currentKey := fmt.Sprintf("%d,%d", current.X, current.Y)

		if current.X == goalX && current.Y == goalY {
			return reconstructPath(cameFrom, current)
		}

		closedSet[currentKey] = true

		for _, dir := range directions {
			neighborX := current.X + dir[0]
			neighborY := current.Y + dir[1]
			neighborKey := fmt.Sprintf("%d,%d", neighborX, neighborY)

			if !grid.IsWalkable(neighborX, neighborY) {
				continue
			}

			if closedSet[neighborKey] {
				continue
			}

			// 获取地形代价
			terrainKey := fmt.Sprintf("%d,%d", neighborX, neighborY)
			baseCost := 1.0
			if cost, exists := terrainCosts[terrainKey]; exists {
				baseCost = cost
			}

			// 计算移动代价
			moveCost := baseCost
			if dir[0] != 0 && dir[1] != 0 {
				moveCost = baseCost * math.Sqrt(2)
			}

			tentativeGScore := gScore[currentKey] + moveCost

			if _, exists := gScore[neighborKey]; !exists {
				neighbor := &Node{
					X:      neighborX,
					Y:      neighborY,
					G:      tentativeGScore,
					H:      heuristic(neighborX, neighborY, goalX, goalY),
					Parent: current,
				}
				neighbor.F = neighbor.G + neighbor.H

				heap.Push(openList, neighbor)
				gScore[neighborKey] = tentativeGScore
				fScore[neighborKey] = neighbor.F
				cameFrom[neighborKey] = current
			} else if tentativeGScore < gScore[neighborKey] {
				gScore[neighborKey] = tentativeGScore
				fScore[neighborKey] = tentativeGScore + heuristic(neighborX, neighborY, goalX, goalY)
				cameFrom[neighborKey] = current

				for i, node := range *openList {
					if node.X == neighborX && node.Y == neighborY {
						node.G = tentativeGScore
						node.F = fScore[neighborKey]
						node.Parent = current
						heap.Fix(openList, i)
						break
					}
				}
			}
		}
	}

	return nil
}

// AStarMultiGoal 多目标A*搜索
// 寻找访问所有目标点的最短路径
func AStarMultiGoal(grid *Grid, startX, startY int, goals [][]int, heuristic Heuristic) []Node {
	if len(goals) == 0 {
		return nil
	}

	// 使用最近邻启发式：总是选择最近的目标
	var path []Node
	currentX, currentY := startX, startY
	remainingGoals := make([][]int, len(goals))
	copy(remainingGoals, goals)

	for len(remainingGoals) > 0 {
		// 找到最近的目标
		closestGoal := 0
		minDistance := math.MaxFloat64
		for i, goal := range remainingGoals {
			dist := heuristic(currentX, currentY, goal[0], goal[1])
			if dist < minDistance {
				minDistance = dist
				closestGoal = i
			}
		}

		// 寻找到最近目标的路径
		goal := remainingGoals[closestGoal]
		subPath := AStar(grid, currentX, currentY, goal[0], goal[1], heuristic)
		if subPath == nil {
			return nil // 无法到达目标
		}

		// 添加路径（除了第一个点，避免重复）
		if len(path) == 0 {
			path = append(path, subPath...)
		} else {
			path = append(path, subPath[1:]...)
		}

		// 更新当前位置
		currentX, currentY = goal[0], goal[1]

		// 移除已访问的目标
		remainingGoals = append(remainingGoals[:closestGoal], remainingGoals[closestGoal+1:]...)
	}

	return path
}

// AStarBidirectional 双向A*搜索
// 从起点和终点同时开始搜索，在中间相遇
func AStarBidirectional(grid *Grid, startX, startY, goalX, goalY int, heuristic Heuristic) []Node {
	if !grid.IsWalkable(startX, startY) || !grid.IsWalkable(goalX, goalY) {
		return nil
	}

	// 正向搜索
	forwardOpen := &PriorityQueue{}
	heap.Init(forwardOpen)
	forwardClosed := make(map[string]bool)
	forwardCameFrom := make(map[string]*Node)
	forwardGScore := make(map[string]float64)

	// 反向搜索
	backwardOpen := &PriorityQueue{}
	heap.Init(backwardOpen)
	backwardClosed := make(map[string]bool)
	backwardCameFrom := make(map[string]*Node)
	backwardGScore := make(map[string]float64)

	// 初始化起始节点
	startNode := &Node{X: startX, Y: startY, G: 0, H: heuristic(startX, startY, goalX, goalY)}
	startNode.F = startNode.G + startNode.H
	heap.Push(forwardOpen, startNode)
	forwardGScore[fmt.Sprintf("%d,%d", startX, startY)] = 0

	goalNode := &Node{X: goalX, Y: goalY, G: 0, H: heuristic(goalX, goalY, startX, startY)}
	goalNode.F = goalNode.G + goalNode.H
	heap.Push(backwardOpen, goalNode)
	backwardGScore[fmt.Sprintf("%d,%d", goalX, goalY)] = 0

	directions := [][]int{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	meetingPoint := ""
	bestCost := math.MaxFloat64

	// 交替进行正向和反向搜索
	for forwardOpen.Len() > 0 && backwardOpen.Len() > 0 {
		// 正向搜索
		if forwardOpen.Len() > 0 {
			current := heap.Pop(forwardOpen).(*Node)
			currentKey := fmt.Sprintf("%d,%d", current.X, current.Y)

			if forwardClosed[currentKey] {
				continue
			}

			forwardClosed[currentKey] = true

			// 检查是否与反向搜索相遇
			if backwardClosed[currentKey] {
				totalCost := forwardGScore[currentKey] + backwardGScore[currentKey]
				if totalCost < bestCost {
					bestCost = totalCost
					meetingPoint = currentKey
				}
			}

			// 扩展邻居
			for _, dir := range directions {
				neighborX := current.X + dir[0]
				neighborY := current.Y + dir[1]
				neighborKey := fmt.Sprintf("%d,%d", neighborX, neighborY)

				if !grid.IsWalkable(neighborX, neighborY) || forwardClosed[neighborKey] {
					continue
				}

				moveCost := 1.0
				if dir[0] != 0 && dir[1] != 0 {
					moveCost = math.Sqrt(2)
				}

				tentativeGScore := forwardGScore[currentKey] + moveCost

				if _, exists := forwardGScore[neighborKey]; !exists || tentativeGScore < forwardGScore[neighborKey] {
					neighbor := &Node{
						X:      neighborX,
						Y:      neighborY,
						G:      tentativeGScore,
						H:      heuristic(neighborX, neighborY, goalX, goalY),
						Parent: current,
					}
					neighbor.F = neighbor.G + neighbor.H

					heap.Push(forwardOpen, neighbor)
					forwardGScore[neighborKey] = tentativeGScore
					forwardCameFrom[neighborKey] = current
				}
			}
		}

		// 反向搜索
		if backwardOpen.Len() > 0 {
			current := heap.Pop(backwardOpen).(*Node)
			currentKey := fmt.Sprintf("%d,%d", current.X, current.Y)

			if backwardClosed[currentKey] {
				continue
			}

			backwardClosed[currentKey] = true

			// 检查是否与正向搜索相遇
			if forwardClosed[currentKey] {
				totalCost := forwardGScore[currentKey] + backwardGScore[currentKey]
				if totalCost < bestCost {
					bestCost = totalCost
					meetingPoint = currentKey
				}
			}

			// 扩展邻居
			for _, dir := range directions {
				neighborX := current.X + dir[0]
				neighborY := current.Y + dir[1]
				neighborKey := fmt.Sprintf("%d,%d", neighborX, neighborY)

				if !grid.IsWalkable(neighborX, neighborY) || backwardClosed[neighborKey] {
					continue
				}

				moveCost := 1.0
				if dir[0] != 0 && dir[1] != 0 {
					moveCost = math.Sqrt(2)
				}

				tentativeGScore := backwardGScore[currentKey] + moveCost

				if _, exists := backwardGScore[neighborKey]; !exists || tentativeGScore < backwardGScore[neighborKey] {
					neighbor := &Node{
						X:      neighborX,
						Y:      neighborY,
						G:      tentativeGScore,
						H:      heuristic(neighborX, neighborY, startX, startY),
						Parent: current,
					}
					neighbor.F = neighbor.G + neighbor.H

					heap.Push(backwardOpen, neighbor)
					backwardGScore[neighborKey] = tentativeGScore
					backwardCameFrom[neighborKey] = current
				}
			}
		}
	}

	// 重建路径
	if meetingPoint != "" {
		return reconstructBidirectionalPath(forwardCameFrom, backwardCameFrom, meetingPoint)
	}

	return nil
}

// 辅助函数
func reconstructPath(cameFrom map[string]*Node, current *Node) []Node {
	var path []Node
	for current != nil {
		path = append([]Node{{X: current.X, Y: current.Y}}, path...)
		current = current.Parent
	}
	return path
}

func reconstructBidirectionalPath(forwardCameFrom, backwardCameFrom map[string]*Node, meetingPoint string) []Node {
	var path []Node

	// 从相遇点向前重建
	current := forwardCameFrom[meetingPoint]
	for current != nil {
		path = append([]Node{{X: current.X, Y: current.Y}}, path...)
		current = current.Parent
	}

	// 从相遇点向后重建
	current = backwardCameFrom[meetingPoint]
	for current != nil {
		path = append(path, Node{X: current.X, Y: current.Y})
		current = current.Parent
	}

	return path
}

func absWithAstar(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func maxWithAstar(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minWithAstar(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 测试函数
func TestAStar() {
	fmt.Println("=== A*搜索算法测试 ===")

	// 创建网格
	grid := NewGrid(10, 10)

	// 添加一些障碍物
	grid.AddObstacle(2, 2)
	grid.AddObstacle(2, 3)
	grid.AddObstacle(2, 4)
	grid.AddObstacle(3, 4)
	grid.AddObstacle(4, 4)
	grid.AddObstacle(5, 4)
	grid.AddObstacle(6, 4)
	grid.AddObstacle(7, 4)

	// 基本A*搜索
	fmt.Println("基本A*搜索:")
	path := AStar(grid, 0, 0, 9, 9, ManhattanDistance)
	if path != nil {
		fmt.Printf("找到路径，长度: %d\n", len(path))
		fmt.Printf("路径: %v\n", path)
	} else {
		fmt.Println("未找到路径")
	}

	// 不同启发式函数比较
	fmt.Println("\n不同启发式函数比较:")
	heuristics := map[string]Heuristic{
		"曼哈顿距离":  ManhattanDistance,
		"欧几里得距离": EuclideanDistance,
		"切比雪夫距离": ChebyshevDistance,
		"八方向距离":  OctileDistance,
	}

	for name, heuristic := range heuristics {
		path := AStar(grid, 0, 0, 9, 9, heuristic)
		if path != nil {
			fmt.Printf("%s: 路径长度 %d\n", name, len(path))
		} else {
			fmt.Printf("%s: 未找到路径\n", name)
		}
	}

	// 带权重的A*搜索
	fmt.Println("\n带权重的A*搜索:")
	terrainCosts := map[string]float64{
		"1,1": 2.0, // 沼泽地
		"1,2": 2.0,
		"2,1": 2.0,
		"3,3": 1.5, // 丘陵
		"4,3": 1.5,
		"5,3": 1.5,
	}

	weightedPath := AStarWithWeights(grid, 0, 0, 9, 9, ManhattanDistance, terrainCosts)
	if weightedPath != nil {
		fmt.Printf("带权重路径长度: %d\n", len(weightedPath))
	} else {
		fmt.Println("未找到带权重路径")
	}

	// 多目标A*搜索
	fmt.Println("\n多目标A*搜索:")
	goals := [][]int{{3, 3}, {7, 7}, {9, 9}}
	multiPath := AStarMultiGoal(grid, 0, 0, goals, ManhattanDistance)
	if multiPath != nil {
		fmt.Printf("多目标路径长度: %d\n", len(multiPath))
	} else {
		fmt.Println("未找到多目标路径")
	}

	// 双向A*搜索
	fmt.Println("\n双向A*搜索:")
	bidirectionalPath := AStarBidirectional(grid, 0, 0, 9, 9, ManhattanDistance)
	if bidirectionalPath != nil {
		fmt.Printf("双向搜索路径长度: %d\n", len(bidirectionalPath))
	} else {
		fmt.Println("未找到双向搜索路径")
	}
}
