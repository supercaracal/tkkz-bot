package command

import (
	"fmt"
)

type suffixTreeNode struct {
	parent   *suffixTreeNode
	r        rune
	depth    int
	count    int
	children map[rune]*suffixTreeNode
}

func newSuffixTreeNode(parent *suffixTreeNode, r rune) *suffixTreeNode {
	depth := 0
	if parent != nil {
		depth = parent.depth + 1
	}

	return &suffixTreeNode{
		parent:   parent,
		r:        r,
		depth:    depth,
		children: map[rune]*suffixTreeNode{},
	}
}

func detectLongestTandemRepeat(s string) (string, int) {
	if s == "" {
		return "", 0
	}

	runes := []rune(s)
	root := newSuffixTreeNode(nil, '0')
	for i := range runes {
		root.add(runes[i:])
	}

	// printSuffixTree(root)

	node := root.dfs()
	if node == nil || node.count < 2 {
		return "", 0
	}

	buf := make([]rune, 0, node.depth)
	for p, n := node, node; p.count == n.count; {
		buf = append(buf, n.r)
		p = n
		n = n.parent
	}

	for i, j := 0, len(buf)-1; i < j; i, j = i+1, j-1 {
		x := buf[j]
		buf[j] = buf[i]
		buf[i] = x
	}

	return string(buf), node.count
}

func (stn *suffixTreeNode) add(runes []rune) {
	if len(runes) == 0 {
		return
	}

	if _, ok := stn.children[runes[0]]; !ok {
		stn.children[runes[0]] = newSuffixTreeNode(stn, runes[0])
	}

	stn.children[runes[0]].count++

	if len(runes) > 1 {
		stn.children[runes[0]].add(runes[1:])
	}
}

func (stn *suffixTreeNode) dfs() *suffixTreeNode {
	longest := stn

	for _, child := range stn.children {
		if child.count < stn.count {
			continue
		}

		if node := child.dfs(); node.depth > longest.depth {
			longest = node
		}
	}

	return longest
}

func printSuffixTree(node *suffixTreeNode) {
	for _, child := range node.children {
		fmt.Printf("%s:%d:%d", string([]rune{child.r}), child.count, child.depth)
		printSuffixTree(child)
		if len(child.children) == 0 {
			fmt.Println()
		}
	}
}
