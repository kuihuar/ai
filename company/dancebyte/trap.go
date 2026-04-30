package dancebyte

import "fmt"

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
