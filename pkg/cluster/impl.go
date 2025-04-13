package cluster

import (
	"github.com/huyouba1/kde/configs"
	"github.com/huyouba1/kde/pkg/k8s"
	"gorm.io/gorm"
)

type ClusterServiceServer struct {
	gdb    *gorm.DB
	client *k8s.Client
}

func NewClusterServiceServer() *ClusterServiceServer {
	gdb, err := configs.C().Database.GetDB()
	if err != nil {
		panic(err)
	}

	return &ClusterServiceServer{gdb: gdb}
}
