package storage

import (
	"fmt"

	"github.com/huyouba1/kde/pkg/storage/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// DB 封装了数据库操作
type DB struct {
	db *gorm.DB
}

// NewDB 创建一个新的数据库连接
func NewDB(dbPath string) (*DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 自动迁移数据库表
	if err := db.AutoMigrate(
		&models.ClusterModel{},
		&models.ClusterInfo{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return &DB{db: db}, nil
}

// GetCluster 获取集群信息
func (db *DB) GetCluster(id string) (*models.ClusterModel, error) {
	var cluster models.ClusterModel
	if err := db.db.First(&cluster, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get cluster: %w", err)
	}
	return &cluster, nil
}

// ListClusters 获取集群列表
func (db *DB) ListClusters() ([]*models.ClusterModel, error) {
	var clusters []*models.ClusterModel
	if err := db.db.Find(&clusters).Error; err != nil {
		return nil, fmt.Errorf("failed to list clusters: %w", err)
	}
	return clusters, nil
}

// CreateCluster 创建集群
func (db *DB) CreateCluster(cluster *models.ClusterModel) error {
	if err := db.db.Create(cluster).Error; err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}
	return nil
}

// UpdateCluster 更新集群
func (db *DB) UpdateCluster(cluster *models.ClusterModel) error {
	if err := db.db.Save(cluster).Error; err != nil {
		return fmt.Errorf("failed to update cluster: %w", err)
	}
	return nil
}

// DeleteCluster 删除集群
func (db *DB) DeleteCluster(id string) error {
	if err := db.db.Delete(&models.ClusterModel{}, "id = ?", id).Error; err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}
	return nil
}

// UpdateClusterInfo 更新集群信息
func (db *DB) UpdateClusterInfo(clusterID string, info *models.ClusterInfo) error {
	info.ClusterID = clusterID
	if err := db.db.Save(info).Error; err != nil {
		return fmt.Errorf("failed to update cluster info: %w", err)
	}
	return nil
}

// GetClusterInfo 获取集群信息
func (db *DB) GetClusterInfo(clusterID string) (*models.ClusterInfo, error) {
	var info models.ClusterInfo
	if err := db.db.First(&info, "cluster_id = ?", clusterID).Error; err != nil {
		return nil, fmt.Errorf("failed to get cluster info: %w", err)
	}
	return &info, nil
}
