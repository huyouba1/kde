package models

import (
	"time"

	"gorm.io/gorm"
)

// ClusterInfo 集群信息
type ClusterInfo struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	ClusterID      string         `gorm:"uniqueIndex;not null" json:"cluster_id"`
	Version        string         `json:"version"`
	NodeCount      int            `json:"node_count"`
	NamespaceCount int            `json:"namespace_count"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
