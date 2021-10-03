package gen

import "strings"

/*
GenRouter:
	1.implement of general match, /static/*filepath
	-> any thing has prefix static
	2.implement of parameter match, /p/:lang/doc
	 -> /p/python/doc p/cplus/doc p/go/doc ...
*/
type trieNode struct {
	pattern string //route pattern to match

	//classical trie struct
	part     string //split from url, value of trie tree
	children []*trieNode
	isWild   bool //if the part shows as ":something" or "*something"
}

//the first parent that successfully match, used by insert process
func (parent *trieNode) matchChild(part string) *trieNode {
	for _, child := range parent.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

//all parents that successfully match, used by search function
func (parent *trieNode) matchChildren(part string) []*trieNode {
	parents := make([]*trieNode, 0)
	for _, child := range parent.children {
		if child.part == part || child.isWild {
			parents = append(parents, child)
		}
	}
	return parents
}

//recursion insert,the order of router will exert effect on search
func (parent *trieNode) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		parent.pattern = pattern
		return
	}

	part := parts[height]
	child := parent.matchChild(part)
	if child == nil {
		child = &trieNode{part: part, isWild: part[0] == ':' || part[0] == '*'}
		parent.children = append(parent.children, child)
	}
	child.insert(pattern, parts, height+1)
}

//recursion search
func (parent *trieNode) search(parts []string, height int) *trieNode {
	if len(parts) == height || strings.HasPrefix(parent.part, "*") {
		if parent.pattern == "" {
			return nil
		}
		return parent
	}

	part := parts[height]
	children := parent.matchChildren(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
