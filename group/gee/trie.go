package gee

import (
	"fmt"
	"strings"
)

// 设计树节点
type node struct {
	//待匹配路由
	pattern string
	//路由中的一部分
	part string
	//子节点
	children []*node
	//是否精确匹配  part含有:或*时为true
	isWild bool
}

// 返回记录当前节点信息的字符串
func (n *node) String() string {
	return fmt.Sprintf("node{pattern=%s,part=%s,isWilkd=%t}", n.pattern, n.part, n.isWild)
}

// 遍历所有节点并将具有pattern的节点添加到列表中
func (n *node) travel(list *([]*node)) {
	//如果当前节点具有pattern
	if n.pattern != "" {
		//*list是一个结构体指针的切片
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}

}

// 第一个匹配成功的节点 用于插入
// 参数part就是待匹配的一部分
func (n *node) matchChild(part string) *node {
	//遍历当前节点的子节点
	for _, child := range n.children {
		//如果当前子节点这部分和part相同 或者是模糊的 则是匹配成功
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点 用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 支持节点的插入
// parts存了我们要匹配总共有几部分
func (n *node) insert(pattern string, parts []string, height int) {
	//插入完成
	if len(parts) == height {
		//只有是一个完整路径的时候pattern才不为空
		n.pattern = pattern
		return
	}
	//我们拿出要匹配的这一部分
	part := parts[height]
	//去找子节点看存不存在这一部分
	child := n.matchChild(part)
	//如果不存在 就新建一个节点
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		//新建完节点后 添加到当前节点的子节点中来
		n.children = append(n.children, child)
	}
	//去查找下一个部分
	child.insert(pattern, parts, height+1)
}

// 支持节点的查询
func (n *node) search(parts []string, height int) *node {
	//判断此时查询完成 要么长度一致 要么n.part中含有前缀字符串*(注意当我们遍历到这里的时候说明父节点已经匹配成功 所以此时如果喊*就是匹配成功)
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//判断这里是否是有效pattern
		if n.pattern == "" {
			return nil
		}
		return n
	}
	//拿出我们下一个要匹配的part
	part := parts[height]
	//找到能匹配到part的所有子节点
	children := n.matchChildren(part)
	//判断所有合法子节点 递归查询看是否有匹配成功的
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
