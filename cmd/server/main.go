package main

import (
	"flag"
	"fmt"
	"github.com/huyouba1/kde/configs"
	"log"
	"os"
	"path/filepath"

	"github.com/huyouba1/kde/pkg/api"
	"github.com/huyouba1/kde/pkg/storage"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "configs/config.yaml", "配置文件路径")
}

func main() {
	flag.Parse()

	// 获取工作目录
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}

	// 解析配置文件路径
	configFile := filepath.Join(wd, configPath)
	fmt.Printf("使用配置文件: %s\n", configFile)

	// 加载配置文件
	err = configs.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	fmt.Println("Kubernetes管理系统服务启动中...")
	// 创建并启动服务器
	server, err := api.NewServer(configs.C())
	if err != nil {
		log.Fatalf("创建服务器失败: %v", err)
	}

	// 初始化SQLite数据库管理器
	sqliteManager, err := storage.NewManager(&configs.C().Database.SQLite)
	if err != nil {
		log.Fatalf("初始化SQLite数据库管理器失败: %v", err)
	}
	defer sqliteManager.Close()

	if err := server.Start(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
