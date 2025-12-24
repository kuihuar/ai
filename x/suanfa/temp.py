def threeSum(self, nums):
    if len(nums) < 3:
        return []
    nums.sort()
    res = set()
    # 外层循环，固定第一个元素，范围是[:-2]，因为至少要留两个元素
    for i, v in enumerate(nums[:-2]):
        # 跳过重复元素
        if i >= 1 and v == nums[i-1]:
            continue
        # 内层循环，固定第二个元素，范围是[i+1:]，因为至少要留一个元素
        d = {}
        for x in nums[i+1:]:
            if x not in d:
                d[-v-x] = 1
            else:
                res.add((v, -v-x, x))
        return map(list, res)



def BFS(graph, start, end):
    queue = []
    queue.append([start])
    visited = set()
    visited.add(start)
    while queue:
        node = queue.pop(0)
        visited.add(node)
        ## 处理当前节点
        process(node)
        ## 获取当前节点的所有相邻节点，找node的后继节点，也就是node的所有子节点
        ## 如果node的后继节点没有访问过，就加入队列
        nodes = generate_related_nodes(node)
        ## 相邻节点加入队列
        queue.push(nodes)
    # other processing work
    ...        

## 递归实现
visited = set()
# 本身递归实现了栈
def DFS(node):
    if node in visited: # terminator
        # already visited
        return
    visited.add(node)
    # process current node here.
   ...
    for next_node in node.children():
        if not next_node in visited:
            DFS(next_node)


## 非递归实现
def DFS1(self, root):
    if root is None:
        return []
   ## 手动维护栈
   visited, stack = [], [root]
   while stack:
       node = stack.pop()
       visited.add(node)
       process(node)
       nodes = generate_related_nodes(node)
       stack.push(nodes)
   # other processing work
   ...
              