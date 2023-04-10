package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 打印堆栈信息
func trace(message string) string {
	var pcs [32]uintptr
	//Callers用来返回调用栈的程序计数器
	//这里第0个是本身 第1个是上一层trace 第2个是defer func为了日志简洁 我们跳过前三个
	//第二个参数是返回调用栈信息的数组 该函数会填充该数组并返回实际的填充数
	n := runtime.Callers(3, pcs[:])
	//这是一种字符串缓存类型 实现了io.Write接口  提供了往缓存中写入字符串的方法
	//可高效的拼接字符串
	var str strings.Builder
	//写入传入消息和traceback
	str.WriteString(message + "\nTraceback:")
	//循环得到所有的调用栈过程
	for _, pc := range pcs[:n] {
		//FuncForPC的作用是返回一个给定的函数指针所在函数的函数值
		fn := runtime.FuncForPC(pc)
		//返回PC对应的文件名和行号  这里就获取调用栈中每个帧的文件名和行号
		file, line := fn.FileLine(pc)
		//写入
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	//返回得到调用栈的信息
	return str.String()
}

// 实现中间件用于错误处理  使用defer挂载上错误恢复的函数 在这个函数调用recover捕获panic
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				//trace函数是用来获取触发panic的堆栈信息
				log.Printf("%s\n\n", trace(message))
				c.Fail(http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		c.Next()
	}
}
