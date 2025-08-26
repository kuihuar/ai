package algorithm

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

// 测试数据
var testCases = []struct {
	name     string
	input    []int
	expected []int
}{
	{
		name:     "空数组",
		input:    []int{},
		expected: []int{},
	},
	{
		name:     "单个元素",
		input:    []int{1},
		expected: []int{1},
	},
	{
		name:     "已排序数组",
		input:    []int{1, 2, 3, 4, 5},
		expected: []int{1, 2, 3, 4, 5},
	},
	{
		name:     "逆序数组",
		input:    []int{5, 4, 3, 2, 1},
		expected: []int{1, 2, 3, 4, 5},
	},
	{
		name:     "重复元素",
		input:    []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5},
		expected: []int{1, 1, 2, 3, 3, 4, 5, 5, 5, 6, 9},
	},
	{
		name:     "负数元素",
		input:    []int{-3, 1, -4, 1, -5, 9, -2, 6, -5, 3, -5},
		expected: []int{-5, -5, -5, -4, -3, -2, 1, 1, 3, 6, 9},
	},
	{
		name:     "大数组",
		input:    []int{64, 34, 25, 12, 22, 11, 90, 88, 76, 54, 32, 21, 19, 8, 7, 6, 5, 4, 3, 2, 1},
		expected: []int{1, 2, 3, 4, 5, 6, 7, 8, 11, 12, 19, 21, 22, 25, 32, 34, 54, 64, 76, 88, 90},
	},
}

// 辅助函数：创建数组副本
func copyArray(arr []int) []int {
	result := make([]int, len(arr))
	copy(result, arr)
	return result
}

// 辅助函数：验证排序结果
func verifySort(t *testing.T, algorithmName string, input []int, result []int, expected []int) {
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("%s 排序失败:\n输入: %v\n期望: %v\n实际: %v",
			algorithmName, input, expected, result)
	}
}

// 测试冒泡排序
func TestBubbleSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := copyArray(tc.input)
			BubbleSort(input)
			verifySort(t, "冒泡排序", tc.input, input, tc.expected)
		})
	}
}

// 测试选择排序
func TestSelectSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := copyArray(tc.input)
			SelectSort(input)
			verifySort(t, "选择排序", tc.input, input, tc.expected)
		})
	}
}

// 测试插入排序
func TestInsertSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := copyArray(tc.input)
			InsertSort(input)
			verifySort(t, "插入排序", tc.input, input, tc.expected)
		})
	}
}

// 测试快速排序
func TestQuickSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := copyArray(tc.input)
			QuickSort(input)
			verifySort(t, "快速排序", tc.input, input, tc.expected)
		})
	}
}

// 测试归并排序（原地排序版本）
func TestMergeSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := copyArray(tc.input)
			MergeSort(input)
			verifySort(t, "归并排序", tc.input, input, tc.expected)
		})
	}
}

// 测试归并排序2（函数式版本）
func TestMergeSort2(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := copyArray(tc.input)
			result := MergeSort2(input)
			verifySort(t, "归并排序2", tc.input, result, tc.expected)

			// 验证原数组没有被修改
			if !reflect.DeepEqual(input, tc.input) {
				t.Errorf("归并排序2 修改了原数组:\n原数组: %v\n修改后: %v", tc.input, input)
			}
		})
	}
}

// 基准测试：比较不同排序算法的性能
func BenchmarkSortingAlgorithms(b *testing.B) {
	// 创建测试数据
	testData := []int{64, 34, 25, 12, 22, 11, 90, 88, 76, 54, 32, 21, 19, 8, 7, 6, 5, 4, 3, 2, 1}

	b.Run("BubbleSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := copyArray(testData)
			BubbleSort(data)
		}
	})

	b.Run("SelectSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := copyArray(testData)
			SelectSort(data)
		}
	})

	b.Run("InsertSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := copyArray(testData)
			InsertSort(data)
		}
	})

	b.Run("QuickSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := copyArray(testData)
			QuickSort(data)
		}
	})

	b.Run("MergeSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := copyArray(testData)
			MergeSort(data)
		}
	})

	b.Run("MergeSort2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := copyArray(testData)
			MergeSort2(data)
		}
	})

	b.Run("GoSort", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			data := copyArray(testData)
			sort.Ints(data)
		}
	})
}

// 测试边界情况
func TestEdgeCases(t *testing.T) {
	t.Run("空数组", func(t *testing.T) {
		empty := []int{}

		// 测试所有排序算法
		BubbleSort(empty)
		SelectSort(empty)
		InsertSort(empty)
		QuickSort(empty)
		MergeSort(empty)
		result := MergeSort2(empty)

		if len(result) != 0 {
			t.Errorf("空数组排序后应该还是空数组，但得到了: %v", result)
		}
	})

	t.Run("单个元素", func(t *testing.T) {
		single := []int{42}
		expected := []int{42}

		BubbleSort(single)
		if !reflect.DeepEqual(single, expected) {
			t.Errorf("单个元素排序失败: %v", single)
		}

		result := MergeSort2(single)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("单个元素排序失败: %v", result)
		}
	})

	t.Run("相同元素", func(t *testing.T) {
		same := []int{1, 1, 1, 1, 1}
		expected := []int{1, 1, 1, 1, 1}

		BubbleSort(same)
		if !reflect.DeepEqual(same, expected) {
			t.Errorf("相同元素排序失败: %v", same)
		}

		result := MergeSort2(same)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("相同元素排序失败: %v", result)
		}
	})
}

// 测试稳定性（对于稳定排序算法）
func TestStability(t *testing.T) {
	// 创建一个包含重复元素的数组，每个元素包含值和索引
	type Item struct {
		value int
		index int
	}

	items := []Item{
		{3, 0}, {1, 1}, {2, 2}, {3, 3}, {1, 4}, {2, 5},
	}

	// 只按值排序，保持索引的相对位置
	sort.Slice(items, func(i, j int) bool {
		return items[i].value < items[j].value
	})

	// 验证稳定性：相同值的元素应该保持原来的相对顺序
	expected := []Item{
		{1, 1}, {1, 4}, {2, 2}, {2, 5}, {3, 0}, {3, 3},
	}

	if !reflect.DeepEqual(items, expected) {
		t.Errorf("稳定性测试失败:\n期望: %v\n实际: %v", expected, items)
	}
}

// ExampleBubbleSort 演示冒泡排序的使用
func ExampleBubbleSort() {
	// 创建一个未排序的数组
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Println("排序前:", arr)

	// 使用冒泡排序
	BubbleSort(arr)
	fmt.Println("排序后:", arr)

	// Output:
	// 排序前: [64 34 25 12 22 11 90]
	// 排序后: [11 12 22 25 34 64 90]
}

// ExampleSelectSort 演示选择排序的使用
func ExampleSelectSort() {
	// 创建一个包含重复元素的数组
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("排序前:", arr)

	// 使用选择排序
	SelectSort(arr)
	fmt.Println("排序后:", arr)

	// Output:
	// 排序前: [3 1 4 1 5 9 2 6]
	// 排序后: [1 1 2 3 4 5 6 9]
}

// TestSelectSortWithSteps 测试选择排序并显示每一步
func TestSelectSortWithSteps(t *testing.T) {
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("=== 选择排序详细步骤 ===")
	fmt.Println("初始数组:", arr)

	// 手动实现选择排序并显示每一步
	for i := 0; i < len(arr)-1; i++ {
		minIndex := i
		for j := i + 1; j < len(arr); j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
			fmt.Printf("第%d轮: 找到最小值%d，与位置%d交换 -> %v\n", i+1, arr[i], i+1, arr)
		} else {
			fmt.Printf("第%d轮: 位置%d已经是最小值%d，无需交换 -> %v\n", i+1, i+1, arr[i], arr)
		}
	}

	fmt.Println("最终结果:", arr)

	// 验证结果
	expected := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("排序结果错误:\n期望: %v\n实际: %v", expected, arr)
	}
}

// TestBubbleSortWithSteps 测试冒泡排序并显示每一步
func TestBubbleSortWithSteps(t *testing.T) {
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6}
	fmt.Println("=== 冒泡排序详细步骤 ===")
	fmt.Println("初始数组:", arr)

	// 手动实现冒泡排序并显示每一步
	for i := 0; i < len(arr)-1; i++ {
		swapped := false
		for j := 0; j < len(arr)-1-i; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
				fmt.Printf("第%d轮第%d次: 交换%d和%d -> %v\n", i+1, j+1, arr[j], arr[j+1], arr)
			}
		}
		if !swapped {
			fmt.Printf("第%d轮: 没有交换，数组已排序\n", i+1)
			break
		}
	}

	fmt.Println("最终结果:", arr)

	// 验证结果
	expected := []int{1, 1, 2, 3, 4, 5, 6, 9}
	if !reflect.DeepEqual(arr, expected) {
		t.Errorf("排序结果错误:\n期望: %v\n实际: %v", expected, arr)
	}
}

// ExampleInsertSort 演示插入排序的使用
func ExampleInsertSort() {
	// 创建一个包含负数的数组
	arr := []int{-3, 1, -4, 1, -5, 9, -2, 6}
	fmt.Println("排序前:", arr)

	// 使用插入排序
	InsertSort(arr)
	fmt.Println("排序后:", arr)

	// Output:
	// 排序前: [-3 1 -4 1 -5 9 -2 6]
	// 排序后: [-5 -4 -3 -2 1 1 6 9]
}

// ExampleQuickSort 演示快速排序的使用
func ExampleQuickSort() {
	// 创建一个逆序数组
	arr := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	fmt.Println("排序前:", arr)

	// 使用快速排序
	QuickSort(arr)
	fmt.Println("排序后:", arr)

	// Output:
	// 排序前: [9 8 7 6 5 4 3 2 1]
	// 排序后: [1 2 3 4 5 6 7 8 9]
}

// ExampleMergeSort 演示归并排序的使用
func ExampleMergeSort() {
	// 创建一个已部分排序的数组
	arr := []int{1, 3, 5, 7, 2, 4, 6, 8}
	fmt.Println("排序前:", arr)

	// 使用归并排序
	MergeSort(arr)
	fmt.Println("排序后:", arr)

	// Output:
	// 排序前: [1 3 5 7 2 4 6 8]
	// 排序后: [1 2 3 4 5 6 7 8]
}

// ExampleMergeSort2 演示函数式归并排序的使用
func ExampleMergeSort2() {
	// 创建一个包含重复元素的数组
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
	fmt.Println("排序前:", arr)

	// 使用函数式归并排序（返回新数组）
	result := MergeSort2(arr)
	fmt.Println("排序后:", result)
	fmt.Println("原数组保持不变:", arr)

	// Output:
	// 排序前: [3 1 4 1 5 9 2 6 5 3 5]
	// 排序后: [1 1 2 3 3 4 5 5 5 6 9]
	// 原数组保持不变: [3 1 4 1 5 9 2 6 5 3 5]
}

func ExampleHeapSort() {
	// 创建一个已部分排序的数组
	arr := []int{1, 3, 5, 7, 2, 4, 6, 8}
	fmt.Println("排序前:", arr)

	// 使用归并排序
	HeapSort(arr)
	fmt.Println("排序后:", arr)

	// Output:
	// 排序前: [1 3 5 7 2 4 6 8]
	// 排序后: [1 2 3 4 5 6 7 8]
}

func TestHeapSort(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			input := copyArray(tc.input)
			fmt.Println("排序前:", input)
			HeapSort(input)
			fmt.Println("排序后:", input)
			verifySort(t, "堆排序", tc.input, input, tc.expected)
		})
	}
}

// TestGetHeapHeight 测试堆高度计算
func TestGetHeapHeight(t *testing.T) {
	tests := []struct {
		n      int
		height int
	}{
		{0, 0},  // 空堆
		{1, 0},  // 只有根节点
		{2, 1},  // 根节点 + 1个子节点
		{3, 1},  // 根节点 + 2个子节点
		{4, 2},  // 根节点 + 2个子节点 + 1个孙子节点
		{5, 2},  // 根节点 + 2个子节点 + 2个孙子节点
		{6, 2},  // 根节点 + 2个子节点 + 3个孙子节点
		{7, 2},  // 根节点 + 2个子节点 + 4个孙子节点
		{8, 3},  // 根节点 + 2个子节点 + 4个孙子节点 + 1个曾孙节点
		{15, 3}, // 完全二叉树，高度为3
		{16, 4}, // 完全二叉树，高度为4
	}

	for _, test := range tests {
		result := GetHeapHeight(test.n)
		if result != test.height {
			t.Errorf("GetHeapHeight(%d) = %d, want %d", test.n, result, test.height)
		}
	}
}

// TestHeapHeightVisualization 可视化堆高度计算
func TestHeapHeightVisualization(t *testing.T) {
	fmt.Println("堆高度计算示例：")
	fmt.Println("节点数 | 高度 | 说明")
	fmt.Println("-------|------|------")

	testCases := []struct {
		n      int
		height int
		desc   string
	}{
		{1, 0, "只有根节点"},
		{2, 1, "根节点 + 1个子节点"},
		{3, 1, "根节点 + 2个子节点"},
		{4, 2, "根节点 + 2个子节点 + 1个孙子节点"},
		{7, 2, "根节点 + 2个子节点 + 4个孙子节点"},
		{8, 3, "根节点 + 2个子节点 + 4个孙子节点 + 1个曾孙节点"},
	}

	for _, tc := range testCases {
		result := GetHeapHeight(tc.n)
		fmt.Printf("%6d | %4d | %s\n", tc.n, result, tc.desc)
		if result != tc.height {
			t.Errorf("GetHeapHeight(%d) = %d, want %d", tc.n, result, tc.height)
		}
	}
}
