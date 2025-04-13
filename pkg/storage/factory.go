package storage

import (
	"fmt"
	"github.com/huyouba1/kde/configs"

	"github.com/huyouba1/kde/pkg/storage/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Factory 存储工厂
type Factory struct {
	db *gorm.DB
}

// NewFactory 创建存储工厂
func NewFactory(cfg *configs.Config) *Factory {
	var db *gorm.DB
	var err error

	switch cfg.Database.Type {
	case "sqlite":
		db, err = initSQLite(cfg.Database.SQLite)
	default:
		panic(fmt.Sprintf("不支持的数据库类型: %s", cfg.Database.Type))
	}

	if err != nil {
		panic(fmt.Sprintf("初始化数据库失败: %v", err))
	}

	// 自动迁移数据库模型
	if err := autoMigrate(db); err != nil {
		panic(fmt.Sprintf("数据库迁移失败: %v", err))
	}

	return &Factory{
		db: db,
	}
}

// GetDB 获取数据库连接
func (f *Factory) GetDB() *gorm.DB {
	return f.db
}

// Close 关闭数据库连接
func (f *Factory) Close() error {
	sqlDB, err := f.db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}
	return sqlDB.Close()
}

// initSQLite 初始化SQLite数据库
func initSQLite(cfg configs.SQLiteConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Path), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
	}
	return db, nil
}

// autoMigrate 自动迁移数据库模型
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.ClusterModel{},
		&models.NodeModel{},
	)
}
