package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

// 给Context对象新增一个属性和方法 从而提供对路由参数的访问
type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	Params     map[string]string
	StatusCode int
}

// 新增方法 用于获取解析的参数
func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

// 初始化上下文结构体
func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

// 获取POST请求表单中指定key的值
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

// 获取URL中指定key的值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 设置响应状态码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置响应头
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 返回字符串响应类型
func (c *Context) String(code int, format string, values ...interface{}) {
	//text/plain表示响应体是纯文本 没有特定格式 可以直接展示给客户阅读
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	//format代表格式
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// 返回JSON响应类型
func (c *Context) JSON(code int, obj interface{}) {
	//application/json表示响应体是JSON格式的数据 常用于API返回数据的格式
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	//创建JSON编码器
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// 返回二进制类型的数据
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	//text/html表示响应体是HTML格式的文本 常用于网页渲染
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
