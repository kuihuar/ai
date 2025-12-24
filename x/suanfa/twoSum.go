package suanfa

import "fmt"

// 暴力枚举
// 时间复杂度：O(N^2)，其中 N 是数组中的元素数量。最坏情况下数组中任意两个数都要被匹配一次。
// 空间复杂度：O(1)。
func TwoSum(nums []int, target int) []int {
	for i, x := range nums {
		fmt.Printf("i: %d and x: %d \n", i, x)
		for j := i + 1; j < len(nums); j++ {
			fmt.Printf("j: %d and nums[j]: %d \n", j, nums[j])
			if x+nums[j] == target {
				//return []int{i, j}
			}
		}
		fmt.Println("======= end =======")
	}
	return []int{}
}

// 哈希表
// 时间复杂度：O(N)，其中 N 是数组中的元素数量。对于每一个元素 x，我们可以 O(1) 地寻找 target - x。
// 空间复杂度：O(N)，其中 N 是数组中的元素数量。主要为哈希表的开销。
// suanfa.TwoSumHashTable([]int{2, 4, 7, 11, 15}, 15)
func TwoSumHashTable(nums []int, target int) []int {
	hashTable := make(map[int]int)

	for i, x := range nums {
		fmt.Printf("i: %d and x: %d \n", i, x)
		if p, ok := hashTable[target-x]; ok {
			fmt.Printf("hashTable: %v, p: %d, ok :%t \n", hashTable, p, ok)
			return []int{p, i}
		}
		hashTable[x] = i

		fmt.Printf("hashTable: %v \n", hashTable)
	}
	return []int{}
	// hashTable := map[int]int{}
	// for i, x := range nums {
	// 	if p, ok := hashTable[target-x]; ok {
	// 		return []int{p, i}
	// 	}
	// 	hashTable[x] = i
	// }
	// return []int{}
}

func TwoSum1(nums []int, target int) []int {
	for i := 0; i < len(nums)-1; i++ {
		fmt.Printf("i: %d and nums[i]: %d \n", i, nums[i])
		for j := i + 1; j < len(nums); j++ {
			fmt.Printf("j: %d and nums[j]: %d \n", j, nums[j])
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
		fmt.Println("======= end =======")
	}
	return []int{}
}
