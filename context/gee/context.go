package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 便于进行处理请求和构造响应
type H map[string]interface{}

// 定义一个请求上下文结构体
// 本结构体包含了请求和响应相关的信息
type Context struct {
	//包含Writer和Req两个原始的HTTP对象
	Writer     http.ResponseWriter //ResponseWriter是一个接口 用于向客户端发送HTTP响应
	Req        *http.Request       //Request 是一个结构体类型 代表一个客户端的HTTP请求
	Path       string              //请求路径
	Method     string              //请求方式  列如GET POST PUT DELETE
	StatusCode int                 //响应状态码
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
