package anagram

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var replacer = strings.NewReplacer(" ", "", "-", "")

// Tree is an anagram tree
type Tree struct {
	Root       *Node
	Leaves     []Leaf
	Letters    []rune
	LettersMap map[rune]int
}

// Node is a node in an AnagramTree
type Node struct {
	Letter   rune
	Children []*Node
	Leaf     int
}

// Leaf is a tree leaf
type Leaf []string

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
		Leaves:     make([]Leaf, 0, 0),
		Letters:    letters,
		LettersMap: lettersMap,
	}
}

// Anagrams finds full anagrams
func (t *Tree) Anagrams(word string) Leaf {
	word = replacer.Replace(word)

	hist := make([]int, len(t.Letters), len(t.Letters))
	Histogram(word, t.LettersMap, false, hist)

	if idx, ok := t.anagrams(hist); ok {
		return t.Leaves[idx]
	}
	return Leaf{}
}

func (t *Tree) anagrams(hist []int) (int, bool) {
	node := t.Root
	for _, cnt := range hist {
		child, ok := node.GetChild(cnt)
		if !ok {
			return -1, false
		}
		node = child
	}
	return node.Leaf, true
}

// AnagramsWithUnknown finds full anagrams
func (t *Tree) AnagramsWithUnknown(word string, minUnknown, maxUnknown uint) []Leaf {
	word = replacer.Replace(word)

	hist := make([]int, len(t.Letters), len(t.Letters))
	Histogram(word, t.LettersMap, false, hist)

	if maxUnknown == 0 {
		if idx, ok := t.anagrams(hist); ok {
			return []Leaf{t.Leaves[idx]}
		}
		return []Leaf{}
	}

	indices := t.anagramsWithUnknown(hist, minUnknown, maxUnknown)
	results := make([]Leaf, len(indices), len(indices))
	for i, idx := range indices {
		results[i] = t.Leaves[idx]
	}

	return results
}

type withUnknown struct {
	Node     *Node
	Unknowns uint
}

func (t *Tree) anagramsWithUnknown(hist []int, minUnknown, maxUnknown uint) []int {
	results := []int{}

	open := []*withUnknown{{t.Root, maxUnknown}}
	for _, cnt := range hist {
		newOpen := []*withUnknown{}

		for _, o := range open {
			for i := cnt; i <= cnt+int(o.Unknowns) && i < len(o.Node.Children); i++ {
				child := o.Node.Children[i]
				if child == nil {
					continue
				}
				newOpen = append(newOpen, &withUnknown{child, o.Unknowns - uint(i-cnt)})
			}
		}
		open = newOpen
	}

	diff := maxUnknown - minUnknown
	for _, o := range open {
		if o.Unknowns <= diff {
			results = append(results, o.Node.Leaf)
		}
	}

	return results
}

// PartialAnagrams finds partial anagrams
func (t *Tree) PartialAnagrams(word string, minLength uint) []Leaf {
	word = replacer.Replace(word)

	hist := make([]int, len(t.Letters), len(t.Letters))
	Histogram(word, t.LettersMap, false, hist)

	indices := t.partialAnagrams(hist, minLength)
	results := make([]Leaf, len(indices), len(indices))
	for i, idx := range indices {
		results[i] = t.Leaves[idx]
	}

	return results
}

func (t *Tree) partialAnagrams(hist []int, minLength uint) []int {
	results := []int{}

	open := []*Node{t.Root}
	for _, cnt := range hist {
		newOpen := []*Node{}

		for _, o := range open {
			for i, child := range o.Children {
				if i > cnt {
					break
				}
				if child == nil {
					continue
				}
				newOpen = append(newOpen, child)
			}
		}
		open = newOpen
	}

	for _, o := range open {
		if minLength == 0 || utf8.RuneCountInString(t.Leaves[o.Leaf][0]) >= int(minLength) {
			results = append(results, o.Leaf)
		}
	}

	return results
}

// PartialAnagramsWithUnknown finds partial anagrams
func (t *Tree) PartialAnagramsWithUnknown(word string, minLength, minUnknown, maxUnknown uint) []Leaf {
	word = replacer.Replace(word)

	hist := make([]int, len(t.Letters), len(t.Letters))
	Histogram(word, t.LettersMap, false, hist)

	var indices []int
	if maxUnknown == 0 {
		indices = t.partialAnagrams(hist, minLength)
	} else {
		indices = t.partialAnagramsWithUnknown(hist, minLength, minUnknown, maxUnknown)
	}
	results := make([]Leaf, len(indices), len(indices))
	for i, idx := range indices {
		results[i] = t.Leaves[idx]
	}

	return results
}

func (t *Tree) partialAnagramsWithUnknown(hist []int, minLength, minUnknown, maxUnknown uint) []int {
	results := []int{}

	open := []*withUnknown{{t.Root, maxUnknown}}
	for _, cnt := range hist {
		newOpen := []*withUnknown{}

		for _, o := range open {
			for i, child := range o.Node.Children {
				if i > cnt+int(o.Unknowns) {
					break
				}
				if child == nil {
					continue
				}
				rem := o.Unknowns
				if i > cnt {
					rem = o.Unknowns - uint(i-cnt)
				}
				childNode := withUnknown{child, rem}
				newOpen = append(newOpen, &childNode)
			}
		}
		open = newOpen
	}

	diff := maxUnknown - minUnknown
	for _, o := range open {
		if o.Unknowns <= diff &&
			(minLength == 0 || utf8.RuneCountInString(t.Leaves[o.Node.Leaf][0]) >= int(minLength)) {

			results = append(results, o.Node.Leaf)
		}
	}

	return results
}

// MultiAnagrams finds combinations of partial anagrams
func (t *Tree) MultiAnagrams(word string, maxWords, minLength uint, permutations bool) [][]Leaf {
	word = replacer.Replace(word)

	hist := make([]int, len(t.Letters), len(t.Letters))
	Histogram(word, t.LettersMap, false, hist)

	tree, indices := t.multiAnagrams(hist, maxWords, minLength, permutations)

	results := make([][]Leaf, len(indices), len(indices))
	for i, ind := range indices {
		row := make([]Leaf, len(ind), len(ind))
		for j, leaf := range ind {
			row[j] = tree.Leaves[leaf]
		}
		results[i] = row
	}

	return results
}

func (t *Tree) multiAnagrams(hist []int, maxWords, minLength uint, permutations bool) (*Tree, [][]int) {
	totalLen := 0
	for _, c := range hist {
		totalLen += c
	}

	partials := t.partialAnagrams(hist, minLength)

	tree := NewTree(t.Letters)
	for _, p := range partials {
		tree.AddWords(t.Leaves[p], nil)
	}

	open := [][]int{}
	closed := [][]int{}

	for i, p := range tree.Leaves {
		if utf8.RuneCountInString(p[0]) == totalLen {
			closed = append(closed, []int{i})
		} else {
			open = append(open, []int{i})
		}
	}

	if maxWords == 1 {
		return &tree, closed
	}

	tempHist := make([]int, len(hist), len(hist))
	for len(open) > 0 {
		curr := open[0]
		open = open[1:]

		for i := 0; i < len(tempHist); i++ {
			tempHist[i] = hist[i]
		}

		strLen := 0
		for _, c := range curr {
			str := tree.Leaves[c][0]
			strLen += utf8.RuneCountInString(str)
			Histogram(str, t.LettersMap, true, tempHist)
		}

		subPartials := tree.partialAnagrams(tempHist, minLength)

		if len(subPartials) == 0 {
			continue
		}

		for _, sub := range subPartials {
			if !permutations && sub < curr[len(curr)-1] {
				continue
			}
			new := []int{}
			new = append(new, curr...)
			new = append(new, sub)

			str := tree.Leaves[sub][0]
			if strLen+utf8.RuneCountInString(str) == totalLen {
				closed = append(closed, new)
			} else {
				if maxWords <= 0 || len(new) < int(maxWords) {
					open = append(open, new)
				}
			}
		}
	}

	return &tree, closed
}

// AddWords adds words to the tree
func (t *Tree) AddWords(words []string, progress chan int) {
	if progress != nil {
		defer close(progress)
	}

	numWords := len(words)
	prog := 0

	result := make([]int, len(t.Letters), len(t.Letters))

	for w, word := range words {
		if len(word) == 0 {
			continue
		}

		for i := 0; i < len(result); i++ {
			result[i] = 0
		}
		Histogram(word, t.LettersMap, false, result)

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

		p := (100 * w) / numWords
		if progress != nil && p > prog {
			progress <- p
			prog = p
		}
	}
	if progress != nil {
		progress <- 100
	}
}

func (t *Tree) addLeaf() int {
	t.Leaves = append(t.Leaves, []string{})
	return len(t.Leaves) - 1
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
