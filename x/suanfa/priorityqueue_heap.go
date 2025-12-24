// 第K个最大元素
// 用小顶堆实现
// 堆的大小为k, 堆顶元素为第k大元素

// 时间复杂度：O(nlogk)，其中 n 是数组 nums 的长度。需要遍历数组 nums 一次，对于数组中的每个元素
// 也，然后排序，时间复杂度为O(klogk)

package suanfa

import (
	"container/heap"
)

type MinHeap []int

func (h MinHeap) Len() int {
	return len(h)
}
func (h MinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type KthLargest struct {
	heap *MinHeap
	k    int
}

func NewKthLargest(k int, nums []int) KthLargest {
	kl := KthLargest{
		heap: &MinHeap{},
		k:    k,
	}
	for _, num := range nums {
		kl.Add(num)
	}
	return kl
}

func (this *KthLargest) Add(val int) int {
	if this.heap.Len() < this.k {
		heap.Push(this.heap, val)
	} else {
		if val > (*this.heap)[0] {
			heap.Pop(this.heap)
			heap.Push(this.heap, val)
		}
	}
	return (*this.heap)[0]
}

// 滑动窗口最大值
// 1. 用大顶堆实现
// 2. 用双端队列实现
// 3. 用单调队列实现

// 大顶堆实现
// 时间复杂度：O(nlogn)，其中 n 是数组 nums 的长度。在最坏情况下，数组 nums
// 中的元素单调递增，那么将数组 nums 中的每个元素加入优先队列中需要 O(nlogn) 的时间。
// 空间复杂度：O(n)，其中 n 是数组 nums 的长度。空间复杂度主要取决于优先队列的大小，
// 而优先队列的大小为 k。

type MaxHeap struct {
	indices []int
	nums    []int
}

func (h MaxHeap) Less(i, j int) bool {
	return h.nums[h.indices[i]] > h.nums[h.indices[j]]
}
func (h MaxHeap) Swap(i, j int) {
	h.indices[i], h.indices[j] = h.indices[j], h.indices[i]
}
func (h MaxHeap) Len() int {
	return len(h.indices)
}

func (h *MaxHeap) Push(x interface{}) {
	h.indices = append(h.indices, x.(int))
}
func (h *MaxHeap) Pop() interface{} {
	old := h.indices
	n := len(old)
	v := old[n-1]
	h.indices = old[:n-1]
	return v
}
func MaxSlidingWindow(nums []int, k int) []int {
	q := &MaxHeap{
		indices: make([]int, 0, k),
		nums:    nums,
	}

	for i := 0; i < k; i++ {
		q.indices = append(q.indices, i)
	}
	heap.Init(q)

	res := []int{nums[q.indices[0]]}

	for i := k; i < len(nums); i++ {
		heap.Push(q, i)
		for q.indices[0] <= i-k {
			heap.Pop(q)
		}
		res = append(res, nums[q.indices[0]])
	}
	return res
}

// 单调队列实现
// 核心思想是维护一个单调递减的队列，保证队首元素为当前窗口的最大值。队列中的元素为数组 nums 中的下标，
// 这样可以方便地判断队首元素是否在当前窗口内。
// 时间复杂度：O(n)，其中 n 是数组 nums 的长度。每个元素最多被入队和出队一次，因此时间复杂度为 O(n)。
// 空间复杂度：O(k)，其中 k 是窗口的大小。队列中最多存储 k 个元素。
func MaxSlidingWindow2(nums []int, k int) []int {
	// 双端队列，存储数组nums的下标，队首元素为当前窗口的最大值，长度为k, 队尾元素为当前窗口的最小值
	deque := make([]int, 0, k)
	// 窗口可以滑动的次数是n - k + 1次，因此结果数组的长度应该是n - k + 1
	res := make([]int, 0, len(nums)-k+1)

	for i := 0; i < len(nums); i++ {
		// 移动超出窗口范围的元素，从头部
		if len(deque) > 0 && deque[0] < i-k+1 {
			deque = deque[1:]
		}
		// 保持队列单调递减，移除队列中小于当前元素的元素，从尾部

		for len(deque) > 0 && nums[i] >= nums[deque[len(deque)-1]] {
			deque = deque[:len(deque)-1]
		}
		// 将当前元素的下标加入队列
		deque = append(deque, i)
		if i >= k-1 {
			// 队首元素为当前窗口的最大值
			res = append(res, nums[deque[0]])
		}
	}
	return res
}

func MaxSlidingWindow3(nums []int, k int) []int {
	// 双端队列，存储数组nums的下标，队首元素为当前窗口的最大值，长度为k, 队尾元素为当前窗口的最小值
	windowIndicesDeque := make([]int, 0, k)
	// 窗口可以滑动的次数是n - k + 1次，因此结果数组的长度应该是n - k + 1
	result := make([]int, 0, len(nums)-k+1)

	maintainQueue := func(i int) {
		// 保持队列单调递减，移除队列中小于当前元素的元素，从尾部
		for len(windowIndicesDeque) > 0 && nums[i] >= nums[windowIndicesDeque[len(windowIndicesDeque)-1]] {
			windowIndicesDeque = windowIndicesDeque[:len(windowIndicesDeque)-1]
		}
		// 将当前元素的下标加入队列
		windowIndicesDeque = append(windowIndicesDeque, i)
	}
	for i := 0; i < k; i++ {
		maintainQueue(i)
	}

	result = append(result, nums[windowIndicesDeque[0]])

	for i := k; i < len(nums); i++ {

		maintainQueue(i)
		// for deque[0] <= i-k {
		for windowIndicesDeque[0] < i-k+1 {
			windowIndicesDeque = windowIndicesDeque[1:]
		}
		result = append(result, nums[windowIndicesDeque[0]])
	}
	return result
}
