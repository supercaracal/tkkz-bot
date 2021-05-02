package command

import (
	"fmt"
	"sync"
)

type suffixTreeNode struct {
	r        rune
	cnt      int
	depth    int
	parent   *suffixTreeNode
	children []*suffixTreeNode
}

const (
	maxTraversalAttempts = 1000
)

var (
	suffixTreeNodePool = sync.Pool{
		New: func() interface{} {
			return &suffixTreeNode{children: make([]*suffixTreeNode, 0, 2)}
		},
	}

	runesPool = sync.Pool{
		New: func() interface{} {
			return make([]rune, 0, 2)
		},
	}
)

func (node *suffixTreeNode) printSuffixTree() {
	for _, child := range node.children {
		fmt.Printf("%s:%d", string([]rune{child.r}), child.cnt)
		child.printSuffixTree()
		if len(child.children) == 0 {
			fmt.Println()
		}
	}
}

func (node *suffixTreeNode) reset() {
	defer suffixTreeNodePool.Put(node)

	for _, child := range node.children {
		child.reset()
	}

	node.r = '0'
	node.cnt = 0
	node.depth = 0
	node.parent = nil
	node.children = node.children[:0]
}

func (node *suffixTreeNode) add(runes []rune) {
	if len(runes) == 0 {
		return
	}

	var child *suffixTreeNode
	for _, nd := range node.children {
		if nd.r == runes[0] {
			child = nd
			child.cnt++
			break
		}
	}

	if child == nil {
		child = suffixTreeNodePool.Get().(*suffixTreeNode)
		child.r = runes[0]
		child.depth = node.depth + 1
		child.parent = node
		node.children = append(node.children, child)
	}

	if len(runes) > 1 {
		child.add(runes[1:])
	}
}

// detect returns the last node of the first pattern part of a tandem repeat such as:
// tandemtandemtandemtan.....
//      ^
func (node *suffixTreeNode) detect() *suffixTreeNode {
	if len(node.children) == 0 || (node.parent != nil && node.cnt < 1) {
		return nil
	}

	longest := node
	for _, child := range node.children {
		if child.cnt < 1 || (child.cnt < node.cnt && node.depth == 1) {
			continue
		}

		n := child.detect()
		if n == nil {
			continue
		}

		if n.depth > longest.depth {
			longest = n
		}
	}

	return longest
}

// count returns a number of the tandem repeat such as:
// tandemtandemtandemfoobarbaz => 3
func (node *suffixTreeNode) count(pattern []rune) (cnt int) {
	if len(pattern) == 0 {
		return
	}

	for n, i := node, 0; i < maxTraversalAttempts; i++ {
		for _, r := range pattern {
			hasNext := false
			for _, child := range n.children {
				if child.r == r {
					hasNext = true
					n = child
					break
				}
			}
			if !hasNext {
				return
			}
		}
		cnt++
	}

	return
}

func detectLongestTandemRepeat(s string) (string, int) {
	if s == "" {
		return "", 0
	}

	root := suffixTreeNodePool.Get().(*suffixTreeNode)
	defer root.reset()

	runes := runesPool.Get().([]rune)
	defer func(r []rune) { r = r[:0]; runesPool.Put(r) }(runes)

	for _, r := range s {
		runes = append(runes, r)
	}

	for i := range runes {
		root.add(runes[i:])
	}

	//root.printSuffixTree()

	node := root.detect()
	if node == nil {
		return "", 0
	}

	subRunes := runesPool.Get().([]rune)
	defer func(r []rune) { r = r[:0]; runesPool.Put(r) }(subRunes)

	for n := node; n.parent != nil; n = n.parent {
		subRunes = append(subRunes, n.r)
	}

	if len(subRunes) < 2 {
		return "", 0
	}

	// reversing
	for i, j := 0, len(subRunes)-1; i < j; i, j = i+1, j-1 {
		r := subRunes[i]
		subRunes[i] = subRunes[j]
		subRunes[j] = r
	}

	if cnt := root.count(subRunes); cnt > 1 {
		return string(subRunes), cnt
	}

	return "", 0
}
