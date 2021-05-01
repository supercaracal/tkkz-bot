package command

import (
	"fmt"
	"sort"
)

type suffixTreeNode struct {
	r        rune
	cnt      int
	depth    int
	parent   *suffixTreeNode
	children []*suffixTreeNode
}

func newSuffixTreeNode(r rune, parent *suffixTreeNode) *suffixTreeNode {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}

	return &suffixTreeNode{
		r:        r,
		depth:    depth,
		parent:   parent,
		children: make([]*suffixTreeNode, 0, 2),
	}
}

func detectLongestTandemRepeat(s string) (string, int) {
	if s == "" {
		return "", 0
	}

	runes := []rune(s)
	root := newSuffixTreeNode('0', nil)
	for i := range runes {
		root.add(runes[i:])
	}

	//root.printSuffixTree()

	node := root.detect()
	if node == nil {
		return "", 0
	}

	subRunes := make([]rune, 0, 8)
	for n := node; n.parent != nil; n = n.parent {
		subRunes = append(subRunes, n.r)
	}
	for i, j := 0, len(subRunes)-1; i < j; i, j = i+1, j-1 {
		r := subRunes[i]
		subRunes[i] = subRunes[j]
		subRunes[j] = r
	}

	return string(subRunes), node.cnt
}

func (node *suffixTreeNode) add(runes []rune) {
	if len(runes) == 0 {
		return
	}

	var child *suffixTreeNode
	for _, nd := range node.children {
		if nd.r == runes[0] {
			child = nd
			break
		}
	}

	if child == nil {
		child = newSuffixTreeNode(runes[0], node)
		node.children = append(node.children, child)
	}

	child.cnt++

	if len(runes) > 1 {
		child.add(runes[1:])
	}
}

func (node *suffixTreeNode) detect() *suffixTreeNode {
	if len(node.children) == 0 {
		return nil
	}

	if len(node.children) > 1 {
		sort.Slice(node.children, node.sort)
	}

	if node.children[0].cnt < 2 {
		return node
	}

	longest := node
	for _, child := range node.children {
		if child.cnt < node.cnt {
			break
		}

		n := child.detect()
		if n == nil {
			continue
		}

		if n.depth > longest.depth {
			longest = n
		}
	}

	if longest.depth == 1 {
		return nil
	}

	return longest
}

func (node *suffixTreeNode) sort(i, j int) bool {
	return node.children[i].cnt > node.children[j].cnt
}

func (node *suffixTreeNode) printSuffixTree() {
	for _, child := range node.children {
		fmt.Printf("%s:%d", string([]rune{child.r}), child.cnt)
		child.printSuffixTree()
		if len(child.children) == 0 {
			fmt.Println()
		}
	}
}
