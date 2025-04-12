package etcd

import (
	"context"
	"fmt"
	"time"

	"github.com/huyouba1/kde/pkg/storage/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// Manager etcd数据库管理器
type Manager struct {
	client *clientv3.Client
	config *config.EtcdConfig
}

// NewManager 创建一个新的etcd管理器
func NewManager(cfg *config.EtcdConfig) (*Manager, error) {
	// 创建etcd客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.Endpoints,
		DialTimeout: time.Duration(cfg.DialTimeout) * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("连接etcd失败: %v", err)
	}

	return &Manager{
		client: client,
		config: cfg,
	}, nil
}

// Client 获取etcd客户端
func (m *Manager) Client() *clientv3.Client {
	return m.client
}

// Close 关闭etcd连接
func (m *Manager) Close() error {
	return m.client.Close()
}

// Put 存储键值对
func (m *Manager) Put(ctx context.Context, key, value string) error {
	_, err := m.client.Put(ctx, key, value)
	return err
}

// Get 获取键值
func (m *Manager) Get(ctx context.Context, key string) (string, error) {
	resp, err := m.client.Get(ctx, key)
	if err != nil {
		return "", err
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return string(resp.Kvs[0].Value), nil
}

// Delete 删除键值
func (m *Manager) Delete(ctx context.Context, key string) error {
	_, err := m.client.Delete(ctx, key)
	return err
}

// Watch 监听键值变化
func (m *Manager) Watch(ctx context.Context, key string) clientv3.WatchChan {
	return m.client.Watch(ctx, key)
}
