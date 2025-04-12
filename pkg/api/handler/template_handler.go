package handler

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
)

//go:embed static/*
var staticFiles embed.FS

// TemplateHandler 处理模板和静态文件
type TemplateHandler struct {
	templates *template.Template
	staticFS  fs.FS
}

// NewTemplateHandler 创建新的模板处理器
func NewTemplateHandler() (*TemplateHandler, error) {
	// 获取模板目录
	templateDir := "pkg/api/templates"

	// 加载所有模板文件
	templates, err := template.ParseGlob(filepath.Join(templateDir, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("加载模板失败: %v", err)
	}

	// 获取静态文件系统
	staticFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		return nil, fmt.Errorf("初始化静态文件系统失败: %v", err)
	}

	return &TemplateHandler{
		templates: templates,
		staticFS:  staticFS,
	}, nil
}

// GetTemplates 获取所有模板
func (h *TemplateHandler) GetTemplates() *template.Template {
	return h.templates
}

// GetTemplate 获取指定名称的模板
func (h *TemplateHandler) GetTemplate(name string) (*template.Template, error) {
	tmpl := h.templates.Lookup(name)
	if tmpl == nil {
		return nil, fmt.Errorf("模板 %s 不存在", name)
	}
	return tmpl, nil
}

// Render 渲染模板
func (h *TemplateHandler) Render(name string, data interface{}) (string, error) {
	tmpl, err := h.GetTemplate(name)
	if err != nil {
		return "", err
	}

	// 创建一个缓冲区来存储渲染结果
	var buf strings.Builder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// ServeStaticFiles 处理静态文件请求
func (h *TemplateHandler) ServeStaticFiles(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.FS(h.staticFS)).ServeHTTP(w, r)
}
