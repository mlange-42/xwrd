package anagram

import (
	"fmt"
	"unicode"
)

// Tree is an anagram tree
type Tree struct {
	Root       *Node
	Leaves     [][]string
	Letters    []rune
	LettersMap map[rune]int
}

// NewTree creates a new Tree
func NewTree(letters []rune) Tree {
	lettersMap := make(map[rune]int, 2*len(letters))
	for i, letter := range letters {
		lettersMap[letter] = i
		lettersMap[unicode.ToUpper(letter)] = i
	}
	root := NewNode(letters[0], -1)
	return Tree{
		Root:       &root,
		Leaves:     make([][]string, 0, 0),
		Letters:    letters,
		LettersMap: lettersMap,
	}
}

// FindAnagrams finds full anagrams
func (t *Tree) FindAnagrams(word string) []string {
	result := make([]int, len(t.Letters), len(t.Letters))
	Histogram(word, t.LettersMap, result)

	node := t.Root
	for _, cnt := range result {
		child, ok := node.GetChild(cnt)
		if !ok {
			return []string{}
		}
		node = child
	}
	return t.Leaves[node.Leaf]
}

// AddWords adds words to the tree
func (t *Tree) AddWords(words []string) {
	result := make([]int, len(t.Letters), len(t.Letters))

	for _, word := range words {
		for i := 0; i < len(result); i++ {
			result[i] = 0
		}
		Histogram(word, t.LettersMap, result)

		node := t.Root
		for i, cnt := range result {
			if child, ok := node.GetChild(cnt); ok {
				node = child
			} else {
				if i < len(result)-1 {
					letter := t.Letters[i+1]
					child := NewNode(letter, -1)
					err := node.AddChild(&child, cnt)
					if err != nil {
						panic(err)
					}
					node = &child
				} else {
					leaf := t.addLeaf()
					child := NewNode('~', leaf)
					err := node.AddChild(&child, cnt)
					if err != nil {
						panic(err)
					}
					node = &child
				}
			}
			if i == len(result)-1 {
				t.Leaves[node.Leaf] = append(t.Leaves[node.Leaf], word)
			}
		}
	}
}

func (t *Tree) addLeaf() int {
	t.Leaves = append(t.Leaves, []string{})
	return len(t.Leaves) - 1
}

// Node is a node in an AnagramTree
type Node struct {
	Letter   rune
	Children []*Node
	Leaf     int
}

// NewNode creates a new Node
func NewNode(letter rune, leaf int) Node {
	return Node{
		Letter:   letter,
		Children: []*Node{},
		Leaf:     leaf,
	}
}

// AddChild adds a child node
func (n *Node) AddChild(node *Node, count int) error {
	if len(n.Children) <= count {
		for len(n.Children) <= count {
			n.Children = append(n.Children, nil)
		}
	} else {
		if n.Children[count] != nil {
			return fmt.Errorf("Node already has a child for count %d", count)
		}
	}
	n.Children[count] = node

	return nil
}

// GetChild returns the child with the given count
func (n *Node) GetChild(count int) (*Node, bool) {
	if len(n.Children) <= count {
		return nil, false
	}
	node := n.Children[count]
	if node == nil {
		return nil, false
	}
	return node, true
}
