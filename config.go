package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ServerConfig 服务器配置结构
type ServerConfig struct {
	Name     string `yaml:"name"`
	IP       string `yaml:"ip"`
	Port     string `yaml:"port"`
	OS       string `yaml:"os"`
	AuthType string `yaml:"auth_type"` // "password" or "key"
	Password string `yaml:"password,omitempty"`
	KeyPath  string `yaml:"key_path,omitempty"`
	Username string `yaml:"username"`
}

// Config 总配置结构
type Config struct {
	Servers []ServerConfig `yaml:"servers"`
}

// getConfigPath 获取配置文件路径
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	
	configDir := filepath.Join(homeDir, ".rexec")
	
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %v", err)
	}
	
	return filepath.Join(configDir, "config.yaml"), nil
}

// loadConfig 加载配置文件
func loadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}
	
	config := &Config{
		Servers: []ServerConfig{},
	}
	
	// 如果配置文件不存在，返回空配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	
	return config, nil
}

// saveConfig 保存配置文件
func saveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}
	
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	
	return nil
}

// findServer 根据名称查找服务器配置
func findServer(name string) (*ServerConfig, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, err
	}
	
	for _, server := range config.Servers {
		if server.Name == name {
			return &server, nil
		}
	}
	
	return nil, fmt.Errorf("server '%s' not found", name)
}

// addServerConfig 添加服务器配置
func addServerConfig(serverConfig ServerConfig) error {
	config, err := loadConfig()
	if err != nil {
		return err
	}
	
	// 检查是否已存在同名服务器
	for _, server := range config.Servers {
		if server.Name == serverConfig.Name {
			return fmt.Errorf("server '%s' already exists", serverConfig.Name)
		}
	}
	
	config.Servers = append(config.Servers, serverConfig)
	
	return saveConfig(config)
}
