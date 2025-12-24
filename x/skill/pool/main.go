package main

import (
	"fmt"

	"github.com/bytedance/gopkg/util/gopool"
)

func main() {
	// go func() {

	// }()

	gopool.Go(func() {
		fmt.Println("ln")
	})

	// pool := gopool.NewPool("test", 10, nil)

}
