package gee

import (
	"log"
	"net/http"
)

type router struct {
	//用来保存每个路由的处理函数
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
	}
}

// 向router里面添加新的路由
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//添加路由时会打印一条日志
	log.Printf("Route %4s - %s,", method, pattern)
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// 为当前请求 获取处理函数 并调用他处理请求
func (r *router) handle(c *Context) {
	key := c.Method + "-" + c.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
