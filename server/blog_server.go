package main

import (
	"blog/server/model"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		blogs, err := model.GetPage("1")
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
		id := c.Param("id")
		c.JSON(200, gin.H{
			"ok":     true,
			"reason": "",
			"id":     id,
		})
	})
	router.Run(":9092")
}
