package models

import (
	"time"

	"gorm.io/gorm"
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

// ClusterModel 集群数据库模型
type ClusterModel struct {
	ID          string         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null" json:"name"`
	Description string         `json:"description"`
	Status      ClusterStatus  `gorm:"not null" json:"status"`
	KubeConfig  string         `gorm:"type:text;not null" json:"kubeconfig"`
	APIServer   string         `gorm:"not null" json:"api_server"`
	Version     string         `json:"version"`
	NodeCount   int            `gorm:"default:0" json:"node_count"`
	CreatedAt   time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// NodeCondition 节点状态
type NodeCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

// NodeModel 节点数据库模型
type NodeModel struct {
	Name       string         `gorm:"primaryKey" json:"name"`
	ClusterID  string         `gorm:"primaryKey" json:"cluster_id"`
	Status     string         `gorm:"not null" json:"status"`
	Role       string         `gorm:"not null" json:"role"`
	IP         string         `gorm:"not null" json:"ip"`
	Version    string         `json:"version"`
	Labels     string         `gorm:"type:text" json:"labels"`
	Conditions string         `gorm:"type:text" json:"conditions"`
	CreatedAt  time.Time      `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"not null" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
