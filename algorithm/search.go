package algorithm

import "fmt"

// LinearSearch 线性搜索算法
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 适用场景: 无序数组、小规模数据、需要找到第一个匹配项
func LinearSearch(arr []int, target int) int {
	for i := 0; i < len(arr); i++ {
		if arr[i] == target {
			return i // 返回第一个匹配的索引
		}
	}
	return -1 // 未找到
}

// LinearSearchAll 线性搜索所有匹配项
// 时间复杂度: O(n)
// 空间复杂度: O(k) - k为匹配项数量
func LinearSearchAll(arr []int, target int) []int {
	var result []int
	for i := 0; i < len(arr); i++ {
		if arr[i] == target {
			result = append(result, i)
		}
	}
	return result
}

// LinearSearchOptimized 优化的线性搜索（提前退出）
// 时间复杂度: O(n) - 最坏情况
// 空间复杂度: O(1)
func LinearSearchOptimized(arr []int, target int) int {
	// 如果数组为空，直接返回-1
	if len(arr) == 0 {
		return -1
	}

	// 如果目标值小于第一个元素或大于最后一个元素，提前退出
	if target < arr[0] || target > arr[len(arr)-1] {
		return -1
	}

	for i := 0; i < len(arr); i++ {
		if arr[i] == target {
			return i
		}
		// 如果当前元素已经大于目标值，提前退出（适用于有序数组）
		if arr[i] > target {
			break
		}
	}
	return -1
}

// LinearSearchWithSentinel 哨兵线性搜索
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 优势: 减少边界检查，提高性能
func LinearSearchWithSentinel(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}

	// 保存最后一个元素
	last := arr[len(arr)-1]

	// 将目标值放在最后作为哨兵
	arr[len(arr)-1] = target

	i := 0
	// 不需要边界检查，因为哨兵保证会找到
	for arr[i] != target {
		i++
	}

	// 恢复原数组
	arr[len(arr)-1] = last

	// 判断是否找到
	if i < len(arr)-1 || last == target {
		return i
	}
	return -1
}

// LinearSearchRecursive 递归线性搜索
// 时间复杂度: O(n)
// 空间复杂度: O(n) - 递归调用栈
func LinearSearchRecursive(arr []int, target int) int {
	return linearSearchRecursiveHelper(arr, target, 0)
}

func linearSearchRecursiveHelper(arr []int, target, index int) int {
	// 基础情况：到达数组末尾
	if index >= len(arr) {
		return -1
	}

	// 找到目标值
	if arr[index] == target {
		return index
	}

	// 递归搜索剩余部分
	return linearSearchRecursiveHelper(arr, target, index+1)
}

// LinearSearchParallel 并行线性搜索（分块处理）
// 时间复杂度: O(n/p) - p为处理器数量
// 空间复杂度: O(1)
func LinearSearchParallel(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}

	// 简单的分块处理（实际应用中可以使用goroutine）
	chunkSize := len(arr) / 4 // 分成4块
	if chunkSize == 0 {
		chunkSize = 1
	}

	for chunk := 0; chunk < len(arr); chunk += chunkSize {
		end := chunk + chunkSize
		if end > len(arr) {
			end = len(arr)
		}

		// 在当前块中搜索
		for i := chunk; i < end; i++ {
			if arr[i] == target {
				return i
			}
		}
	}

	return -1
}

// LinearSearchWithComparator 带比较器的线性搜索
// 时间复杂度: O(n)
// 空间复杂度: O(1)
type Comparator func(a, b int) bool

func LinearSearchWithComparator(arr []int, target int, comparator Comparator) int {
	for i := 0; i < len(arr); i++ {
		if comparator(arr[i], target) {
			return i
		}
	}
	return -1
}

// BinarySearch 二分搜索算法
// 时间复杂度: O(log n)
// 空间复杂度: O(1)
// 前提条件: 数组必须是有序的
func BinarySearch(arr []int, target int) int {
	left := 0
	right := len(arr) - 1

	for left <= right {
		mid := left + (right-left)/2 // 避免整数溢出

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

// BinarySearchRecursive 递归二分搜索
// 时间复杂度: O(log n)
// 空间复杂度: O(log n) - 递归调用栈
func BinarySearchRecursive(arr []int, target int) int {
	return binarySearchRecursiveHelper(arr, target, 0, len(arr)-1)
}

func binarySearchRecursiveHelper(arr []int, target, left, right int) int {
	if left > right {
		return -1
	}

	mid := left + (right-left)/2

	if arr[mid] == target {
		return mid
	} else if arr[mid] < target {
		return binarySearchRecursiveHelper(arr, target, mid+1, right)
	} else {
		return binarySearchRecursiveHelper(arr, target, left, mid-1)
	}
}

// BinarySearchFirst 二分搜索第一个匹配项
// 时间复杂度: O(log n)
// 空间复杂度: O(1)
func BinarySearchFirst(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			result = mid
			right = mid - 1 // 继续向左搜索
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}

// BinarySearchLast 二分搜索最后一个匹配项
// 时间复杂度: O(log n)
// 空间复杂度: O(1)
func BinarySearchLast(arr []int, target int) int {
	left := 0
	right := len(arr) - 1
	result := -1

	for left <= right {
		mid := left + (right-left)/2

		if arr[mid] == target {
			result = mid
			left = mid + 1 // 继续向右搜索
		} else if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return result
}

// BinarySearchRange 二分搜索目标值的范围
// 时间复杂度: O(log n)
// 空间复杂度: O(1)
func BinarySearchRange(arr []int, target int) (int, int) {
	first := BinarySearchFirst(arr, target)
	if first == -1 {
		return -1, -1
	}

	last := BinarySearchLast(arr, target)
	return first, last
}

// JumpSearch 跳跃搜索算法
// 时间复杂度: O(√n)
// 空间复杂度: O(1)
// 前提条件: 数组必须是有序的
func JumpSearch(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}

	// 计算跳跃步长
	step := int(float64(len(arr)) * 0.5)
	if step == 0 {
		step = 1
	}

	// 跳跃阶段
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

// InterpolationSearch 插值搜索算法
// 时间复杂度: O(log log n) - 平均情况，O(n) - 最坏情况
// 空间复杂度: O(1)
// 前提条件: 数组必须是有序的，且元素分布均匀
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

// ExponentialSearch 指数搜索算法
// 时间复杂度: O(log n)
// 空间复杂度: O(1)
// 前提条件: 数组必须是有序的
func ExponentialSearch(arr []int, target int) int {
	if len(arr) == 0 {
		return -1
	}

	// 如果第一个元素就是目标值
	if arr[0] == target {
		return 0
	}

	// 找到范围
	i := 1
	for i < len(arr) && arr[i] <= target {
		i = i * 2
	}

	// 在找到的范围内进行二分搜索
	return binarySearchRecursiveHelper(arr, target, i/2, min(i, len(arr)-1))
}

// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// 测试函数
func TestSearchAlgorithms() {
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	target := 22

	fmt.Println("=== 线性搜索测试 ===")
	fmt.Println("原始数组:", arr)
	fmt.Println("搜索目标:", target)

	// 基本线性搜索
	result := LinearSearch(arr, target)
	fmt.Printf("基本线性搜索结果: 索引 %d\n", result)

	// 搜索所有匹配项
	arrWithDuplicates := []int{1, 2, 3, 2, 4, 2, 5}
	allResults := LinearSearchAll(arrWithDuplicates, 2)
	fmt.Printf("搜索所有2的结果: 索引 %v\n", allResults)

	fmt.Println("\n=== 二分搜索测试 ===")
	sortedArr := []int{11, 12, 22, 25, 34, 64, 90}
	fmt.Println("有序数组:", sortedArr)

	// 基本二分搜索
	binaryResult := BinarySearch(sortedArr, target)
	fmt.Printf("二分搜索结果: 索引 %d\n", binaryResult)

	// 递归二分搜索
	recursiveResult := BinarySearchRecursive(sortedArr, target)
	fmt.Printf("递归二分搜索结果: 索引 %d\n", recursiveResult)

	// 搜索范围
	rangeArr := []int{1, 2, 2, 2, 3, 4, 5}
	first, last := BinarySearchRange(rangeArr, 2)
	fmt.Printf("目标值2的范围: [%d, %d]\n", first, last)

	fmt.Println("\n=== 其他搜索算法测试 ===")

	// 跳跃搜索
	jumpResult := JumpSearch(sortedArr, target)
	fmt.Printf("跳跃搜索结果: 索引 %d\n", jumpResult)

	// 插值搜索
	interpolationResult := InterpolationSearch(sortedArr, target)
	fmt.Printf("插值搜索结果: 索引 %d\n", interpolationResult)

	// 指数搜索
	exponentialResult := ExponentialSearch(sortedArr, target)
	fmt.Printf("指数搜索结果: 索引 %d\n", exponentialResult)
}
