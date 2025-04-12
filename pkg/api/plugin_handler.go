package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huyouba1/kde/pkg/plugin"
	pluginconfig "github.com/huyouba1/kde/pkg/plugin/config"
	"github.com/huyouba1/kde/pkg/plugin/registry"
)

// PluginHandler 插件API处理器
type PluginHandler struct {
	pluginManager  *plugin.Manager
	pluginRegistry *registry.Registry
	configManager  *pluginconfig.Manager
}

// NewPluginHandler 创建一个新的插件API处理器
func NewPluginHandler(manager *plugin.Manager, registry *registry.Registry, configManager *pluginconfig.Manager) *PluginHandler {
	return &PluginHandler{
		pluginManager:  manager,
		pluginRegistry: registry,
		configManager:  configManager,
	}
}

// RegisterRoutes 注册插件API路由
func (h *PluginHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/", h.listPlugins)
	router.GET("/:id", h.getPlugin)
	router.POST("/", h.installPlugin)
	router.DELETE("/:id", h.uninstallPlugin)
	router.PUT("/:id/enable", h.enablePlugin)
	router.PUT("/:id/disable", h.disablePlugin)
	router.GET("/:id/config", h.getPluginConfig)
	router.PUT("/:id/config", h.updatePluginConfig)
}

// PluginResponse 插件响应结构
type PluginResponse struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Description  string              `json:"description"`
	Version      string              `json:"version"`
	Author       string              `json:"author"`
	Status       plugin.PluginStatus `json:"status"`
	Enabled      bool                `json:"enabled"`
	AutoStart    bool                `json:"auto_start"`
	Type         string              `json:"type,omitempty"`
	Capabilities []string            `json:"capabilities,omitempty"`
}

// listPlugins 获取插件列表
func (h *PluginHandler) listPlugins(c *gin.Context) {
	ctx := context.Background()

	// 获取插件列表
	plugins, err := h.pluginManager.ListPlugins(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取插件列表失败: %v", err)})
		return
	}

	// 构建响应
	response := make([]PluginResponse, 0, len(plugins))
	for _, p := range plugins {
		// 获取插件配置
		config, _ := h.configManager.GetConfig(p.ID)
		enabled := false
		autoStart := false
		if config != nil {
			enabled = config.Enabled
			autoStart = config.AutoStart
		}

		// 添加到响应
		response = append(response, PluginResponse{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Version:     p.Version,
			Author:      p.Author,
			Status:      p.Status,
			Enabled:     enabled,
			AutoStart:   autoStart,
		})
	}

	c.JSON(http.StatusOK, response)
}

// getPlugin 获取插件详情
func (h *PluginHandler) getPlugin(c *gin.Context) {
	ctx := context.Background()
	pluginID := c.Param("id")

	// 获取插件列表
	plugins, err := h.pluginManager.ListPlugins(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("获取插件列表失败: %v", err)})
		return
	}

	// 查找指定插件
	var pluginInfo *plugin.PluginInfo
	for _, p := range plugins {
		if p.ID == pluginID {
			pluginInfo = p
			break
		}
	}

	if pluginInfo == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "插件不存在"})
		return
	}

	// 获取插件配置
	config, _ := h.configManager.GetConfig(pluginID)
	enabled := false
	autoStart := false
	if config != nil {
		enabled = config.Enabled
		autoStart = config.AutoStart
	}

	// 构建响应
	response := PluginResponse{
		ID:          pluginInfo.ID,
		Name:        pluginInfo.Name,
		Description: pluginInfo.Description,
		Version:     pluginInfo.Version,
		Author:      pluginInfo.Author,
		Status:      pluginInfo.Status,
		Enabled:     enabled,
		AutoStart:   autoStart,
	}

	// 如果是能力插件，添加类型和能力信息
	if p, ok := h.pluginManager.GetPlugin(pluginID); ok {
		if cp, ok := p.(plugin.CapabilityPlugin); ok {
			response.Type = string(cp.GetType())
			capabilities := cp.GetCapabilities()
			strCaps := make([]string, len(capabilities))
			for i, cap := range capabilities {
				strCaps[i] = string(cap)
			}
			response.Capabilities = strCaps
		}
	}

	c.JSON(http.StatusOK, response)
}

// InstallPluginRequest 安装插件请求
type InstallPluginRequest struct {
	Path      string `json:"path" binding:"required"`
	AutoStart bool   `json:"auto_start"`
}

// installPlugin 安装插件
func (h *PluginHandler) installPlugin(c *gin.Context) {
	ctx := context.Background()

	// 解析请求
	var req InstallPluginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的请求: %v", err)})
		return
	}

	// 安装插件
	pluginInfo, err := h.pluginManager.InstallPlugin(ctx, req.Path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("安装插件失败: %v", err)})
		return
	}

	// 创建插件配置
	config := &pluginconfig.PluginConfig{
		ID:        pluginInfo.ID,
		Enabled:   true,
		AutoStart: req.AutoStart,
		Settings:  make(map[string]interface{}),
	}

	// 保存配置
	if err := h.configManager.SaveConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("保存插件配置失败: %v", err)})
		return
	}

	// 如果设置了自动启动，则加载并启动插件
	if req.AutoStart {
		// 加载插件
		if err := h.pluginRegistry.LoadPlugin(ctx, pluginInfo); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("加载插件失败: %v", err)})
			return
		}

		// 启动插件
		if err := h.pluginRegistry.StartPlugin(ctx, pluginInfo.ID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("启动插件失败: %v", err)})
			return
		}
	}

	// 构建响应
	response := PluginResponse{
		ID:          pluginInfo.ID,
		Name:        pluginInfo.Name,
		Description: pluginInfo.Description,
		Version:     pluginInfo.Version,
		Author:      pluginInfo.Author,
		Status:      pluginInfo.Status,
		Enabled:     true,
		AutoStart:   req.AutoStart,
	}

	c.JSON(http.StatusOK, response)
}

// uninstallPlugin 卸载插件
func (h *PluginHandler) uninstallPlugin(c *gin.Context) {
	ctx := context.Background()
	pluginID := c.Param("id")

	// 停止并卸载插件
	if err := h.pluginManager.UninstallPlugin(ctx, pluginID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("卸载插件失败: %v", err)})
		return
	}

	// 删除插件配置
	if err := h.configManager.DeleteConfig(pluginID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("删除插件配置失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "插件已成功卸载"})
}

// enablePlugin 启用插件
func (h *PluginHandler) enablePlugin(c *gin.Context) {
	ctx := context.Background()
	pluginID := c.Param("id")

	// 获取插件配置
	config, exists := h.configManager.GetConfig(pluginID)
	if !exists {
		// 创建新配置
		config = &pluginconfig.PluginConfig{
			ID:       pluginID,
			Settings: make(map[string]interface{}),
		}
	}

	// 更新配置
	config.Enabled = true
	if err := h.configManager.SaveConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("保存插件配置失败: %v", err)})
		return
	}

	// 启用插件
	if err := h.pluginManager.EnablePlugin(ctx, pluginID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("启用插件失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "插件已成功启用"})
}

// disablePlugin 禁用插件
func (h *PluginHandler) disablePlugin(c *gin.Context) {
	ctx := context.Background()
	pluginID := c.Param("id")

	// 获取插件配置
	config, exists := h.configManager.GetConfig(pluginID)
	if !exists {
		// 创建新配置
		config = &pluginconfig.PluginConfig{
			ID:       pluginID,
			Settings: make(map[string]interface{}),
		}
	}

	// 更新配置
	config.Enabled = false
	if err := h.configManager.SaveConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("保存插件配置失败: %v", err)})
		return
	}

	// 禁用插件
	if err := h.pluginManager.DisablePlugin(ctx, pluginID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("禁用插件失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "插件已成功禁用"})
}

// getPluginConfig 获取插件配置
func (h *PluginHandler) getPluginConfig(c *gin.Context) {
	pluginID := c.Param("id")

	// 获取插件配置
	config, exists := h.configManager.GetConfig(pluginID)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "插件配置不存在"})
		return
	}

	c.JSON(http.StatusOK, config)
}

// UpdatePluginConfigRequest 更新插件配置请求
type UpdatePluginConfigRequest struct {
	AutoStart    *bool                  `json:"auto_start,omitempty"`
	Settings     map[string]interface{} `json:"settings,omitempty"`
	Dependencies []string               `json:"dependencies,omitempty"`
}

// updatePluginConfig 更新插件配置
func (h *PluginHandler) updatePluginConfig(c *gin.Context) {
	pluginID := c.Param("id")

	// 解析请求
	var req UpdatePluginConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("无效的请求: %v", err)})
		return
	}

	// 获取插件配置
	config, exists := h.configManager.GetConfig(pluginID)
	if !exists {
		// 创建新配置
		config = &pluginconfig.PluginConfig{
			ID:       pluginID,
			Enabled:  true,
			Settings: make(map[string]interface{}),
		}
	}

	// 更新配置
	if req.AutoStart != nil {
		config.AutoStart = *req.AutoStart
	}

	if req.Settings != nil {
		if config.Settings == nil {
			config.Settings = make(map[string]interface{})
		}
		for k, v := range req.Settings {
			config.Settings[k] = v
		}
	}

	if req.Dependencies != nil {
		config.Dependencies = req.Dependencies
	}

	// 保存配置
	if err := h.configManager.SaveConfig(config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("保存插件配置失败: %v", err)})
		return
	}

	c.JSON(http.StatusOK, config)
}
