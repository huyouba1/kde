package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client 封装了 Kubernetes 客户端
type Client struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
}

// NewClient 创建一个新的 Kubernetes 客户端
func NewClient(kubeconfigPath string) (*Client, error) {
	var config *rest.Config
	var err error

	// 如果提供了 kubeconfig 路径，使用它
	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
		}
	} else {
		// 否则尝试使用 in-cluster 配置
		config, err = rest.InClusterConfig()
		if err != nil {
			// 如果 in-cluster 配置失败，尝试使用默认的 kubeconfig 位置
			home, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get user home directory: %w", err)
			}
			kubeconfig := filepath.Join(home, ".kube", "config")
			config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
			if err != nil {
				return nil, fmt.Errorf("failed to build config from default kubeconfig: %w", err)
			}
		}
	}

	// 创建 clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	return &Client{
		clientset: clientset,
		config:    config,
	}, nil
}

// GetClientSet 返回 Kubernetes clientset
func (c *Client) GetClientSet() *kubernetes.Clientset {
	return c.clientset
}

// GetConfig 返回 Kubernetes 配置
func (c *Client) GetConfig() *rest.Config {
	return c.config
}

// TestConnection 测试与 Kubernetes 集群的连接
func (c *Client) TestConnection(ctx context.Context) error {
	_, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to connect to Kubernetes cluster: %w", err)
	}
	return nil
}

// GetClusterInfo 获取集群信息
func (c *Client) GetClusterInfo(ctx context.Context) (*ClusterInfo, error) {
	version, err := c.clientset.ServerVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to get server version: %w", err)
	}

	nodes, err := c.clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list nodes: %w", err)
	}

	namespaces, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	return &ClusterInfo{
		Version:        version.String(),
		NodeCount:      len(nodes.Items),
		NamespaceCount: len(namespaces.Items),
	}, nil
}

// ClusterInfo 包含集群的基本信息
type ClusterInfo struct {
	Version        string
	NodeCount      int
	NamespaceCount int
}
