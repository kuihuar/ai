1. reverse a linked list
2. swap nodes in pairs
3. linked list cycle


#### 链表反转
1. 问题描述
给定一个链表，反转它并返回新链表的头节点。
示例：
- 输入：1 → 2 → 3 → nil
- 输出：3 → 2 → 1 → nil
2. 解题思路
- 方法一：迭代法
  - 步骤一：初始化三个指针：prev 指向 nil，curr 指向头节点，next 用于暂存 curr.Next。
  - 步骤二：遍历链表，直到 curr 为 nil。
  - 步骤三：反转当前节点的指针方向，将 curr.Next 指向 prev。
  - 步骤四：移动指针，prev 指向 curr，curr 指向 next。
  - 步骤五：重复步骤三到步骤四，直到 curr 为 nil。
  - 步骤六：返回 prev，即新链表的头节点。
- 方法二：递归法
  - 步骤一：递归终止条件：如果链表为空或只有一个节点，直接返回头节点。
  - 步骤二：递归反转当前节点的下一个节点，并将返回的新头节点赋值给 curr.Next。
  - 步骤三：将当前节点的 Next 指向 prev。
  - 步骤四：返回新头节点。
3. 代码实现
方法一：迭代法
```go
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode
    curr := head
    for curr != nil {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    return prev
}
```
方法二：
```go
func reverseList(head *ListNode) *ListNode {
    // 递归终止条件：空链表或单节点链表无需反转，直接返回
    if head == nil || head.Next == nil {
        return head
    }

    // 递归调用：先反转 head.Next 之后的子链表
    newHead := reverseList(head.Next)

    // 核心反转逻辑
    head.Next.Next = head  // ✅ 将后驱节点的 Next 指向自己（建立反向链接）
    head.Next = nil        // ✅ 断开原正向链接（避免成环）

    // 返回新链表的头节点（即原链表的尾节点）
    return newHead
} 

```
1. 反转流程
- 步骤一：保存 curr.Next 到 next，防止后续操作丢失链表信息。

- 步骤二：将 curr.Next 指向 prev，反转当前节点的指针方向。

- 步骤三：prev 和 curr 向前移动，处理下一个节点。

- 步骤四：重复步骤一到步骤三，直到 curr 为 nil，完成链表反转。
- 步骤五：返回 prev，即新链表的头节点。
2. 示例演示
以链表 1 → 2 → 3 → nil 为例：

|循环次数	|curr	|prev	|操作后链表状态|
|--|--|--|--|
|初始	|1	|nil	| |
|第一次	|2	|1	|1 → nil|
|第二次	|3	|2	| 2 → 1 → nil|
|结束	|nil	|3	|3 → 2 → 1 → nil|

最终返回 prev = 3，即新链表的头节点。



3. 复杂度分析

|指标|值|说明|
|--|--|--|
|时间复杂度|O(n)|遍历一次链表，n 为节点数|
|空间复杂度|O(1)|仅使用固定指针变量|


4. 总结
核心逻辑：通过三指针（prev、curr、next）逐步反转节点指针方向。


#### 递归过程分步演示
1->2->3->4->5->nil
##### 第一步：递归终止条件
- 当链表为空或单节点时，直接返回自身。

- 递归到最后一个节点 5 时，触发终止条件，返回 5。

##### 第二步：递归回退到节点 4
- 当前状态：head=4, newHead=5

- 操作：

    1. head.Next.Next = head → 4.Next.Next = 4（即 5.Next = 4，建立反向链接）

    2. head.Next = nil → 4.Next = nil（断开原正向链接）

- 结果：子链表变为 5→4，4 的 Next 断开。

##### 第三步：递归回退到节点 3
- 当前状态：head=3, newHead=5

- 操作：

    1. head.Next.Next = head → 3.Next.Next = 3（即 4.Next = 3）

    2. head.Next = nil → 3.Next = nil

- 结果：子链表变为 5→4→3，3 的 Next 断开。
##### 第四步：递归回退到节点 2
- 当前状态：head=2, newHead=5

- 操作：

    1. head.Next.Next = head → 2.Next.Next = 2（即 3.Next = 2）

    2. head.Next = nil → 2.Next = nil

- 结果：子链表变为 5→4→3→2，2 的 Next 断开。

##### 第五步：递归回退到节点 1
- 当前状态：head=1, newHead=5

- 操作：

    1. head.Next.Next = head → 1.Next.Next = 1（即 2.Next = 1）

    2. head.Next = nil → 1.Next = nil

- 最终结果：链表完全反转为 5→4→3→2→1。

##### 复杂度分析
- 时间复杂度：O(n)，遍历所有节点。

- 空间复杂度：O(n)，递归栈深度与链表长度相关。

##### 递归调用栈示意图
reverseList(1)                    → 最终返回 newHead=5
│
└─ reverseList(2)                 → 返回 newHead=5
   │
   └─ reverseList(3)              → 返回 newHead=5
      │
      └─ reverseList(4)           → 返回 newHead=5
         │
         └─ reverseList(5) → 触发终止，返回 5



##### 总结
递归法通过深度优先遍历到链表末尾，再回溯过程中逐步反转每个节点的指向。核心在于：

- 先处理子问题（反转 head.Next 后的链表）；

- 再处理当前节点（调整指针方向，避免成环）。

- 最终返回的 newHead 始终是原链表的尾节点，即反转后的头节点。         

---

|特性|	make(map[...]...)|	map[...]struct{}{}|
|--|--|--|
|初始化方式|	显式构造|	字面量语法|
|容量控制|	支持预分配（make(..., cap)）|	无法指定初始容量|
|编译器优化|	直接调用 makemap|	可能被优化为 makemap|
|代码可读性|	更明确的创建意图|	更简洁|


### 两两交换链表中的节点
### 判断链表是否有环