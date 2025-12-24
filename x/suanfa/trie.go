package suanfa

// 208 实现 Trie (前缀树)

// Trie 是一个前缀树

type Trie struct {
	// 存储子节点
	children [26]*Trie
	// 标记是否为单词结尾
	isEnd bool
}

// Constructor 初始化 Trie
func ConstructorTrie() Trie {
	return Trie{}
}

// Insert 插入单词·
func (t *Trie) Insert(word string) {
	node := t
	for _, ch := range word {
		// 计算字符在 children 数组中的索引
		ch -= 'a'
		// 如果子节点不存在，创建一个新的子节点
		if node.children[ch] == nil {
			node.children[ch] = &Trie{}
		}
		// 移动到下一个子节点
		node = node.children[ch]
	}
	node.isEnd = true
}

// Search 搜索单词
func (t *Trie) Search(word string) bool {
	node := t.searchPrefix(word)
	return node != nil && node.isEnd
}

// StartsWith 前缀搜索
func (t *Trie) StartsWith(prefix string) bool {
	return t.searchPrefix(prefix) != nil
}

// searchPrefix 前缀搜索
func (t *Trie) searchPrefix(prefix string) *Trie {
	node := t
	for _, ch := range prefix {
		// 计算字符在 children 数组中的索引
		ch -= 'a'
		// 如果子节点不存在，返回 nil
		if node.children[ch] == nil {
			return nil
		}
		// 移动到下一个子节点
		node = node.children[ch]
	}
	return node
}
