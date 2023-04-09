package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		t := time.Now()
		//表示错误了
		c.Fail(500, "Internal Server Error")
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
func main() {
	r := gee.New()
	//为所有请求添加中间件
	r.Use(gee.Logger())
	//原始方法可以addroute
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2")
	//为该group添加中间件
	//这个中间件会阻止handler运行
	v2.Use(onlyForV2())

	v2.GET("/hello/:name", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s,you're at %s\n", c.Param("name"), c.Path)
	})

	r.Run(":9999")
}
