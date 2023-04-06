package gee

import (
	"net/http"
	"strings"
)

type router struct {
	//存储每种请求方式的Trie树根节点  使用前缀数存储动态路由
	roots map[string]*node
	//存储每种（请求方法+路由路径）的HandlerFunc
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 本函数的作用是将pattern切割成parts
func parsePattern(pattern string) []string {
	//将pattern按照/分割
	vs := strings.Split(pattern, "/")
	//存放parts
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			//当加入到有*时 停止
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// 添加路由  将路由路径添加到前缀树中，并映射到对应的处理函数
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//先将请求划分成parts
	parts := parsePattern(pattern)
	//构造完整请求
	key := method + "-" + pattern
	//查看当前请求方式是否有根节点 如果没有申请一个
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	//以当前请求方式为根节点 进行插入
	r.roots[method].insert(pattern, parts, 0)
	//将对应请求的函数记录
	r.handlers[key] = handler
}

// 函数根据请求的方法method和路径path查找路由
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	//解析请求路径 searchParts是一个parts字符串数组
	searchParts := parsePattern(path)
	//用于存储匹配到的路由参数 参数映射 比如说/:name  匹配的是/liu  那么name=liu params就是用来存这个的
	params := make(map[string]string)
	//查找该请求方法的根节点 如果没有则表示该请求方法没有被注册
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	//从根节点开始搜索路径
	n := root.search(searchParts, 0)
	//如果找到匹配路径 就解析路径中的参数并返回处理函数
	if n != nil {
		//我们将匹配到的路径分解  注意这里是匹配路径 而不是我们的path 不要混在一起了
		parts := parsePattern(n.pattern)
		//遍历我们分解的匹配路径的每一部分
		for index, part := range parts {
			//如果该部分是路由参数 我们就把他添加到参数映射中
			if part[0] == ':' {
				//这里就可以理解为params[name]=liu 后面是path
				params[part[1:]] = searchParts[index]
			}
			//若该部分是通配符 则将剩余部分添加到参数映射中
			if part[0] == '*' && len(part) > 1 {
				//这里就是个添加多个 并且用/分割
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		//注意这里的n.是匹配路径 而非真实路径
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
