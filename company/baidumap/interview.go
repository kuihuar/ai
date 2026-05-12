package baidumap

import (
	"net"
	"os"
)

func loopDefer() {

	files := []string{"a.txt", "b.txt", "c.txt"}

	for _, file := range files {
		err := func(fileName string) error {

			f, err := os.Open(fileName)
			if err != nil {
				return err
			}
			defer f.Close()
			return nil
		}(file)
		if err != nil {
			panic(err)
		}
	}

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			panic(err)
		}

		f.Close()
	}
}

// 正确写法（匿名函数 + 立即执行）：

func loopDeferCorrect() {

	for i := 0; i < 10; i++ {
		go func(idx int) {
			conn, err := net.Dial("tcp", ":8080")
			if err != nil {
				return

			}
			defer conn.Close()
		}(i)
	}
}
