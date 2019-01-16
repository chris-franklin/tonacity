package tonacity

import "fmt"

// TODO: Node removal

// TrieNode A data structure for searching for values associated with a sequence of numbers, where the maximum range of those numbers is known up front.
type TrieNode struct {
	min      HalfSteps     // Minimum legal path value
	max      HalfSteps     // Maximum legal path value
	children []*TrieNode   // Slice of child nodes
	values   []interface{} // Values for paths that terminate at this node
}

// NewTrie Creates a new trie where each value in a path must be between min and max inclusive.
func NewTrie(min HalfSteps, max HalfSteps) *TrieNode {
	if max <= min {
		return nil
	}
	l := (max - min) + 1
	return &TrieNode{min, max, make([]*TrieNode, l, l), make([]interface{}, 0)}
}

func (node *TrieNode) String() string {
	return fmt.Sprintf("Trie Node (%d to %d) with %d value(s)", node.min, node.max, len(node.values))
}

// AddValue Adds a value to this trie node at the given path
func (node *TrieNode) AddValue(path []HalfSteps, value interface{}) {
	if len(path) == 0 {
		node.values = append(node.values, value)
	} else {
		child := node.children[path[0]-node.min]
		if child == nil {
			child = NewTrie(node.min, node.max)
			node.children[path[0]-node.min] = child
		}
		child.AddValue(path[1:], value)
	}
}

// FindValue Find the value for the given path, or nil if the path has no value
func (node *TrieNode) FindValue(path []HalfSteps) []interface{} {
	if len(path) == 0 {
		return node.values
	}
	if path[0] < node.min || path[0] > node.max {
		return nil
	}
	child := node.children[path[0]-node.min]
	if child == nil {
		return nil
	}
	return child.FindValue(path[1:])
}
