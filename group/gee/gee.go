package gee

import (
	"log"
	"net/http"
	"strings"
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
	//先构造一个自己
	engine := &Engine{router: newRouter()}
	//注意这个是创建根路由组  且我们所有路由组共享一个engine
	engine.RouterGroup = &RouterGroup{engine: engine}
	//将根路由组加入到路由组中
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

// comp是路径
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

// 将中间件加入到改group组
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

// 修改过后的ServerHTTP会先获取该请求要用到的中间件
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	//遍历所有group组
	for _, group := range engine.groups {
		//当我们的请求和该组group有相同前缀的时候 表示该请求要应用这个group的中间组
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			//获得group的中间组
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	//存储最c.handlers里面
	c.handlers = middlewares
	engine.router.handle(c)
}
