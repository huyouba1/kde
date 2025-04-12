package cluster

import (
	"context"
	"fmt"
	"time"

	"github.com/huyouba1/kde/pkg/storage"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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
	ID          string        `json:"id" gorm:"primaryKey"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Status      ClusterStatus `json:"status"`
	KubeConfig  string        `json:"kubeconfig" gorm:"type:text"`
	APIServer   string        `json:"api_server"`
	Version     string        `json:"version"`
	NodeCount   int           `json:"node_count"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// Node 表示一个Kubernetes节点
type Node struct {
	Name       string            `json:"name"`
	ClusterID  string            `json:"cluster_id"`
	Status     string            `json:"status"`
	Role       string            `json:"role"`
	IP         string            `json:"ip"`
	Version    string            `json:"version"`
	Labels     map[string]string `json:"labels"`
	Conditions []NodeCondition   `json:"conditions"`
}

// NodeCondition 节点状态条件
type NodeCondition struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Manager 集群管理器
type Manager struct {
	storageFactory storage.Factory
	clients        map[string]*kubernetes.Clientset
}

// NewManager 创建一个新的集群管理器
func NewManager(factory storage.Factory) *Manager {
	return &Manager{
		storageFactory: factory,
		clients:        make(map[string]*kubernetes.Clientset),
	}
}

// GetCluster 获取集群信息
func (m *Manager) GetCluster(ctx context.Context, id string) (*Cluster, error) {
	// TODO: 从数据库获取集群信息
	return nil, fmt.Errorf("未实现")
}

// ListClusters 获取集群列表
func (m *Manager) ListClusters(ctx context.Context) ([]*Cluster, error) {
	// TODO: 从数据库获取集群列表
	return nil, fmt.Errorf("未实现")
}

// CreateCluster 创建集群
func (m *Manager) CreateCluster(ctx context.Context, cluster *Cluster) error {
	// TODO: 保存集群信息到数据库
	return fmt.Errorf("未实现")
}

// UpdateCluster 更新集群
func (m *Manager) UpdateCluster(ctx context.Context, cluster *Cluster) error {
	// TODO: 更新集群信息到数据库
	return fmt.Errorf("未实现")
}

// DeleteCluster 删除集群
func (m *Manager) DeleteCluster(ctx context.Context, id string) error {
	// TODO: 从数据库删除集群信息
	return fmt.Errorf("未实现")
}

// GetClient 获取Kubernetes客户端
func (m *Manager) GetClient(ctx context.Context, id string) (*kubernetes.Clientset, error) {
	// 检查缓存
	if client, ok := m.clients[id]; ok {
		return client, nil
	}

	// 获取集群信息
	cluster, err := m.GetCluster(ctx, id)
	if err != nil {
		return nil, err
	}

	// 创建客户端
	client, err := m.createClient(cluster.KubeConfig)
	if err != nil {
		return nil, err
	}

	// 缓存客户端
	m.clients[id] = client

	return client, nil
}

// createClient 创建Kubernetes客户端
func (m *Manager) createClient(kubeconfig string) (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error

	if kubeconfig == "" {
		// 使用集群内部配置
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("获取集群内部配置失败: %v", err)
		}
	} else {
		// 使用提供的kubeconfig
		config, err = clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
		if err != nil {
			return nil, fmt.Errorf("解析kubeconfig失败: %v", err)
		}
	}

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	return clientset, nil
}

// GetNodes 获取集群节点列表
func (m *Manager) GetNodes(ctx context.Context, clusterID string) ([]Node, error) {
	// 获取客户端
	client, err := m.GetClient(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	// 获取节点列表
	nodeList, err := client.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("获取节点列表失败: %v", err)
	}

	// 转换为内部节点结构
	nodes := make([]Node, 0, len(nodeList.Items))
	for _, n := range nodeList.Items {
		// 确定节点角色
		role := "worker"
		for label := range n.Labels {
			if label == "node-role.kubernetes.io/master" || label == "node-role.kubernetes.io/control-plane" {
				role = "master"
				break
			}
		}

		// 获取节点IP
		var nodeIP string
		for _, addr := range n.Status.Addresses {
			if addr.Type == corev1.NodeInternalIP {
				nodeIP = addr.Address
				break
			}
		}

		// 转换节点状态条件
		conditions := make([]NodeCondition, 0, len(n.Status.Conditions))
		for _, c := range n.Status.Conditions {
			conditions = append(conditions, NodeCondition{
				Type:    string(c.Type),
				Status:  string(c.Status),
				Message: c.Message,
			})
		}

		// 创建节点对象
		node := Node{
			Name:       n.Name,
			ClusterID:  clusterID,
			Status:     getNodeStatus(n),
			Role:       role,
			IP:         nodeIP,
			Version:    n.Status.NodeInfo.KubeletVersion,
			Labels:     n.Labels,
			Conditions: conditions,
		}

		nodes = append(nodes, node)
	}

	return nodes, nil
}

// getNodeStatus 获取节点状态
func getNodeStatus(node corev1.Node) string {
	for _, condition := range node.Status.Conditions {
		if condition.Type == corev1.NodeReady {
			if condition.Status == corev1.ConditionTrue {
				return "Ready"
			} else {
				return "NotReady"
			}
		}
	}
	return "Unknown"
}
