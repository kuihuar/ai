# 背包问题 (Knapsack Problem)

## 问题概述

背包问题是组合优化中的一个经典问题，描述为：给定一组物品，每个物品都有自己的重量和价值，在限定的总重量内，选择物品使得总价值最大。这个问题是NP完全问题，在实际应用中有广泛的应用场景。

## 问题分类

### 1. 0-1背包问题
- **特点**：每个物品只能选择一次（放入或不放入）
- **约束**：总重量不超过背包容量
- **目标**：最大化总价值
- **复杂度**：NP完全问题

### 2. 完全背包问题
- **特点**：每个物品可以选择无限次
- **约束**：总重量不超过背包容量
- **目标**：最大化总价值
- **复杂度**：伪多项式时间

### 3. 多重背包问题
- **特点**：每个物品有数量限制
- **约束**：总重量不超过背包容量，每种物品不超过指定数量
- **目标**：最大化总价值
- **复杂度**：NP完全问题

### 4. 分组背包问题
- **特点**：物品分组，每组最多选择一个物品
- **约束**：总重量不超过背包容量，每组最多选一个
- **目标**：最大化总价值
- **复杂度**：NP完全问题

## 算法变种

### 1. 动态规划解法
```go
// 0-1背包问题动态规划
func Knapsack01DP(weights []int, values []int, capacity int) int

// 完全背包问题动态规划
func KnapsackUnboundedDP(weights []int, values []int, capacity int) int

// 多重背包问题动态规划
func KnapsackMultipleDP(weights []int, values []int, counts []int, capacity int) int

// 分组背包问题动态规划
func KnapsackGroupDP(groups [][]Item, capacity int) int
```

### 2. 贪心算法
```go
// 按价值密度贪心
func KnapsackGreedyByDensity(weights []int, values []int, capacity int) int

// 按价值贪心
func KnapsackGreedyByValue(weights []int, values []int, capacity int) int

// 按重量贪心
func KnapsackGreedyByWeight(weights []int, values []int, capacity int) int
```

### 3. 回溯算法
```go
// 0-1背包回溯
func Knapsack01Backtrack(weights []int, values []int, capacity int) int

// 带剪枝的回溯
func KnapsackBacktrackWithPruning(weights []int, values []int, capacity int) int

// 分支限界
func KnapsackBranchAndBound(weights []int, values []int, capacity int) int
```

### 4. 近似算法
```go
// 贪心近似
func KnapsackGreedyApproximation(weights []int, values []int, capacity int) int

// 动态规划近似
func KnapsackDPApproximation(weights []int, values []int, capacity int, epsilon float64) int

// 随机化算法
func KnapsackRandomized(weights []int, values []int, capacity int, iterations int) int
```

## 应用场景

### 1. 资源分配
- **投资组合**：在有限资金下选择最优投资项目
- **时间管理**：在有限时间内安排最有价值的任务
- **空间利用**：在有限空间内放置最有价值的物品
- **能源分配**：在有限能源下分配最优用途

### 2. 生产制造
- **生产计划**：在有限资源下安排最优生产计划
- **库存管理**：在有限仓库空间下存储最有价值商品
- **设备选择**：在有限预算下选择最优设备组合
- **原材料采购**：在有限预算下采购最优原材料

### 3. 计算机科学
- **内存管理**：在有限内存下加载最有价值的程序
- **缓存优化**：在有限缓存空间下存储最有价值数据
- **网络带宽**：在有限带宽下传输最有价值信息
- **CPU调度**：在有限CPU时间下执行最有价值任务

### 4. 游戏开发
- **装备选择**：在有限负重下选择最优装备组合
- **技能点分配**：在有限技能点下分配最优技能
- **道具携带**：在有限道具栏下携带最有价值道具
- **资源收集**：在有限时间下收集最有价值资源

## 算法特点

### 1. 动态规划优势
- **最优性保证**：保证找到全局最优解
- **完备性**：能够处理所有可能的输入
- **理论基础**：有坚实的数学理论基础
- **可扩展性**：容易添加约束条件

### 2. 贪心算法优势
- **效率高**：时间复杂度通常为O(n log n)
- **实现简单**：代码简洁易懂
- **内存友好**：空间复杂度低
- **实时性好**：适合实时应用

### 3. 回溯算法优势
- **灵活性**：容易添加各种约束条件
- **剪枝优化**：通过剪枝大幅提高效率
- **解空间探索**：能够探索完整的解空间
- **启发式指导**：可以使用启发式函数指导搜索

## 算法劣势

### 1. 动态规划劣势
- **空间复杂度高**：需要存储完整的DP表
- **时间复杂度高**：对于大容量问题效率较低
- **内存消耗大**：在内存受限环境中可能不可行
- **初始化复杂**：需要正确初始化DP表

### 2. 贪心算法劣势
- **不保证最优**：通常只能得到局部最优解
- **问题依赖**：贪心策略需要针对特定问题设计
- **质量不稳定**：解的质量可能波动很大
- **理论基础弱**：缺乏统一的理论基础

### 3. 回溯算法劣势
- **指数复杂度**：在最坏情况下是指数级复杂度
- **剪枝依赖**：效率高度依赖剪枝策略
- **参数敏感**：对参数设置比较敏感
- **调试困难**：递归调用栈调试复杂

## 优化技巧

### 1. 动态规划优化
```go
// 空间优化 - 滚动数组
func Knapsack01Optimized(weights []int, values []int, capacity int) int {
    dp := make([]int, capacity+1)
    
    for i := 0; i < len(weights); i++ {
        for j := capacity; j >= weights[i]; j-- {
            dp[j] = max(dp[j], dp[j-weights[i]]+values[i])
        }
    }
    
    return dp[capacity]
}

// 二进制优化 - 多重背包
func KnapsackMultipleOptimized(weights []int, values []int, counts []int, capacity int) int {
    dp := make([]int, capacity+1)
    
    for i := 0; i < len(weights); i++ {
        // 二进制分解
        for k := 1; counts[i] > 0; k *= 2 {
            if counts[i] < k {
                k = counts[i]
            }
            counts[i] -= k
            
            for j := capacity; j >= k*weights[i]; j-- {
                dp[j] = max(dp[j], dp[j-k*weights[i]]+k*values[i])
            }
        }
    }
    
    return dp[capacity]
}
```

### 2. 贪心优化
```go
// 价值密度排序
func KnapsackGreedyOptimized(weights []int, values []int, capacity int) int {
    n := len(weights)
    items := make([]Item, n)
    
    for i := 0; i < n; i++ {
        items[i] = Item{
            weight: weights[i],
            value:  values[i],
            density: float64(values[i]) / float64(weights[i]),
        }
    }
    
    // 按价值密度降序排序
    sort.Slice(items, func(i, j int) bool {
        return items[i].density > items[j].density
    })
    
    totalValue := 0
    remainingCapacity := capacity
    
    for _, item := range items {
        if item.weight <= remainingCapacity {
            totalValue += item.value
            remainingCapacity -= item.weight
        }
    }
    
    return totalValue
}
```

### 3. 回溯优化
```go
// 带剪枝的回溯
func KnapsackBacktrackOptimized(weights []int, values []int, capacity int) int {
    n := len(weights)
    bestValue := 0
    
    // 计算价值密度
    densities := make([]float64, n)
    for i := 0; i < n; i++ {
        densities[i] = float64(values[i]) / float64(weights[i])
    }
    
    var backtrack func(index int, currentWeight int, currentValue int, upperBound float64)
    backtrack = func(index int, currentWeight int, currentValue int, upperBound float64) {
        // 剪枝：如果上界小于当前最优解，则剪枝
        if float64(currentValue)+upperBound <= float64(bestValue) {
            return
        }
        
        if index == n {
            if currentValue > bestValue {
                bestValue = currentValue
            }
            return
        }
        
        // 不选择当前物品
        backtrack(index+1, currentWeight, currentValue, upperBound-densities[index])
        
        // 选择当前物品
        if currentWeight+weights[index] <= capacity {
            backtrack(index+1, currentWeight+weights[index], currentValue+values[index], upperBound-densities[index])
        }
    }
    
    // 计算初始上界
    upperBound := 0.0
    for i := 0; i < n; i++ {
        upperBound += densities[i]
    }
    
    backtrack(0, 0, 0, upperBound)
    return bestValue
}
```

## 代码示例

### 0-1背包问题动态规划
```go
type Item struct {
    weight int
    value  int
    density float64
}

func Knapsack01DP(weights []int, values []int, capacity int) int {
    // 获取物品数量
    n := len(weights)
    
    // 创建二维DP表，dp[i][j]表示前i个物品在容量j下的最大价值
    // dp[i][j] = 前i个物品中，在容量j的限制下能获得的最大价值
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, capacity+1)
    }
    
    // 外层循环：遍历每个物品（从1开始，因为dp[0][j] = 0）
    for i := 1; i <= n; i++ {
        // 内层循环：遍历所有可能的容量（从0到capacity）
        for j := 0; j <= capacity; j++ {
            // 如果当前物品的重量小于等于当前容量
            if weights[i-1] <= j {
                // 状态转移方程：
                // dp[i][j] = max(不选择物品i, 选择物品i)
                // 不选择物品i: dp[i-1][j]
                // 选择物品i: dp[i-1][j-weights[i-1]] + values[i-1]
                dp[i][j] = max(dp[i-1][j], dp[i-1][j-weights[i-1]]+values[i-1])
            } else {
                // 如果当前物品重量超过容量，只能不选择
                dp[i][j] = dp[i-1][j]
            }
        }
    }
    
    // 返回最终结果：前n个物品在容量capacity下的最大价值
    return dp[n][capacity]
}

// 空间优化版本
func Knapsack01Optimized(weights []int, values []int, capacity int) int {
    dp := make([]int, capacity+1)
    
    for i := 0; i < len(weights); i++ {
        for j := capacity; j >= weights[i]; j-- {
            dp[j] = max(dp[j], dp[j-weights[i]]+values[i])
        }
    }
    
    return dp[capacity]
}
```

### 完全背包问题
```go
func KnapsackUnboundedDP(weights []int, values []int, capacity int) int {
    dp := make([]int, capacity+1)
    
    for i := 0; i < len(weights); i++ {
        for j := weights[i]; j <= capacity; j++ {
            dp[j] = max(dp[j], dp[j-weights[i]]+values[i])
        }
    }
    
    return dp[capacity]
}
```

### 多重背包问题
```go
func KnapsackMultipleDP(weights []int, values []int, counts []int, capacity int) int {
    dp := make([]int, capacity+1)
    
    for i := 0; i < len(weights); i++ {
        for j := capacity; j >= weights[i]; j-- {
            for k := 1; k <= counts[i] && k*weights[i] <= j; k++ {
                dp[j] = max(dp[j], dp[j-k*weights[i]]+k*values[i])
            }
        }
    }
    
    return dp[capacity]
}
```

### 贪心算法实现
```go
func KnapsackGreedyByDensity(weights []int, values []int, capacity int) int {
    n := len(weights)
    items := make([]Item, n)
    
    for i := 0; i < n; i++ {
        items[i] = Item{
            weight:  weights[i],
            value:   values[i],
            density: float64(values[i]) / float64(weights[i]),
        }
    }
    
    // 按价值密度降序排序
    sort.Slice(items, func(i, j int) bool {
        return items[i].density > items[j].density
    })
    
    totalValue := 0
    remainingCapacity := capacity
    
    for _, item := range items {
        if item.weight <= remainingCapacity {
            totalValue += item.value
            remainingCapacity -= item.weight
        }
    }
    
    return totalValue
}
```

## 问题变种

### 1. 子集和问题
- **描述**：给定一个数组，判断是否存在子集的和等于目标值
- **解法**：可以转化为0-1背包问题
- **应用**：资源分配、目标达成

### 2. 分割等和子集
- **描述**：将数组分成两个子集，使两个子集的和相等
- **解法**：转化为背包容量为总和一半的背包问题
- **应用**：负载均衡、资源分配

### 3. 目标和问题
- **描述**：在数组中添加+或-号，使结果等于目标值
- **解法**：转化为背包问题
- **应用**：表达式求值、符号分配

### 4. 硬币兑换问题
- **描述**：用最少数量的硬币凑出目标金额
- **解法**：完全背包问题的变种
- **应用**：货币兑换、零钱找零

## 性能比较

| 算法 | 时间复杂度 | 空间复杂度 | 最优性 | 适用场景 |
|------|------------|------------|--------|----------|
| 动态规划 | O(nW) | O(nW) | 保证最优 | 小规模问题 |
| 贪心算法 | O(n log n) | O(1) | 不保证最优 | 实时应用 |
| 回溯算法 | O(2^n) | O(n) | 保证最优 | 精确解 |
| 分支限界 | O(2^n) | O(n) | 保证最优 | 中等规模 |
| 近似算法 | O(n log n) | O(n) | 近似最优 | 大规模问题 |

## 总结

背包问题是组合优化中的经典问题，具有重要的理论价值和实际应用意义。不同的解法各有优劣：

- **动态规划**：保证最优解，适合小规模问题
- **贪心算法**：效率高，适合实时应用
- **回溯算法**：灵活性好，适合精确解
- **近似算法**：平衡效率和精度，适合大规模问题

选择合适的算法需要根据具体问题的规模、精度要求和实时性需求来决定。在实际应用中，往往需要结合多种算法的优点，设计出适合特定场景的解决方案。

## 0-1背包问题详细示例

### 输入数据
```go
weights = [2, 1, 3, 2]  // 物品重量
values  = [12, 10, 20, 15]  // 物品价值  
capacity = 5  // 背包容量
```

### DP表构建过程

| 物品\容量 | 0 | 1 | 2 | 3 | 4 | 5 |
|-----------|---|---|---|---|---|---|
| 0个物品   | 0 | 0 | 0 | 0 | 0 | 0 |
| 物品1(w=2,v=12) | 0 | 0 | 12 | 12 | 12 | 12 |
| 物品2(w=1,v=10) | 0 | 10 | 12 | 22 | 22 | 22 |
| 物品3(w=3,v=20) | 0 | 10 | 12 | 22 | 22 | 32 |
| 物品4(w=2,v=15) | 0 | 10 | 15 | 25 | 37 | 37 |

### 详细计算过程

**第1行（物品1，重量2，价值12）：**
- 容量0: dp[1][0] = dp[0][0] = 0
- 容量1: dp[1][1] = dp[0][1] = 0 (重量2 > 容量1)
- 容量2: dp[1][2] = max(dp[0][2], dp[0][0]+12) = max(0, 12) = 12
- 容量3: dp[1][3] = max(dp[0][3], dp[0][1]+12) = max(0, 12) = 12
- 容量4: dp[1][4] = max(dp[0][4], dp[0][2]+12) = max(0, 12) = 12
- 容量5: dp[1][5] = max(dp[0][5], dp[0][3]+12) = max(0, 12) = 12

**第2行（物品2，重量1，价值10）：**
- 容量0: dp[2][0] = dp[1][0] = 0
- 容量1: dp[2][1] = max(dp[1][1], dp[1][0]+10) = max(0, 10) = 10
- 容量2: dp[2][2] = max(dp[1][2], dp[1][1]+10) = max(12, 10) = 12
- 容量3: dp[2][3] = max(dp[1][3], dp[1][2]+10) = max(12, 22) = 22
- 容量4: dp[2][4] = max(dp[1][4], dp[1][3]+10) = max(12, 22) = 22
- 容量5: dp[2][5] = max(dp[1][5], dp[1][4]+10) = max(12, 22) = 22

**第3行（物品3，重量3，价值20）：**
- 容量0: dp[3][0] = dp[2][0] = 0
- 容量1: dp[3][1] = dp[2][1] = 10 (重量3 > 容量1)
- 容量2: dp[3][2] = dp[2][2] = 12 (重量3 > 容量2)
- 容量3: dp[3][3] = max(dp[2][3], dp[2][0]+20) = max(22, 20) = 22
- 容量4: dp[3][4] = max(dp[2][4], dp[2][1]+20) = max(22, 30) = 30
- 容量5: dp[3][5] = max(dp[2][5], dp[2][2]+20) = max(22, 32) = 32

**第4行（物品4，重量2，价值15）：**
- 容量0: dp[4][0] = dp[3][0] = 0
- 容量1: dp[4][1] = dp[3][1] = 10 (重量2 > 容量1)
- 容量2: dp[4][2] = max(dp[3][2], dp[3][0]+15) = max(12, 15) = 15
- 容量3: dp[4][3] = max(dp[3][3], dp[3][1]+15) = max(22, 25) = 25
- 容量4: dp[4][4] = max(dp[3][4], dp[3][2]+15) = max(30, 27) = 30
- 容量5: dp[4][5] = max(dp[3][5], dp[3][3]+15) = max(32, 37) = 37

### 最终结果
- **最大价值**：37
- **最优选择**：物品2(重量1,价值10) + 物品3(重量3,价值20) + 物品4(重量2,价值15) = 总重量6 > 容量5 ❌
- **实际最优**：物品1(重量2,价值12) + 物品3(重量3,价值20) = 总重量5, 总价值32

### 算法复杂度
- **时间复杂度**：O(n × capacity) = O(4 × 5) = O(20)
- **空间复杂度**：O(n × capacity) = O(4 × 5) = O(20)

### 最优解回溯
通过DP表可以回溯出最优解的选择：
```go
func GetOptimalItems(weights []int, values []int, capacity int) []int {
    n := len(weights)
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, capacity+1)
    }
    
    // 构建DP表（同上）
    for i := 1; i <= n; i++ {
        for j := 0; j <= capacity; j++ {
            if weights[i-1] <= j {
                dp[i][j] = max(dp[i-1][j], dp[i-1][j-weights[i-1]]+values[i-1])
            } else {
                dp[i][j] = dp[i-1][j]
            }
        }
    }
    
    // 回溯最优解
    selected := []int{}
    i, j := n, capacity
    
    for i > 0 && j > 0 {
        if dp[i][j] != dp[i-1][j] {
            selected = append(selected, i-1) // 选择了物品i-1
            j -= weights[i-1]
        }
        i--
    }
    
    return selected
}
```

对于示例数据，回溯结果：选择物品1和物品3，总价值32。 