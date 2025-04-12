package storage

import (
	"fmt"

	"github.com/huyouba1/kde/pkg/storage/config"
	"github.com/huyouba1/kde/pkg/storage/etcd"
	"github.com/huyouba1/kde/pkg/storage/sqlite"
)

// Factory 存储工厂接口
type Factory interface {
	// GetSQLiteManager 获取SQLite管理器
	GetSQLiteManager() (*sqlite.Manager, error)
	// GetEtcdManager 获取etcd管理器
	GetEtcdManager() (*etcd.Manager, error)
	// Close 关闭所有连接
	Close() error
}

// factory 存储工厂实现
type factory struct {
	config   *config.Config
	sqliteDB *sqlite.Manager
	etcdDB   *etcd.Manager
}

// NewFactory 创建一个新的存储工厂
func NewFactory(cfg *config.Config) Factory {
	return &factory{
		config: cfg,
	}
}

// GetSQLiteManager 获取SQLite管理器
func (f *factory) GetSQLiteManager() (*sqlite.Manager, error) {
	if f.sqliteDB != nil {
		return f.sqliteDB, nil
	}

	db, err := sqlite.NewManager(&f.config.Database.SQLite)
	if err != nil {
		return nil, err
	}

	f.sqliteDB = db
	return db, nil
}

// GetEtcdManager 获取etcd管理器
func (f *factory) GetEtcdManager() (*etcd.Manager, error) {
	if f.etcdDB != nil {
		return f.etcdDB, nil
	}

	db, err := etcd.NewManager(&f.config.Database.Etcd)
	if err != nil {
		return nil, err
	}

	f.etcdDB = db
	return db, nil
}

// GetStorageManager 根据配置获取存储管理器
func (f *factory) GetStorageManager() (interface{}, error) {
	switch f.config.Database.Type {
	case "sqlite":
		return f.GetSQLiteManager()
	case "etcd":
		return f.GetEtcdManager()
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", f.config.Database.Type)
	}
}

// Close 关闭所有连接
func (f *factory) Close() error {
	var lastErr error

	if f.sqliteDB != nil {
		if err := f.sqliteDB.Close(); err != nil {
			lastErr = err
		}
	}

	if f.etcdDB != nil {
		if err := f.etcdDB.Close(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
