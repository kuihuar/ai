package suanfa

import "cmp"

func BubbleSort(arr []int) {
	n := len(arr)

	// 外层循环：控制冒泡的轮数，一共进行 n-1 轮
	for i := 0; i < n-1; i++ {
		// 内层循环：逐对比较相邻元素
		// 注意：每经过一轮，末尾就多了一个确定排好序的元素，
		// 所以比较终止位置是 n - 1 - i
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				// 如果前一个元素大于后一个元素，交换它们（让大元素向右“冒泡”）
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func SelectionSort(arr []int) {
	n := len(arr)

	// 外层循环：决定当前要确定位置的索引 i
	// 只需要遍历到 n-1，因为剩下最后一个元素自然是最大的
	for i := 0; i < n-1; i++ {
		// 假设当前未排序区域的第一个元素就是最小的
		minIndex := i
		// 内层循环：从 i+1 开始往后查找未排序区域中真正的最小值索引
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				// 找到了更小的元素，更新最小值的索引
				minIndex = j
			}
		}
		// 如果找到了比 arr[i] 更小的元素，则交换它们的位置
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
	}
}

// 插入排序的基本思想是：
//
// 1. 从第一个元素开始，该元素可以认为已经被排序
// 2. 取出下一个元素，在已经排序的元素序列中从后向前扫描
// 3. 如果该元素（已排序）大于新元素，将该元素移到下一位置
// 4. 重复步骤3，直到找到已排序的元素小于或者等于新元素的位置
// 5. 将新元素插入到该位置后
func InsertionSort(arr []int) {
	n := len(arr)
	if n <= 1 {
		return
	}
	// 从第 2 个元素（索引 1）开始，把它插入到左边已排序的区域中
	for i := 1; i < n; i++ {
		// 取出当前元素
		key := arr[i] // 当前需要寻找位置并插入的元素
		// 从当前元素的前一个元素开始向前扫描
		// 从已排序区域的最右侧开始比较
		j := i - 1
		// 将大于key的元素向后移动

		// 从右向左遍历已排序区域：
		// 只要已排序区的元素 arr[j] 比 key 大，就把 arr[j] 往右挪一位
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j] // 元素右移，腾出位置
			j--
		}
		// 将key插入到正确的位置
		// 找到了合适的位置（第一个 <= key 的元素右边），把 key 插入进去
		arr[j+1] = key
	}
}

// InsertionSort 是插入排序的通用实现。
// 支持任意实现了 cmp.Ordered 的类型（int, float, string 等）。
func InsertionSortGeneric[T cmp.Ordered](arr []T) {
	n := len(arr)
	if n <= 1 {
		return
	}

	// 从第 2 个元素（索引 1）开始，把它插入到左边已排序的区域中
	for i := 1; i < n; i++ {
		key := arr[i] // 当前需要寻找位置并插入的元素
		j := i - 1    // 从已排序区域的最右侧开始比较

		// 从右向左遍历已排序区域：
		// 只要已排序区的元素 arr[j] 比 key 大，就把 arr[j] 往右挪一位
		for j >= 0 && arr[j] > key {
			//第一次执行 arr[j+1] = arr[j] 时（即 j = i - 1）
			//此时 j + 1 刚好就是 i
			arr[j+1] = arr[j] // 元素右移，腾出位置（）
			j--
		}

		// 找到了合适的位置（第一个 <= key 的元素右边），把 key 插入进去
		arr[j+1] = key
	}
}

// 快速排序步骤说明
// 选择基准值（Pivot）

// 通常选择数组中间、首尾元素或随机位置的元素作为基准，此处以中间元素为例。

// 分区操作（Partitioning）

// 将数组分为两部分：小于基准的元素移到左侧，大于基准的元素移到右侧。

// 使用左右双指针向中间扫描，交换不符合条件的元素。

// 递归排序子数组

// 对左右两个子数组递归执行上述步骤，直到子数组长度为1或0。
func QuickSortInt(arr []int, low, high int) {
	if low < high {
		pivotIndex := partitionInt(arr, low, high)
		QuickSortInt(arr, low, pivotIndex-1)
		QuickSortInt(arr, pivotIndex+1, high)
	}
}

func partitionInt(arr []int, low, high int) int {
	// 选择基准值
	pivot := arr[(low+high)/2]
	left, right := low, high
	for left <= right {
		// 从左向右找到第一个大于等于基准值的元素
		for arr[left] < pivot {
			left++
		}
		// 从右向左找到第一个小于等于基准值的元素
		for arr[right] > pivot {
			right--
		}
		if left <= right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
			right--
		}
	}
	return left
}

func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	mid := len(arr) / 2
	left := MergeSort(arr[:mid])
	right := MergeSort(arr[mid:])
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

// 是暴露给外部调用的快速排序入口函数
func QuickSortGeneric[T cmp.Ordered](arr []T) {
	// 边界条件判断：如果切片为空或只有一个元素，本身就是有序的，直接返回
	if len(arr) <= 1 {
		return
	}
	// 调用递归辅助函数，初始区间为整个切片的索引范围 [0, len(arr)-1]
	quickSortGenericHelper(arr, 0, len(arr)-1)
}

// quickSortGenericHelper 是递归执行快排的核心函数。
func quickSortGenericHelper[T cmp.Ordered](arr []T, low, high int) {
	// 当区间有效（至少包含 2 个元素）时才继续递归排序
	if low < high {
		// 执行分区，获取左半部分的结束/右半部分的起始分界线
		pivotIndex := partitionGeneric(arr, low, high)
		// 【关键注意点】：
		// 采用双指针相向移动（Hoare 变体）时，Pivot 并没有固定在 splitIndex 位置，
		// 而是把数组分成了 [low, splitIndex-1] 和 [splitIndex, high] 两部分。
		// 因此递归区间必须是 splitIndex-1 和 splitIndex，不能写成 splitIndex+1！
		quickSortGenericHelper(arr, low, pivotIndex-1)  // 递归排序左半部分
		quickSortGenericHelper(arr, pivotIndex+1, high) // 递归排序右半部分
	}
}

// partitionGeneric 采用双指针相向移动法（Hoare 变体）进行分区。
// 它把切片划分成两部分：左侧的所有元素都 <= pivot，右侧的所有元素都 >= pivot。
// 返回值 left 是右半部分区间的起始索引。
func partitionGeneric[T cmp.Ordered](arr []T, low, high int) int {
	// 1. 取中间节点的值作为基准值 (Pivot)，能有效规避已排序数组导致性能退化的问题
	pivot := arr[(low+high)/2]
	// 2. 初始化双指针，分别指向区间的头和尾
	left, right := low, high

	// 3. 循环直到左右指针相遇或交错
	for left <= right {
		// 从左往右找：寻找第一个大于或等于 pivot 的元素
		for arr[left] < pivot {
			left++
		}
		// 从右往左找：寻找第一个小于或等于 pivot 的元素
		for arr[right] > pivot {
			right--
		}
		// 如果左指针仍在右指针左侧（或重合），说明找到了需要交换的一对元素
		if left <= right {
			// 交换元素，将较小的换到左边，较大的换到右边
			arr[left], arr[right] = arr[right], arr[left]
			// 交换后，两指针继续向中间靠拢
			left++
			right--
		}
	}
	// 4. 返回分界点 left。
	// 循环结束后，left 左侧（不含 left）的元素都 <= pivot，left 及右侧的元素都 >= pivot
	return left
}
