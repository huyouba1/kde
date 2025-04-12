# Kubernetes Dashboard Engine (KDE)

Kubernetes Dashboard Engine 是一个基于 Go 语言开发的 Kubernetes 集群管理平台，提供直观的 Web 界面来管理和监控 Kubernetes 集群。

## 功能特点

- 集群管理：支持多集群的添加、配置和管理
- 部署管理：提供应用部署、扩缩容和更新功能
- 资源监控：实时监控集群资源使用情况
- 系统设置：支持系统配置和插件管理
- 用户友好：提供直观的 Web 界面

## 目录结构

```
kde/
├── cmd/                    # 命令行入口
│   └── server/            # 服务器主程序
│       └── main.go        # 程序入口文件
│
├── configs/               # 配置文件目录
│   └── config.yaml        # 主配置文件
│
├── pkg/                   # 核心功能包
│   ├── api/              # API 相关代码
│   │   ├── handler/      # 请求处理器
│   │   │   ├── template_handler.go  # 模板处理逻辑
│   │   │   └── static/   # 静态资源文件
│   │   │       ├── css/  # CSS 样式文件
│   │   │       ├── js/   # JavaScript 文件
│   │   │       └── index.html  # 基础 HTML 模板
│   │   └── templates/    # HTML 模板文件
│   │       ├── base.html            # 基础布局模板
│   │       ├── index.html           # 首页模板
│   │       ├── clusters.html        # 集群列表模板
│   │       ├── cluster.html         # 集群详情模板
│   │       ├── deployments.html     # 部署管理模板
│   │       └── settings.html        # 系统设置模板
│   │
│   ├── storage/          # 数据存储相关
│   └── deploy/           # 部署相关功能
│
├── scripts/              # 脚本文件
├── web/                  # Web 前端资源
├── data/                 # 数据存储目录
├── go.mod               # Go 模块定义
└── README.md            # 项目说明文档
```

### 目录说明

#### cmd/
- 包含项目的主要入口点
- `server/main.go` 是应用程序的启动入口，负责初始化配置和启动服务器

#### configs/
- 存放项目的配置文件
- `config.yaml` 包含服务器配置、数据库配置、部署配置等

#### pkg/
- 核心功能包，包含项目的主要代码
- `api/`: API 相关代码
  - `handler/`: 请求处理器
    - `template_handler.go`: 负责 HTML 模板的加载、渲染和静态文件服务
    - `static/`: 存放静态资源文件
      - `css/`: 样式文件
      - `js/`: JavaScript 文件
      - `index.html`: 基础 HTML 模板
  - `templates/`: HTML 模板文件，使用 Go 的模板引擎
    - 包含各个页面的模板文件
    - 使用 `base.html` 作为基础布局
- `storage/`: 数据存储相关代码
- `deploy/`: 部署相关功能代码

#### scripts/
- 存放各种辅助脚本
- 如部署脚本、构建脚本等

#### web/
- 存放 Web 前端资源
- 如 JavaScript、CSS、图片等静态文件

#### data/
- 应用程序数据存储目录
- 如数据库文件、缓存文件等

## 快速开始

1. 克隆项目
```bash
git clone https://github.com/yourusername/kde.git
cd kde
```

2. 安装依赖
```bash
go mod download
```

3. 配置
- 复制 `configs/config.yaml.example` 为 `configs/config.yaml`
- 根据需求修改配置文件

4. 运行
```bash
go run cmd/server/main.go
```

5. 访问
- 打开浏览器访问 `http://localhost:8080`

## 开发指南

### 环境要求
- Go 1.16 或更高版本
- Kubernetes 集群
- SQLite 或 etcd（用于数据存储）

### 开发流程
1. 创建新的功能分支
2. 实现功能
3. 编写测试
4. 提交代码
5. 创建 Pull Request

## 贡献指南

欢迎提交 Issue 和 Pull Request 来帮助改进项目。

## 许可证

MIT License