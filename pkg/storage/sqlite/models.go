package sqlite

import (
	"time"

	"github.com/huyouba1/kde/pkg/cluster"
	"gorm.io/gorm"
)

// ClusterModel 集群数据库模型
type ClusterModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	Status      string `gorm:"not null"`
	KubeConfig  string `gorm:"type:text;not null"`
	APIServer   string `gorm:"not null"`
	Version     string
	NodeCount   int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

// NodeModel 节点数据库模型
type NodeModel struct {
	Name       string `gorm:"primaryKey"`
	ClusterID  string `gorm:"primaryKey"`
	Status     string `gorm:"not null"`
	Role       string `gorm:"not null"`
	IP         string `gorm:"not null"`
	Version    string
	Labels     string    `gorm:"type:text"`
	Conditions string    `gorm:"type:text"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

// ToCluster 转换为集群对象
func (m *ClusterModel) ToCluster() *cluster.Cluster {
	return &cluster.Cluster{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Status:      cluster.ClusterStatus(m.Status),
		KubeConfig:  m.KubeConfig,
		APIServer:   m.APIServer,
		Version:     m.Version,
		NodeCount:   m.NodeCount,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

// FromCluster 从集群对象创建模型
func FromCluster(c *cluster.Cluster) *ClusterModel {
	return &ClusterModel{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Status:      string(c.Status),
		KubeConfig:  c.KubeConfig,
		APIServer:   c.APIServer,
		Version:     c.Version,
		NodeCount:   c.NodeCount,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// ToNode 转换为节点对象
func (m *NodeModel) ToNode() *cluster.Node {
	return &cluster.Node{
		Name:       m.Name,
		ClusterID:  m.ClusterID,
		Status:     m.Status,
		Role:       m.Role,
		IP:         m.IP,
		Version:    m.Version,
		Labels:     make(map[string]string),          // TODO: 实现 JSON 解析
		Conditions: make([]cluster.NodeCondition, 0), // TODO: 实现 JSON 解析
	}
}

// FromNode 从节点对象创建模型
func FromNode(n *cluster.Node) *NodeModel {
	return &NodeModel{
		Name:       n.Name,
		ClusterID:  n.ClusterID,
		Status:     n.Status,
		Role:       n.Role,
		IP:         n.IP,
		Version:    n.Version,
		Labels:     "", // TODO: 实现 JSON 序列化
		Conditions: "", // TODO: 实现 JSON 序列化
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

// AutoMigrate 自动迁移数据库模型
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&ClusterModel{}, &NodeModel{})
}
