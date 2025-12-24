作用：将数据处理拆分为多个阶段，通过 Channel 连接。

```go
// 生成数据
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// 平方处理
func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// 使用
for n := range sq(gen(2, 3)) {
    fmt.Println(n) // 输出: 4 9
}
```
文件处理 Pipeline
假设我们有一个需求，要读取一个文件，统计文件中每行的单词数，然后找出单词数最多的行。可以通过 Pipeline 模式实现这个过程
```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// readFile 读取文件内容，将每行数据发送到通道中
func readFile(filePath string) <-chan string {
    out := make(chan string)
    go func() {
        file, err := os.Open(filePath)
        if err != nil {
            fmt.Println("Error opening file:", err)
            close(out)
            return
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            out <- scanner.Text()
        }
        if err := scanner.Err(); err != nil {
            fmt.Println("Error reading file:", err)
        }
        close(out)
    }()
    return out
}

// countWords 统计每行的单词数，将结果发送到通道中
func countWords(in <-chan string) <-chan int {
    out := make(chan int)
    go func() {
        for line := range in {
            words := strings.Fields(line)
            out <- len(words)
        }
        close(out)
    }()
    return out
}

// findMax 找出单词数最多的行
func findMax(in <-chan int) int {
    max := 0
    for num := range in {
        if num > max {
            max = num
        }
    }
    return max
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: go run main.go <file_path>")
        return
    }
    filePath := os.Args[1]

    // 构建 Pipeline
    lines := readFile(filePath)
    wordCounts := countWords(lines)
    maxWordCount := findMax(wordCounts)

    fmt.Println("单词数最多的行的单词数:", maxWordCount)
}
```