package helm

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/huyouba1/kde/pkg/delivery"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Manager Helm交付管理器
type Manager struct {
	workdir   string
	cachePath string
}

// NewManager 创建一个新的Helm交付管理器
func NewManager(workdir, cachePath string) *Manager {
	return &Manager{
		workdir:   workdir,
		cachePath: cachePath,
	}
}

// Deploy 部署Helm Chart
func (m *Manager) Deploy(ctx context.Context, options *delivery.HelmOptions) error {
	// 创建工作目录
	deployDir := filepath.Join(m.workdir, options.ClusterID, "helm", options.Name)
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return fmt.Errorf("创建工作目录失败: %v", err)
	}

	// 获取Kubernetes客户端配置
	kubeconfig, err := m.getKubeconfigPath(options.ClusterID)
	if err != nil {
		return fmt.Errorf("获取kubeconfig失败: %v", err)
	}

	// 准备Helm命令参数
	args := []string{}

	// 检查是否已安装
	installed, err := m.isReleaseInstalled(kubeconfig, options.Name, options.Namespace)
	if err != nil {
		return fmt.Errorf("检查Release状态失败: %v", err)
	}

	if installed {
		// 升级已存在的Release
		args = append(args, "upgrade")
	} else {
		// 安装新的Release
		args = append(args, "install")
	}

	// 添加Release名称
	args = append(args, options.Name)

	// 添加Chart来源
	if options.ChartPath != "" {
		// 使用本地Chart路径
		args = append(args, options.ChartPath)
	} else {
		// 使用Chart仓库
		chartRef := options.ChartName
		if options.ChartRepo != "" {
			// 添加仓库
			repoName := strings.Split(options.ChartRepo, "/")[0]
			if err := m.addHelmRepo(repoName, options.ChartRepo); err != nil {
				return fmt.Errorf("添加Helm仓库失败: %v", err)
			}
			chartRef = fmt.Sprintf("%s/%s", repoName, options.ChartName)
		}

		// 添加版本信息
		if options.Version != "" {
			chartRef = fmt.Sprintf("%s --version %s", chartRef, options.Version)
		}

		args = append(args, chartRef)
	}

	// 添加命名空间
	if options.Namespace != "" {
		args = append(args, "--namespace", options.Namespace, "--create-namespace")
	}

	// 添加值覆盖
	if len(options.Values) > 0 {
		// 创建values文件
		valuesFile := filepath.Join(deployDir, "values.yaml")
		if err := m.createValuesFile(valuesFile, options.Values); err != nil {
			return fmt.Errorf("创建values文件失败: %v", err)
		}
		args = append(args, "-f", valuesFile)
	}

	// 设置kubeconfig
	args = append(args, "--kubeconfig", kubeconfig)

	// 执行Helm命令
	cmd := exec.CommandContext(ctx, "helm", args...)
	cmd.Dir = deployDir

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行Helm命令失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// Uninstall 卸载Helm Release
func (m *Manager) Uninstall(ctx context.Context, clusterID, name, namespace string) error {
	// 获取Kubernetes客户端配置
	kubeconfig, err := m.getKubeconfigPath(clusterID)
	if err != nil {
		return fmt.Errorf("获取kubeconfig失败: %v", err)
	}

	// 准备Helm命令参数
	args := []string{"uninstall", name}

	// 添加命名空间
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	// 设置kubeconfig
	args = append(args, "--kubeconfig", kubeconfig)

	// 执行Helm命令
	cmd := exec.CommandContext(ctx, "helm", args...)

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("执行Helm命令失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// GetReleaseStatus 获取Release状态
func (m *Manager) GetReleaseStatus(ctx context.Context, clusterID, name, namespace string) (string, error) {
	// 获取Kubernetes客户端配置
	kubeconfig, err := m.getKubeconfigPath(clusterID)
	if err != nil {
		return "", fmt.Errorf("获取kubeconfig失败: %v", err)
	}

	// 准备Helm命令参数
	args := []string{"status", name, "--output", "json"}

	// 添加命名空间
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	// 设置kubeconfig
	args = append(args, "--kubeconfig", kubeconfig)

	// 执行Helm命令
	cmd := exec.CommandContext(ctx, "helm", args...)

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("执行Helm命令失败: %v, 输出: %s", err, string(output))
	}

	return string(output), nil
}

// isReleaseInstalled 检查Release是否已安装
func (m *Manager) isReleaseInstalled(kubeconfig, name, namespace string) (bool, error) {
	// 准备Helm命令参数
	args := []string{"list", "--filter", name, "--output", "json"}

	// 添加命名空间
	if namespace != "" {
		args = append(args, "--namespace", namespace)
	}

	// 设置kubeconfig
	args = append(args, "--kubeconfig", kubeconfig)

	// 执行Helm命令
	cmd := exec.Command("helm", args...)

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("执行Helm命令失败: %v, 输出: %s", err, string(output))
	}

	// 检查输出是否包含Release名称
	return strings.Contains(string(output), name), nil
}

// addHelmRepo 添加Helm仓库
func (m *Manager) addHelmRepo(name, url string) error {
	// 执行Helm命令添加仓库
	cmd := exec.Command("helm", "repo", "add", name, url)

	// 捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("添加Helm仓库失败: %v, 输出: %s", err, string(output))
	}

	// 更新仓库
	cmd = exec.Command("helm", "repo", "update")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("更新Helm仓库失败: %v, 输出: %s", err, string(output))
	}

	return nil
}

// createValuesFile 创建values文件
func (m *Manager) createValuesFile(filePath string, values map[string]string) error {
	// 创建YAML格式的values文件
	content := ""
	for k, v := range values {
		content += fmt.Sprintf("%s: %s\n", k, v)
	}

	// 写入文件
	return os.WriteFile(filePath, []byte(content), 0644)
}

// getKubeconfigPath 获取kubeconfig路径
func (m *Manager) getKubeconfigPath(clusterID string) (string, error) {
	// TODO: 从存储中获取集群的kubeconfig并保存到临时文件
	// 这里应该实现从数据库获取kubeconfig并写入临时文件的逻辑

	// 临时返回默认kubeconfig路径
	return filepath.Join(os.Getenv("HOME"), ".kube", "config"), nil
}

// getKubernetesClient 获取Kubernetes客户端
func (m *Manager) getKubernetesClient(clusterID string) (*kubernetes.Clientset, error) {
	// 获取kubeconfig路径
	kubeconfig, err := m.getKubeconfigPath(clusterID)
	if err != nil {
		return nil, err
	}

	// 加载kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("加载kubeconfig失败: %v", err)
	}

	// 创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("创建Kubernetes客户端失败: %v", err)
	}

	return clientset, nil
}
