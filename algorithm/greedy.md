# 贪心搜索算法 (Greedy Search Algorithm)

## 算法概述

贪心搜索算法是一种启发式搜索算法，它在每一步都选择当前看起来最优的选择，而不考虑全局最优解。贪心算法总是选择启发式函数值最小的节点进行扩展，即选择h(n)最小的节点，其中h(n)是从当前节点到目标节点的启发式估计代价。

## 核心特点

### 1. 搜索策略
- **局部最优选择**：每一步都选择当前最优的选项
- **启发式驱动**：完全基于启发式函数h(n)进行决策
- **不回溯**：一旦做出选择就不再改变
- **快速决策**：决策过程简单快速

### 2. 评估函数
- **f(n) = h(n)**：只考虑启发式估计代价
- **忽略实际代价**：不考虑从起点到当前节点的实际代价g(n)
- **纯启发式**：完全依赖启发式函数的质量

### 3. 数据结构
- **优先队列**：按h(n)值排序的开放列表
- **访问标记**：避免重复访问节点
- **路径记录**：记录搜索路径

### 4. 时间复杂度
- **最坏情况**: O(V²) - V为顶点数
- **平均情况**: O(V log V) - 使用优先队列
- **实际性能**: 通常比A*算法更快

### 5. 空间复杂度
- **O(V)**：需要存储访问过的节点
- **优先队列**: O(V) - 开放列表大小
- **路径记录**: O(V) - 父节点映射

## 算法变种

### 1. 基础贪心搜索
```go
// 基本贪心搜索
func GreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64) []int

// 带约束贪心搜索
func ConstrainedGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64, constraints []Constraint) []int

// 多目标贪心搜索
func MultiObjectiveGreedySearch(graph map[int][]Edge, start, target int, heuristics []func(int, int) float64) []int
```

### 2. 启发式函数
```go
// 距离启发式
func DistanceHeuristic(current, target Point) float64

// 代价启发式
func CostHeuristic(current, target int, costMap map[int]float64) float64

// 组合启发式
func CombinedHeuristic(current, target int, heuristics []func(int, int) float64, weights []float64) float64
```

### 3. 优化变种
```go
// 随机贪心搜索
func RandomizedGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64, randomness float64) []int

// 自适应贪心搜索
func AdaptiveGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64) []int

// 分层贪心搜索
func HierarchicalGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64) []int
```

### 4. 特殊应用
```go
// 背包问题贪心
func GreedyKnapsack(items []Item, capacity float64) []Item

// 调度问题贪心
func GreedyScheduling(jobs []Job) []Job

// 最小生成树贪心
func GreedyMST(graph map[int][]Edge) []Edge
```

## 应用场景

### 1. 路径规划
- **快速路径查找**：需要快速找到可行路径的场景
- **实时导航**：实时系统中的路径规划
- **游戏AI**：简单角色的移动路径
- **机器人导航**：快速避障路径规划

### 2. 优化问题
- **背包问题**：物品选择优化
- **调度问题**：任务调度优化
- **分配问题**：资源分配优化
- **排序问题**：元素排序优化

### 3. 图论问题
- **最小生成树**：Kruskal和Prim算法
- **最短路径**：Dijkstra算法的变种
- **网络设计**：网络拓扑优化
- **聚类分析**：数据聚类

### 4. 组合优化
- **集合覆盖**：最小集合覆盖问题
- **顶点覆盖**：最小顶点覆盖问题
- **匹配问题**：最大匹配问题
- **着色问题**：图着色问题

## 算法优势

### 1. 效率优势
- **快速决策**：每一步的决策都很简单快速
- **低计算复杂度**：通常比全局搜索算法快
- **内存效率**：空间复杂度相对较低
- **实时性能**：适合实时系统

### 2. 实现简单
- **算法逻辑**：逻辑简单，易于理解和实现
- **调试方便**：决策过程清晰，便于调试
- **扩展性好**：容易添加约束和优化

### 3. 适用性广
- **问题类型**：适用于多种优化问题
- **启发式灵活**：可以使用各种启发式函数
- **约束处理**：容易处理各种约束条件

## 算法劣势

### 1. 最优性不保证
- **局部最优**：可能陷入局部最优解
- **全局视角缺失**：不考虑全局最优性
- **解的质量**：解的质量可能不如全局搜索算法

### 2. 启发式依赖
- **启发式质量**：解的质量高度依赖启发式函数
- **设计困难**：好的启发式函数可能难以设计
- **问题特定**：启发式函数通常针对特定问题

### 3. 不可回溯
- **决策不可逆**：一旦做出选择就不能改变
- **错误累积**：早期错误会影响后续决策
- **适应性差**：难以适应动态变化的环境

## 启发式函数设计

### 1. 距离启发式
```go
// 欧几里得距离
func EuclideanHeuristic(current, target Point) float64 {
    dx := current.x - target.x
    dy := current.y - target.y
    return math.Sqrt(dx*dx + dy*dy)
}

// 曼哈顿距离
func ManhattanHeuristic(current, target Point) float64 {
    return float64(abs(current.x-target.x) + abs(current.y-target.y))
}

// 切比雪夫距离
func ChebyshevHeuristic(current, target Point) float64 {
    return float64(max(abs(current.x-target.x), abs(current.y-target.y)))
}
```

### 2. 代价启发式
```go
// 基于代价的启发式
func CostBasedHeuristic(current, target int, costMap map[int]float64) float64 {
    if cost, exists := costMap[target]; exists {
        return cost
    }
    return 0.0
}

// 基于权重的启发式
func WeightBasedHeuristic(current, target int, weightMap map[int]float64) float64 {
    if weight, exists := weightMap[target]; exists {
        return 1.0 / weight // 权重越大，启发式值越小
    }
    return 1.0
}
```

### 3. 问题特定启发式
```go
// 8数码问题的启发式
func EightPuzzleHeuristic(current, target []int) float64 {
    misplaced := 0
    for i := 0; i < 9; i++ {
        if current[i] != target[i] && current[i] != 0 {
            misplaced++
        }
    }
    return float64(misplaced)
}

// 背包问题的启发式
func KnapsackHeuristic(item Item, remainingCapacity float64) float64 {
    if item.weight > remainingCapacity {
        return 0.0 // 不可选
    }
    return item.value / item.weight // 价值密度
}
```

## 优化技巧

### 1. 随机化优化
```go
// 随机贪心搜索
func RandomizedGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64, randomness float64) []int {
    // 在每一步决策时引入随机性
    // 避免陷入局部最优
}

// ε-贪心策略
func EpsilonGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64, epsilon float64) []int {
    // 以ε概率选择随机动作，以1-ε概率选择贪心动作
}
```

### 2. 自适应优化
```go
// 自适应贪心搜索
func AdaptiveGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64) []int {
    // 根据搜索历史动态调整启发式函数
    // 避免重复访问相同的局部最优
}

// 学习贪心搜索
func LearningGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64) []int {
    // 使用机器学习方法学习更好的启发式函数
}
```

### 3. 约束处理
```go
// 带约束的贪心搜索
func ConstrainedGreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64, constraints []Constraint) []int {
    // 在贪心选择时考虑约束条件
    // 确保解满足所有约束
}

// 多目标贪心搜索
func MultiObjectiveGreedySearch(graph map[int][]Edge, start, target int, heuristics []func(int, int) float64, weights []float64) []int {
    // 考虑多个目标函数
    // 使用加权和或帕累托最优
}
```

## 代码示例

### 基本贪心搜索实现
```go
type Node struct {
    id       int
    h        float64 // 启发式估计代价
    parent   int     // 父节点
}

func GreedySearch(graph map[int][]Edge, start, target int, heuristic func(int, int) float64) []int {
    openList := make(map[int]*Node)
    closedList := make(map[int]bool)
    
    // 初始化起始节点
    startNode := &Node{
        id:     start,
        h:      heuristic(start, target),
        parent: -1,
    }
    openList[start] = startNode
    
    for len(openList) > 0 {
        // 找到h值最小的节点
        current := findMinHNode(openList)
        
        if current.id == target {
            return reconstructPath(current, openList)
        }
        
        // 从开放列表移除，加入关闭列表
        delete(openList, current.id)
        closedList[current.id] = true
        
        // 检查所有邻居
        for _, edge := range graph[current.id] {
            neighborID := edge.to
            
            if closedList[neighborID] {
                continue
            }
            
            neighbor, exists := openList[neighborID]
            if !exists {
                neighbor = &Node{
                    id:     neighborID,
                    h:      heuristic(neighborID, target),
                    parent: current.id,
                }
                openList[neighborID] = neighbor
            }
        }
    }
    
    return []int{} // 没有找到路径
}
```

### 贪心背包算法
```go
type Item struct {
    id     int
    weight float64
    value  float64
}

func GreedyKnapsack(items []Item, capacity float64) []Item {
    // 按价值密度排序（价值/重量）
    sort.Slice(items, func(i, j int) bool {
        return items[i].value/items[i].weight > items[j].value/items[j].weight
    })
    
    selected := []Item{}
    remainingCapacity := capacity
    
    for _, item := range items {
        if item.weight <= remainingCapacity {
            selected = append(selected, item)
            remainingCapacity -= item.weight
        }
    }
    
    return selected
}
```

### 贪心调度算法
```go
type Job struct {
    id       int
    duration float64
    deadline float64
    priority float64
}

func GreedyScheduling(jobs []Job) []Job {
    // 按优先级排序
    sort.Slice(jobs, func(i, j int) bool {
        return jobs[i].priority > jobs[j].priority
    })
    
    scheduled := []Job{}
    currentTime := 0.0
    
    for _, job := range jobs {
        if currentTime + job.duration <= job.deadline {
            scheduled = append(scheduled, job)
            currentTime += job.duration
        }
    }
    
    return scheduled
}
```

## 与其他算法的比较

| 特性 | 贪心搜索 | A*搜索 | Dijkstra | BFS |
|------|----------|--------|----------|-----|
| 启发式使用 | ✅ (仅h(n)) | ✅ (g(n)+h(n)) | ❌ | ❌ |
| 最优性保证 | ❌ | ✅ (可接受启发式) | ✅ | ✅ (无权图) |
| 效率 | 最高 | 高 | 中等 | 低 |
| 解的质量 | 可能次优 | 最优 | 最优 | 最优 |
| 实现复杂度 | 简单 | 中等 | 简单 | 简单 |
| 适用场景 | 快速近似解 | 最优路径 | 最短路径 | 无权图 |

## 实际应用案例

### 1. 网络路由
```go
// 贪心路由算法
func GreedyRouting(network map[int][]Link, start, target int) []int {
    // 在每一步选择延迟最小的链路
    // 快速找到可行路径
}
```

### 2. 资源分配
```go
// 贪心资源分配
func GreedyResourceAllocation(resources []Resource, tasks []Task) map[int]Resource {
    // 按任务优先级分配资源
    // 最大化资源利用率
}
```

### 3. 任务调度
```go
// 贪心任务调度
func GreedyTaskScheduling(tasks []Task, processors []Processor) map[int][]Task {
    // 将任务分配给负载最小的处理器
    // 平衡处理器负载
}
```

## 总结

贪心搜索算法是一种简单而高效的启发式搜索算法，特别适合需要快速找到可行解的场景。它通过局部最优选择策略，在保证效率的同时提供合理的解决方案。

贪心算法的核心优势在于其简单性和高效性，这使得它在实时系统、快速原型开发和近似算法中表现出色。虽然不能保证全局最优解，但在许多实际应用中，贪心算法提供的解已经足够好。

通过合理设计启发式函数和优化策略，贪心算法可以适应各种不同的应用场景，从简单的路径规划到复杂的组合优化问题。它是算法工具箱中重要的基础工具，为更复杂的算法提供了重要的基础。 