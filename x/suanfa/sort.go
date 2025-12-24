package suanfa

func BubbleSort(arr []int) {
	n := len(arr)

	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
		}
	}
}

func SelectionSort(arr []int) {
	n := len(arr)

	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		arr[i], arr[minIndex] = arr[minIndex], arr[i]
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
	for i := 1; i < n; i++ {
		// 取出当前元素
		key := arr[i]
		// 从当前元素的前一个元素开始向前扫描
		j := i - 1
		// 将大于key的元素向后移动
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		// 将key插入到正确的位置
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
func QuickSort(arr []int, low, high int) {
	if low < high {
		pivotIndex := partition(arr, low, high)
		QuickSort(arr, low, pivotIndex-1)
		QuickSort(arr, pivotIndex+1, high)
	}
}
func partition(arr []int, low, high int) int {
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
