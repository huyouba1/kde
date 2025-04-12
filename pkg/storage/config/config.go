package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Config 表示系统的配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Deploy   DeployConfig   `yaml:"deploy"`
	Delivery DeliveryConfig `yaml:"delivery"`
	Log      LogConfig      `yaml:"log"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Host string `yaml:"host"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type   string       `yaml:"type"`
	SQLite SQLiteConfig `yaml:"sqlite"`
	Etcd   EtcdConfig   `yaml:"etcd"`
}

// SQLiteConfig SQLite数据库配置
type SQLiteConfig struct {
	Path string `yaml:"path"`
}

// EtcdConfig etcd配置
type EtcdConfig struct {
	Endpoints   []string `yaml:"endpoints"`
	DialTimeout int      `yaml:"dialTimeout"`
}

// DeployConfig 部署配置
type DeployConfig struct {
	Ansible   AnsibleConfig   `yaml:"ansible"`
	Container ContainerConfig `yaml:"container"`
}

// AnsibleConfig Ansible配置
type AnsibleConfig struct {
	InventoryPath string `yaml:"inventoryPath"`
	PlaybooksPath string `yaml:"playbooksPath"`
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	Registry  string `yaml:"registry"`
	Namespace string `yaml:"namespace"`
}

// DeliveryConfig 应用交付配置
type DeliveryConfig struct {
	Helm    HelmConfig `yaml:"helm"`
	Workdir string     `yaml:"workdir"`
}

// HelmConfig Helm配置
type HelmConfig struct {
	RepoUrl   string `yaml:"repoUrl"`
	CachePath string `yaml:"cachePath"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"`
	Path       string `yaml:"path"`
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"maxSize"`
	MaxBackups int    `yaml:"maxBackups"`
	MaxAge     int    `yaml:"maxAge"`
	Compress   bool   `yaml:"compress"`
}

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", path)
	}

	// 读取配置文件
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return cfg, nil
}
