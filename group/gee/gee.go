package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type (
	//用于实现路由分组功能
	RouterGroup struct {
		prefix      string        //路由组的前缀
		middlewares []HandlerFunc //中间件支持
		parent      *RouterGroup  //支持嵌套  parent可以理解为父亲
		engine      *Engine       //所有路由公用一个Engine实例
	}

	Engine struct {
		*RouterGroup                //根路由组
		router       *router        //路由器
		groups       []*RouterGroup //所有的路由组
	}
)

// Engine的构造函数
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// 用于创建一个新的路由组 可以理解为生儿子
// 注意所有路由组公用一个Engine实例
// prefix表示新建的RouterGroup的路由前缀
// RouterGroup支持嵌套  一个RouterGroup可以以多个子RouterGroup 通过parent关联
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
