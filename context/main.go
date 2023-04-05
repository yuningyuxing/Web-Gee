package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()
	//下面创造三个路由类型
	//返回一个HTML页面
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	//返回一个包含请求参数中name值的字符串和请求路径的字符串
	r.GET("/hello", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})
	//返回一个JSON格式的响应
	r.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})
	r.Run(":9999")
}
