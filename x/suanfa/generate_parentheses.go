package suanfa

import "fmt"

// 22 括号生成

// 1. 数学归纳法
// 2. 递归（深度优先）
// - 1. 定义子问题，子树包信PQ节点
// - 2. 子问题的解
// - 3. 合并子问题的解
// - 4. 递归终止条件

// 长度为2n的合法括号组合，可以分解为以下子问题：
// - 1. 长度为2n-1的合法括号组合，加上一个左括号
// - 2. 长度为2n-1的合法括号组合，加上一个右括号
// - 3. 递归次数为2n次，每次递归的时间复杂度为O(1)，所以总的时间复杂度为O(2^n)
// - 4. 递归终止条件：n=0时，返回空字符串
// 2. 递归实现
// 剪枝：
// - 1. 左括号的数量小于n
// - 2. 右括号的数量小于左括号的数量
// - 3. 左右括号的数量都等于n
// - 4. 递归终止条件：左右括号的数量都等于n
// - 5. 递归终止条件：左括号的数量小于n
// - 6. 递归终止条件：右括号的数量小于左括号的数量
// - 7. 递归终止条件：左右括号的数量都等于n
// 局部不合法，不再递归
// 时间复杂度：O(2^n)，空间复杂度：O(n)
// 递归实现的时间复杂度为O(2^n)，空间复杂度为O(n)
// https://google.github.io/styleguide/go/

// 时间复杂度：O(2^n)，空间复杂度：O(n)
// 每个有效组合需要O(2n)的时间来生成，并且有2^n个有效组合，所以总的时间复杂度为O(2^n)。

// left < 3
// t
// left =0 right = 0 path = ""
// left =1 right = 0 path = "("
// left =2 right = 0 path = "(("
// left =3 right = 0 path = "((("
// left =4 right = 0 path = "((()"
// left =4 right = 1 path = "((())"
// left =4 right = 2 path = "((()))"
// left =4 right = 3 path = "((()))"
// left =4 right = 4 path = "((()))"
// left =3 right = 1 path = "(()"

func GenerateParenthesis(n int) []string {
	var result []string
	// left 表示已使用的左括号数量
	// right 表示已使用的右括号数量
	// path 表示当前的括号组合
	var dfs func(left, right int, path string)
	dfs = func(left, right int, path string) {
		// 递归终止条件
		// 左右括号的数量都等于n
		if left == n && right == n {
			result = append(result, path)
			return
		}
		// 左括号未用完，继续使用左括号

		if left < n {
			dfs(left+1, right, path+"(")
		}
		// 右括号未用完，且右括号数量小于左括号数量，继续使用右括号
		// 右括号数量小于左括号数量，才能保证括号的合法性
		if right < left {
			// if right < n && right < left {
			dfs(left, right+1, path+")")
		}
	}
	dfs(0, 0, "")
	return result
}

func GenerateParenthesis1(n int) []string {
	var result []string
	var dfs func(left, right int, path string)
	dfs = func(left, right int, path string) {
		if left == n && right == n {
			fmt.Printf("dfs(%d, %d, \"%s\")返回，添加到结果集\n", left, right, path)
			result = append(result, path)
			return
		}

		if left < n {
			fmt.Printf("dfs(%d+1, %d, \"%s\"+\"(\")\n", left, right, path)
			dfs(left+1, right, path+"(")
		}
		if right < left {
			fmt.Printf("dfs(%d, %d+1, \"%s\"+\")\")\n", left, right, path)
			dfs(left, right+1, path+")")
		}
	}
	dfs(0, 0, "")
	return result
}
