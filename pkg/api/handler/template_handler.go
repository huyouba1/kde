package handler

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed static
var staticFS embed.FS

// TemplateHandler 处理模板和静态文件
type TemplateHandler struct {
	templates *template.Template
	staticFS  embed.FS
}

// NewTemplateHandler 创建新的模板处理器
func NewTemplateHandler() (*TemplateHandler, error) {
	// 使用工作目录相对路径加载模板
	pattern := "./pkg/api/templates/*.html"
	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	return &TemplateHandler{
		templates: tmpl,
		staticFS:  staticFS,
	}, nil
}

// GetTemplates 获取模板
func (h *TemplateHandler) GetTemplates() *template.Template {
	return h.templates
}

// Render 渲染模板
func (h *TemplateHandler) Render(w http.ResponseWriter, name string, data interface{}) error {
	return h.templates.ExecuteTemplate(w, name, data)
}

// ServeStatic 提供静态文件服务
func (h *TemplateHandler) ServeStatic(c *gin.Context) {
	// 移除 /static/ 前缀
	path := strings.TrimPrefix(c.Request.URL.Path, "/static/")

	// 构建完整的文件路径
	filePath := filepath.Join("static", path)

	file, err := h.staticFS.ReadFile(filePath)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	// 设置正确的 Content-Type
	ext := filepath.Ext(path)
	switch ext {
	case ".css":
		c.Header("Content-Type", "text/css")
	case ".js":
		c.Header("Content-Type", "application/javascript")
	case ".png", ".jpg", ".jpeg", ".gif":
		c.Header("Content-Type", "image/"+strings.TrimPrefix(ext, "."))
	default:
		c.Header("Content-Type", "text/plain")
	}

	c.Data(http.StatusOK, c.Writer.Header().Get("Content-Type"), file)
}
