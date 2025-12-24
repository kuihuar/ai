package suanfa

// 212 单词搜索 II
// 给定一个 m x n 二维字符网格 board 和一个单词（字符串）列表 words，找出所有同时在二维网格和字典中出现的单词。
// 单词必须按照字母顺序，通过 相邻的单元格 内的字母构成，其中“相邻”单元格是那些水平相邻或垂直相邻的单元格。同一个单元格内的字母在一个单词中不允许被重复使用。

func FindWords(board [][]byte, words []string) []string {
	// 构建 Trie
	trie := &Trie{}
	for _, word := range words {
		trie.Insert(word)
	}
	// 结果集
	result := make([]string, 0)
	var dirs = [][]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	// 深度优先搜索
	var dfs func(node *Trie, i, j int, path string)
	dfs = func(node *Trie, i, j int, path string) {
		// 边界条件
		if i < 0 || i >= len(board) || j < 0 || j >= len(board[0]) {
			return
		}
		// 获取当前字符
		c := board[i][j]
		// 如果当前字符不在 Trie 中，返回

		if c == '#' || node.children[c-'a'] == nil {
			return
		}

		// 进入下一层
		node = node.children[c-'a']
		path += string(c)
		// 如果当前节点是单词结尾，添加到结果集
		if node.isEnd {
			result = append(result, path)
			// 避免重复添加单词
			node.isEnd = false
		}
		// 标记当前字符已经访问过
		board[i][j] = '#'
		for _, dir := range dirs {
			dfs(node, i+dir[0], j+dir[1], path)
		}
		// 继续搜索
		// dfs(node, i+1, j, path)
		// dfs(node, i-1, j, path)
		// dfs(node, i, j+1, path)
		// dfs(node, i, j-1, path)
		// 回溯
		board[i][j] = c

	}
	// 遍历整个二维网格
	for i := 0; i < len(board); i++ {
		for j := 0; j < len(board[0]); j++ {
			dfs(trie, i, j, "")
		}
	}
	return result
}
