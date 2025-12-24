package suanfa

// 69 x 的平方根
// 1. 二分查找
// 1e-6 表示 10 的 -6 次方，即 0.000001
// 1e6 表示 10 的 6 次方，即 1000000

func MySqrt(x int) int {
	if x == 0 || x == 1 {
		return x
	}
	left, right, res := 0, x, -1
	for left <= right {
		mid := left + (right-left)/2

		// 不用乘号，以防越界
		if mid == x/mid {
			return mid
		} else if mid < x/mid {
			res = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return res
}
func MySqrt2(x int) int {
	if x == 0 || x == 1 {
		return x
	}
	left, right, res := 0, x, -1
	for left <= right {
		mid := left + (right-left)/2
		if mid*mid <= x {
			res = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return res
}

func MySqrt1(x int) int {
	left, right := 1, x
	for left <= right {
		mid := left + (right-left)/2
		if mid*mid > x {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return right
}

// 2. 牛顿迭代法
// 3. 袖珍计算器算法
