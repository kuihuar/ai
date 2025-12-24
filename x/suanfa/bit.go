package suanfa

// 191 位1的个数
// 编写一个函数，输入是一个无符号整数（以二进制串的形式），返回其二进制表达式中数字位数为 '1' 的个数（也被称为汉明重量）。
// 提示：

// 1. x %2 == 0 为偶数， x %2 == 1 为奇数 (32位)
// 2. x >> 1 右移一位，相当于 x / 2
// 3. x & 1 == 1 为奇数， x & 1 == 0 为偶数
// 4. x & (x - 1) 清零最低位的 1 （最佳1的个数）

func HammingWeight(num uint32) int {
	var count int
	for num > 0 {
		if num&1 == 1 {
			count++
		}
		num = num >> 1
	}
	return count
}
func HammingWeight1(num uint32) int {
	var count int
	for num > 0 {
		count++
		num = num & (num - 1)
	}
	return count

}

// 331 是否是2的倍数 power of two
// 1. mod 2 取余
// 2. 二进制位运算
// 特点，只有一个1，其他位都是0
// 00000001
// 00000010
// 判断方法 x & (x - 1) == 0 为 true 为2的倍数 且 x > 0
func IsPowerOfTwo(n int) bool {
	// 给最后一位清零，如果是0，说明是2的倍数
	return n > 0 && n&(n-1) == 0
}

// 338 比特位计数
// 给你一个整数 n ，对于 0 <= i <= n 中的每个 i ，计算其二进制表示中 1 的个数 ，
// 返回一个长度为 n + 1 的数组 ans 作为答案。
// 1.  for i := 1; i <= n; i++ {
// 2.  ans[i] = ans[i&(i-1)] + 1
// 3.  }
// 4.  return ans
// 5.  }

func CountBits(n int) []int {
	ans := make([]int, n+1)
	for i := 1; i <= n; i++ {
		ans[i] = ans[i&(i-1)] + 1
	}
	return ans
}

// for i :=0; i <= n; i++ {
// 分别计算每个数字的二进制表示中1的个数，然后生数组
