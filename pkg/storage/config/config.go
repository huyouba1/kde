package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Deploy     DeployConfig     `mapstructure:"deploy"`
	Delivery   DeliveryConfig   `mapstructure:"delivery"`
	Log        LogConfig        `mapstructure:"log"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Kubernetes KubernetesConfig `mapstructure:"kubernetes"`
	Cache      CacheConfig      `mapstructure:"cache"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type   string       `mapstructure:"type"`
	SQLite SQLiteConfig `mapstructure:"sqlite"`
	Etcd   EtcdConfig   `mapstructure:"etcd"`
}

// SQLiteConfig SQLite数据库配置
type SQLiteConfig struct {
	Path string `mapstructure:"path"`
}

// EtcdConfig etcd配置
type EtcdConfig struct {
	Endpoints   []string `mapstructure:"endpoints"`
	DialTimeout int      `mapstructure:"dialTimeout"`
}

// DeployConfig 部署配置
type DeployConfig struct {
	Ansible   AnsibleConfig   `mapstructure:"ansible"`
	Container ContainerConfig `mapstructure:"container"`
}

// AnsibleConfig Ansible配置
type AnsibleConfig struct {
	InventoryPath string `mapstructure:"inventoryPath"`
	PlaybooksPath string `mapstructure:"playbooksPath"`
}

// ContainerConfig 容器配置
type ContainerConfig struct {
	Registry  string `mapstructure:"registry"`
	Namespace string `mapstructure:"namespace"`
}

// DeliveryConfig 应用交付配置
type DeliveryConfig struct {
	Helm    HelmConfig `mapstructure:"helm"`
	Workdir string     `mapstructure:"workdir"`
}

// HelmConfig Helm配置
type HelmConfig struct {
	RepoUrl   string `mapstructure:"repoUrl"`
	CachePath string `mapstructure:"cachePath"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret   string `mapstructure:"jwt_secret"`
	TokenExpiry string `mapstructure:"token_expiry"`
}

// KubernetesConfig Kubernetes配置
type KubernetesConfig struct {
	ConfigPath string `mapstructure:"config_path"`
	InCluster  bool   `mapstructure:"in_cluster"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Type string `mapstructure:"type"`
	Size int    `mapstructure:"size"`
	TTL  string `mapstructure:"ttl"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置默认值
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	if config.Server.Host == "" {
		config.Server.Host = "0.0.0.0"
	}
	if config.Database.Type == "" {
		config.Database.Type = "sqlite"
	}
	if config.Database.SQLite.Path == "" {
		config.Database.SQLite.Path = "data/kde.db"
	}
	if config.Log.Level == "" {
		config.Log.Level = "info"
	}
	if config.Log.Format == "" {
		config.Log.Format = "text"
	}
	if config.Log.Output == "" {
		config.Log.Output = "stdout"
	}

	return &config, nil
}
