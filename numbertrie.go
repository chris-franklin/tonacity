package tonacity

// TODO: Node removal, handling of negative pattern values (offset to make positive), feedback if pattern already in trie

// TrieNode A data structure for searching for values associated with a sequence of numbers, where the maximum range of those numbers is known up front.
type TrieNode struct {
	children []*TrieNode
	value    interface{}
}

// NewTrie Creates a new trie where each value in a path is between 0 and maxChildren
func NewTrie(maxChildren int) *TrieNode {
	return &TrieNode{make([]*TrieNode, maxChildren, maxChildren), nil}
}

// AddValue Adds a value to this trie node at the given path
func (node *TrieNode) AddValue(path []HalfSteps, value interface{}) {
	if len(path) == 0 {
		node.value = value
	} else {
		child := node.children[path[0]]
		if child == nil {
			child = NewTrie(len(node.children))
			node.children[path[0]] = child
		}
		child.AddValue(path[1:], value)
	}
}

// FindValue Find the value for the given path, or nil if the path has no value
func (node *TrieNode) FindValue(path []HalfSteps) interface{} {
	if len(path) == 0 {
		return node.value
	}
	if path[0] < 0 || int(path[0]) >= len(node.children) {
		return nil
	}
	child := node.children[path[0]]
	if child == nil {
		return nil
	}
	return child.FindValue(path[1:])
}
