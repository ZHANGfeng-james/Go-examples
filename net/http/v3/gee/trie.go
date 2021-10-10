package gee

import "strings"

// node constructor of router trie tree
type node struct {
	pattern  string  // 完整匹配路径
	part     string  // 当前节点的匹配内容
	children []*node // 每个节点下的子节点
	isWild   bool    // 是否包含通配符（* 和 :）
}

// insert trie tree node with pattern
func (n *node) insert(pattern string, parts []string, height int) {
	//TEST CASE: /p/:name/join [p, :name, join] 0
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// TDD
	// 0 --> p
	// 1 --> :name
	// 2 --> join

	//FIXME /p/:name/join /p/:time/sell
	//FIXME /p/:name /p/michoi
	part := parts[height]
	child := n.matchChild(part)
	if child == nil {
		// :name *filepath 存入 node.part
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// just for * only once
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			// middle path of route，并不是一个路由
			return nil
		}
		return n
	}

	part := parts[height]
	child := n.matchChildren(part)
	// children 为 []*node，是否存在多个匹配情况
	for _, item := range child {
		result := item.search(parts, height+1)
		if child != nil {
			return result
		}
	}
	return nil
}

// matchChild matches children of node to find match one
func (n *node) matchChild(path string) *node {
	for _, ele := range n.children {
		if ele.part == path || ele.isWild {
			return ele
		}
	}
	return nil
}

// matchChildren matches all children, and return all nodes
func (n *node) matchChildren(path string) []*node {
	nodes := make([]*node, 0)
	for _, ele := range n.children {
		if ele.part == path || ele.isWild {
			nodes = append(nodes, ele)
		}
	}
	return nodes
}
