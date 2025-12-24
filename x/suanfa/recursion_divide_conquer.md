1. Recursion
2. Divide & Conquer
3. Backtracking
4. Traversal



### 递归模板
```python

def recursion(level, param1, param2, ...):
    # 终止条件
    # recursion terminator
    if level > MAX_LEVEL:
        process_result
        return
    # 业务处理逻辑    
    # process logic in current level
    process(level, data...)
    # drill down
    # 下探到下一层
    self.recursion(level + 1, p1, ...)
    # 如有必要，清理当前层的状态
    # reverse the current level status if needed
    reverse_state(level)
```


```go 
func Factorial(n int) int {
	if n <= 1 {
		return 1
	}
	return n * Factorial(n-1)
}

```
Factorial(5)
5 * Factorial(4)
5 * 4 * Factorial(3)
5 * 4 * 3 * Factorial(2)
5 * 4 * 3 * 2 * Factorial(1)
5 * 4 * 3 * 2 * 1
5 * 4 * 3 * 2
5 * 4 * 6
5 * 24
120
   


Divide & Conquer
    1. 分治（Divide）：将问题分解为更小的子问题。
    2. 解决（Conquer）：递归地解决每个子问题。
    3. 合并（Combine）：将子问题的解合并为原问题的解。
    4. 终止条件：当问题足够小，可以直接求解时，终止递归。
    分治法的核心思想是将一个大问题分解为若干个小问题，然后递归地解决这些小问题，最后将它们的解合并为原问题的解。
    分治法通常用于解决以下类型的问题：
    - 排序问题：如快速排序、归并排序等。
    - 搜索问题：如二分搜索、线性搜索等。

分治的子问问题没有相关性，可以并行处理，没有重复计算（中间结果）

### 分冶模板
```python
def divide_conquer(problem, param1, param2,...):
    # 终止条件
    # recursion terminator
    if problem is None:
        print_result
        return
    # 准备数据
    # prepare data
    data = prepare_data(problem)
    subproblems = split_problem(problem, data)
    # 分治
    # conquer subproblems
    subresult1 = self.divide_conquer(subproblems[0], p1,...)
    subresult2 = self.divide_conquer(subproblems[1], p1,...)
    subresult3 = self.divide_conquer(subproblems[2], p1,...)
    ...
    # 合并
    # process and generate the final result
    result = process_result(subresult1, subresult2, subresult3,...)
    # revert the current level states
```
### 回溯模板
```python
def backtrack(路径, 选择列表):
    if 满足结束条件:
        result.add(路径)
        return
    for 选择 in 选择列表:
        做选择
        backtrack(路径, 选择列表)
        撤销选择
```
   