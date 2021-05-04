package command

import (
	"fmt"
	"strings"
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

// using suffix tree
func detectLongestTandemRepeat1(s string) (string, int) {
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

// sliding and comparing
func detectLongestTandemRepeat2(s string) (string, int) {
	if s == "" {
		return "", 0
	}

	runes := runesPool.Get().([]rune)
	defer func(r []rune) { r = r[:0]; runesPool.Put(r) }(runes)

	for _, r := range s {
		runes = append(runes, r)
	}

	var (
		// ex) fooabcabcabcbar
		//        ^^^^^^^^^
		// first=3, size=3
		first, size int
	)

	for i, j, tmp := 1, 0, len(runes); i <= len(runes)/2; i++ {
		for j = 0; i+j < len(runes); j++ {
			if runes[i+j] != runes[j] {
				if j-tmp > size {
					first = tmp
					size = j - tmp
					//fmt.Printf("1: max updated: first=%d size=%d: %s\n", first, size, s)
				}
				tmp = len(runes)
			} else {
				if tmp == len(runes) {
					tmp = j
				}
			}
			//fmt.Printf("%d: %s, %d: %s\n", i+j, string([]rune{runes[i+j]}), j, string([]rune{runes[j]}))
		}
		if j-tmp > size {
			first = tmp
			size = j - tmp
			//fmt.Printf("2: max updated: first=%d size=%d: %s\n", first, size, s)
		}
	}

	mono := true
	for i := first + 1; i < first+size; i++ {
		if runes[i-1] != runes[i] {
			mono = false
			break
		}
	}

	if mono || size < 2 {
		return "", 0
	}

	return string(runes[first : first+size]), func() (cnt int) {
		for i := first; i < len(runes) && runes[i] == runes[first+((i-first)%size)]; i++ {
			//fmt.Printf("%d: %s, %d: %s\n", i, string([]rune{runes[i]}), first+((i-first)%size), string([]rune{runes[first+((i-first)%size)]}))
			if (i-first+1)%size == 0 {
				cnt++
			}
		}
		return
	}()
}

func detectLongestTandemRepeat3(s string) (string, int) {
	if s == "" {
		return "", 0
	}

	runes := runesPool.Get().([]rune)
	defer func(r []rune) { r = r[:0]; runesPool.Put(r) }(runes)

	for _, r := range s {
		runes = append(runes, r)
	}

	l := len(runes)
	if l < 3 {
		return "", 0
	}

	var (
		w   []rune
		cnt int
	)

	for i, end := 1, l-1; i < l/2; i++ {
		start := end - i
		w = runes[start:end]
		c := 0
		for j := 1; j*i <= l; j++ {
			tandem := strings.Repeat(s, j)
			if !strings.HasSuffix(s, tandem) {
				c = j - 1
				break
			}
		}
		if c > cnt {
			cnt = c
		} else {
			break
		}
	}

	return string(w), cnt
}
