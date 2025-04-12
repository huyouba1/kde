package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

// PluginConfig 插件配置
type PluginConfig struct {
	ID           string                 `json:"id"`
	Enabled      bool                   `json:"enabled"`
	AutoStart    bool                   `json:"auto_start"`
	Settings     map[string]interface{} `json:"settings"`
	Dependencies []string               `json:"dependencies"`
}

// Manager 插件配置管理器
type Manager struct {
	mu          sync.RWMutex
	configDir   string
	configs     map[string]*PluginConfig
	configFiles map[string]string
}

// NewManager 创建一个新的插件配置管理器
func NewManager(configDir string) (*Manager, error) {
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %v", err)
	}

	return &Manager{
		configDir:   configDir,
		configs:     make(map[string]*PluginConfig),
		configFiles: make(map[string]string),
	}, nil
}

// LoadAllConfigs 加载所有插件配置
func (m *Manager) LoadAllConfigs() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 清空现有配置
	m.configs = make(map[string]*PluginConfig)
	m.configFiles = make(map[string]string)

	// 遍历配置目录
	files, err := ioutil.ReadDir(m.configDir)
	if err != nil {
		return fmt.Errorf("读取配置目录失败: %v", err)
	}

	// 加载每个配置文件
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(m.configDir, file.Name())
		config, err := m.loadConfigFile(filePath)
		if err != nil {
			fmt.Printf("加载配置文件 %s 失败: %v\n", filePath, err)
			continue
		}

		m.configs[config.ID] = config
		m.configFiles[config.ID] = filePath
	}

	return nil
}

// loadConfigFile 加载单个配置文件
func (m *Manager) loadConfigFile(filePath string) (*PluginConfig, error) {
	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析JSON
	config := &PluginConfig{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return config, nil
}

// GetConfig 获取插件配置
func (m *Manager) GetConfig(pluginID string) (*PluginConfig, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, ok := m.configs[pluginID]
	return config, ok
}

// SaveConfig 保存插件配置
func (m *Manager) SaveConfig(config *PluginConfig) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 更新内存中的配置
	m.configs[config.ID] = config

	// 确定配置文件路径
	filePath, ok := m.configFiles[config.ID]
	if !ok {
		filePath = filepath.Join(m.configDir, fmt.Sprintf("%s.json", config.ID))
		m.configFiles[config.ID] = filePath
	}

	// 序列化为JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	// 写入文件
	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	return nil
}

// DeleteConfig 删除插件配置
func (m *Manager) DeleteConfig(pluginID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查配置是否存在
	filePath, ok := m.configFiles[pluginID]
	if !ok {
		return nil // 配置不存在，视为成功
	}

	// 删除文件
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除配置文件失败: %v", err)
	}

	// 从内存中删除
	delete(m.configs, pluginID)
	delete(m.configFiles, pluginID)

	return nil
}

// GetAutoStartPlugins 获取自动启动的插件列表
func (m *Manager) GetAutoStartPlugins() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var pluginIDs []string
	for id, config := range m.configs {
		if config.Enabled && config.AutoStart {
			pluginIDs = append(pluginIDs, id)
		}
	}

	return pluginIDs
}

// GetPluginSetting 获取插件设置
func (m *Manager) GetPluginSetting(pluginID, key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	config, ok := m.configs[pluginID]
	if !ok || config.Settings == nil {
		return nil, false
	}

	val, ok := config.Settings[key]
	return val, ok
}

// SetPluginSetting 设置插件设置
func (m *Manager) SetPluginSetting(pluginID, key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 获取配置，如果不存在则创建
	config, ok := m.configs[pluginID]
	if !ok {
		config = &PluginConfig{
			ID:       pluginID,
			Enabled:  true,
			Settings: make(map[string]interface{}),
		}
		m.configs[pluginID] = config
	}

	// 确保Settings不为nil
	if config.Settings == nil {
		config.Settings = make(map[string]interface{})
	}

	// 设置值
	config.Settings[key] = value

	// 保存配置
	return m.SaveConfig(config)
}
