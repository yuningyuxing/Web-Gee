package main

import (
	"fmt"
	"gee"
	"net/http"
	"time"
)

// 定义一个结构体 表示一个学生
type student struct {
	Name string
	Age  int8
}

// 将时间格式化为字符串
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	r := gee.Default()
	r.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "Hello yuning\n")
	})
	r.GET("/panic", func(c *gee.Context) {
		names := []string{"yuning"}
		//故意引起越界错误从而测试
		c.String(http.StatusOK, names[100])
	})
	r.Run(":9999")
	//r := gee.New()
	//r.Use(gee.Logger())
	////将上面函数添加到FuncMap
	////设置模板函数映射表
	//r.SetFuncMap(template.FuncMap{
	//	"FormatAsDate": FormatAsDate,
	//})
	////加载指定目录下所有的HTML模板的文件
	////加载到htmlTemplates
	//r.LoadHTMLGlob("templates/*")
	////将指定目录 也就是./static(我们自己建立的目录)下的所有静态文件映射到/assets路由上
	////这样我们就可以通过assets访问到我磁盘中static中的文件了
	//r.Static("/assets", "./static")
	//
	////创建两个学生对象
	//stu1 := &student{Name: "Geektutu", Age: 20}
	//stu2 := &student{Name: "Jack", Age: 22}
	//
	//r.GET("/", func(c *gee.Context) {
	//	//使用HTML方法渲染名为css.tmpl的模板文件 并将渲染结果返回给客户端
	//	c.HTML(http.StatusOK, "css.tmpl", nil)
	//})
	//
	//r.GET("/students", func(c *gee.Context) {
	//	//同理渲染arr.tmpl模板文件  并传递两个参数
	//	c.HTML(http.StatusOK, "arr.tmpl", gee.H{
	//		"title":  "gee",
	//		"stuArr": [2]*student{stu1, stu2},
	//	})
	//})
	//
	//r.GET("/date", func(c *gee.Context) {
	//	//与上面同理
	//	c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
	//		"title": "gee",
	//		"now":   time.Date(2023, 4, 10, 0, 0, 0, 0, time.UTC),
	//	})
	//})
	//
	//r.Run(":9999")
}
