package cluster

const (
	AppName = "cluster"
)

// 集群管理接口
type ClusterInterface interface {
	CreateCluster(*Cluster) (*Cluster, error)
	UpdateCluster(*Cluster) (*Cluster, error)
	DeleteCluster(*Cluster) (*Cluster, error)
}
