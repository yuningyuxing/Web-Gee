package gee

import (
	"log"
	"time"
)

// 这个中间件的作用是计时
func Logger() HandlerFunc {
	return func(c *Context) {
		t := time.Now()
		//在本样例中我们先得到t然后开始执行下一个中间件也就是handler 然后再回来执行log
		c.Next()
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}
