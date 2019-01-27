package main

import (
	"blog/server/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"html/template"
)

func Format(input string) string {
	return input
}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("../client/info.html", "../client/info_frame.html")

	// router.Delims("{[{", "}]}")
	// 自定制模板转换函数
	// router.SetFuncMap(template.FuncMap{
	// 	"Format": Format,
	// })

	router.GET("/page/:page", func(c *gin.Context) {
		page := c.Param("page")
		blogs, err := model.GetPage(page)
		if err != nil {
			fmt.Println(err)
			c.JSON(404, gin.H{
				"ok":     false,
				"reason": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"ok":     true,
				"reason": "",
				"blogs":  blogs,
			})
		}
	})

	router.GET("/blog/:id", func(c *gin.Context) {
		name := c.Param("id")
		c.HTML(200, "info_frame.html", gin.H{
			"name": "/blog_content/" + name,
		})
	})

	router.GET("/blog_content/:id", func(c *gin.Context) {
		name := c.Param("id")
		blog, err := model.GetBlog(name)
		if err != nil {
			fmt.Println(err)
			c.JSON(404, gin.H{
				"ok":     false,
				"reason": err.Error(),
			})
		} else {
			c.HTML(200, "info.html", gin.H{
				"content": template.HTML(blog.Content), // 关键性的 template.HTML 操作, 避免html标签被转义
			})
		}
	})
	router.Static("/static/", "../client/")
	// router.StaticFile("index.html", "../client/index.html")
	router.Run(":9092")
}
