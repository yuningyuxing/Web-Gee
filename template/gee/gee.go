package gee

import (
	"html/template"
	"log"
	"net/http"
	"path"
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
		//将所有的模板加载进内存
		//它用于存储和管理一个或多个HTML模板 同时提供方法来执行这些模板 从而渲染出最终的HTML输出
		htmlTemplates *template.Template
		//表示所有自定义模板渲染函数
		//FuncMap类型定义了字符串到函数的映射
		funcMap template.FuncMap
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
	//新增实例化时要给c.engine赋值
	c.engine = engine
	engine.router.handle(c)
}

// 该方法通过传入相对路径和文件系统  返回一个可以处理请求并提供静态文件服务的函数
// FileSystem是一个接口  它实现了对一系列命名文件的访问，文件路径的分割符为‘/’ 不管主机操作系统的惯例如何
// 可以理解FileSystem一般存的是静态文件目录或者是打包后的文件 因为是接口 所以可以存储不同类型的文件系统(真实地址)
// 解析请求的地址 并将其映射到服务器上文件的真实地址
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	//计算文件的绝对路径
	//path是实现了对斜杠分隔的路径的实用操作函数
	//Join函数可以将任意数量的路径元素放入一个单一路径里 会根据需要添加斜杠
	absolutePath := path.Join(group.prefix, relativePath)
	//创建一个http.FileServer实例 并移除绝对路径前缀
	//StripPrefix会返回一个Handler 该处理器会将请求的URL.Path字段中给定前缀的prefix去除后交给实现了Handler的ServeHTTP处理
	//StripPrefix会向URL.Path字段中没有给定前缀的请求回复404 page not found
	//http.FileServer可以通过fs访问到指定静态文件 它会将请求的URL路径映射到本地磁盘上的静态文件 并向客户端发送该文件的内容
	//这里他就会把fs的absolutePath去掉后交给FileServer返回的Handler(文件服务器)处理
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	//返回处理函数  该函数尝试打开并服务于请求的文件
	return func(c *Context) {
		//获取文件路径
		file := c.Param("filepath")
		//检查文件是否存在以及我们是否能打开
		if _, err := fs.Open(file); err != nil {
			//如果不能打开
			c.Status(http.StatusNotFound)
			return
		}
		//如果能打开调用ServeHTTP服务于请求的文件
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// 这个方法是暴露给用户的 用户可以将磁盘上的某个文件夹root映射到路由relativePath上
// static方法用于注册静态文件处理函数 支持绝对和相对路径
// 注意root是指定了静态文件所在的磁盘路径  relativePath是指定了要映射到的路由相对路径
func (group *RouterGroup) Static(relativePath string, root string) {
	//创建一个处理静态文件的HandlerFunc
	//http.Dir可以将一个目录转化为http.FileSystem接口
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	//计算出静态文件的路径模式 并将其注册到group的get方法里面
	urlPattern := path.Join(relativePath, "/*filepath")
	group.GET(urlPattern, handler)
}

// 设置自定义渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// 加载模板函数的方法
// 通过传入文件模式参数pattern从而将指定路径下的所有符合该模式的HTML文件加载到Engine实例的htmlTemplates字段中
func (engine *Engine) LoadHTMLGlob(pattern string) {
	//Must用于检查模板解析过程中是否出现了错误 如果有错误会引起panic
	//New创建一个名字为空的模板 这个模板用于解析服务端渲染模板中的内容
	//用Funcs方法将其与模板相关联(funcMap)
	//ParseGlob用于解析一个匹配pattern的模板文件 此时pattern是一个包含统配符的模板文件路径 函数会将匹配到的所有文件都解析为模板
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}
