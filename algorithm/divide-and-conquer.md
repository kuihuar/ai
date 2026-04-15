# 分治算法 (Divide and Conquer)

## 📖 概述

分治算法是一种重要的算法设计策略，它将一个复杂的问题分解为若干个相同或相似的子问题，递归地解决这些子问题，然后将子问题的解合并得到原问题的解。

## 🎯 核心思想

### 分治三步曲
1. **分解 (Divide)** - 将原问题分解为若干个子问题
2. **解决 (Conquer)** - 递归地解决各个子问题
3. **合并 (Combine)** - 将子问题的解合并为原问题的解

### 递归基础
```go
// 递归的基本结构
func recursiveFunction(n int) int {
    // 1. 基本情况 (Base Case) - 递归的终止条件
    if n <= 1 {
        return 1
    }
    
    // 2. 递归情况 (Recursive Case) - 问题分解
    subResult := recursiveFunction(n - 1)
    
    // 3. 合并结果 (Combine) - 将子问题结果合并
    return n * subResult
}
```

## 🛠️ 经典分治算法

### 1. 归并排序 (Merge Sort)

#### 算法思想
将数组分成两半，分别排序，然后合并两个有序数组。

#### Go 实现
```go
package main

import "fmt"

// MergeSort 归并排序
func MergeSort(arr []int) []int {
    n := len(arr)
    
    // 基本情况：数组长度为0或1时直接返回
    if n <= 1 {
        return arr
    }
    
    // 分解：将数组分成两半
    mid := n / 2
    left := MergeSort(arr[:mid])
    right := MergeSort(arr[mid:])
    
    // 合并：合并两个有序数组
    return merge(left, right)
}

// merge 合并两个有序数组
func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    // 比较两个数组的元素，选择较小的放入结果
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    // 将剩余元素添加到结果
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    
    return result
}

func main() {
    arr := []int{64, 34, 25, 12, 22, 11, 90}
    fmt.Println("原始数组:", arr)
    
    sorted := MergeSort(arr)
    fmt.Println("排序后:", sorted)
}
```

#### 复杂度分析
- **时间复杂度**: O(n log n)
- **空间复杂度**: O(n)
- **稳定性**: 稳定排序

### 2. 快速排序 (Quick Sort)

#### 算法思想
选择一个基准元素，将数组分为小于基准和大于基准的两部分，递归排序。

#### Go 实现
```go
// QuickSort 快速排序
func QuickSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    // 选择基准元素（这里选择第一个元素）
    pivot := arr[0]
    
    // 分解：将数组分为小于基准和大于基准的两部分
    var left, right []int
    for i := 1; i < len(arr); i++ {
        if arr[i] <= pivot {
            left = append(left, arr[i])
        } else {
            right = append(right, arr[i])
        }
    }
    
    // 递归排序左右两部分
    left = QuickSort(left)
    right = QuickSort(right)
    
    // 合并：左部分 + 基准 + 右部分
    result := append(left, pivot)
    result = append(result, right...)
    
    return result
}

// QuickSortInPlace 原地快速排序（优化版本）
func QuickSortInPlace(arr []int, low, high int) {
    if low < high {
        // 分区并获取基准位置
        pivotIndex := partition(arr, low, high)
        
        // 递归排序基准左右两部分
        QuickSortInPlace(arr, low, pivotIndex-1)
        QuickSortInPlace(arr, pivotIndex+1, high)
    }
}

// partition 分区函数
func partition(arr []int, low, high int) int {
    pivot := arr[high]
    i := low - 1
    
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}
```

### 3. 二分搜索 (Binary Search)

#### 算法思想
在有序数组中，通过比较中间元素来缩小搜索范围。

#### Go 实现
```go
// BinarySearch 二分搜索（递归版本）
func BinarySearch(arr []int, target, left, right int) int {
    if left > right {
        return -1 // 未找到
    }
    
    mid := left + (right-left)/2
    
    if arr[mid] == target {
        return mid
    } else if arr[mid] > target {
        return BinarySearch(arr, target, left, mid-1)
    } else {
        return BinarySearch(arr, target, mid+1, right)
    }
}

// BinarySearchIterative 二分搜索（迭代版本）
func BinarySearchIterative(arr []int, target int) int {
    left, right := 0, len(arr)-1
    
    for left <= right {
        mid := left + (right-left)/2
        
        if arr[mid] == target {
            return mid
        } else if arr[mid] > target {
            right = mid - 1
        } else {
            left = mid + 1
        }
    }
    
    return -1
}

func leftBound(nums []int, target int)int{
    left, right := 0, len(nums) -1
    for left <= right{
        mid := left + (right - left) / 2
        if nums[mid] >= target {
            right = mid -1
        }else {
            left = mid +1
        }
    }
    if left >= len(nums) || nums[left] != target {
        return -1
    }
    return left
}

func rightBound(nums []int, target int) {

    left, right := 0,len(nums) -1

    for left <=right {
        mid := left + (right - left) /2
        if nums[mid] <=target {
            left = mid +1
        }else{
            right = mid -1
        }
    }
    if right <0 || nums[mid]!= target {
        return -1
    }
    return right
}


func firtGreaterOrEqual(nums []int, target int) {
    left, right:=0, len(nums)-1

    for left <=right {
        mid := left+(right-left) /2

        if nums[mid] >=target {
            right = mid -1
        }else{
            left =mid +1
        }
    }
    return left
}


func lasetLessEqual(nums []int, target int) int {
    left, right :=0, len(nums)-1


    for left <= right {

        mid := left + (right - left) /2

        if nums[mid] <=target {
            left = mid +1
        }else{
            right = mid -1
        }
    }
    return right
}
```

### 4. 大整数乘法 (Karatsuba Algorithm)

#### 算法思想
将大整数乘法分解为更小的乘法问题，减少乘法次数。

#### Go 实现
```go
// KaratsubaMultiply Karatsuba大整数乘法
func KaratsubaMultiply(x, y int) int {
    // 基本情况：如果数字很小，直接相乘
    if x < 10 || y < 10 {
        return x * y
    }
    
    // 计算数字的位数
    n := max(getDigits(x), getDigits(y))
    m := n / 2
    
    // 分解数字
    a := x / pow(10, m)
    b := x % pow(10, m)
    c := y / pow(10, m)
    d := y % pow(10, m)
    
    // 递归计算三个乘法
    ac := KaratsubaMultiply(a, c)
    bd := KaratsubaMultiply(b, d)
    ad_plus_bc := KaratsubaMultiply(a+b, c+d) - ac - bd
    
    // 合并结果
    return ac*pow(10, 2*m) + ad_plus_bc*pow(10, m) + bd
}

// 辅助函数
func getDigits(n int) int {
    if n == 0 {
        return 1
    }
    count := 0
    for n != 0 {
        n /= 10
        count++
    }
    return count
}

func pow(base, exp int) int {
    result := 1
    for i := 0; i < exp; i++ {
        result *= base
    }
    return result
}
```

## 🧮 主定理 (Master Theorem)

### 主定理公式
对于递归关系：T(n) = aT(n/b) + f(n)

其中：
- a ≥ 1：子问题数量
- b > 1：问题规模缩小因子
- f(n)：分解和合并的代价

### 三种情况
1. **情况1**: 如果 f(n) = O(n^(log_b(a) - ε))，则 T(n) = Θ(n^(log_b(a)))
2. **情况2**: 如果 f(n) = Θ(n^(log_b(a)) * log^k(n))，则 T(n) = Θ(n^(log_b(a)) * log^(k+1)(n))
3. **情况3**: 如果 f(n) = Ω(n^(log_b(a) + ε))，则 T(n) = Θ(f(n))

### 应用示例
```go
// 归并排序的复杂度分析
// T(n) = 2T(n/2) + O(n)
// a = 2, b = 2, f(n) = O(n)
// log_b(a) = log_2(2) = 1
// f(n) = O(n) = O(n^1) = O(n^(log_b(a)))
// 属于情况2，k = 0
// 因此 T(n) = Θ(n * log n)
```

## 🚀 递归优化技巧

### 1. 记忆化递归 (Memoization)
```go
// 斐波那契数列的记忆化递归
var memo = make(map[int]int)

func FibonacciMemo(n int) int {
    if n <= 1 {
        return n
    }
    
    // 检查是否已经计算过
    if val, exists := memo[n]; exists {
        return val
    }
    
    // 计算并存储结果
    result := FibonacciMemo(n-1) + FibonacciMemo(n-2)
    memo[n] = result
    
    return result
}
```

### 2. 尾递归优化
```go
// 尾递归版本的阶乘计算
func FactorialTailRec(n, acc int) int {
    if n <= 1 {
        return acc
    }
    return FactorialTailRec(n-1, n*acc)
}

// 调用方式
func Factorial(n int) int {
    return FactorialTailRec(n, 1)
}
```

### 3. 自底向上 (Bottom-Up)
```go
// 斐波那契数列的自底向上解法
func FibonacciBottomUp(n int) int {
    if n <= 1 {
        return n
    }
    
    dp := make([]int, n+1)
    dp[0], dp[1] = 0, 1
    
    for i := 2; i <= n; i++ {
        dp[i] = dp[i-1] + dp[i-2]
    }
    
    return dp[n]
}
```

## 🎯 经典分治问题

### 1. 最大子数组和 (Maximum Subarray Sum)
```go
// MaxSubArraySum 最大子数组和（分治解法）
func MaxSubArraySum(arr []int) int {
    return maxSubArrayHelper(arr, 0, len(arr)-1)
}

func maxSubArrayHelper(arr []int, left, right int) int {
    if left == right {
        return arr[left]
    }
    
    mid := left + (right-left)/2
    
    // 递归求解左右两部分
    leftMax := maxSubArrayHelper(arr, left, mid)
    rightMax := maxSubArrayHelper(arr, mid+1, right)
    
    // 求解跨越中点的最大子数组和
    crossMax := maxCrossingSum(arr, left, mid, right)
    
    return max3(leftMax, rightMax, crossMax)
}

func maxCrossingSum(arr []int, left, mid, right int) int {
    // 向左扩展
    leftSum := 0
    leftMax := arr[mid]
    for i := mid; i >= left; i-- {
        leftSum += arr[i]
        if leftSum > leftMax {
            leftMax = leftSum
        }
    }
    
    // 向右扩展
    rightSum := 0
    rightMax := arr[mid+1]
    for i := mid + 1; i <= right; i++ {
        rightSum += arr[i]
        if rightSum > rightMax {
            rightMax = rightSum
        }
    }
    
    return leftMax + rightMax
}
```

### 2. 最近点对问题 (Closest Pair of Points)
```go
type Point struct {
    x, y float64
}

// ClosestPair 最近点对问题
func ClosestPair(points []Point) (Point, Point, float64) {
    // 按x坐标排序
    sort.Slice(points, func(i, j int) bool {
        return points[i].x < points[j].x
    })
    
    return closestPairHelper(points)
}

func closestPairHelper(points []Point) (Point, Point, float64) {
    n := len(points)
    
    if n <= 3 {
        return bruteForceClosest(points)
    }
    
    mid := n / 2
    midX := points[mid].x
    
    // 递归求解左右两部分
    leftP1, leftP2, leftDist := closestPairHelper(points[:mid])
    rightP1, rightP2, rightDist := closestPairHelper(points[mid:])
    
    // 取较小的距离
    var minDist float64
    var p1, p2 Point
    if leftDist < rightDist {
        minDist = leftDist
        p1, p2 = leftP1, leftP2
    } else {
        minDist = rightDist
        p1, p2 = rightP1, rightP2
    }
    
    // 检查跨越中线的点对
    stripP1, stripP2, stripDist := closestInStrip(points, midX, minDist)
    
    if stripDist < minDist {
        return stripP1, stripP2, stripDist
    }
    
    return p1, p2, minDist
}

func distance(p1, p2 Point) float64 {
    dx := p1.x - p2.x
    dy := p1.y - p2.y
    return math.Sqrt(dx*dx + dy*dy)
}
```

## 🧪 测试和练习

### 1. 基础测试
```go
func TestMergeSort(t *testing.T) {
    testCases := []struct {
        input    []int
        expected []int
    }{
        {[]int{64, 34, 25, 12, 22, 11, 90}, []int{11, 12, 22, 25, 34, 64, 90}},
        {[]int{1}, []int{1}},
        {[]int{}, []int{}},
        {[]int{3, 3, 3}, []int{3, 3, 3}},
    }
    
    for _, tc := range testCases {
        result := MergeSort(tc.input)
        if !reflect.DeepEqual(result, tc.expected) {
            t.Errorf("MergeSort(%v) = %v, 期望 %v", tc.input, result, tc.expected)
        }
    }
}
```

### 2. 性能测试
```go
func BenchmarkMergeSort(b *testing.B) {
    // 生成随机数组
    arr := make([]int, 1000)
    for i := range arr {
        arr[i] = rand.Intn(10000)
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        arrCopy := make([]int, len(arr))
        copy(arrCopy, arr)
        MergeSort(arrCopy)
    }
}
```

## 📚 学习建议

### 1. 学习步骤
```markdown
1. 🧠 理解分治思想：分解、解决、合并
2. 📝 掌握递归基础：基本情况、递归情况
3. 🔢 学习主定理：分析递归复杂度
4. 💻 实现经典算法：归并排序、快速排序
5. 🎯 解决实际问题：最大子数组、最近点对
6. ⚡ 优化技巧：记忆化、尾递归、自底向上
```

### 2. 练习题目
```markdown
基础练习：
- 实现归并排序
- 实现快速排序
- 实现二分搜索
- 计算斐波那契数列

进阶练习：
- 最大子数组和问题
- 最近点对问题
- Karatsuba大整数乘法
- Strassen矩阵乘法

高级练习：
- 分治优化动态规划
- 分治解决几何问题
- 并行分治算法
```

### 3. 常见错误
```go
// 错误1：忘记基本情况
func wrongRecursion(n int) int {
    return n + wrongRecursion(n-1) // 缺少基本情况，会无限递归
}

// 错误2：递归深度过深
func deepRecursion(n int) int {
    if n <= 1 {
        return 1
    }
    return deepRecursion(n-1) + deepRecursion(n-2) // 指数级复杂度
}

// 错误3：没有正确合并结果
func wrongMerge(left, right []int) []int {
    // 没有正确合并两个有序数组
    return append(left, right...) // 这样会保持原顺序，不是合并排序
}
```

## 🎯 实际应用

### 1. 数据库查询优化
```go
// 分治在数据库查询中的应用
func parallelQuery(data []Record, query Query) []Record {
    if len(data) < 1000 {
        return sequentialQuery(data, query)
    }
    
    mid := len(data) / 2
    
    // 并行处理两部分
    var leftResults, rightResults []Record
    var wg sync.WaitGroup
    
    wg.Add(2)
    go func() {
        defer wg.Done()
        leftResults = parallelQuery(data[:mid], query)
    }()
    
    go func() {
        defer wg.Done()
        rightResults = parallelQuery(data[mid:], query)
    }()
    
    wg.Wait()
    
    // 合并结果
    return mergeQueryResults(leftResults, rightResults)
}
```

### 2. 图像处理
```go
// 分治在图像处理中的应用
func divideAndConquerImageProcessing(image [][]Pixel) [][]Pixel {
    height := len(image)
    width := len(image[0])
    
    if height <= 64 && width <= 64 {
        return processSmallImage(image)
    }
    
    // 将图像分成四个象限
    midH := height / 2
    midW := width / 2
    
    topLeft := divideAndConquerImageProcessing(image[:midH][:midW])
    topRight := divideAndConquerImageProcessing(image[:midH][midW:])
    bottomLeft := divideAndConquerImageProcessing(image[midH:][:midW])
    bottomRight := divideAndConquerImageProcessing(image[midH:][midW:])
    
    // 合并四个象限的结果
    return mergeImageQuadrants(topLeft, topRight, bottomLeft, bottomRight)
}
```

---

**分治算法是算法设计中的核心思想，掌握它对于理解更复杂的算法非常重要！** 🎉

通过系统学习分治算法，您将能够：
- 理解递归的本质和优化技巧
- 掌握经典的分治算法实现
- 学会分析算法的复杂度
- 应用分治思想解决实际问题 