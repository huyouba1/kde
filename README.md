# Kubernetes Deployment Engine (KDE)

KDE 是一个专注于 Kubernetes 应用交付的部署引擎，旨在简化客户环境的持续交付流程。通过统一的界面和标准化的部署流程，帮助团队快速、可靠地将应用部署到各种 Kubernetes 环境中。

## 架构设计

### 技术栈
- **后端**: 使用 Go 语言开发，采用单体架构
- **前端**: 使用 Go 模板引擎，Bootstrap 框架
- **数据库**: SQLite（支持 etcd 作为备选）
- **部署工具**: Ansible，支持 RKE 部署方式

### 系统架构
```
KDE
├── 前端界面 (Go Template + Bootstrap)
├── API 服务层
├── 业务逻辑层
├── 数据访问层 (SQLite/etcd)
└── 部署引擎 (Ansible)
```
![系统架构](images/arch.png)

![功能模块](images/module.png)

## 核心功能

### 1. 环境部署
- 支持客户环境快速部署
  - 单机部署
  - 集群部署
- 多环境适配
  - 支持不同操作系统版本
  - 支持离线环境部署
  - 支持容器化部署
- 环境配置管理
  - 环境变量管理
  - 配置文件管理
  - 密钥管理

### 2. 应用交付
- 多格式应用支持
  - Kubernetes YAML 文件
  - Helm Charts
  - Kustomize 配置
- 交付流程管理
  - 版本控制
  - 环境隔离
  - 回滚机制
- 模板化部署
  - 变量替换
  - 条件部署
  - 依赖管理

### 3. 持续交付
- 部署流水线
  - 自动化部署
  - 环境检查
  - 健康检查
- 部署策略
  - 蓝绿部署
  - 金丝雀发布
  - 滚动更新
- 部署监控
  - 部署状态
  - 资源使用
  - 性能指标

### 4. 扩展性
- 插件系统
  - 自定义部署策略
  - 自定义资源管理
  - 自定义检查规则
- 集成能力
  - CI/CD 工具集成
  - 监控系统集成
  - 日志系统集成

## 快速开始

### 环境要求
- Go 1.24+
- SQLite3
- Ansible
- Docker (可选，用于容器化部署)

### 安装步骤
1. 克隆项目
```bash
git clone https://github.com/huyouba1/kde.git
cd kde
```

2. 安装依赖
```bash
go mod tidy
```

3. 配置
```bash
cp configs/config.example.yaml configs/config.yaml
# 编辑配置文件
vim configs/config.yaml
```

4. 运行
```bash
go run cmd/server/main.go
```

5. 访问
打开浏览器访问 http://localhost:8080

## 项目结构
```
kde/
├── cmd/                # 命令行入口
├── configs/           # 配置文件
├── pkg/               # 核心代码
│   ├── api/          # API 服务
│   ├── deploy/       # 部署引擎
│   ├── storage/      # 数据存储
│   └── utils/        # 工具函数
├── scripts/          # 部署脚本
├── web/              # 前端资源
└── data/             # 数据目录
```

## 开发指南

### 代码规范
- 遵循 Go 标准代码规范
- 使用 gofmt 格式化代码
- 编写单元测试

### 贡献流程
1. Fork 项目
2. 创建特性分支
3. 提交更改
4. 发起 Pull Request

## 许可证
MIT License - 详见 [LICENSE](LICENSE) 文件