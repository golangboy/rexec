package server

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"golang.org/x/term"
	"rexec/internal/config"
)

// AddServer 添加服务器配置
func AddServer(ip, port, name, osType string) error {
	// 验证端口号
	if _, err := strconv.Atoi(port); err != nil {
		return fmt.Errorf("invalid port number: %s", port)
	}
	
	reader := bufio.NewReader(os.Stdin)
	
	// 获取用户名
	fmt.Print("Enter username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read username: %v", err)
	}
	username = strings.TrimSpace(username)
	
	if username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	
	serverConfig := config.ServerConfig{
		Name:     name,
		IP:       ip,
		Port:     port,
		OS:       osType,
		Username: username,
	}
	
	// 选择认证方式（所有系统都支持SSH认证）
	fmt.Println("Select authentication method:")
	fmt.Println("1. Password")
	fmt.Println("2. Private Key")
	fmt.Print("Enter choice (1 or 2): ")
	
	choice, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read choice: %v", err)
	}
	choice = strings.TrimSpace(choice)
	
	// 为Windows用户提供SSH提示
	if osType == "windows" {
		fmt.Println("\nNote: Windows server must have OpenSSH Server installed and enabled.")
		if port != "22" {
			fmt.Printf("Note: Port %s is not the standard SSH port (22).\n", port)
		}
	}
	
	switch choice {
	case "1":
		// 密码认证
		serverConfig.AuthType = "password"
		password, err := readPassword("Enter password: ")
		if err != nil {
			return fmt.Errorf("failed to read password: %v", err)
		}
		serverConfig.Password = password
		
	case "2":
		// 密钥认证
		serverConfig.AuthType = "key"
		fmt.Print("Enter private key file path: ")
		keyPath, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read key path: %v", err)
		}
		keyPath = strings.TrimSpace(keyPath)
		
		if keyPath == "" {
			return fmt.Errorf("key path cannot be empty")
		}
		
		// 验证密钥文件是否存在
		if _, err := os.Stat(keyPath); os.IsNotExist(err) {
			return fmt.Errorf("private key file does not exist: %s", keyPath)
		}
		
		serverConfig.KeyPath = keyPath
		
	default:
		return fmt.Errorf("invalid choice: %s", choice)
	}
	
	// 保存配置
	if err := config.AddServerConfig(serverConfig); err != nil {
		return err
	}
	
	return nil
}

// readPassword 安全地读取密码（不显示在终端）
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	
	// 获取终端文件描述符
	fd := int(syscall.Stdin)
	
	// 读取密码（不回显）
	password, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}
	
	fmt.Println() // 换行
	
	return string(password), nil
}
