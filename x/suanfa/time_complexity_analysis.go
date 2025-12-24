package suanfa

import (
	"fmt"
	"math"
)

func timeComplexity() {
	n := 1000

	// O(1)
	fmt.Printf("O(1) is : %d \n", n)
	fmt.Println("-----------")
	// O(1)
	fmt.Printf("O(1) is : %d \n", n)
	fmt.Printf("O(1) is : %d \n", n)
	fmt.Println("-----------")
	// O(n)
	for i := 1; i <= n; i++ {
		fmt.Printf("O(n)  i: %d \n", i)
	}
	fmt.Println("-----------")
	// O(n^2)
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			fmt.Printf("O(n^2) i: %d and j: %d \n", i, j)
		}
	}
	fmt.Println("-----------")
	// O(log(n)), 对数，终止条件
	for i := 1; i <= n; i *= 2 {
		fmt.Printf("O(log(n)) i: %d \n", i)
	}
	fmt.Println("-----------")
	// O(k^n)
	for j := 1.0; j <= math.Pow(2, float64(n)); j += 1 {
		fmt.Printf("O(k^n) j: %f \n", j)
	}
	fmt.Println("-----------")

	// O(n!)
	for i := 1; i <= FibonacciRecursive(n); i++ {
		fmt.Printf("O(n!)  i: %d \n", i)
	}
	fmt.Println("-----------")
}
