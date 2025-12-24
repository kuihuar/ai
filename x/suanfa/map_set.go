package suanfa

import (
	"fmt"
	"sort"
	"unsafe"
)

// 排序实现
// 1. 将字符串转换为字符数组
// 2. 对字符数组进行排序
// 3. 比较排序后的字符数组是否相同
// 排序实现的时间复杂度：O(nlogn)，空间复杂度：O(n)
// 其中n是字符串的长度。排序的时间复杂度为O(nlogn)，比较的时间复杂度为O(n)。总时间复杂度为O(nlogn+n)=O(nlogn)。
// 空间复杂度为O(n)，因为我们需要一个额外的字符数组来存储排序后的字符。
func IsAnagram(s string, t string) bool {
	runes := []rune(s)
	runet := []rune(t)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})

	sort.Slice(runet, func(i, j int) bool {
		return runet[i] < runet[j]
	})
	// 遍历比较
	// for i := 0; i < len(runes); i++ {
	// 	if runes[i] != runet[i] {
	// 		return false
	// 	}
	// }
	// return true
	// 转换比较
	// return string(runes) == string(runet)
	// 内存块比较
	// 比较内存块是否相同，相同则为true，不同则为false
	// 注意：这种方法只适用于Go语言的字符串类型，因为Go语言的字符串类型是不可变的，所以内存块是固定的
	//  利用 unsafe 包强制转换内存布局，将 []rune 切片直接映射为 string 类型的内存结构进行比较。

	// (1) unsafe.Pointer(&runes) 的作用
	//     &runes：获取 []rune 类型变量 runes 的内存地址（即指向 runes 切片的指针，类型为 *[]rune）
	//     unsafe.Pointer(&runes)：将 *[]rune 类型的指针强制转换为 unsafe.Pointer 类型，这是 Go 允许的跨类型指针转换的桥梁
	// (2) (*string)(...) 的类型转换
	//     (*string)(unsafe.Pointer(&runes))：将 unsafe.Pointer 类型的指针转换为 *string 类型（即指向字符串的指针）
	//     *string：表示指向字符串的指针类型
	//     这一步的本质是内存布局的强制映射：将原本表示 []rune 切片的内存结构，强行解释为 string 类型的内存结构
	// (3) *(*string)(...) 的取值操作
	//     * 解引用的必要性
	//     只有解引用后，才能比较两个字符串的实际内容，而非指针地址。
	//     *(*string)(unsafe.Pointer(&runes))：解引用 *string 类型的指针，获取其指向的 string 值
	//     这一步的目的是将 *string 类型的指针转换为 string 类型的值，以便进行后续的字符串比较
	// 总的来说，这段代码的作用是将 []rune 类型的切片 runes 转换为 string 类型的值，然后比较这两个 string 值是否相等。
	// 这种转换方式是通过 unsafe 包提供的跨类型指针转换功能来实现的，虽然在 Go 语言中不推荐直接操作内存，但在某些特定的场景下，
	// 这种方式可以提供更高效的内存操作方式。
	return *(*string)(unsafe.Pointer(&runes)) == *(*string)(unsafe.Pointer(&runet))
}

// 哈希表实现
// 1. 遍历字符串，将每个字符作为键，出现的次数作为值，存入哈希表
// 2. 遍历另一个字符串，将每个字符作为键，出现的次数作为值，存入哈希表
// 3. 比较两个哈希表是否相同
// 哈希表实现的时间复杂度：O(n)，
// 其中n是字符串的长度。遍历字符串的时间复杂度为O(n)，比较哈希表的时间复杂度为O(n)。
// 总时间复杂度为O(n+n)=O(n)。
// 空间复杂度：O(n)
// 其中n是字符串的长度。哈希表的空间复杂度为O(n)。

func IsAnagram2(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}
	// 定义哈希表
	cnt := make(map[rune]int)
	for _, char := range s {
		cnt[char]++
	}
	for _, char := range t {
		cnt[char]--
		if cnt[char] < 0 {
			return false
		}
	}
	return true
}

// 数组实现
// 1. 定义一个长度为26的数组，用于存储每个字符出现的次数
// 2. 遍历字符串，将每个字符作为索引，出现的次数作为值，存入数组
// 3. 遍历另一个字符串，将每个字符作为索引，出现的次数作为值，存入数组
// 4. 比较两个数组是否相同
// 注意：这种方法只适用于小写字母，如果字符串中包含大写字母，则需要定义一个长度为52的数组
// 时间复杂度：O(n)，
// 其中n是字符串的长度，S是字符集的大小。在这个例子中，字符集是小写字母，所以S=26。
// 数组的空间复杂度为O(S)。
func IsAnagram3(s string, t string) bool {
	var c1, c2 [26]int
	for _, ch := range s {
		c1[ch-'a']++
	}
	for _, ch := range s {
		c2[ch-'a']++
	}
	return c1 == c2
}

// 两数之和也是这里的应用，题号1

// 15. threeSum

// 暴力枚举
// 时间复杂度：O(N^3)，其中 N 是数组中的元素数量。最坏情况下数组中任意三个数都要被匹配一次。
// 空间复杂度：O(1)。
func ThreeSum(nums []int) [][]int {
	// 定义一个二维数组，用于存储结果
	var result [][]int
	sort.Ints(nums)

	n := len(nums)
	fmt.Println("nums.len: ", n)
	unique := make(map[[3]int]struct{}, 0)
	for i, x := range nums {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if x+nums[j]+nums[k] == 0 {
					if _, ok := unique[[3]int{x, nums[j], nums[k]}]; ok {
						// fmt.Printf("unique: %v, ok: %t \n", unique, ok)
						continue
					}
					unique[[3]int{x, nums[j], nums[k]}] = struct{}{}
					result = append(result, []int{x, nums[j], nums[k]})
				}
			}
		}
	}
	return result
}

// 两层循环
// 时间复杂度：O(N^2)，其中 N 是数组中的元素数量。最坏情况下数组中任意两个数都要被匹配一次。
// 空间复杂度：O(n)。
func ThreeSum1(nums []int) [][]int {
	// 定义一个二维数组，用于存储结果
	var result [][]int
	sort.Ints(nums)

	n := len(nums)
	fmt.Println("nums.len: ", n)
	seen := make(map[[3]int]struct{}, 0)
	for i := 0; i < n; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		minus := make(map[int]struct{}, 0)
		for j := i + 1; j < n; j++ {
			complement := -nums[i] - nums[j]

			if _, ok := minus[complement]; ok {
				triplet := [3]int{nums[i], complement, nums[j]}
				if _, ok := seen[triplet]; !ok {
					result = append(result, triplet[:])
					seen[triplet] = struct{}{}
				}
			}
			minus[nums[j]] = struct{}{}
		}
	}
	return result
}

// 双指针
// 时间复杂度：O(N^2)，其中 N 是数组中的元素数量。最坏情况下数组中任意两个数都要被匹配一次。
// 空间复杂度：O(n)。
func ThreeSum2(nums []int) [][]int {
	// 定义一个二维数组，用于存储结果
	var result [][]int
	sort.Ints(nums)
	n := len(nums)

	// 遍历数组，固定第一个数字，然后使用双指针法找到另外两个数字
	// 注意：为了避免重复，需要跳过相同的数字
	// 时间复杂度：O(N^2)，其中 N 是数组中的元素数量。最坏情况下数组中任意两个数都要被匹配一次。
	// 空间复杂度：O(n)。
	// 双指针法的时间复杂度：O(N)，其中 N 是数组中的元素数量。在最坏情况下，左右指针分别移动了 N 次。
	// 双指针法的空间复杂度：O(1)。
	// 总的时间复杂度：O(N^2)，其中 N 是数组中的元素数量。
	// 总的空间复杂度：O(n)。
	for i := 0; i < n-2; i++ {
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		left, right := i+1, n-1
		for left < right {
			sum := nums[i] + nums[left] + nums[right]
			if sum == 0 {
				result = append(result, []int{nums[i], nums[left], nums[right]})

				for left < right && nums[left] == nums[left+1] {
					left++
				}
				for left < right && nums[right] == nums[right-1] {
					right--
				}

				left++
				right--

			} else if sum < 0 {
				left++
			} else {
				right--
			}
		}
	}
	return result
}

func ThreeSum3(nums []int) [][]int {
	// 定义一个二维数组，用于存储结果
	var result [][]int
	sort.Ints(nums)
	seen := make(map[[3]int]struct{}, 0)
	// 外层循环遍历到倒数第二个元素（len(nums)-2），因为需要三个元素
	for i := 0; i < len(nums)-2; i++ {
		// 如果当前元素和前一个元素相同，跳过，避免重复计算
		if i > 0 && nums[i] == nums[i-1] {
			continue
		}
		//使用一个映射 complements 来存储当前外层循环的元素（-nums[i]-x）。
		complements := make(map[int]struct{})
		// 内层循环从 i+1 开始遍历，寻找和为 -nums[i] 的元素。
		for j := i + 1; j < len(nums); j++ {
			x := nums[j]
			if _, exists := complements[-nums[i]-x]; exists {
				// 如果存在，将这三个元素添加到结果数组中。
				triplet := [3]int{nums[i], -nums[i] - x, x}
				if _, ok := seen[triplet]; !ok {
					result = append(result, triplet[:])
					seen[triplet] = struct{}{}
				}
			} else {
				complements[x] = struct{}{}
			}
		}
	}
	return result
}
