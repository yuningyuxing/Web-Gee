package gee

//框架入口
//这部分代码 实现了请求处理函数和路由器的注册 查找 匹配以及处理功能
import (
	"log"
	"net/http"
)

// 定义了处理HTTP请求的函数类型 *Context包含了请求和响应的所有信息
type HandlerFunc func(*Context)

// 他是Web框架Gee的核心结构体
type Engine struct {
	//包含了一个路由器
	router *router
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

// 将路由规则和处理函数注册到路由表handler
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %s", method, pattern)
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 启动HTTP服务器 监听指定的端口
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// 实现ServerHTTP接口 接管了所有的HTTP请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
