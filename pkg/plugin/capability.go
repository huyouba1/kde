package plugin

// PluginType 插件类型
type PluginType string

const (
	// TypeClusterManager 集群管理插件
	TypeClusterManager PluginType = "cluster_manager"
	// TypeDeployment 部署插件
	TypeDeployment PluginType = "deployment"
	// TypeDelivery 交付插件
	TypeDelivery PluginType = "delivery"
	// TypeMonitoring 监控插件
	TypeMonitoring PluginType = "monitoring"
	// TypeSecurity 安全插件
	TypeSecurity PluginType = "security"
	// TypeBackup 备份插件
	TypeBackup PluginType = "backup"
	// TypeGeneral 通用插件
	TypeGeneral PluginType = "general"
)

// PluginCapability 插件能力
type PluginCapability string

const (
	// CapabilityClusterCreate 创建集群能力
	CapabilityClusterCreate PluginCapability = "cluster_create"
	// CapabilityClusterMonitor 监控集群能力
	CapabilityClusterMonitor PluginCapability = "cluster_monitor"
	// CapabilityClusterBackup 备份集群能力
	CapabilityClusterBackup PluginCapability = "cluster_backup"
	// CapabilityDeployK8s 部署Kubernetes能力
	CapabilityDeployK8s PluginCapability = "deploy_k8s"
	// CapabilityDeployOffline 离线部署能力
	CapabilityDeployOffline PluginCapability = "deploy_offline"
	// CapabilityDeliveryHelm Helm交付能力
	CapabilityDeliveryHelm PluginCapability = "delivery_helm"
	// CapabilityDeliveryKustomize Kustomize交付能力
	CapabilityDeliveryKustomize PluginCapability = "delivery_kustomize"
	// CapabilityDeliveryYaml YAML交付能力
	CapabilityDeliveryYaml PluginCapability = "delivery_yaml"
)

// CapabilityPlugin 具有能力的插件接口
type CapabilityPlugin interface {
	Plugin
	// GetType 获取插件类型
	GetType() PluginType
	// GetCapabilities 获取插件支持的能力
	GetCapabilities() []PluginCapability
	// HasCapability 检查插件是否具有指定能力
	HasCapability(capability PluginCapability) bool
}

// BaseCapabilityPlugin 基础能力插件实现
type BaseCapabilityPlugin struct {
	pluginType   PluginType
	capabilities []PluginCapability
}

// GetType 获取插件类型
func (p *BaseCapabilityPlugin) GetType() PluginType {
	return p.pluginType
}

// GetCapabilities 获取插件支持的能力
func (p *BaseCapabilityPlugin) GetCapabilities() []PluginCapability {
	return p.capabilities
}

// HasCapability 检查插件是否具有指定能力
func (p *BaseCapabilityPlugin) HasCapability(capability PluginCapability) bool {
	for _, cap := range p.capabilities {
		if cap == capability {
			return true
		}
	}
	return false
}

// SetType 设置插件类型
func (p *BaseCapabilityPlugin) SetType(pluginType PluginType) {
	p.pluginType = pluginType
}

// AddCapability 添加插件能力
func (p *BaseCapabilityPlugin) AddCapability(capability PluginCapability) {
	p.capabilities = append(p.capabilities, capability)
}

// SetCapabilities 设置插件能力列表
func (p *BaseCapabilityPlugin) SetCapabilities(capabilities []PluginCapability) {
	p.capabilities = capabilities
}
