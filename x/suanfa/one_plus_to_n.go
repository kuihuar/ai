package suanfa

import "fmt"

func onePlusToN(n int) {
	res := 0
	for i := 1; i <= n; i++ {
		res = res + i
	}
	fmt.Printf("res: %d \n", res)
}
func onePlusToN_1(n int) {
	res := n * (n + 1) / 2
	fmt.Printf("res: %d \n", res)
}
