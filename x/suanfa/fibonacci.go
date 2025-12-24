package suanfa

import "fmt"

// 1,1,2,3,5,7,13,21
// f(n) = f(n-1) + f(n-2)
// O(2^n)
func FibonacciRecursive(n int) int {
	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-1) + FibonacciRecursive(n-2)
}

func PrintFibonacci(n int) {

	for i := 0; i < n; i++ {
		fmt.Printf("PrintFibonacci(%d)= %d \n", i, FibonacciRecursive(i))
	}
}

// O(n)
func FibonacciMemo(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}

	if val, ok := memo[n]; ok {
		return val
	}
	memo[n] = FibonacciMemo(n-1, memo) + FibonacciMemo(n-2, memo)
	return memo[n]
}

func PrintFibonacciMemo(n int) {
	memo := make(map[int]int)

	for i := 0; i <= n; i++ {

		val := FibonacciMemo(i, memo)
		fmt.Printf("PrintFibonacciMemo(%d)= %d \n", i, val)
	}
}

// 动态规划
func FibonacciIterative(n int) []int {

	series := make([]int, n+1)

	series[0] = 0
	series[1] = 1
	for i := 2; i <= n; i++ {
		series[i] = series[i-1] + series[i-2]
	}
	return series
}

func PrintFibonacciIterative(n int) {
	series := FibonacciIterative(n)

	for i := 0; i <= n; i++ {
		fmt.Printf("PrintFibonacciIterative(%d)= %d \n", i, series[i])
	}
}

// 迭代法
// O(n)
func FibonacciNth(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1

	for i := 2; i <= n; i++ {
		a, b = b, a+b
	}
	return b
}
func PrintFibonacciNth(n int) {
	res := FibonacciNth(n)
	fmt.Printf("FibonacciNth(%d)= %d \n", n, res)

}

// 动态规划（自下而上）
func fibonacciDP(n int) int {
	if n <= 1 {
		return n
	}
	fib := make([]int, n+1)
	fib[0], fib[1] = 0, 1
	for i := 2; i <= n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
	return fib[n]
}

// 动态规划（自上而下，带备忘录）
var memo = make(map[int]int)

func fibonacciMemo(n int) int {
	if n <= 1 {
		return n
	}
	if val, found := memo[n]; found {
		return val
	}
	memo[n] = fibonacciMemo(n-1) + fibonacciMemo(n-2)
	return memo[n]
}
