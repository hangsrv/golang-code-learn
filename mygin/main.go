package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 设置前端静态文件路由
	router.Static("/web", "./public")

	// 设置后端接口路由
	router.GET("/api/test", func(c *gin.Context) {
		// 处理后端接口逻辑
		c.JSON(200, gin.H{
			"message": "Hello from backend API!",
		})
	})

	router.Run(":18080")
}
