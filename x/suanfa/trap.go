package suanfa

import "fmt"

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

// 接雨水，暴力解法
func trap(height []int) int {
	if len(height) == 0 {
		return 0
	}
	var res int
	for i := 1; i < len(height)-1; i++ {
		fmt.Printf("i:%d\n=======\n", i)
		leftMax := 0
		for l := 0; l < i; l++ {
			leftMax = max(leftMax, height[l])
		}
		fmt.Println("leftMax:", leftMax)
		rightMax := 0
		for r := i + 1; r < len(height); r++ {
			rightMax = max(rightMax, height[r])
		}
		fmt.Println("rightMax:", rightMax)
		trap := min(leftMax, rightMax) - height[i]
		fmt.Printf("min(%d,%d) - height[i]:%d\n", leftMax, rightMax, height[i])
		if trap > 0 {
			res += trap
		}
		fmt.Printf("trap:%d\n", trap)
		fmt.Printf("res:%d\n=======\n", res)
	}
	return res
}
func trap4(height []int) int {
	if len(height) == 0 {
		return 0
	}
	var res int
	n := len(height)
	leftMax := make([]int, n)
	rightMax := make([]int, n)

	// 1. 从左往右算左最大值（从 i=1 开始！）
	leftMax[0] = height[0]
	for i := 1; i < n; i++ {
		leftMax[i] = max(leftMax[i-1], height[i])
	}

	// 2. 从右往左算右最大值（从 i=n-2 开始！）
	rightMax[n-1] = height[n-1]
	for i := n - 2; i >= 0; i-- {
		rightMax[i] = max(rightMax[i+1], height[i])
	}

	// 3. 计算总雨水量
	for i := 1; i < n-1; i++ {
		res += min(leftMax[i], rightMax[i]) - height[i]
	}

	return res
}

// 接雨水，双指针解法
func trap2(height []int) int {
	if len(height) == 0 {
		return 0
	}
	var res int
	left, right := 0, len(height)-1
	leftMax, rightMax := height[left], height[right]
	fmt.Printf("初始状态:left=%d,right=%d,leftMax=%d,rightMax=%d\n", left, right, leftMax, rightMax)
	round := 1
	resround := 0
	for left < right {
		fmt.Printf("第%d轮,left=%d,right=%d,", round, left, right)
		round++
		if height[left] < height[right] {
			fmt.Printf("左边柱子高(%d)<右边柱子高(%d):处理左边,", height[left], height[right])
			if height[left] >= leftMax {
				fmt.Printf("左边柱子高(%d)大于等于leftMax(%d),", height[left], leftMax)
				leftMax = height[left]
				fmt.Printf("更新leftMax=%d", leftMax)
			} else {
				fmt.Printf("左边柱子高(%d)小于leftMax(%d),", height[left], leftMax)
				res += leftMax - height[left]
				resround++
				fmt.Printf("接雨水(%d),res:%d, leftMax(%d)-height[%d](%d)=%d", resround, res, leftMax, left, height[left], leftMax-height[left])
			}
			left++
			fmt.Printf("移动左边:left++ %d\n", left)
		} else {
			fmt.Printf("右边柱子高(%d)>=左边柱子高(%d):处理右边,", height[right], height[left])
			if height[right] >= rightMax {
				fmt.Printf("右边柱子高(%d)大于等于rightMax(%d),", height[right], rightMax)
				rightMax = height[right]
				fmt.Printf("更新rightMax=%d", rightMax)
			} else {
				fmt.Printf("右边柱子高(%d)小于rightMax(%d),", height[right], rightMax)
				res += rightMax - height[right]
				resround++
				fmt.Printf("接雨水(%d),res:%d, rightMax(%d)-height[%d](%d)=%d", resround, res, rightMax, right, height[right], rightMax-height[right])
			}
			right--
			fmt.Printf("移动右边:right-- %d \n", right)
		}

	}
	return res
}

func trapX(height []int) int {
	if len(height) == 0 {
		return 0
	}
	left, right := 0, len(height)-1
	leftMax, rightMax := 0, 0
	var res int

	for left < right {
		if height[left] < height[right] {
			if height[left] >= leftMax {
				leftMax = height[left]
			} else {
				res += leftMax - height[left]
			}
			left++
		} else {
			if height[right] >= rightMax {
				rightMax = height[right]
			} else {
				res += rightMax - height[right]
			}
			right--
		}
	}
	return res
}

// 接雨水，栈解法
func trap3(height []int) int {
	if len(height) == 0 {
		return 0
	}
	var res int
	stack := []int{}
	for i := 0; i < len(height); i++ {
		for len(stack) > 0 && height[i] > height[stack[len(stack)-1]] {
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			if len(stack) == 0 {
				break
			}
			res += (i - stack[len(stack)-1] - 1) * (height[i] - height[top])
		}
		stack = append(stack, i)
	}
	return res
}

// 接雨水，动态规划解法
func trap4(height []int) int {
	if len(height) == 0 {
		return 0
	}
	var res int
	n := len(height)
	leftMax := make([]int, n)
	rightMax := make([]int, n)

	// 1. 从左往右算左最大值（从 i=1 开始！）
	leftMax[0] = height[0]
	for i := 1; i < n; i++ {
		leftMax[i] = max(leftMax[i-1], height[i])
	}

	// 2. 从右往左算右最大值（从 i=n-2 开始！）
	rightMax[n-1] = height[n-1]
	for i := n - 2; i >= 0; i-- {
		rightMax[i] = max(rightMax[i+1], height[i])
	}

	// 3. 计算总雨水量
	for i := 1; i < n-1; i++ {
		res += min(leftMax[i], rightMax[i]) - height[i]
	}

	return res
}

// 必须带上这两个函数

func max(a, b int) int {
	if a > b {
		return a
	}
	return b

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// trapBruteForce 接雨水 - 暴力解法
// 时间复杂度: O(n²) | 空间复杂度: O(1)
func trapBruteForce(height []int) int {
	if len(height) == 0 {
		return 0
	}

	n := len(height)
	var res int

	// 最左和最右两个边界一定存不住水，直接从 1 遍历到 n-2
	for i := 1; i < n-1; i++ {
		leftMax := 0
		rightMax := 0

		// 1. 现场向左扫描，寻找当前位置左侧（包含自身）的最大高度
		for l := i; l >= 0; l-- {
			leftMax = max(leftMax, height[l])
		}

		// 2. 现场向右扫描，寻找当前位置右侧（包含自身）的最大高度
		for r := i; r < n; r++ {
			rightMax = max(rightMax, height[r])
		}

		// 3. 当前柱子能接的水 = min(左侧最高, 右侧最高) - 当前柱子高度
		res += min(leftMax, rightMax) - height[i]
	}

	return res
}

// trap 计算给定高度图能接的雨水总量
// 核心思想：对于索引 i 处的柱子，它上方能装的水取决于 min(左边最高柱子, 右边最高柱子) - 当前柱子高度
func trapDP(height []int) int {
	// 边界条件判断：如果数组为空，无法接水，直接返回 0
	if len(height) == 0 {
		return 0
	}

	n := len(height)
	var res int

	// leftMax[i] 表示索引 i 及其左边所有柱子中的最大高度
	leftMax := make([]int, n)
	// rightMax[i] 表示索引 i 及其右边所有柱子中的最大高度
	rightMax := make([]int, n)

	// 初始化边界值：
	// 最左边的柱子，其左侧最大高度就是它自己
	leftMax[0] = height[0]
	// 最右边的柱子，其右侧最大高度就是它自己
	rightMax[n-1] = height[n-1]

	// 从左向右遍历，填充 leftMax 数组
	// 状态转移方程：leftMax[i] = max(左边前一个位置的最大值, 当前柱子高度)
	for left := 1; left < n; left++ {
		leftMax[left] = max(leftMax[left-1], height[left])
	}

	// 从右向左遍历，填充 rightMax 数组
	// 状态转移方程：rightMax[i] = max(右边后一个位置的最大值, 当前柱子高度)
	for right := n - 2; right >= 0; right-- {
		rightMax[right] = max(rightMax[right+1], height[right])
	}

	// 计算每个柱子上方能接的水，并累加到结果中
	// 注意：最左边 (i=0) 和最右边 (i=n-1) 的柱子边界无法留水，所以从 1 遍历到 n-2
	for i := 1; i < n-1; i++ {
		// 当前位置能装的水 = min(左边最高, 右边最高) - 当前柱子高度
		res += min(leftMax[i], rightMax[i]) - height[i]
	}

	return res
}

// trapDynamicProgramming 接雨水 - 动态规划解法
// 时间复杂度: O(n) | 空间复杂度: O(n)
func trapDynamicProgramming(height []int) int {
	if len(height) == 0 {
		return 0
	}

	n := len(height)
	var res int

	// leftMax[i] 表示位置 i 及其左边的最高柱子
	// rightMax[i] 表示位置 i 及其右边的最高柱子
	leftMax := make([]int, n)
	rightMax := make([]int, n)

	leftMax[0] = height[0]
	rightMax[n-1] = height[n-1]

	// 从左向右预处理左侧最大值
	for left := 1; left < n; left++ {
		leftMax[left] = max(leftMax[left-1], height[left])
	}

	// 从右向左预处理右侧最大值
	for right := n - 2; right >= 0; right-- {
		rightMax[right] = max(rightMax[right+1], height[right])
	}

	// 汇总每个位置的接水量
	for i := 1; i < n-1; i++ {
		res += min(leftMax[i], rightMax[i]) - height[i]
	}

	return res
}



// trapTwoPointers 接雨水 - 双指针解法
// 时间复杂度: O(n) | 空间复杂度: O(1)
func trapTwoPointers(height []int) int {
	if len(height) == 0 {
		return 0
	}

	// 1. 初始化双指针分别指向首尾
	left, right := 0, len(height)-1

	// 2. 维护左右两侧已知的最大值
	leftMax, rightMax := 0, 0

	var res int

	// 3. 当左右指针未相遇时循环
	for left < right {
		// 更新左侧和右侧的最大高度
		leftMax = max(leftMax, height[left])
		rightMax = max(rightMax, height[right])

		// 4. 谁小就结算谁，因为较小的一方是木桶的“短板”
		if leftMax < rightMax {
			// 左侧是瓶颈，计算 left 位置的水量并右移
			res += leftMax - height[left]
			left++
		} else {
			// 右侧是瓶颈，计算 right 位置的水量并左移
			res += rightMax - height[right]
			right--
		}
	}

	return res
}


3. 三种解法全方位对比
维度		暴力解法 (trapBruteForce)	动态规划 (trapDynamicProgramming)	双指针 (trapTwoPointers)
时间复杂度	$\mathcal{O}(n^2)$			$\mathcal{O}(n)$					$\mathcal{O}(n)$
空间复杂度	$\mathcal{O}(1)$			$\mathcal{O}(n)$					$\mathcal{O}(1)$