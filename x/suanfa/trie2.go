package suanfa

// 使用动态字符集 + 索引函数

type TrieWithIndexFun struct {
	children []*TrieWithIndexFun
	isEnd    bool
	indexFun func(rune) int
}

func ConstructorTrieWithIndexFun(indexFun func(rune) int) TrieWithIndexFun {
	return TrieWithIndexFun{
		children: make([]*TrieWithIndexFun, 0),
		isEnd:    false,
		indexFun: func(r rune) int {
			return int(r)
		},
	}
}
func (t *TrieWithIndexFun) Insert(word string) {
	node := t
	for _, ch := range word {
		// 计算字符在 children 数组中的索引
		chIndex := t.indexFun(ch)
		// 如果子节点不存在，创建一个新的子节点
		if node.children[chIndex] == nil {
			node.children[chIndex] = &TrieWithIndexFun{}
		}
		// 移动到下一个子节点
		node = node.children[chIndex]

	}
	node.isEnd = true
}
func (t *TrieWithIndexFun) Search(word string) bool {
	node := t.searchPrefix(word)
	return node != nil && node.isEnd
}
func (t *TrieWithIndexFun) StartsWith(prefix string) bool {
	return t.searchPrefix(prefix) != nil
}
func (t *TrieWithIndexFun) searchPrefix(prefix string) *TrieWithIndexFun {
	node := t
	for _, ch := range prefix {
		// 计算字符在 children 数组中的索引
		chIndex := t.indexFun(ch)
		// 如果子节点不存在，返回 nil
		if node.children[chIndex] == nil {
			return nil
		}
		// 移动到下一个子节点
		node = node.children[chIndex]
	}
	return node
}
