package suanfa

// []rune{}：用于需要 字符级操作 的场景（尤其是多语言文本）。

// []byte{}：用于需要 字节级操作 的场景（如二进制数据或纯 ASCII 文本）。

func isValidParentheses(s string) bool {
	stack := make([]byte, 0)
	// stack := []byte{}

	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	for _, char := range s {
		// 当前是个左括号，入栈
		if _, ok := pairs[byte(char)]; !ok {
			stack = append(stack, byte(char))
		} else {
			//栈顶元素与当前元素不匹配，返回false；
			if len(stack) != 0 && pairs[byte(char)] != stack[len(stack)-1] {
				return false
			} else {
				//栈顶元素与当前元素匹配，出栈
				stack = stack[:len(stack)-1]
			}
		}
	}
	// 栈为空，返回true；
	// 栈不为空，返回false；
	return len(stack) == 0
}

func isValidParenthesesRune(s string) bool {
	stack := make([]rune, 0)
	// stack := []byte{}

	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}
	for _, char := range s {
		if _, ok := pairs[char]; !ok { // 当前是个左括号，入栈
			stack = append(stack, char)
		} else { // 当前是个右括号
			//栈顶元素与当前元素不匹配，返回false；或者栈为空
			if len(stack) == 0 || pairs[char] != stack[len(stack)-1] {
				return false
			} else {
				//栈不空，并且栈顶元素与当前元素匹配，出栈
				stack = stack[:len(stack)-1]
			}
		}
	}

	// 栈为空，返回true；
	// 栈不为空，返回false；
	return len(stack) == 0
}

func isValidParenthesesRune1(s string) bool {
	stack := make([]byte, 0)
	// stack := []byte{}

	pairs := map[byte]byte{
		')': '(',
		']': '[',
		'}': '{',
	}
	for i := 0; i < len(s); i++ {
		char := s[i]
		if _, ok := pairs[char]; !ok { // 当前是个左括号，入栈
			stack = append(stack, char)
		} else { // 当前是个右括号
			//栈顶元素与当前元素不匹配，返回false；或者栈为空
			if len(stack) == 0 || pairs[char] != stack[len(stack)-1] {
				return false
			} else {
				//栈不空，并且栈顶元素与当前元素匹配，出栈
				stack = stack[:len(stack)-1]
			}
		}
	}

	// 栈为空，返回true；
	// 栈不为空，返回false；
	return len(stack) == 0
}

// 用栈实现队列
// 输入栈：用于存储输入的元素
// 输出栈：用于存储输出的元素
// 当输出栈为空时，将输入栈的元素全部出栈，入栈到输出栈
// 当输出栈不为空时，直接输出栈顶元素
// 当输出栈为空且输入栈为空时，返回true；否则返回false
// 时间复杂度：push、empty 为 O(1)，pop 和 peek 为均摊 O(1)。
// 对于每个元素，至多入栈和出栈各两次，故均摊复杂度为 O(1)。
// 空间复杂度：O(n)，其中 n 是队列的容量。
// 栈存储队列中的元素，需要 O(n) 的额外空间。
type ImplementQueueUsingStacks struct {
	InputStack  []int
	OutputStack []int
}

func Constructor() ImplementQueueUsingStacks {
	return ImplementQueueUsingStacks{
		InputStack:  make([]int, 0),
		OutputStack: make([]int, 0),
	}
}

func (this *ImplementQueueUsingStacks) Push(x int) {
	this.InputStack = append(this.InputStack, x)
}

func (this *ImplementQueueUsingStacks) Pop() int {

	// 输出栈为空，将输入栈的元素全部出栈，入栈到输出栈
	if len(this.OutputStack) == 0 {
		for len(this.InputStack) > 0 {
			this.OutputStack = append(this.OutputStack, this.InputStack[len(this.InputStack)-1])
			this.InputStack = this.InputStack[:len(this.InputStack)-1]
		}
	}
	res := this.OutputStack[len(this.OutputStack)-1]
	this.OutputStack = this.OutputStack[:len(this.OutputStack)-1]
	return res
}
func (this *ImplementQueueUsingStacks) Peek() int {
	// 输出栈为空，将输入栈的元素全部出栈，入栈到输出栈
	if len(this.OutputStack) == 0 {
		for len(this.InputStack) > 0 {
			this.OutputStack = append(this.OutputStack, this.InputStack[len(this.InputStack)-1])
			this.InputStack = this.InputStack[:len(this.InputStack)-1]
		}
	}
	return this.OutputStack[len(this.OutputStack)-1]
}
func (this *ImplementQueueUsingStacks) Empty() bool {
	// 输出栈为空且输入栈为空，返回true；否则返回false
	return len(this.OutputStack) == 0 && len(this.InputStack) == 0

}

// 用队列实现栈
// 存储队列：用于存储元素
// 辅助队列：用于辅助存储元素
// push: 将元素入辅助队列，将存储队列的元素全部出队，入队到辅助队列，将存储队列和辅助队列交换
// pop: 将存储队列的元素出队
// top: 返回存储队列的第一个元素
// empty: 判断存储队列是否为空
// 时间复杂度：push 和 empty 为 O(1)，pop 和 top 为 O(n)。
// 对于 pop 操作，需要将存储队列的元素全部出队，入队到辅助队列，时间复杂度为 O(n)，其中 n 是队列的容量。
// 空间复杂度：O(n)，其中 n 是队列的容量。
// 存储队列和辅助队列需要 O(n) 的额外空间。
type ImplementStackUsingQueues struct {
	StoreQueue  []int
	AssistQueue []int
}

func NewImplementStackUsingQueues() ImplementStackUsingQueues {
	return ImplementStackUsingQueues{
		StoreQueue:  make([]int, 0),
		AssistQueue: make([]int, 0),
	}
}

func (this *ImplementStackUsingQueues) Push(x int) {
	this.AssistQueue = append(this.AssistQueue, x)
	for len(this.StoreQueue) > 0 {
		this.AssistQueue = append(this.AssistQueue, this.StoreQueue[0])
		this.StoreQueue = this.StoreQueue[1:]
	}
	this.StoreQueue, this.AssistQueue = this.AssistQueue, this.StoreQueue
}

func (this *ImplementStackUsingQueues) Pop() int {
	res := this.StoreQueue[0]
	this.StoreQueue = this.StoreQueue[1:]
	return res
}
func (this *ImplementStackUsingQueues) Top() int {
	return this.StoreQueue[0]
}

func (this *ImplementStackUsingQueues) Empty() bool {
	return len(this.StoreQueue) == 0
}

// 使用一个队列实现栈
type ImplementStackUsingQueues1 struct {
	Queue []int
}

func NewImplementStackUsingQueues1() ImplementStackUsingQueues1 {
	return ImplementStackUsingQueues1{
		Queue: make([]int, 0),
	}
}
func (this *ImplementStackUsingQueues1) Push(x int) {

	n := len(this.Queue)

	this.Queue = append(this.Queue, x)
	for i := 0; i < n; i++ {
		this.Queue = append(this.Queue, this.Queue[0])
		this.Queue = this.Queue[1:]
	}
}

func (this *ImplementStackUsingQueues1) Pop() int {
	res := this.Queue[0]
	this.Queue = this.Queue[1:]
	return res
}
func (this *ImplementStackUsingQueues1) Top() int {
	return this.Queue[0]
}
func (this *ImplementStackUsingQueues1) Empty() bool {
	return len(this.Queue) == 0
}
