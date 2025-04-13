package api

import (
	"context"
	"fmt"
	"github.com/huyouba1/kde/configs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huyouba1/kde/pkg/api/handler"
	"github.com/huyouba1/kde/pkg/storage"
	"github.com/huyouba1/kde/pkg/storage/models"
)

// Server API服务器
type Server struct {
	config          *configs.Config
	router          *gin.Engine
	httpServer      *http.Server
	storageFactory  *storage.Factory
	templateHandler *handler.TemplateHandler
}

// NewServer 创建一个新的API服务器
// func NewServer(cfg *configs.Config) (*Server, error) {
func NewServer(cfg *configs.Config) (*Server, error) {
	// 创建Gin路由器
	router := gin.Default()

	// 创建存储工厂
	storageFactory := storage.NewFactory(cfg)

	// 创建模板处理器
	templateHandler, err := handler.NewTemplateHandler()
	if err != nil {
		return nil, fmt.Errorf("failed to create template handler: %w", err)
	}

	// 设置静态文件服务
	router.Static("/static", "pkg/api/handler/static")

	// 设置模板渲染
	router.SetHTMLTemplate(templateHandler.GetTemplates())

	// 创建服务器
	server := &Server{
		config:          cfg,
		router:          router,
		storageFactory:  storageFactory,
		templateHandler: templateHandler,
	}

	// 初始化路由
	server.initRoutes()

	return server, nil
}

// initRoutes 初始化API路由
func (s *Server) initRoutes() {
	// 前端页面路由
	s.router.GET("/", s.handleIndex)
	s.router.GET("/clusters", s.handleClusters)
	s.router.GET("/deployments", s.handleDeployments)
	s.router.GET("/settings", s.handleSettings)

	// API版本前缀
	api := s.router.Group("/api/v1")

	// 健康检查
	api.GET("/health", s.healthCheck)

	// 仪表板数据API
	api.GET("/dashboard", s.handleDashboardData)

	// 集群管理API
	cluster := api.Group("/clusters")
	{
		cluster.GET("/", s.listClusters)
		cluster.POST("/", s.createCluster)
		cluster.GET("/:id", s.getCluster)
		cluster.PUT("/:id", s.updateCluster)
		cluster.DELETE("/:id", s.deleteCluster)
	}

	// 部署API
	deploy := api.Group("/deploy")
	{
		deploy.POST("/single", s.deploySingleNode)
		deploy.POST("/cluster", s.deployCluster)
		deploy.GET("/packages", s.listOfflinePackages)
		deploy.POST("/packages", s.createOfflinePackage)
	}

	// 应用交付API
	delivery := api.Group("/delivery")
	{
		delivery.POST("/yaml", s.deployYaml)
		delivery.POST("/helm", s.deployHelm)
		delivery.POST("/kustomize", s.deployKustomize)
	}

	// 插件API
	plugin := api.Group("/plugins")
	{
		plugin.GET("/", s.listPlugins)
		plugin.POST("/", s.installPlugin)
		plugin.DELETE("/:name", s.uninstallPlugin)
	}
}

// Start 启动API服务器
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)

	s.httpServer = &http.Server{
		Addr:    addr,
		Handler: s.router,
	}

	fmt.Printf("API服务器启动在 %s\n", addr)
	return s.httpServer.ListenAndServe()
}

// Stop 停止API服务器
func (s *Server) Stop() error {
	// 创建一个5秒超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if err := s.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	// 关闭存储连接
	return s.storageFactory.Close()
}

// 页面处理函数
func (s *Server) handleIndex(c *gin.Context) {
	// 获取集群数量
	var clusterCount int64
	if err := s.storageFactory.GetDB().Model(&models.ClusterModel{}).Count(&clusterCount).Error; err != nil {
		clusterCount = 0
	}

	// 获取部署数量（TODO: 实现实际的部署计数）
	deploymentCount := 0

	// 检查系统状态
	systemHealthy := true // TODO: 实现实际的健康检查

	c.HTML(http.StatusOK, "index.html", gin.H{
		"ClusterCount":    clusterCount,
		"DeploymentCount": deploymentCount,
		"SystemHealthy":   systemHealthy,
	})
}

func (s *Server) handleClusters(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Active": "clusters",
	})
}

func (s *Server) handleDeployments(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Active": "deployments",
	})
}

func (s *Server) handleSettings(c *gin.Context) {
	c.HTML(http.StatusOK, "base.html", gin.H{
		"Active": "settings",
	})
}

// healthCheck 健康检查处理函数
func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// API处理函数
func (s *Server) handleDashboardData(c *gin.Context) {
	// TODO: 从存储中获取实际数据
	c.JSON(http.StatusOK, gin.H{
		"clusterCount":    0,
		"deploymentCount": 0,
		"systemStatus":    "正常",
	})
}

// 集群管理处理函数
func (s *Server) listClusters(c *gin.Context) {
	// TODO: 实现集群列表获取
	clusters := []map[string]interface{}{
		{
			"id":        "cluster-1",
			"status":    "Running",
			"nodeCount": 3,
			"createdAt": time.Now().Add(-24 * time.Hour),
		},
		{
			"id":        "cluster-2",
			"status":    "Creating",
			"nodeCount": 1,
			"createdAt": time.Now().Add(-1 * time.Hour),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"clusters": clusters,
	})
}

func (s *Server) createCluster(c *gin.Context) {
	var data struct {
		Name       string `json:"name"`
		NodeCount  int    `json:"nodeCount"`
		DeployType string `json:"deployType"`
	}

	if err := c.BindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的请求数据",
		})
		return
	}

	// TODO: 实现集群创建
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("集群 %s 创建请求已提交", data.Name),
	})
}

func (s *Server) getCluster(c *gin.Context) {
	id := c.Param("id")
	// TODO: 实现获取集群详情
	c.JSON(http.StatusOK, gin.H{
		"id":        id,
		"status":    "Running",
		"nodeCount": 3,
		"createdAt": time.Now().Add(-24 * time.Hour),
	})
}

func (s *Server) updateCluster(c *gin.Context) {
	id := c.Param("id")
	// TODO: 实现更新集群
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("集群 %s 更新成功", id),
	})
}

func (s *Server) deleteCluster(c *gin.Context) {
	id := c.Param("id")
	// TODO: 实现删除集群
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("集群 %s 删除请求已提交", id),
	})
}

// 部署相关处理函数
func (s *Server) deploySingleNode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "单节点部署请求已提交",
	})
}

func (s *Server) deployCluster(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "集群部署请求已提交",
	})
}

func (s *Server) listOfflinePackages(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"packages": []string{},
	})
}

func (s *Server) createOfflinePackage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "离线包创建请求已提交",
	})
}

// 应用交付相关处理函数
func (s *Server) deployYaml(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "YAML部署请求已提交",
	})
}

func (s *Server) deployHelm(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Helm部署请求已提交",
	})
}

func (s *Server) deployKustomize(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Kustomize部署请求已提交",
	})
}

// 插件相关处理函数
func (s *Server) listPlugins(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"plugins": []string{},
	})
}

func (s *Server) installPlugin(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "插件安装请求已提交",
	})
}

func (s *Server) uninstallPlugin(c *gin.Context) {
	name := c.Param("name")
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("插件 %s 卸载请求已提交", name),
	})
}
