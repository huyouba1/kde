package cluster

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/huyouba1/kde/pkg/k8s"
	"github.com/huyouba1/kde/pkg/storage"
	"github.com/huyouba1/kde/pkg/storage/models"
)

// ClusterStatus 集群状态
type ClusterStatus string

const (
	// StatusActive 活跃状态
	StatusActive ClusterStatus = "active"
	// StatusInactive 非活跃状态
	StatusInactive ClusterStatus = "inactive"
	// StatusCreating 创建中
	StatusCreating ClusterStatus = "creating"
	// StatusDeleting 删除中
	StatusDeleting ClusterStatus = "deleting"
	// StatusError 错误状态
	StatusError ClusterStatus = "error"
)

// Cluster 表示一个Kubernetes集群
type Cluster struct {
	ID          string
	Name        string
	Description string
	Status      models.ClusterStatus
	KubeConfig  string
	APIServer   string
	Version     string
	NodeCount   int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Node 表示一个Kubernetes节点
type Node struct {
	Name       string
	ClusterID  string
	Status     string
	Role       string
	IP         string
	Version    string
	Labels     map[string]string
	Conditions []models.NodeCondition
}

// ClusterManager 管理集群
type ClusterManager struct {
	db        *storage.DB
	clients   map[string]*k8s.Client
	clientsMu sync.RWMutex
}

// NewClusterManager 创建一个新的集群管理器
func NewClusterManager(db *storage.DB) *ClusterManager {
	return &ClusterManager{
		db:      db,
		clients: make(map[string]*k8s.Client),
	}
}

// GetClient 获取指定集群的 Kubernetes 客户端
func (m *ClusterManager) GetClient(clusterID string) (*k8s.Client, error) {
	m.clientsMu.RLock()
	client, exists := m.clients[clusterID]
	m.clientsMu.RUnlock()

	if exists {
		return client, nil
	}

	// 从数据库获取集群配置
	cluster, err := m.db.GetCluster(clusterID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster: %w", err)
	}

	// 创建新的 Kubernetes 客户端
	client, err = k8s.NewClient(cluster.KubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	// 测试连接
	if err := client.TestConnection(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to connect to cluster: %w", err)
	}

	// 缓存客户端
	m.clientsMu.Lock()
	m.clients[clusterID] = client
	m.clientsMu.Unlock()

	return client, nil
}

// GetClusterInfo 获取集群信息
func (m *ClusterManager) GetClusterInfo(clusterID string) (*models.ClusterInfo, error) {
	client, err := m.GetClient(clusterID)
	if err != nil {
		return nil, err
	}

	info, err := client.GetClusterInfo(context.Background())
	if err != nil {
		return nil, err
	}

	return &models.ClusterInfo{
		Version:        info.Version,
		NodeCount:      info.NodeCount,
		NamespaceCount: info.NamespaceCount,
	}, nil
}

// RemoveClient 移除指定集群的客户端
func (m *ClusterManager) RemoveClient(clusterID string) {
	m.clientsMu.Lock()
	delete(m.clients, clusterID)
	m.clientsMu.Unlock()
}

// GetCluster 获取集群信息
func (m *ClusterManager) GetCluster(ctx context.Context, id string) (*models.ClusterModel, error) {
	cluster, err := m.db.GetCluster(id)
	if err != nil {
		return nil, fmt.Errorf("查询集群失败: %v", err)
	}
	return cluster, nil
}

// ListClusters 获取集群列表
func (m *ClusterManager) ListClusters(ctx context.Context) ([]*models.ClusterModel, error) {
	clusters, err := m.db.ListClusters()
	if err != nil {
		return nil, fmt.Errorf("查询集群列表失败: %v", err)
	}
	return clusters, nil
}

// CreateCluster 创建集群
func (m *ClusterManager) CreateCluster(ctx context.Context, cluster *models.ClusterModel) error {
	// 设置创建和更新时间
	now := time.Now()
	cluster.CreatedAt = now
	cluster.UpdatedAt = now

	if err := m.db.CreateCluster(cluster); err != nil {
		return fmt.Errorf("创建集群失败: %v", err)
	}

	return nil
}

// UpdateCluster 更新集群
func (m *ClusterManager) UpdateCluster(ctx context.Context, cluster *models.ClusterModel) error {
	// 设置更新时间
	cluster.UpdatedAt = time.Now()

	if err := m.db.UpdateCluster(cluster); err != nil {
		return fmt.Errorf("更新集群失败: %v", err)
	}

	return nil
}

// DeleteCluster 删除集群
func (m *ClusterManager) DeleteCluster(ctx context.Context, id string) error {
	// 删除集群记录
	if err := m.db.DeleteCluster(id); err != nil {
		return fmt.Errorf("删除集群失败: %v", err)
	}

	// 移除客户端缓存
	m.RemoveClient(id)

	return nil
}
