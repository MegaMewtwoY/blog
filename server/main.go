package main

import (
	"blog/server/model"
	_ "fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ok": true,
		})
	})

	router.GET("/page/:page", func(c *gin.Context) {
		page := c.Param("page")
		c.JSON(200, gin.H{
			"ok":    true,
			"blogs": model.GetPage(page),
		})
	})

	router.GET("/blog/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{
			"ok": true,
			"id": id,
		})
	})
	router.Run(":9092")
}
