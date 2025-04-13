package configs

import (
	"fmt"
	"github.com/huyouba1/kde/pkg/k8s"
	"github.com/huyouba1/kde/pkg/storage/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"path/filepath"
	"sync"
)

// 全局Confifg实例对象

// 为了不被程序在运行时恶意修改,设置成私有变量
var global *Config

// 全局Mysql 客户端实例
var (
	gdb    *gorm.DB
	client *k8s.Client
)

// 全局Config对象获取函数

func C() *Config {
	if global == nil {
		panic("Load Config first")
	}
	return global
}

func NewDefaultConfig() *Config {
	return &Config{
		Server:   NewServerConfig(),
		Database: NewDatabaseConfig(),
		Log:      NewLogConfig(),
	}
}

// Config 应用配置
type Config struct {
	Server   *ServerConfig   `mapstructure:"server"`
	Database *DatabaseConfig `mapstructure:"database"`
	//Deploy     *DeployConfig     `mapstructure:"deploy"`
	//Delivery   *DeliveryConfig   `mapstructure:"delivery"`
	Log *LogConfig `mapstructure:"log"`
	//Auth       *AuthConfig       `mapstructure:"auth"`
	//Kubernetes *KubernetesConfig `mapstructure:"kubernetes"`
	//Cache      *CacheConfig      `mapstructure:"cache"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: 8080,
		Host: "0.0.0.0",
	}
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type: "sqlite",
		SQLite: SQLiteConfig{
			Path: "data/kde.db",
		},
		Etcd: EtcdConfig{
			Endpoints:   []string{"127.0.0.1:2379"},
			DialTimeout: 5,
		},
	}
}

func (d *DatabaseConfig) GetDB() (*gorm.DB, error) {
	// 直接加锁，锁住临界区
	d.lock.Lock()
	// 函数退出时，释放锁
	defer d.lock.Unlock()
	if gdb == nil {
		conn, err := d.GetDBCon()
		if err != nil {
			return nil, err
		}
		gdb = conn
	}
	// 全局变量gdb就一定存在
	return gdb, nil
}

func (d *DatabaseConfig) GetDBCon() (*gorm.DB, error) {
	// 确保数据库目录存在
	dbDir := filepath.Dir(d.SQLite.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %v", err)
	}
	// 连接数据库
	db, err := gorm.Open(sqlite.Open(d.SQLite.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
	}

	// 自动迁移数据库模型
	if err := db.AutoMigrate(&models.ClusterModel{}, &models.NodeModel{}); err != nil {
		return nil, fmt.Errorf("迁移数据库模型失败: %v", err)
	}
	return db, nil
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type   string       `mapstructure:"type"`
	SQLite SQLiteConfig `mapstructure:"sqlite"`
	Etcd   EtcdConfig   `mapstructure:"etcd"`
	lock   sync.Mutex
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

func NewLogConfig() *LogConfig {
	return &LogConfig{
		Level:  "info",
		Format: "text",
		Output: "stdout",
	}
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
