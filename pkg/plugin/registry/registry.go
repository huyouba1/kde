package registry

import (
	"context"
	"fmt"
	"sync"

	"github.com/huyouba1/kde/pkg/plugin"
	"github.com/huyouba1/kde/pkg/storage"
)

// Registry 插件注册表
type Registry struct {
	mu             sync.RWMutex
	storageFactory storage.Factory
	pluginManager  *plugin.Manager
	pluginHooks    map[string][]PluginHook
}

// PluginHook 插件钩子函数类型
type PluginHook func(plugin.Plugin) error

// HookType 钩子类型
type HookType string

const (
	// HookBeforeInit 初始化前钩子
	HookBeforeInit HookType = "before_init"
	// HookAfterInit 初始化后钩子
	HookAfterInit HookType = "after_init"
	// HookBeforeStart 启动前钩子
	HookBeforeStart HookType = "before_start"
	// HookAfterStart 启动后钩子
	HookAfterStart HookType = "after_start"
	// HookBeforeStop 停止前钩子
	HookBeforeStop HookType = "before_stop"
	// HookAfterStop 停止后钩子
	HookAfterStop HookType = "after_stop"
)

// NewRegistry 创建一个新的插件注册表
func NewRegistry(factory storage.Factory, manager *plugin.Manager) *Registry {
	return &Registry{
		storageFactory: factory,
		pluginManager:  manager,
		pluginHooks:    make(map[string][]PluginHook),
	}
}

// RegisterHook 注册插件钩子
func (r *Registry) RegisterHook(hookType HookType, hook PluginHook) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := string(hookType)
	r.pluginHooks[key] = append(r.pluginHooks[key], hook)
}

// ExecuteHooks 执行指定类型的所有钩子
func (r *Registry) ExecuteHooks(hookType HookType, p plugin.Plugin) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := string(hookType)
	hooks, ok := r.pluginHooks[key]
	if !ok {
		return nil
	}

	for _, hook := range hooks {
		if err := hook(p); err != nil {
			return err
		}
	}

	return nil
}

// LoadPlugin 加载单个插件并执行钩子
func (r *Registry) LoadPlugin(ctx context.Context, info *plugin.PluginInfo) error {
	// 执行初始化前钩子
	if err := r.ExecuteHooks(HookBeforeInit, nil); err != nil {
		return fmt.Errorf("执行初始化前钩子失败: %v", err)
	}

	// 加载插件
	err := r.pluginManager.LoadPlugin(ctx, info)
	if err != nil {
		return err
	}

	// 获取加载的插件实例
	p, ok := r.pluginManager.GetPlugin(info.ID)
	if !ok {
		return fmt.Errorf("插件加载后无法获取实例: %s", info.ID)
	}

	// 执行初始化后钩子
	if err := r.ExecuteHooks(HookAfterInit, p); err != nil {
		return fmt.Errorf("执行初始化后钩子失败: %v", err)
	}

	return nil
}

// StartPlugin 启动插件并执行钩子
func (r *Registry) StartPlugin(ctx context.Context, pluginID string) error {
	// 获取插件实例
	p, ok := r.pluginManager.GetPlugin(pluginID)
	if !ok {
		return fmt.Errorf("插件未加载: %s", pluginID)
	}

	// 执行启动前钩子
	if err := r.ExecuteHooks(HookBeforeStart, p); err != nil {
		return fmt.Errorf("执行启动前钩子失败: %v", err)
	}

	// 启动插件
	if err := p.Start(); err != nil {
		return err
	}

	// 执行启动后钩子
	if err := r.ExecuteHooks(HookAfterStart, p); err != nil {
		return fmt.Errorf("执行启动后钩子失败: %v", err)
	}

	return nil
}

// StopPlugin 停止插件并执行钩子
func (r *Registry) StopPlugin(ctx context.Context, pluginID string) error {
	// 获取插件实例
	p, ok := r.pluginManager.GetPlugin(pluginID)
	if !ok {
		return fmt.Errorf("插件未加载: %s", pluginID)
	}

	// 执行停止前钩子
	if err := r.ExecuteHooks(HookBeforeStop, p); err != nil {
		return fmt.Errorf("执行停止前钩子失败: %v", err)
	}

	// 停止插件
	if err := p.Stop(); err != nil {
		return err
	}

	// 执行停止后钩子
	if err := r.ExecuteHooks(HookAfterStop, p); err != nil {
		return fmt.Errorf("执行停止后钩子失败: %v", err)
	}

	return nil
}

// GetPluginByType 根据插件类型获取插件
func (r *Registry) GetPluginByType(pluginType string) ([]plugin.Plugin, error) {
	// TODO: 实现根据类型获取插件的逻辑
	return nil, nil
}

// GetPluginByCapability 根据插件能力获取插件
func (r *Registry) GetPluginByCapability(capability string) ([]plugin.Plugin, error) {
	// TODO: 实现根据能力获取插件的逻辑
	return nil, nil
}
