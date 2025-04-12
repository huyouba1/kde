package main

import (
	"fmt"
	"time"

	"github.com/huyouba1/kde/pkg/plugin"
)

// ExamplePlugin 示例插件实现
type ExamplePlugin struct {
	info    plugin.PluginInfo
	running bool
}

// Plugin 导出的插件变量，必须命名为Plugin
var Plugin = &ExamplePlugin{
	info: plugin.PluginInfo{
		ID:          "example-plugin",
		Name:        "示例插件",
		Description: "这是一个示例插件，展示如何实现插件接口",
		Version:     "1.0.0",
		Author:      "KDE Team",
		Path:        "example.so",
	},
	running: false,
}

// GetInfo 获取插件信息
func (p *ExamplePlugin) GetInfo() plugin.PluginInfo {
	return p.info
}

// Init 初始化插件
func (p *ExamplePlugin) Init() error {
	fmt.Println("示例插件初始化中...")
	return nil
}

// Start 启动插件
func (p *ExamplePlugin) Start() error {
	fmt.Println("示例插件启动中...")
	p.running = true

	// 这里可以启动一个goroutine执行插件的主要功能
	go func() {
		for p.running {
			// 插件的主要逻辑
			fmt.Println("示例插件正在运行...")
			time.Sleep(5 * time.Second)
		}
	}()

	return nil
}

// Stop 停止插件
func (p *ExamplePlugin) Stop() error {
	fmt.Println("示例插件停止中...")
	p.running = false
	return nil
}

// 插件可以定义自己的功能函数
func (p *ExamplePlugin) DoSomething() string {
	return "示例插件执行了一些操作"
}

// 必须有main函数，即使为空
func main() {}
