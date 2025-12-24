package suanfa

import (
	"fmt"
	"sort"
)

func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

func Fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return Fibonacci(n-1) + Fibonacci(n-2)
}

// 暴力法：
// 1. 暴力法的时间复杂度为O(n)，空间复杂度为O(1)
func Power0(base, exponent int) int {
	res := 1
	for i := 0; i < exponent; i++ {
		res *= base
	}
	return res
}

// 分冶法：
// 1. 分冶法的时间复杂度为O(nlogn)，空间复杂度为O(logn)，因为用到了递归栈
func Power(base, exponent int) int {
	if exponent == 0 {
		return 1
	}
	if exponent < 0 {
		return 1 / Power(base, -exponent)
	}
	// 偶数条件下，base^exponent = (base*base)^(exponent/2)
	if exponent%2 == 0 {
		return Power(base*base, exponent/2)
	}
	// 奇数条件下，base^exponent = base * (base*base)^((exponent-1)/2)
	return base * Power(base*base, (exponent-1)/2)
}

// 非递归法：
// 1. 分解指数：

// 当指数为偶数时，𝑥𝑛=(𝑥𝑛/2)2

// 当指数为奇数时，𝑥𝑛=𝑥×𝑥𝑛−1

// 位运算的性质：

// 可以通过 exponent % 2 检查当前指数是奇数还是偶数。
// 将指数除以 2 相当于右移一位，这样可以在每次迭代中快速减少计算量。
func Power1(base float64, exponent int) float64 {
	result := 1.0
	if exponent < 0 {
		base = 1 / base
		exponent = -exponent
	}
	// 当前迭代的指数
	currentBase := base
	// 剩余的指数
	remainingExponent := exponent

	for remainingExponent > 0 {
		fmt.Printf("%d%%2: %v \t", remainingExponent, remainingExponent%2)
		if remainingExponent%2 == 1 {
			result *= currentBase
		}
		// 基数值平方
		currentBase *= currentBase
		// 右移指数（等价于除以2）
		remainingExponent /= 2
	}
	return result

}
func Power2(base, exponent int) int {
	result := 1

	for exponent > 0 {
		fmt.Printf("%d%%2: %v \t", exponent, exponent%2)
		if exponent%2 == 1 {
			result *= base
		}
		// 基数值平方
		base *= base
		// 右移指数（等价于除以2）
		// exponent /= 2
		exponent >>= 1
	}
	return result

}

// 1. 暴力法：
// 2. 哈希表法：
// 3. 排序法：
// 4. 分冶法：

// 排序法
// 排序之后，中间的元素就是出现次数超过一半的元素。
// 时间复杂度为O(nlogn)，空间复杂度为O(1)。
// 排序算法的时间复杂度为O(nlogn)，空间复杂度为O(1)。
// 排序法的核心原理是：多数元素必定占据排序后数组的中间位置
//索引: 0 1 2 3 4 5 6 7 8 9 10 11 12 13 ... 24
// 元素: 1 1 1 1 1 1 1 2 2 2 2 3 3 3 ... 3
//                    ↑
//                 中间位置

func MajorityElement(nums []int) int {
	sort.Ints(nums)
	return nums[len(nums)/2]
}

// 暴力法
// 遍历数组，统计每个元素出现的次数。
// 时间复杂度为O(n^2)，空间复杂度为O(1)。
// 暴力法的时间复杂度为O(n^2)，空间复杂度为O(1)。
func MajorityElement1(nums []int) int {
	n := len(nums)

	for i := 0; i < n; i++ {
		count := 0
		for j := 0; j < len(nums); j++ {
			if nums[i] == nums[j] {
				count++
			}
		}
		if count > n/2 {
			return nums[i]
		}
	}
	return -1
}

// 哈希表法
// 遍历数组，统计每个元素出现的次数。
// 时间复杂度为O(n)，空间复杂度为O(n)。
// 哈希表法的时间复杂度为O(n)，空间复杂度为O(n)。
func MajorityElement2(nums []int) int {
	n := len(nums)
	countMap := make(map[int]int)
	for i := 0; i < n; i++ {
		countMap[nums[i]]++
		if countMap[nums[i]] > n/2 {
			return nums[i]
		}
	}
	return -1
}

// 分冶法
// 分冶法的时间复杂度为O(nlogn)，空间复杂度为O(logn)。
func MajorityElement3(nums []int) int {
	return majorityElementRecursive(nums, 0, len(nums)-1)
}
func majorityElementRecursive(nums []int, left, right int) int {
	// 当只有一个元素时，直接返回该元素
	if left == right {
		return nums[left]
	}
	// 分治：
	mid := left + (right-left)/2
	// 查找左半部分的众数
	leftMajority := majorityElementRecursive(nums, left, mid)
	// 查找右半部分的众数
	rightMajority := majorityElementRecursive(nums, mid+1, right)

	// 合并结果

	// 如果 leftMajority 和 rightMajority 相同，则返回其中一个（因为它们是相同的众数）。
	// 如果 leftMajority 和 rightMajority 不同，则分别统计它们在整个数组中的出现次数，返回出现次数较多的那个。
	if leftMajority == rightMajority {
		return leftMajority
	}
	leftCount := 0
	rightCount := 0
	// 循环遍历当前范围
	// 通过这种方式，我们可以统计出两个候选众数在当前范围内的真实出现次数
	for i := left; i <= right; i++ {
		if nums[i] == leftMajority {
			leftCount++
		} else if nums[i] == rightMajority {
			rightCount++
		}
	}
	if leftCount > rightCount {
		return leftMajority
	}
	return rightMajority
}

// 摩尔投票法
// 摩尔投票法的时间复杂度为O(n)，空间复杂度为O(1)。
// 摩尔投票法的时间复杂度为O(n)，空间复杂度为O(1)。
// 摩尔投票法的核心思想是通过不断消除不同的元素，最终剩下的元素就是出现次数超过一半的元素。
// 具体步骤如下：
// 1. 初始化候选元素candidate和计数器count为0。
// 2. 遍历数组中的每个元素num：
// 2.1 如果计数器count为0，将当前元素num赋值给candidate。
// 2.2 如果当前元素num与candidate相同，计数器count加1。
// 2.3 如果当前元素num与candidate不同，计数器count减1。
// 3. 遍历结束后，candidate就是出现次数超过一半的元素。
// 4. 由于题目保证存在多数元素，所以最终的candidate就是答案。
func MajorityElement4(nums []int) int {
	candidate := 0 // 候选多数元素
	count := 0     // 计数器

	for _, num := range nums {
		if count == 0 {
			// 当计数器归零时，重新选择候选元素
			candidate = num
		}
		// 根据当前元素是否与候选元素相同，更新计数器
		if num == candidate {
			count++
		} else {
			count--
		}
	}
	return candidate // 题目保证存在多数元素，无需验证
}
