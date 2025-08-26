# 深度优先搜索 (Depth-First Search, DFS)

## 算法概述

深度优先搜索是一种用于遍历或搜索树或图的算法。它沿着树的深度遍历树的节点，尽可能深的搜索树的分支。当节点v的所在边都已被探寻过，搜索将回溯到发现节点v的那条边的起始节点。

## 核心特点

### 1. 搜索策略
- **深度优先**：沿着一条路径一直走到底，直到无法继续前进
- **回溯机制**：当到达死胡同或目标节点时，回退到上一个节点，尝试其他路径
- **递归实现**：自然符合DFS的思维模式，代码简洁易懂

### 2. 数据结构
- **递归调用栈**：递归版本使用系统调用栈
- **显式栈**：迭代版本使用栈数据结构
- **访问标记**：避免重复访问节点

### 3. 时间复杂度
- **图遍历**: O(V + E) - V为顶点数，E为边数
- **树遍历**: O(n) - n为节点数
- **回溯算法**: 取决于具体问题，通常是指数级

### 4. 空间复杂度
- **递归版本**: O(V) - 递归调用栈深度
- **迭代版本**: O(V) - 显式栈空间
- **树遍历**: O(h) - h为树的高度

## 算法变种

### 1. 图遍历
```go
// 递归DFS
func DFS(graph map[int][]int, start int) []int

// 迭代DFS
func DFSIterative(graph map[int][]int, start int) []int

// 带路径记录
func DFSWithPath(graph map[int][]int, start, target int) []int

// 查找所有路径
func DFSAllPaths(graph map[int][]int, start, target int) [][]int
```

### 2. 树遍历
```go
// 前序遍历：根 -> 左 -> 右
func DFSTree(root *TreeNode) []int

// 中序遍历：左 -> 根 -> 右
func DFSInorder(root *TreeNode) []int

// 后序遍历：左 -> 右 -> 根
func DFSPostorder(root *TreeNode) []int
```

### 3. 图算法
```go
// 连通分量查找
func DFSConnectedComponents(graph map[int][]int) [][]int

// 环检测
func DFSCycleDetection(graph map[int][]int) bool

// 拓扑排序
func DFSTopologicalSort(graph map[int][]int) []int
```

### 4. 回溯算法
```go
// 通用回溯框架
func DFSBacktracking(n int) [][]int
```

## 应用场景

### 1. 图论问题
- **连通性检测**：判断图中两个节点是否连通
- **路径查找**：寻找从起点到终点的路径
- **环检测**：检测图中是否存在环
- **拓扑排序**：对有向无环图进行排序

### 2. 树相关问题
- **树遍历**：前序、中序、后序遍历
- **树的高度**：计算树的最大深度
- **路径和**：查找从根到叶子的路径和
- **最近公共祖先**：查找两个节点的LCA

### 3. 回溯问题
- **全排列**：生成所有可能的排列
- **组合问题**：从n个元素中选择k个
- **子集问题**：生成所有可能的子集
- **N皇后**：在N×N棋盘上放置N个皇后

### 4. 游戏和AI
- **迷宫求解**：寻找迷宫的出口路径
- **数独求解**：填充数独格子
- **八数码问题**：移动数字块到目标状态
- **状态空间搜索**：探索问题的所有可能状态

## 算法优势

### 1. 内存效率
- **空间复杂度低**：只需要存储当前路径上的节点
- **适合深度搜索**：在深度较大的问题中表现良好
- **避免重复访问**：通过访问标记避免循环

### 2. 实现简单
- **递归实现**：代码简洁，逻辑清晰
- **易于理解**：符合人类思维模式
- **调试方便**：递归调用栈便于调试

### 3. 适用性广
- **通用性强**：适用于各种图结构
- **变种丰富**：可以适应不同的问题需求
- **扩展性好**：容易添加约束条件和优化

## 算法劣势

### 1. 可能陷入深度
- **无限深度**：在有向图中可能陷入无限循环
- **局部最优**：可能找到局部最优解而非全局最优
- **路径长度**：找到的路径可能不是最短路径

### 2. 性能问题
- **重复计算**：可能重复访问相同的子问题
- **指数复杂度**：在某些问题中时间复杂度很高
- **栈溢出**：递归深度过大时可能导致栈溢出

## 优化技巧

### 1. 剪枝优化
```go
// 提前返回
if condition {
    return
}

// 边界检查
if node == nil || visited[node] {
    return
}

// 约束条件
if !isValid(path) {
    return
}
```

### 2. 记忆化搜索
```go
// 使用map缓存结果
memo := make(map[string]bool)

// 检查缓存
if result, exists := memo[key]; exists {
    return result
}
```

### 3. 迭代优化
```go
// 使用显式栈避免递归
stack := list.New()
stack.PushBack(start)

for stack.Len() > 0 {
    node := stack.Remove(stack.Back()).(int)
    // 处理节点
}
```

## 代码示例

### 基本DFS实现
```go
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
```

### 回溯框架
```go
func DFSBacktracking(n int) [][]int {
    var result [][]int
    
    var backtrack func(path []int, used []bool)
    backtrack = func(path []int, used []bool) {
        if len(path) == n {
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
```

## 总结

深度优先搜索是一种强大而灵活的算法，特别适合需要探索所有可能性的问题。它的递归实现简洁优雅，迭代实现高效实用。通过合理的剪枝和优化，DFS可以解决从简单的图遍历到复杂的组合优化等各种问题。

DFS的核心思想是"一条路走到底"，这种策略在需要深度探索的问题中表现出色，是算法工具箱中不可或缺的重要工具。 