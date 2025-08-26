# 广度优先搜索 (Breadth-First Search, BFS)

## 算法概述

广度优先搜索是一种用于遍历或搜索树或图的算法。它从根节点开始，逐层遍历所有相邻节点，先访问所有相邻节点，再访问下一层的节点。BFS保证在访问下一层节点之前，当前层的所有节点都已被访问。

## 核心特点

### 1. 搜索策略
- **广度优先**：逐层遍历，先访问当前层的所有节点
- **队列实现**：使用队列数据结构保证先进先出的访问顺序
- **最短路径**：在无权图中能找到从起点到目标的最短路径

### 2. 数据结构
- **队列**：使用队列存储待访问的节点
- **访问标记**：避免重复访问节点
- **距离记录**：记录从起点到每个节点的距离

### 3. 时间复杂度
- **图遍历**: O(V + E) - V为顶点数，E为边数
- **树遍历**: O(n) - n为节点数
- **最短路径**: O(V + E) - 在无权图中

### 4. 空间复杂度
- **队列空间**: O(V) - 最坏情况下队列中存储所有节点
- **访问标记**: O(V) - 存储所有节点的访问状态
- **树遍历**: O(w) - w为树的最大宽度

## 算法变种

### 1. 图遍历
```go
// 基本BFS
func BFS(graph map[int][]int, start int) []int

// 带距离记录
func BFSWithDistance(graph map[int][]int, start int) map[int]int

// 最短路径
func BFSShortestPath(graph map[int][]int, start, target int) []int

// 层次遍历
func BFSLevelOrder(graph map[int][]int, start int) [][]int
```

### 2. 树遍历
```go
// 层次遍历
func BFSLevelOrderTree(root *TreeNode) [][]int

// 锯齿形遍历
func BFSZigzagLevelOrder(root *TreeNode) [][]int

// 自底向上遍历
func BFSBottomUpLevelOrder(root *TreeNode) [][]int
```

### 3. 图算法
```go
// 连通分量
func BFSConnectedComponents(graph map[int][]int) [][]int

// 二分图检测
func BFSBipartiteCheck(graph map[int][]int) bool

// 拓扑排序
func BFSTopologicalSort(graph map[int][]int) []int
```

### 4. 特殊应用
```go
// 多源BFS
func MultiSourceBFS(grid [][]int, sources [][]int) [][]int

// 双向BFS
func BidirectionalBFS(graph map[int][]int, start, target int) []int
```

## 应用场景

### 1. 图论问题
- **最短路径**：在无权图中找到最短路径
- **连通性检测**：判断图中两个节点是否连通
- **层次分析**：分析图的层次结构
- **网络分析**：分析社交网络、通信网络等

### 2. 树相关问题
- **层次遍历**：按层次访问树节点
- **树的高度**：计算树的最小高度
- **最近公共祖先**：查找两个节点的LCA
- **树的序列化**：将树转换为字符串

### 3. 矩阵和网格问题
- **岛屿数量**：计算矩阵中岛屿的数量
- **最短路径**：在网格中寻找最短路径
- **包围区域**：标记被包围的区域
- **机器人路径**：机器人寻路问题

### 4. 游戏和AI
- **迷宫求解**：寻找迷宫的最短路径
- **八数码问题**：移动数字块到目标状态
- **状态空间搜索**：探索问题的状态空间
- **游戏AI**：游戏中的路径规划

## 算法优势

### 1. 最短路径保证
- **最优性**：在无权图中保证找到最短路径
- **完整性**：能够找到所有可达的节点
- **层次性**：按层次组织搜索结果

### 2. 实现简单
- **队列操作**：使用标准队列数据结构
- **逻辑清晰**：访问顺序直观易懂
- **调试方便**：层次结构便于调试

### 3. 适用性广
- **通用性强**：适用于各种图结构
- **扩展性好**：容易添加约束条件
- **并行友好**：可以并行处理同一层的节点

## 算法劣势

### 1. 内存消耗
- **空间复杂度高**：需要存储整层的节点
- **队列大小**：在最坏情况下队列可能很大
- **不适合深度搜索**：在深度较大的问题中效率较低

### 2. 性能问题
- **可能访问不必要节点**：在目标较深时可能访问过多节点
- **不适合启发式搜索**：无法利用问题特定的启发信息
- **层次限制**：在层次较多时效率下降

## 优化技巧

### 1. 双向BFS
```go
// 从起点和终点同时开始搜索
func BidirectionalBFS(graph map[int][]int, start, target int) []int {
    if start == target {
        return []int{start}
    }
    
    // 从起点开始的队列
    startQueue := list.New()
    startQueue.PushBack(start)
    startVisited := make(map[int]int)
    startVisited[start] = 0
    
    // 从终点开始的队列
    targetQueue := list.New()
    targetQueue.PushBack(target)
    targetVisited := make(map[int]int)
    targetVisited[target] = 0
    
    // 交替搜索
    for startQueue.Len() > 0 && targetQueue.Len() > 0 {
        // 从起点搜索一层
        if result := searchLevel(graph, startQueue, startVisited, targetVisited); result != nil {
            return result
        }
        
        // 从终点搜索一层
        if result := searchLevel(graph, targetQueue, targetVisited, startVisited); result != nil {
            return result
        }
    }
    
    return []int{}
}
```

### 2. 多源BFS
```go
// 从多个源点同时开始搜索
func MultiSourceBFS(grid [][]int, sources [][]int) [][]int {
    queue := list.New()
    visited := make([][]bool, len(grid))
    
    // 初始化访问矩阵
    for i := range visited {
        visited[i] = make([]bool, len(grid[0]))
    }
    
    // 将所有源点加入队列
    for _, source := range sources {
        queue.PushBack(source)
        visited[source[0]][source[1]] = true
    }
    
    // BFS搜索
    for queue.Len() > 0 {
        current := queue.Remove(queue.Front()).([]int)
        // 处理当前节点
    }
    
    return grid
}
```

### 3. 层次优化
```go
// 记录层次信息
func BFSWithLevel(graph map[int][]int, start int) map[int]int {
    queue := list.New()
    visited := make(map[int]bool)
    level := make(map[int]int)
    
    queue.PushBack(start)
    visited[start] = true
    level[start] = 0
    
    for queue.Len() > 0 {
        node := queue.Remove(queue.Front()).(int)
        currentLevel := level[node]
        
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                level[neighbor] = currentLevel + 1
                queue.PushBack(neighbor)
            }
        }
    }
    
    return level
}
```

## 代码示例

### 基本BFS实现
```go
func BFS(graph map[int][]int, start int) []int {
    queue := list.New()
    visited := make(map[int]bool)
    result := []int{}
    
    queue.PushBack(start)
    visited[start] = true
    
    for queue.Len() > 0 {
        node := queue.Remove(queue.Front()).(int)
        result = append(result, node)
        
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue.PushBack(neighbor)
            }
        }
    }
    
    return result
}
```

### 最短路径BFS
```go
func BFSShortestPath(graph map[int][]int, start, target int) []int {
    if start == target {
        return []int{start}
    }
    
    queue := list.New()
    visited := make(map[int]bool)
    parent := make(map[int]int)
    
    queue.PushBack(start)
    visited[start] = true
    
    for queue.Len() > 0 {
        node := queue.Remove(queue.Front()).(int)
        
        if node == target {
            // 重建路径
            path := []int{}
            for node != start {
                path = append([]int{node}, path...)
                node = parent[node]
            }
            return append([]int{start}, path...)
        }
        
        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                parent[neighbor] = node
                queue.PushBack(neighbor)
            }
        }
    }
    
    return []int{} // 没有找到路径
}
```

### 层次遍历
```go
func BFSLevelOrder(graph map[int][]int, start int) [][]int {
    queue := list.New()
    visited := make(map[int]bool)
    result := [][]int{}
    
    queue.PushBack(start)
    visited[start] = true
    
    for queue.Len() > 0 {
        levelSize := queue.Len()
        level := []int{}
        
        for i := 0; i < levelSize; i++ {
            node := queue.Remove(queue.Front()).(int)
            level = append(level, node)
            
            for _, neighbor := range graph[node] {
                if !visited[neighbor] {
                    visited[neighbor] = true
                    queue.PushBack(neighbor)
                }
            }
        }
        
        result = append(result, level)
    }
    
    return result
}
```

## 与DFS的比较

| 特性 | BFS | DFS |
|------|-----|-----|
| 搜索策略 | 广度优先，逐层搜索 | 深度优先，一条路走到底 |
| 数据结构 | 队列 | 栈（递归调用栈） |
| 最短路径 | 保证找到最短路径 | 不保证最短路径 |
| 空间复杂度 | O(V) - 存储整层节点 | O(h) - 存储当前路径 |
| 适用场景 | 最短路径、层次分析 | 深度探索、回溯问题 |
| 实现复杂度 | 简单 | 简单（递归） |

## 总结

广度优先搜索是一种强大而实用的算法，特别适合需要找到最短路径或进行层次分析的问题。它的队列实现保证了访问顺序的正确性，在无权图中能够保证找到最优解。

BFS的核心思想是"逐层探索"，这种策略在需要广度覆盖的问题中表现出色，是图论和树论中不可或缺的重要工具。通过合理的优化和变种，BFS可以解决从简单的图遍历到复杂的网络分析等各种问题。 