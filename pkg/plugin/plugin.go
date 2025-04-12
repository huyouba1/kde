package plugin

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"time"

	"github.com/huyouba1/kde/pkg/storage"
)

// PluginStatus 插件状态
type PluginStatus string

const (
	// StatusEnabled 启用状态
	StatusEnabled PluginStatus = "enabled"
	// StatusDisabled 禁用状态
	StatusDisabled PluginStatus = "disabled"
	// StatusError 错误状态
	StatusError PluginStatus = "error"
)

// PluginInfo 插件信息
type PluginInfo struct {
	ID          string       `json:"id" gorm:"primaryKey"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Version     string       `json:"version"`
	Author      string       `json:"author"`
	Path        string       `json:"path"`
	Status      PluginStatus `json:"status"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// Plugin 插件接口
type Plugin interface {
	// GetInfo 获取插件信息
	GetInfo() PluginInfo
	// Init 初始化插件
	Init() error
	// Start 启动插件
	Start() error
	// Stop 停止插件
	Stop() error
}

// Manager 插件管理器
type Manager struct {
	storageFactory storage.Factory
	pluginDir      string
	plugins        map[string]Plugin
}

// NewManager 创建一个新的插件管理器
func NewManager(factory storage.Factory, pluginDir string) *Manager {
	return &Manager{
		storageFactory: factory,
		pluginDir:      pluginDir,
		plugins:        make(map[string]Plugin),
	}
}

// LoadPlugins 加载所有插件
func (m *Manager) LoadPlugins(ctx context.Context) error {
	// TODO: 从数据库获取插件列表
	pluginInfos := []*PluginInfo{}

	for _, info := range pluginInfos {
		if info.Status == StatusEnabled {
			if err := m.LoadPlugin(ctx, info); err != nil {
				// 记录错误但继续加载其他插件
				fmt.Printf("加载插件 %s 失败: %v\n", info.Name, err)
			}
		}
	}

	return nil
}

// LoadPlugin 加载单个插件
func (m *Manager) LoadPlugin(ctx context.Context, info *PluginInfo) error {
	// 构建插件路径
	plugPath := info.Path
	if !filepath.IsAbs(plugPath) {
		plugPath = filepath.Join(m.pluginDir, plugPath)
	}

	// 打开插件
	plug, err := plugin.Open(plugPath)
	if err != nil {
		return fmt.Errorf("打开插件失败: %v", err)
	}

	// 查找插件符号
	sym, err := plug.Lookup("Plugin")
	if err != nil {
		return fmt.Errorf("查找插件符号失败: %v", err)
	}

	// 类型断言
	var p Plugin
	p, ok := sym.(Plugin)
	if !ok {
		return fmt.Errorf("插件不实现Plugin接口")
	}

	// 初始化插件
	if err := p.Init(); err != nil {
		return fmt.Errorf("初始化插件失败: %v", err)
	}

	// 启动插件
	if err := p.Start(); err != nil {
		return fmt.Errorf("启动插件失败: %v", err)
	}

	// 存储插件实例
	m.plugins[info.ID] = p

	return nil
}

// GetPlugin 获取插件实例
func (m *Manager) GetPlugin(id string) (Plugin, bool) {
	p, ok := m.plugins[id]
	return p, ok
}

// ListPlugins 获取所有插件信息
func (m *Manager) ListPlugins(ctx context.Context) ([]*PluginInfo, error) {
	// TODO: 从数据库获取插件列表
	return nil, fmt.Errorf("未实现")
}

// InstallPlugin 安装插件
func (m *Manager) InstallPlugin(ctx context.Context, path string) (*PluginInfo, error) {
	// 打开插件以验证
	plug, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("打开插件失败: %v", err)
	}

	// 查找插件符号
	sym, err := plug.Lookup("Plugin")
	if err != nil {
		return nil, fmt.Errorf("查找插件符号失败: %v", err)
	}

	// 类型断言
	p, ok := sym.(Plugin)
	if !ok {
		return nil, fmt.Errorf("插件不实现Plugin接口")
	}

	// 获取插件信息
	info := p.GetInfo()

	// 设置插件状态和时间
	info.Status = StatusEnabled
	info.CreatedAt = time.Now()
	info.UpdatedAt = time.Now()

	// TODO: 保存插件信息到数据库

	return &info, nil
}

// UninstallPlugin 卸载插件
func (m *Manager) UninstallPlugin(ctx context.Context, id string) error {
	// 检查插件是否已加载
	p, ok := m.plugins[id]
	if ok {
		// 停止插件
		if err := p.Stop(); err != nil {
			return fmt.Errorf("停止插件失败: %v", err)
		}

		// 从内存中删除
		delete(m.plugins, id)
	}

	// TODO: 从数据库删除插件信息

	return nil
}

// EnablePlugin 启用插件
func (m *Manager) EnablePlugin(ctx context.Context, id string) error {
	// TODO: 从数据库获取插件信息
	info := &PluginInfo{}

	// 加载并启动插件
	if err := m.LoadPlugin(ctx, info); err != nil {
		return err
	}

	// 更新插件状态
	info.Status = StatusEnabled
	info.UpdatedAt = time.Now()

	// TODO: 更新插件信息到数据库

	return nil
}

// DisablePlugin 禁用插件
func (m *Manager) DisablePlugin(ctx context.Context, id string) error {
	// 检查插件是否已加载
	p, ok := m.plugins[id]
	if !ok {
		return fmt.Errorf("插件未加载")
	}

	// 停止插件
	if err := p.Stop(); err != nil {
		return fmt.Errorf("停止插件失败: %v", err)
	}

	// 从内存中删除
	delete(m.plugins, id)

	// TODO: 更新插件状态到数据库

	return nil
}
