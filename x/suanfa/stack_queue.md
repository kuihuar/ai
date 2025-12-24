1. Stack - Fist In Last Out
    - Array or Linked List
2. Queue - Fist In First Out
    - Array or Doubly Linked List
3. Priority Queue - Fist In First Out, but the first one is the one with the highest priority
    - heap(Binary, Binomial,Fibonacci) 二叉堆 二项堆 斐波那契堆
    - Binary Search Tree 二叉搜索树
    - 正常入，按照优先级出

[复杂度](https://www.bigocheatsheet.com/)



1. backspace string compare
2. implement queue using stacks（用栈实现队列）232
3. implement stack using queues（用队列实现栈）225
4. valid parentheses




#### valid parentheses

#####  时间复杂度：O(n)
- 原因：算法需要遍历输入字符串中的每个字符一次，其中 n 是字符串的长度。

- 操作细节：

    - 压栈（Push）：遇到左括号时，执行一次 O(1) 操作。

    - 弹栈（Pop）和匹配检查：遇到右括号时，执行一次 O(1) 的弹栈和一次 O(1) 的匹配检查。

- 总操作次数：每个字符仅处理一次，总时间为线性增长。

##### 空间复杂度：O(n)
- 原因：栈的最大深度与输入字符串的长度相关。

- 最坏情况：当所有字符均为左括号（如 ((((...）时，栈的大小达到 n。

- 一般情况：即使括号正确匹配，栈的深度可能为 n/2（例如 ()()()...），但大 O 表示法以最坏情况为准。

- 关键点总结
    - 栈的必要性：需维护括号嵌套顺序，无法用计数器优化（仅限多种括号时）。

    - 操作效率：栈的 push 和 pop 均为 O(1)，哈希表查映射也为 O(1)。

- 边界情况：无效字符直接跳过，不影响复杂度。

- 结论
使用栈实现括号匹配的时间复杂度和空间复杂度均为 O(n)，适用于任意括号类型的匹配场景。


[Comparison of theoretic bounds for variants](https://en.wikipedia.org/wiki/Heap_(data_structure))
##### Priority Queue
- 以小顶堆为例 Min Heap
    1. 越小的元素越靠近堆顶
    2. 父亲节点小于子节点
    3. 最小的元素永远在堆顶
- 以大顶堆为例 Max Heap
    1. 最大的元素永远在堆顶
    2. 父亲节点大于子节点
    3. 越大的元素越靠近堆顶

