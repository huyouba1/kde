package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/huyouba1/kde/pkg/api/handlers"
)

// SetupRoutes 设置所有路由
func SetupRoutes(router *gin.Engine) {
	// 静态文件服务
	router.Static("/static", "static")

	// API 路由
	api := router.Group("/api/v1")

	// 集群管理路由
	api.GET("/clusters", handlers.GetClusters)
	api.POST("/deploy/cluster", handlers.CreateCluster)
	api.DELETE("/clusters/:id/delete", handlers.DeleteCluster)

	// 首页路由
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/clusters")
	})

	// 集群管理页面
	router.GET("/clusters", func(c *gin.Context) {
		c.File("pkg/api/templates/clusters.html")
	})
}
