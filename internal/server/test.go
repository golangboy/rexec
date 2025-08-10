package server

import (
	"fmt"
	"strings"

	"rexec/internal/config"
	"rexec/internal/ssh"
)

// TestConnection 测试与指定服务器的连接
func TestConnection(serverName string) error {
	// 查找服务器配置
	server, err := config.FindServer(serverName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Testing connection to server '%s' (%s:%s)...\n", serverName, server.IP, server.Port)
	
	// 创建SSH连接
	client, err := ssh.NewClient(server)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer client.Close()
	
	fmt.Println("✓ Connection successful!")
	fmt.Println()
	
	// 获取系统信息
	fmt.Println("Gathering system information...")
	sysInfo, err := client.GetSystemInfo(server.OS)
	if err != nil {
		return fmt.Errorf("failed to get system info: %v", err)
	}
	
	// 显示系统信息
	displaySystemInfo(sysInfo, server.OS)
	
	return nil
}

// displaySystemInfo 格式化显示系统信息
func displaySystemInfo(info map[string]string, osType string) {
	if osType == "windows" {
		fmt.Println("Windows System Information:")
		fmt.Println("==========================")
		
		// Windows字段显示顺序
		fields := []struct {
			key   string
			label string
		}{
			{"hostname", "Computer Name"},
			{"os", "Operating System"},
			{"version", "OS Version"},
			{"architecture", "Architecture"},
			{"uptime", "Uptime"},
			{"cpu", "CPU"},
			{"memory", "Total Memory"},
			{"disk", "C: Drive"},
		}
		
		for _, field := range fields {
			if value, exists := info[field.key]; exists {
				value = strings.TrimSpace(value)
				if value != "" {
					fmt.Printf("%-16s: %s\n", field.label, value)
				}
			}
		}
	} else {
		fmt.Println("Linux System Information:")
		fmt.Println("=========================")
		
		// Linux字段显示顺序
		fields := []struct {
			key   string
			label string
		}{
			{"hostname", "Hostname"},
			{"os", "Operating System"},
			{"kernel", "Kernel Version"},
			{"uptime", "Uptime"},
			{"cpu", "CPU"},
			{"memory", "Total Memory"},
			{"disk", "Root Disk"},
			{"load", "Load Average"},
		}
		
		for _, field := range fields {
			if value, exists := info[field.key]; exists {
				value = strings.TrimSpace(value)
				if value != "" {
					fmt.Printf("%-16s: %s\n", field.label, value)
				}
			}
		}
	}
}
