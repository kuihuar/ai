package algorithm

// BubbleSort 冒泡排序
// 时间复杂度: O(n²) - 最坏和平均情况
// 空间复杂度: O(1) - 原地排序
// 稳定性: 稳定排序
// 算法思想: 重复遍历数组，每次比较相邻元素，如果顺序错误则交换
func BubbleSort(arr []int) {
	n := len(arr)
	// 外层循环控制排序轮数
	for i := 0; i < n-1; i++ {
		// 内层循环进行相邻元素比较和交换
		// 每轮排序后，最大的元素会"冒泡"到数组末尾
		for j := 0; j < n-i-1; j++ {
			// 如果前一个元素大于后一个元素，则交换它们
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

// SelectSort 选择排序
// 时间复杂度: O(n²) - 最坏、平均和最好情况都是
// 空间复杂度: O(1) - 原地排序
// 稳定性: 不稳定排序
// 算法思想: 每次从未排序区间选择最小的元素，放到已排序区间的末尾
func SelectSort(arr []int) {
	n := len(arr)
	// 外层循环控制已排序区间的边界
	for i := 0; i < n-1; i++ {
		// 假设当前位置的元素是最小的
		minIndex := i
		// 在未排序区间中寻找真正的最小元素
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		// 将找到的最小元素与当前位置的元素交换
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
	}
}

// InsertSort 插入排序
// 时间复杂度: O(n²) - 最坏和平均情况，O(n) - 最好情况（已排序）
// 空间复杂度: O(1) - 原地排序
// 稳定性: 稳定排序
// 算法思想: 将数组分为已排序和未排序两部分，每次从未排序部分取一个元素插入到已排序部分的正确位置
func InsertSort(arr []int) {
	n := len(arr)
	// 从第二个元素开始，逐个插入到已排序区间
	for i := 1; i < n; i++ {
		// 保存当前要插入的元素
		temp := arr[i]
		// 从已排序区间的末尾开始，向前查找插入位置
		j := i - 1
		// 如果已排序区间的元素大于待插入元素，则向后移动
		for j >= 0 && arr[j] > temp {
			arr[j+1] = arr[j]
			j--
		}
		// 在正确位置插入元素
		arr[j+1] = temp
	}
}

// QuickSort 快速排序
// 时间复杂度: O(n log n) - 平均情况，O(n²) - 最坏情况
// 空间复杂度: O(log n) - 平均情况，O(n) - 最坏情况
// 稳定性: 不稳定排序
// 算法思想: 分治法，选择一个基准元素，将数组分为两部分，左边都小于基准，右边都大于基准，然后递归排序
func QuickSort(arr []int) {
	quickSort(arr, 0, len(arr)-1)
}

// quickSort 快速排序的递归实现
func quickSort(arr []int, left, right int) {
	// 递归终止条件：左边界小于右边界
	if left < right {
		// 选择基准元素并进行分区
		pivot := partition(arr, left, right)
		// 递归排序左半部分
		quickSort(arr, left, pivot-1)
		// 递归排序右半部分
		quickSort(arr, pivot+1, right)
	}
}

// partition 分区函数，选择基准元素并将数组分为两部分
func partition(arr []int, left, right int) int {
	// 选择最后一个元素作为基准
	pivot := arr[right]

	// i 指向小于基准元素的最后一个位置
	i := left - 1

	// 遍历数组，将小于基准的元素移到左边
	for j := left; j < right; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	// 将基准元素放到正确的位置
	arr[i+1], arr[right] = arr[right], arr[i+1]

	// 返回基准元素的位置
	return i + 1
}

// MergeSort 归并排序
// 时间复杂度: O(n log n) - 最坏、平均和最好情况都是
// 空间复杂度: O(n) - 需要额外的临时数组
// 稳定性: 稳定排序
// 算法思想: 分治法，将数组分成两半，递归排序，然后合并两个有序数组
func MergeSort(arr []int) {
	n := len(arr)
	// 如果数组长度小于2，已经有序
	if n < 2 {
		return
	}
	// 创建临时数组用于合并过程
	temp := make([]int, n)
	mergeSort(arr, temp, 0, n-1)
}

// mergeSort 归并排序的递归实现
func mergeSort(arr []int, temp []int, left, right int) {
	// 递归终止条件：左边界小于右边界
	if left < right {
		// 计算中间位置
		mid := (left + right) / 2
		// 递归排序左半部分
		mergeSort(arr, temp, left, mid)
		// 递归排序右半部分
		mergeSort(arr, temp, mid+1, right)
		// 合并两个有序数组
		merge(arr, temp, left, mid, right)
	}
}

// merge 合并两个有序数组
func merge(arr []int, temp []int, left, mid, right int) {
	// 初始化三个指针：左半部分、右半部分、临时数组
	i, j, k := left, mid+1, left // ✅ k 从 left 开始，而不是 0

	// 比较两个有序数组的元素，将较小的放入临时数组
	for i <= mid && j <= right {
		if arr[i] <= arr[j] {
			temp[k] = arr[i]
			i++
		} else {
			temp[k] = arr[j]
			j++
		}
		k++
	}

	// 将左半部分剩余的元素复制到临时数组
	for i <= mid {
		temp[k] = arr[i]
		i++
		k++
	}

	// 将右半部分剩余的元素复制到临时数组
	for j <= right {
		temp[k] = arr[j]
		j++
		k++
	}

	// ✅ 正确复制：从 temp[left] 到 temp[k-1]
	// 将临时数组中的有序元素复制回原数组
	for idx := left; idx < k; idx++ {
		arr[idx] = temp[idx]
	}

	// ✅ 使用 temp[0:k] 来复制正确范围
	// copy(arr[left:right+1], temp[0:k])
}

// MergeSort2 归并排序的另一种实现（函数式风格）
// 时间复杂度: O(n log n) - 最坏、平均和最好情况都是
// 空间复杂度: O(n) - 需要额外的临时数组
// 稳定性: 稳定排序
// 算法思想: 分治法，将数组分成两半，递归排序，然后合并两个有序数组
// 与 MergeSort 的区别：返回新的排序数组，不修改原数组
func MergeSort2(arr []int) []int {
	// 递归终止条件：数组长度小于等于1时已经有序
	if len(arr) <= 1 {
		return arr
	}
	// 计算中间位置，将数组分成两半
	mid := len(arr) / 2
	// 递归排序左半部分
	left := MergeSort2(arr[:mid])
	// 递归排序右半部分
	right := MergeSort2(arr[mid:])
	// 合并两个有序数组并返回结果
	return merge2(left, right)
}

// merge2 合并两个有序数组的函数式实现
// 参数: left, right - 两个已排序的数组
// 返回: 合并后的有序数组
func merge2(left, right []int) []int {
	// 预分配结果数组的容量，避免频繁的内存重新分配
	result := make([]int, 0, len(left)+len(right))
	// 初始化两个数组的指针
	i, j := 0, 0
	// 比较两个数组的元素，将较小的添加到结果数组中
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	// 将左数组剩余的元素添加到结果中
	result = append(result, left[i:]...)
	// 将右数组剩余的元素添加到结果中
	result = append(result, right[j:]...)
	return result
}

// HeapSort 堆排序
// 时间复杂度: O(n log n) - 最坏、平均和最好情况都是
// 空间复杂度: O(1) - 原地排序
// 稳定性: 不稳定排序
// 算法思想: 利用堆的特性，先构建最大堆，然后逐个提取堆顶元素
// 正序：构建最大堆，然后逐个提取堆顶元素到末尾
// 倒序：构建最小堆，然后逐个提取堆顶元素到末尾
func HeapSort(arr []int) {
	n := len(arr)

	// 第一步：构建最大堆
	// 从最后一个非叶子节点开始，自底向上进行堆化
	for i := n/2 - 1; i >= 0; i-- {
		heapifyMax(arr, n, i)
	}
	// 第二步：逐个提取堆顶元素
	// 将堆顶元素（最大值）与末尾元素交换，然后重新堆化
	// 从后往前放置最大值
	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		heapifyMax(arr, i, 0)
	}

}

// heapifyMax 堆化函数
// 参数: arr - 待堆化的数组, n - 堆的大小, i - 要堆化的节点索引
// 功能: 将以 i 为根的子树调整为最大堆

func heapifyMax(arr []int, n, i int) {
	// 初始化最大值为根节点
	largest := i
	// 计算左子节点索引
	left := 2*i + 1
	// 计算右子节点索引
	right := 2*i + 2

	// 如果左子节点存在且大于根节点，更新最大值
	if left < n && arr[left] > arr[largest] {
		largest = left
	}

	// 如果右子节点存在且大于当前最大值，更新最大值
	if right < n && arr[right] > arr[largest] {
		largest = right
	}

	// 如果最大值不是根节点，则交换并继续堆化
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		// 递归堆化被交换的子树
		heapifyMax(arr, n, largest)
	}
}

func heapifyMin(arr []int, n, i int) {
	smallest := i
	left := 2*i + 1
	right := 2*i + 2

	if left < n && arr[left] < arr[smallest] {
		smallest = left
	}

	if right < n && arr[right] < arr[smallest] {
		smallest = right
	}

	if smallest != i {
		arr[i], arr[smallest] = arr[smallest], arr[i]
		heapifyMin(arr, n, smallest)
	}
}

func HeapsortDescending(arr []int) {
	n := len(arr)

	for i := n/2 - 1; i >= 0; i-- {
		heapifyMin(arr, n, i)
	}
	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		heapifyMin(arr, i, 0)
	}
}

// 堆相关的辅助函数

// IsMaxHeap 检查数组是否满足最大堆性质
func IsMaxHeap(arr []int) bool {
	n := len(arr)
	for i := 0; i < n/2; i++ {
		left := 2*i + 1
		right := 2*i + 2

		// 检查左子节点
		if left < n && arr[i] < arr[left] {
			return false
		}
		// 检查右子节点
		if right < n && arr[i] < arr[right] {
			return false
		}
	}
	return true
}

// GetHeapHeight 获取堆的高度
// 堆的高度 = ⌊log₂(n)⌋，其中 n 是堆中节点的数量
// 高度定义：从根节点到最远叶子节点的路径长度
func GetHeapHeight(n int) int {
	if n == 0 {
		return 0
	}
	// 使用位运算计算 log₂(n)
	height := 0
	temp := n
	for temp > 1 {
		temp = temp >> 1 // 等价于 temp = temp / 2
		height++
	}
	return height
}

// ShellSort 希尔排序
// 时间复杂度: O(n^1.3) - 平均情况，取决于增量序列
// 空间复杂度: O(1) - 原地排序
// 稳定性: 不稳定排序
// 算法思想: 插入排序的改进版，通过设置不同的间隔来分组排序

func ShellSort(arr []int) {
	n := len(arr)
	// 使用希尔增量序列: gap = gap/2，外层循环：控制间隔序列
	for gap := n / 2; gap > 0; gap /= 2 {
		// 中层循环：遍历每个子序列
		for i := gap; i < n; i++ {
			temp := arr[i] // 保存当前需要插入的元素
			j := i
			// 在子序列中进行插入排序
			// 内层循环：将 temp 插入到正确的位置
			for j >= gap && arr[j-gap] > temp {
				arr[j] = arr[j-gap] // 将较大的元素向后移动
				j -= gap            // 更新 j 以继续向前查找插入位置
			}
			arr[j] = temp // 插入 temp 到正确的位置
		}
	}
}

func ShellSortDescending(arr []int) {
	n := len(arr)

	for gap := n / 2; gap > 0; gap /= 2 {

		for i := gap; i < n; i++ {
			temp := arr[i]
			j := i
			// 降序：移动小的元素向后移动,区别在于这个比较符号
			// 在子序列中寻找 temp 的正确位置
			for j >= gap && arr[j-gap] < temp {
				arr[j] = arr[j-gap]
				j -= gap
			}
			arr[j] = temp
		}

	}
}
