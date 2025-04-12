package yaml

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/huyouba1/kde/pkg/delivery"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Manager YAML交付管理器
type Manager struct {
	workdir string
}

// NewManager 创建一个新的YAML交付管理器
func NewManager(workdir string) *Manager {
	return &Manager{
		workdir: workdir,
	}
}

// Deploy 部署YAML资源
func (m *Manager) Deploy(ctx context.Context, options *delivery.YAMLOptions) error {
	// 创建工作目录
	deployDir := filepath.Join(m.workdir, options.ClusterID, "yaml", options.Name)
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return fmt.Errorf("创建工作目录失败: %v", err)
	}

	// 保存YAML内容到文件
	yamlFile := filepath.Join(deployDir, "resources.yaml")
	if err := ioutil.WriteFile(yamlFile, []byte(options.Content), 0644); err != nil {
		return fmt.Errorf("保存YAML文件失败: %v", err)
	}

	// 获取Kubernetes客户端
	clientset, dynamicClient, err := m.getKubernetesClient(options.ClusterID)
	if err != nil {
		return fmt.Errorf("获取Kubernetes客户端失败: %v", err)
	}

	// 应用YAML资源
	if err := m.applyYAML(ctx, dynamicClient, options.Namespace, options.Content); err != nil {
		return fmt.Errorf("应用YAML资源失败: %v", err)
	}

	return nil
}

// getKubernetesClient 获取Kubernetes客户端
func (m *Manager) getKubernetesClient(clusterID string) (*kubernetes.Clientset, dynamic.Interface, error) {
	// TODO: 从存储中获取集群的kubeconfig
	kubeconfig := ""

	// 创建客户端配置
	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfig))
	if err != nil {
		return nil, nil, fmt.Errorf("创建REST配置失败: %v", err)
	}

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("创建clientset失败: %v", err)
	}

	// 创建dynamic客户端
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("创建dynamic客户端失败: %v", err)
	}

	return clientset, dynamicClient, nil
}

// applyYAML 应用YAML资源
func (m *Manager) applyYAML(ctx context.Context, dynamicClient dynamic.Interface, namespace, content string) error {
	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

	// 分割多文档YAML
	documents := splitYAMLDocuments(content)

	for _, doc := range documents {
		if len(doc) == 0 {
			continue
		}

		// 解析YAML为非结构化对象
		obj := &unstructured.Unstructured{}
		_, gvk, err := decoder.Decode([]byte(doc), nil, obj)
		if err != nil {
			return fmt.Errorf("解析YAML文档失败: %v", err)
		}

		// 设置命名空间
		if len(namespace) > 0 && obj.GetNamespace() == "" {
			obj.SetNamespace(namespace)
		}

		// 获取资源的GVR
		gvr := schema.GroupVersionResource{
			Group:    gvk.Group,
			Version:  gvk.Version,
			Resource: pluralizeKind(gvk.Kind),
		}

		// 确定是否为命名空间级别的资源
		var resourceClient dynamic.ResourceInterface
		if obj.GetNamespace() != "" {
			resourceClient = dynamicClient.Resource(gvr).Namespace(obj.GetNamespace())
		} else {
			resourceClient = dynamicClient.Resource(gvr)
		}

		// 检查资源是否存在
		existing, err := resourceClient.Get(ctx, obj.GetName(), metav1.GetOptions{})
		if err == nil {
			// 资源存在，更新它
			obj.SetResourceVersion(existing.GetResourceVersion())
			_, err = resourceClient.Update(ctx, obj, metav1.UpdateOptions{})
		} else {
			// 资源不存在，创建它
			_, err = resourceClient.Create(ctx, obj, metav1.CreateOptions{})
		}

		if err != nil {
			return fmt.Errorf("应用资源 %s/%s 失败: %v", gvk.Kind, obj.GetName(), err)
		}
	}

	return nil
}

// splitYAMLDocuments 分割多文档YAML
func splitYAMLDocuments(content string) []string {
	// TODO: 实现YAML文档分割
	// 这里应该实现一个简单的YAML文档分割器
	return []string{content}
}

// pluralizeKind 将Kind转换为复数形式的资源名称
func pluralizeKind(kind string) string {
	// 简单的复数规则
	switch kind {
	case "Endpoints":
		return "endpoints"
	case "ConfigMap":
		return "configmaps"
	case "Pod":
		return "pods"
	case "Service":
		return "services"
	case "Deployment":
		return "deployments"
	case "StatefulSet":
		return "statefulsets"
	case "DaemonSet":
		return "daemonsets"
	case "Ingress":
		return "ingresses"
	case "Secret":
		return "secrets"
	case "PersistentVolume":
		return "persistentvolumes"
	case "PersistentVolumeClaim":
		return "persistentvolumeclaims"
	case "Namespace":
		return "namespaces"
	default:
		// 默认规则：添加s或es
		if kind[len(kind)-1] == 's' || kind[len(kind)-1] == 'x' ||
			(len(kind) > 1 && kind[len(kind)-2:] == "sh") ||
			(len(kind) > 1 && kind[len(kind)-2:] == "ch") {
			return kind + "es"
		}
		return kind + "s"
	}
}
