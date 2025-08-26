package algorithm

import (
	"strconv"
	"strings"
)

// ListNode 链表节点定义
type ListNode struct {
	Val  int
	Next *ListNode
}

// 链表相关操作

// ReverseList 链表反转 - 迭代方法
// 时间复杂度: O(n)
// 空间复杂度: O(1)
// 核心思想: 逐个节点反转，维护三个指针：前驱、当前、后继
func ReverseList(head *ListNode) *ListNode {
	var prev *ListNode = nil
	curr := head

	for curr != nil {
		// 保存下一个节点
		next := curr.Next
		// 反转当前节点的指针
		curr.Next = prev
		// 移动指针
		prev = curr
		curr = next
	}

	return prev
}

// ReverseListRecursive 链表反转 - 递归方法
// 时间复杂度: O(n)
// 空间复杂度: O(n) - 递归调用栈
// 核心思想: 先递归到末尾，再逐层返回时反转
func ReverseListRecursive(head *ListNode) *ListNode {
	// 基础情况：空链表或只有一个节点
	if head == nil || head.Next == nil {
		return head
	}

	// 递归反转剩余部分
	newHead := ReverseListRecursive(head.Next)

	// 反转当前节点
	head.Next.Next = head
	head.Next = nil

	return newHead
}

// ReverseBetween 反转链表指定区间 [left, right]
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func ReverseBetween(head *ListNode, left, right int) *ListNode {
	if head == nil || left == right {
		return head
	}

	// 创建虚拟头节点
	dummy := &ListNode{Next: head}
	prev := dummy

	// 移动到left位置的前一个节点
	for i := 0; i < left-1; i++ {
		prev = prev.Next
	}

	// 开始反转
	curr := prev.Next
	for i := 0; i < right-left; i++ {
		next := curr.Next
		curr.Next = next.Next
		next.Next = prev.Next
		prev.Next = next
	}

	return dummy.Next
}

// ReverseKGroup K个一组反转链表
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func ReverseKGroup(head *ListNode, k int) *ListNode {
	if head == nil || k == 1 {
		return head
	}

	// 检查是否有k个节点
	count := 0
	curr := head
	for curr != nil && count < k {
		curr = curr.Next
		count++
	}

	if count < k {
		return head // 不足k个，不反转
	}

	// 反转前k个节点
	prev := ReverseKGroup(curr, k) // 递归处理剩余部分
	curr = head

	for i := 0; i < k; i++ {
		next := curr.Next
		curr.Next = prev
		prev = curr
		curr = next
	}

	return prev
}

// 链表工具函数

// CreateList 根据数组创建链表
func CreateList(arr []int) *ListNode {
	if len(arr) == 0 {
		return nil
	}

	head := &ListNode{Val: arr[0]}
	curr := head

	for i := 1; i < len(arr); i++ {
		curr.Next = &ListNode{Val: arr[i]}
		curr = curr.Next
	}

	return head
}

// ListToArray 链表转数组
func ListToArray(head *ListNode) []int {
	var result []int
	curr := head

	for curr != nil {
		result = append(result, curr.Val)
		curr = curr.Next
	}

	return result
}

// PrintList 打印链表
func PrintList(head *ListNode) string {
	var result []string
	curr := head

	for curr != nil {
		result = append(result, strconv.Itoa(curr.Val))
		curr = curr.Next
	}

	return "[" + strings.Join(result, " -> ") + "]"
}

// GetListLength 获取链表长度
func GetListLength(head *ListNode) int {
	count := 0
	curr := head

	for curr != nil {
		count++
		curr = curr.Next
	}

	return count
}

// 链表检测操作

// HasCycle 检测链表是否有环 - 快慢指针法
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func HasCycle(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return false
	}

	slow := head
	fast := head.Next

	for slow != fast {
		if fast == nil || fast.Next == nil {
			return false
		}
		slow = slow.Next
		fast = fast.Next.Next
	}

	return true
}

// DetectCycle 检测环的起始节点
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func DetectCycle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return nil
	}

	// 第一步：找到相遇点
	slow := head
	fast := head

	for fast != nil && fast.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next

		if slow == fast {
			break
		}
	}

	// 没有环
	if fast == nil || fast.Next == nil {
		return nil
	}

	// 第二步：找到环的起始点
	slow = head
	for slow != fast {
		slow = slow.Next
		fast = fast.Next
	}

	return slow
}

// 链表合并操作

// MergeTwoLists 合并两个有序链表
// 时间复杂度: O(n + m)
// 空间复杂度: O(1)
func MergeTwoLists(l1, l2 *ListNode) *ListNode {
	dummy := &ListNode{}
	curr := dummy

	for l1 != nil && l2 != nil {
		if l1.Val <= l2.Val {
			curr.Next = l1
			l1 = l1.Next
		} else {
			curr.Next = l2
			l2 = l2.Next
		}
		curr = curr.Next
	}

	// 连接剩余节点
	if l1 != nil {
		curr.Next = l1
	}
	if l2 != nil {
		curr.Next = l2
	}

	return dummy.Next
}

// 链表查找操作

// FindMiddle 找到链表的中间节点
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func FindMiddle(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	slow := head
	fast := head

	for fast.Next != nil && fast.Next.Next != nil {
		slow = slow.Next
		fast = fast.Next.Next
	}

	return slow
}

// FindNthFromEnd 找到倒数第n个节点
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func FindNthFromEnd(head *ListNode, n int) *ListNode {
	if head == nil || n <= 0 {
		return nil
	}

	fast := head
	slow := head

	// 快指针先走n步
	for i := 0; i < n; i++ {
		if fast == nil {
			return nil // n大于链表长度
		}
		fast = fast.Next
	}

	// 快慢指针同时移动
	for fast != nil {
		slow = slow.Next
		fast = fast.Next
	}

	return slow
}

// 链表删除操作

// RemoveNthFromEnd 删除倒数第n个节点
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func RemoveNthFromEnd(head *ListNode, n int) *ListNode {
	dummy := &ListNode{Next: head}
	fast := dummy
	slow := dummy

	// 快指针先走n+1步
	for i := 0; i <= n; i++ {
		if fast == nil {
			return head // n大于链表长度
		}
		fast = fast.Next
	}

	// 快慢指针同时移动
	for fast != nil {
		slow = slow.Next
		fast = fast.Next
	}

	// 删除节点
	slow.Next = slow.Next.Next

	return dummy.Next
}

// DeleteDuplicates 删除排序链表中的重复元素
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func DeleteDuplicates(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}

	curr := head

	for curr.Next != nil {
		if curr.Val == curr.Next.Val {
			curr.Next = curr.Next.Next
		} else {
			curr = curr.Next
		}
	}

	return head
}

// 链表排序操作

// SortList 链表排序 - 归并排序
// 时间复杂度: O(n log n)
// 空间复杂度: O(log n) - 递归调用栈
func SortList(head *ListNode) *ListNode {
	if head == nil || head.Next == nil {
		return head
	}

	// 找到中间节点
	mid := FindMiddle(head)
	right := mid.Next
	mid.Next = nil

	// 递归排序左右两部分
	left := SortList(head)
	right = SortList(right)

	// 合并两个有序链表
	return MergeTwoLists(left, right)
}

// 链表回文检测

// IsPalindrome 检测链表是否为回文
// 时间复杂度: O(n)
// 空间复杂度: O(1)
func IsPalindrome(head *ListNode) bool {
	if head == nil || head.Next == nil {
		return true
	}

	// 找到中间节点
	mid := FindMiddle(head)

	// 反转后半部分
	second := ReverseList(mid.Next)
	mid.Next = nil

	// 比较前后两部分
	p1 := head
	p2 := second

	for p2 != nil {
		if p1.Val != p2.Val {
			return false
		}
		p1 = p1.Next
		p2 = p2.Next
	}

	// 恢复原链表结构（可选）
	mid.Next = ReverseList(second)

	return true
}

// 链表相交检测

// GetIntersectionNode 找到两个链表的相交节点
// 时间复杂度: O(n + m)
// 空间复杂度: O(1)
func GetIntersectionNode(headA, headB *ListNode) *ListNode {
	if headA == nil || headB == nil {
		return nil
	}

	pA := headA
	pB := headB

	// 当两个指针相遇时，就是相交节点
	// 如果没有相交，两个指针都会变成nil
	for pA != pB {
		if pA == nil {
			pA = headB
		} else {
			pA = pA.Next
		}

		if pB == nil {
			pB = headA
		} else {
			pB = pB.Next
		}
	}

	return pA
}
