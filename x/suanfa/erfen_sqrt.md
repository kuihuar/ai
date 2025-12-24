// 二分查找

1. sorted 单调递增或递减
2. bounded 存在上下界
3. accessible by index 索引访问

left, right :=0, len(nums) - 1
// 这里必须是小于等于，因为当left==right时，区间[left, right]依然有效
for left <= right {
	mid := (left + right) / 2
	if nums[mid] == target {
		return mid
        // 去数组的右边查找
	} else if nums[mid] < target {
		left = mid + 1
	} else {
        // 去数组的左边查找
		right = mid - 1
	}
}