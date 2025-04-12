package kustomize

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/huyouba1/kde/pkg/delivery"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

// Manager Kustomize交付管理器
type Manager struct {
	workdir string
}

// NewManager 创建一个新的Kustomize交付管理器
func NewManager(workdir string) *Manager {
	return &Manager{
		workdir: workdir,
	}
}

// Deploy 部署Kustomize配置
func (m *Manager) Deploy(ctx context.Context, options *delivery.KustomizeOptions) error {
	// 创建工作目录
	deployDir := filepath.Join(m.workdir, options.ClusterID, "kustomize", options.Name)
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return fmt.Errorf("创建工作目录失败: %v", err)
	}

	// 获取Kubernetes客户端
	clientset, dynamicClient, err := m.getKubernetesClient(options.ClusterID)
	if err != nil {
		return fmt.Errorf("获取Kubernetes客户端失败: %v", err)
	}

	// 构建Kustomize资源
	manifests, err := m.buildKustomize(options.BasePath, options.OverlayPath)
	if err != nil {
		return fmt.Errorf("构建Kustomize资源失败: %v", err)
	}

	// 应用资源到集群
	if err := m.applyManifests(ctx, dynamicClient, clientset, options.Namespace, manifests); err != nil {
		return fmt.Errorf("应用资源到集群失败: %v", err)
	}

	return nil
}

// buildKustomize 构建Kustomize资源
func (m *Manager) buildKustomize(basePath, overlayPath string) (string, error) {
	// 确定kustomization路径
	kustomizationPath := basePath
	if overlayPath != "" {
		kustomizationPath = overlayPath
	}

	// 执行kustomize build命令
	cmd := exec.Command("kustomize", "build", kustomizationPath)

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行kustomize build命令失败: %v, 输出: %s", err, string(output))
	}

	return string(output), nil
}

// applyManifests 应用资源到集群
func (m *Manager) applyManifests(ctx context.Context, dynamicClient dynamic.Interface, clientset *kubernetes.Clientset, namespace, manifests string) error {
	// 创建RESTMapper
	dc, err := discovery.NewDiscoveryClientForConfig(clientset.RESTClient().GetConfig())
	if err != nil {
		return fmt.Errorf("创建Discovery客户端失败: %v", err)
	}

	mapper := restmapper.NewDeferredDiscoveryRESTMapper(
		memory.NewMemCacheClient(dc),
	)

	// 分割多文档YAML
	documents := splitYAMLDocuments(manifests)

	decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

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
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return fmt.Errorf("获取资源映射失败: %v", err)
		}

		// 确定是否为命名空间级别的资源
		var resourceClient dynamic.ResourceInterface
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			resourceClient = dynamicClient.Resource(mapping.Resource).Namespace(obj.GetNamespace())
		} else {
			resourceClient = dynamicClient.Resource(mapping.Resource)
		}

		// 应用资源
		_, err = resourceClient.Apply(ctx, obj.GetName(), obj, metav1.ApplyOptions{
			FieldManager: "kustomize-manager",
			Force:        true,
		})

		if err != nil {
			return fmt.Errorf("应用资源 %s/%s 失败: %v", gvk.Kind, obj.GetName(), err)
		}
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

// getKubeconfigPath 获取kubeconfig路径
func (m *Manager) getKubeconfigPath(clusterID string) (string, error) {
	// TODO: 从存储中获取集群的kubeconfig并保存到临时文件
	// 这里应该实现从数据库获取kubeconfig并写入临时文件的逻辑

	// 临时返回默认kubeconfig路径
	return filepath.Join(os.Getenv("HOME"), ".kube", "config"), nil
}

// splitYAMLDocuments 分割多文档YAML
func splitYAMLDocuments(content string) []string {
	// 使用---作为分隔符分割YAML文档
	docs := strings.Split(content, "---\n")
	result := make([]string, 0, len(docs))

	for _, doc := range docs {
		// 去除空白
		doc = strings.TrimSpace(doc)
		if doc != "" {
			result = append(result, doc)
		}
	}

	return result
}
