package gee

import (
	"fmt"
	"net/http"
)

// 定义某种函数类型
// 这是提供给框架用户的 用来定义路由映射的处理方法
type HandlerFunc func(http.ResponseWriter, *http.Request)

// 给Engine添加了一张路由映射表router
type Engine struct {
	//这里可以理解为每种请求对应一个处理函数HandlerFunc
	//键值key由请求方法和静态路由地址构成 这样我们针对相同路由 如果请求方法不同 可以给出不同的处理方法
	//Value是用户映射的处理方法
	router map[string]HandlerFunc
}

// 创建一个Engine实例
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// 增加一种处理方法针对某种哦请求
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// 这个方法会将GET请求和路由和处理方法 放到映射表router
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// 这个方法会将POST请求和路由和处理方法 放到映射表router
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 启动Web服务
func (engine *Engine) Run(add string) (err error) {
	return http.ListenAndServe(add, engine)
}

// 解析请求的路径，查找路由映射表 如果找到就执行注册的处理方法 如果找不到 返回404 NOT FOUND
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//req.Method表示请求 req.URL.Path表示静态路由地址
	key := req.Method + "-" + req.URL.Path
	//在map中查找
	if handler, ok := engine.router[key]; ok {
		//如果查找到就执行对应的处理方法
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}
