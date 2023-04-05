package main

import (
	"fmt"
	"gee"
	"net/http"
)

//我们通过创建gee实例 来将之前过程封装 此时直接调用gee的方法即可
//此时只支持静态路由
func main() {
	r := gee.New()
	//添加路由
	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	})
	r.GET("/hello", func(w http.ResponseWriter, req *http.Request) {
		for k, v := range req.Header {
			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
		}
	})
	//启动Web服务
	r.Run(":9999")
}
