package mianshiskill

import (
	"fmt"
	"time"
)

//设计一个带有超时机制的channel 操作， 当超过500毫秒时，返回超时错误

func work(result chan string) {
	result <- "done"
}

func HandleTimeoutChannel() (error, string) {

	resultChan := make(chan string)
	go work(resultChan)

	timeoutChan := time.After(500 * time.Millisecond)
	select {
	case result := <-resultChan:
		fmt.Println("result: ", result)
		return nil, result
	case <-timeoutChan:
		fmt.Println("timeout")
		return fmt.Errorf("timeout"), ""
	}
}
