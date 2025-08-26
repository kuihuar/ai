# 搜索算法详解

## 概述

搜索算法是计算机科学中最基础和最重要的算法之一，用于在数据集合中查找特定元素。根据数据结构和应用场景的不同，有多种搜索算法可供选择。

## 线性搜索 (Linear Search)

### 基本概念
线性搜索是最简单直观的搜索算法，通过逐个检查数组中的每个元素来查找目标值。

### 算法特点
- **时间复杂度**: O(n)
- **空间复杂度**: O(1)
- **适用场景**: 无序数组、小规模数据
- **稳定性**: 稳定，总能找到目标值（如果存在）

### 实现变种

#### 1. 基本线性搜索
```go
func LinearSearch(arr []int, target int) int {
    for i := 0; i < len(arr); i++ {
        if arr[i] == target {
            return i
        }
    }
    return -1
}
```

#### 2. 搜索所有匹配项
```go
func LinearSearchAll(arr []int, target int) []int {
    var result []int
    for i := 0; i < len(arr); i++ {
        if arr[i] == target {
            result = append(result, i)
        }
    }
    return result
}
```

#### 3. 哨兵线性搜索
**优势**: 减少边界检查，提高性能
```go
func LinearSearchWithSentinel(arr []int, target int) int {
    if len(arr) == 0 {
        return -1
    }
    
    last := arr[len(arr)-1]
    arr[len(arr)-1] = target  // 设置哨兵
    
    i := 0
    for arr[i] != target {
        i++
    }
    
    arr[len(arr)-1] = last  // 恢复原数组
    
    if i < len(arr)-1 || last == target {
        return i
    }
    return -1
}
```

### 线性搜索的优缺点

**优点**:
- 简单易懂，容易实现
- 适用于有序和无序数组
- 内存友好，空间复杂度低
- 稳定可靠

**缺点**:
- 效率较低，时间复杂度高
- 不适合大规模数据
- 无优化空间

## 二分搜索 (Binary Search)

### 基本概念
二分搜索是一种高效的搜索算法，通过将搜索范围减半来快速定位目标值。

### 算法特点
- **时间复杂度**: O(log n)
- **空间复杂度**: O(1) - 迭代版本，O(log n) - 递归版本
- **前提条件**: 数组必须是有序的
- **适用场景**: 大规模有序数据

### 核心思想
1. 比较目标值与中间元素
2. 如果相等，返回索引
3. 如果目标值小于中间元素，在左半部分搜索
4. 如果目标值大于中间元素，在右半部分搜索
5. 重复直到找到目标值或搜索范围为空

### 实现变种

#### 1. 基本二分搜索
```go
func BinarySearch(arr []int, target int) int {
    left := 0
    right := len(arr) - 1
    
    for left <= right {
        mid := left + (right-left)/2  // 避免整数溢出
        
        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return -1
}
```

#### 2. 递归二分搜索
```go
func BinarySearchRecursive(arr []int, target int) int {
    return binarySearchHelper(arr, target, 0, len(arr)-1)
}

func binarySearchHelper(arr []int, target, left, right int) int {
    if left > right {
        return -1
    }
    
    mid := left + (right-left)/2
    
    if arr[mid] == target {
        return mid
    } else if arr[mid] < target {
        return binarySearchHelper(arr, target, mid+1, right)
    } else {
        return binarySearchHelper(arr, target, left, mid-1)
    }
}
```

#### 3. 搜索第一个匹配项
```go
func BinarySearchFirst(arr []int, target int) int {
    left := 0
    right := len(arr) - 1
    result := -1
    
    for left <= right {
        mid := left + (right-left)/2
        
        if arr[mid] == target {
            result = mid
            right = mid - 1  // 继续向左搜索
        } else if arr[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    
    return result
}
```

### 二分搜索的优缺点

**优点**:
- 效率高，时间复杂度低
- 适合大规模数据
- 内存友好

**缺点**:
- 要求数据有序
- 实现相对复杂
- 不适合频繁插入/删除的场景

## 跳跃搜索 (Jump Search)

### 基本概念
跳跃搜索是线性搜索和二分搜索的折中方案，通过跳跃式前进来减少比较次数。

### 算法特点
- **时间复杂度**: O(√n)
- **空间复杂度**: O(1)
- **前提条件**: 数组必须是有序的
- **适用场景**: 有序数组，比线性搜索快，比二分搜索简单

### 核心思想
1. 设定跳跃步长（通常为√n）
2. 跳跃式前进，直到找到大于目标值的元素
3. 在上一跳和当前跳之间进行线性搜索

### 实现
```go
func JumpSearch(arr []int, target int) int {
    if len(arr) == 0 {
        return -1
    }
    
    step := int(float64(len(arr)) * 0.5)
    if step == 0 {
        step = 1
    }
    
    prev := 0
    for i := 0; i < len(arr); i += step {
        if arr[i] == target {
            return i
        }
        
        if arr[i] > target {
            // 在上一跳和当前跳之间进行线性搜索
            for j := prev; j < i && j < len(arr); j++ {
                if arr[j] == target {
                    return j
                }
            }
            return -1
        }
        
        prev = i
    }
    
    // 在最后一段进行线性搜索
    for j := prev; j < len(arr); j++ {
        if arr[j] == target {
            return j
        }
    }
    
    return -1
}
```

## 插值搜索 (Interpolation Search)

### 基本概念
插值搜索是二分搜索的改进版本，通过插值公式来估计目标值的位置。

### 算法特点
- **时间复杂度**: O(log log n) - 平均情况，O(n) - 最坏情况
- **空间复杂度**: O(1)
- **前提条件**: 数组必须是有序的，且元素分布均匀
- **适用场景**: 均匀分布的有序数据

### 核心思想
使用插值公式来估计目标值的位置：
```
pos = left + ((right-left) * (target-arr[left])) / (arr[right]-arr[left])
```

### 实现
```go
func InterpolationSearch(arr []int, target int) int {
    if len(arr) == 0 {
        return -1
    }
    
    left := 0
    right := len(arr) - 1
    
    for left <= right && target >= arr[left] && target <= arr[right] {
        if left == right {
            if arr[left] == target {
                return left
            }
            return -1
        }
        
        // 插值公式
        pos := left + ((right-left)*(target-arr[left]))/(arr[right]-arr[left])
        
        if arr[pos] == target {
            return pos
        } else if arr[pos] < target {
            left = pos + 1
        } else {
            right = pos - 1
        }
    }
    
    return -1
}
```

## 指数搜索 (Exponential Search)

### 基本概念
指数搜索用于在无界排序数组中搜索元素，通过指数增长来找到搜索范围。

### 算法特点
- **时间复杂度**: O(log n)
- **空间复杂度**: O(1)
- **前提条件**: 数组必须是有序的
- **适用场景**: 无界排序数组，流式数据

### 核心思想
1. 找到包含目标值的范围
2. 在该范围内进行二分搜索

### 实现
```go
func ExponentialSearch(arr []int, target int) int {
    if len(arr) == 0 {
        return -1
    }
    
    if arr[0] == target {
        return 0
    }
    
    // 找到范围
    i := 1
    for i < len(arr) && arr[i] <= target {
        i = i * 2
    }
    
    // 在找到的范围内进行二分搜索
    return binarySearchHelper(arr, target, i/2, min(i, len(arr)-1))
}
```

## 算法比较

| 算法 | 时间复杂度 | 空间复杂度 | 前提条件 | 适用场景 |
|------|------------|------------|----------|----------|
| 线性搜索 | O(n) | O(1) | 无 | 小规模数据，无序数组 |
| 二分搜索 | O(log n) | O(1) | 有序 | 大规模有序数据 |
| 跳跃搜索 | O(√n) | O(1) | 有序 | 中等规模有序数据 |
| 插值搜索 | O(log log n) | O(1) | 有序且均匀分布 | 均匀分布的有序数据 |
| 指数搜索 | O(log n) | O(1) | 有序 | 无界排序数组 |

## 选择指南

### 选择线性搜索的情况：
- 数据规模小（n < 100）
- 数组无序
- 需要找到所有匹配项
- 内存极度受限

### 选择二分搜索的情况：
- 数据规模大（n > 1000）
- 数组有序
- 需要频繁搜索
- 对性能要求高

### 选择跳跃搜索的情况：
- 中等规模数据
- 数组有序
- 需要比线性搜索快但实现简单

### 选择插值搜索的情况：
- 数据分布均匀
- 数组有序
- 对性能要求极高

### 选择指数搜索的情况：
- 无界排序数组
- 流式数据
- 不知道数组大小

## 实际应用

### 1. 数据库查询
- 索引查找使用二分搜索
- 全文搜索使用线性搜索

### 2. 文件系统
- 文件查找使用二分搜索
- 目录遍历使用线性搜索

### 3. 网络路由
- IP地址查找使用二分搜索
- 域名解析使用哈希搜索

### 4. 游戏开发
- 碰撞检测使用空间分割
- 路径查找使用A*算法

## 性能优化技巧

### 1. 缓存友好
- 使用局部性原理
- 减少内存访问

### 2. 分支预测
- 减少条件分支
- 使用哨兵技术

### 3. 并行化
- 分块处理
- 多线程搜索

### 4. 数据结构优化
- 使用跳表
- 使用B树

## 总结

搜索算法是算法设计的基础，选择合适的搜索算法对系统性能至关重要。在实际应用中，需要根据数据特征、规模大小、性能要求等因素综合考虑，选择最适合的算法。
