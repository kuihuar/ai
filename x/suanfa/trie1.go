package suanfa

type TrieUnicode struct {
	children map[rune]*TrieUnicode
	isEnd    bool
}

func ConstructorTrieUnicode() TrieUnicode {
	return TrieUnicode{
		children: make(map[rune]*TrieUnicode),
		isEnd:    false,
	}
}
func (t *TrieUnicode) Insert(word string) {
	node := t
	for _, ch := range word {
		if node.children == nil {
			node.children = make(map[rune]*TrieUnicode)
		}

		if node.children[ch] == nil {
			node.children[ch] = &TrieUnicode{}
		}
		node = node.children[ch]
	}
	node.isEnd = true
}

func (t *TrieUnicode) Search(word string) bool {
	node := t.searchPrefix(word)
	return node != nil && node.isEnd
}
func (t *TrieUnicode) StartsWith(prefix string) bool {
	return t.searchPrefix(prefix) != nil
}
func (t *TrieUnicode) searchPrefix(prefix string) *TrieUnicode {
	node := t

	for _, ch := range prefix {
		if node.children == nil {
			return nil
		}
		if node.children[ch] == nil {
			return nil
		}
		node = node.children[ch]
	}
	return node
}
