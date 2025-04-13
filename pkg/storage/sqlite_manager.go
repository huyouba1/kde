package storage

import (
	"fmt"
	"github.com/huyouba1/kde/configs"
	"os"
	"path/filepath"

	"github.com/huyouba1/kde/pkg/storage/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Manager SQLite数据库管理器
type Manager struct {
	db     *gorm.DB
	config *configs.SQLiteConfig
}

// NewManager 创建一个新的SQLite管理器
func NewManager(cfg *configs.SQLiteConfig) (*Manager, error) {
	// 确保数据库目录存在
	dbDir := filepath.Dir(cfg.Path)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("创建数据库目录失败: %v", err)
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
	}

	// 自动迁移数据库模型
	if err := db.AutoMigrate(&models.ClusterModel{}, &models.NodeModel{}); err != nil {
		return nil, fmt.Errorf("迁移数据库模型失败: %v", err)
	}

	return &Manager{
		db:     db,
		config: cfg,
	}, nil
}

// DB 获取数据库连接
func (m *Manager) DB() *gorm.DB {
	return m.db
}

// Close 关闭数据库连接
func (m *Manager) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// AutoMigrate 自动迁移数据库模型
func (m *Manager) AutoMigrate(models ...interface{}) error {
	return m.db.AutoMigrate(models...)
}
