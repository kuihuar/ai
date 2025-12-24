// reverse a linked list
// swap nodes in pairs
// linked list cycle

// reverse nodes in k-group
// remove nth node from end of list
// merge two sorted lists

package suanfa

// Input : 1->2->3->4->5->NULL
// Output : 5->4->3->2->1->NULL
type ListNode struct {
	Val  int
	Next *ListNode
}

// 把每一节点的next指向它的前一个节点，这样就完成了反转。
// 我们同时也需要一个指针来保存当前节点的前一个节点。
// 最后，我们返回新的头引用，即原链表的最后一个节点。
//
// 反转流程
// 步骤一：保存 curr.Next 到 next，防止后续操作丢失链表信息。

// 步骤二：将 curr.Next 指向 prev，反转当前节点的指针方向。

// 步骤三：prev 和 curr 向前移动，处理下一个节点。

// 步骤四：重复步骤一到步骤三，直到 curr 为 nil，遍历完整个链表。

// 步骤五：返回 prev，即反转后的链表头节点。

// 在Go语言中，多重赋值的右侧表达式会在赋值前先全部求值，然后再依次赋值给左侧变量。
// 原代码的一行赋值等价于四行拆分，因为需要保留原始curr.Next的值，避免修改指针后丢失。
func ReverseLinkedList(head *ListNode) *ListNode {

	var prev *ListNode
	curr := head
	for curr != nil {
		// 1. 保存当前节点的下一个节点
		next := curr.Next
		// 2. 将当前节点的next指向它的前一个节点
		curr.Next = prev
		// 3. 将当前节点的前一个节点指向当前节点
		prev = curr
		// 4. 将当前节点的下一个节点作为当前节点
		curr = next

		// 5. 或者可以这样写
		curr.Next, prev, curr = prev, curr, curr.Next

	}
	return prev
}

// 创建dummy的原因，是最后需要返回dummy.Next，因为dummy是一个虚拟头节点，它的Next才是真正的头节点。
// 这样做的好处是，我们可以统一处理头节点的交换，而不需要单独处理头节点的交换。
// 创建 prev 的原因，是为了连接交换后的节点。
// 这样做的好处是，是为了循环和移动节点的方便，因为我们需要知道当前节点的前一个节点，以便连接交换后的节点。
func SwapPairs(head *ListNode) *ListNode {
	// 创建虚拟头节点， 简楷头节点交换操作
	// 1. 创建一个虚拟头节点 dummy，将其 Next 指向 head，即 dummy.Next = head。
	// 2. 创建一个指针 prev，初始化为 dummy。
	// 3. 进入循环，条件为 prev.Next != nil && prev.Next.Next != nil。
	// 4. 保存第一个节点 first 和第二个节点 second，以及 second 的下一个节点 next。
	// 5. 将 prev.Next 指向 second，将 second.Next 指向 first，将 first.Next 指向 next。
	// 6. 将 prev 移动到 first，继续下一轮交换。
	// 7. 返回 dummy.Next，即交换后的链表头节点。
	dummy := &ListNode{Next: head}
	// 前驱节点，用于连接交换后的节点
	//dummy(prev)->1->2->3->4->5
	prev := dummy

	for prev.Next != nil && prev.Next.Next != nil {

		// 保存第一个节点和第二个节点

		current := prev.Next
		nextNode := current.Next

		// 交换节点（规律，左值都是.Next）
		prev.Next = nextNode         // 前驱节点连接到第二个节点
		current.Next = nextNode.Next // 第一个节点连接到第三个节点
		nextNode.Next = current      // 第二个节点连接到第一个节点

		// 一行代码交换节点
		prev.Next, current.Next, nextNode.Next = nextNode, nextNode.Next, current

		// 移动到下一对节点
		prev = current

	}
	return dummy.Next

	// round 1:
	// 1. 保存第一个节点 current = prev.Next = 1
	// 2. 保存第二个节点 nextNode = current.Next = 2
	// 3. 交换节点
	//    prev.Next = nextNode
	//    current.Next = nextNode.Next = 3
	//    nextNode.Next = current = 1
	// 4. 移动到下一对节点
	//    prev = current = 1
	// round 2:
	// 1. 保存第一个节点 current = prev.Next = 3
	// 2. 保存第二个节点 nextNode = current.Next = 4
	// 3. 交换节点
	//    prev.Next = nextNode
	//    current.Next = nextNode.Next = 5
	//    nextNode.Next = current = 3
	// 4. 移动到下一对节点
	//    prev = current = 3
	// round 3: prev.Next == 5 && prev.Next.Next == nil , 退出循环
}

func SwapPairs1(head *ListNode) *ListNode {
	dummy := &ListNode{Next: head}
	prev := dummy

	for prev.Next != nil && prev.Next.Next != nil {

		// 保存第一个节点和第二个节点
		current := prev.Next
		nextNode := current.Next

		// 交换节点（规律，左值都是.Next）
		// prev.Next = nextNode         // 前驱节点连接到第二个节点
		// current.Next = nextNode.Next // 第一个节点连接到第三个节点
		// nextNode.Next = current      // 第二个节点连接到第一个节点

		prev.Next, current.Next, nextNode.Next = nextNode, nextNode.Next, current
		// 移动到下一对节点
		prev = current

	}
	return dummy.Next
}

// 哈希表解法
// 1. 创建一个哈希表 seen，用于存储已经访问过的节点。
// 2. 遍历链表，每次访问一个节点时，检查该节点是否已经在哈希表中。
// 3. 如果该节点已经在哈希表中，说明链表有环，返回 true。
// 4. 如果该节点不在哈希表中，将该节点添加到哈希表中，并继续遍历下一个节点。
// 5. 如果遍历完整个链表，都没有发现环，返回 false。
// 时间复杂度：O(n)，其中 n 是链表的长度。
// 空间复杂度：O(n)。
func HasCycle(head *ListNode) bool {

	seen := make(map[*ListNode]struct{})

	for head != nil {
		if _, ok := seen[head]; ok {
			return true
		}
		seen[head] = struct{}{}
		head = head.Next
	}
	return false
}

// 快慢指针解法
// 1. 创建两个指针，slow 和 fast，初始时都指向链表的头节点。
// 2. 进入循环，每次循环中，slow 指针移动一步，fast 指针移动两步。
// 3. 如果 fast 指针到达链表的末尾（即 fast == nil 或 fast.Next == nil），说明链表没有环，返回 false。
// 4. 如果 slow 指针和 fast 指针相遇（即 slow == fast），说明链表有环，返回 true。
// 5. 如果循环结束后都没有返回 true，说明链表没有环，返回 false。
// 时间复杂度：O(n)，其中 n 是链表的长度。
// 空间复杂度：O(1)。

// 循环条件 fast != nil && fast.Next != nil 严格限制 fast 的移动，确保 fast.Next.Next 不会触发空指针异常。
// 这样做的好处是，我们可以安全地访问 fast.Next.Next，而不需要担心空指针异常。
func HasCycleSlowFastPointer(head *ListNode) bool {
	slow, fast := head, head
	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
		if slow == fast {
			return true
		}
	}
	return false
}
