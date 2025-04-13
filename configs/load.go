package configs

import (
	"fmt"
	"github.com/spf13/viper"
)

// 如何把配置映射成config对象

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	config := NewDefaultConfig()
	if err := viper.Unmarshal(&config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	global = config
	return nil
}

// 从环境变量加载配置
//func LoadConfigFromEnv() error {
//	config := NewDefaultConfig()
//	err := env.Parse(config)
//	if err != nil {
//		return err
//	}
//	global = config
//	return nil
//}
