package ssh

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"rexec/internal/config"
)

// Client SSH客户端封装
type Client struct {
	*ssh.Client
}

// NewClient 创建SSH客户端连接
func NewClient(server *config.ServerConfig) (*Client, error) {
	var auth []ssh.AuthMethod
	
	// 根据认证类型设置认证方法
	switch server.AuthType {
	case "password":
		auth = []ssh.AuthMethod{
			ssh.Password(server.Password),
		}
	case "key":
		keyAuth, err := createKeyAuth(server.KeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to create key auth: %v", err)
		}
		auth = []ssh.AuthMethod{keyAuth}
	default:
		return nil, fmt.Errorf("unsupported auth type: %s", server.AuthType)
	}
	
	sshConfig := &ssh.ClientConfig{
		User:            server.Username,
		Auth:            auth,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // 注意：生产环境中应该验证主机密钥
		Timeout:         30 * time.Second,
	}
	
	address := net.JoinHostPort(server.IP, server.Port)
	client, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}
	
	return &Client{Client: client}, nil
}

// createKeyAuth 创建密钥认证
func createKeyAuth(keyPath string) (ssh.AuthMethod, error) {
	keyData, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read private key file: %v", err)
	}
	
	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}
	
	return ssh.PublicKeys(signer), nil
}

// ExecuteCommand 在远程服务器上执行命令
func (c *Client) ExecuteCommand(command string) (string, error) {
	session, err := c.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	
	output, err := session.CombinedOutput(command)
	if err != nil {
		return "", fmt.Errorf("failed to execute command: %v", err)
	}
	
	return string(output), nil
}

// GetSystemInfo 获取系统信息
func (c *Client) GetSystemInfo(osType string) (map[string]string, error) {
	info := make(map[string]string)
	var commands map[string]string
	
	if osType == "windows" {
		// Windows系统命令
		commands = map[string]string{
			"hostname":     "hostname",
			"os":          "powershell \"(Get-WmiObject -Class Win32_OperatingSystem).Caption\"",
			"version":     "powershell \"(Get-WmiObject -Class Win32_OperatingSystem).Version\"",
			"uptime":      "powershell \"((Get-Date) - (Get-WmiObject Win32_OperatingSystem).ConvertToDateTime((Get-WmiObject Win32_OperatingSystem).LastBootUpTime)).Days.ToString() + ' days'\"",
			"cpu":         "powershell \"(Get-WmiObject -Class Win32_Processor | Select-Object -First 1).Name\"",
			"memory":      "powershell \"[math]::Round((Get-WmiObject -Class Win32_ComputerSystem).TotalPhysicalMemory/1GB, 2).ToString() + ' GB'\"",
			"disk":        "powershell \"$disk = Get-WmiObject -Class Win32_LogicalDisk -Filter 'DriveType=3' | Where-Object {$_.DeviceID -eq 'C:'}; [math]::Round($disk.Size/1GB, 2).ToString() + ' GB (Free: ' + [math]::Round($disk.FreeSpace/1GB, 2).ToString() + ' GB)'\"",
			"architecture": "powershell \"$env:PROCESSOR_ARCHITECTURE\"",
		}
	} else {
		// Linux系统命令
		commands = map[string]string{
			"hostname":     "hostname",
			"os":          "cat /etc/os-release | grep PRETTY_NAME | cut -d'=' -f2 | tr -d '\"'",
			"kernel":      "uname -r",
			"uptime":      "uptime -p",
			"cpu":         "cat /proc/cpuinfo | grep 'model name' | head -1 | cut -d':' -f2 | sed 's/^ *//'",
			"memory":      "free -h | grep Mem | awk '{print $2}'",
			"disk":        "df -h / | tail -1 | awk '{print $2 \" (\" $5 \" used)\"}'",
			"load":        "cat /proc/loadavg | awk '{print $1, $2, $3}'",
		}
	}
	
	for key, command := range commands {
		output, err := c.ExecuteCommand(command)
		if err != nil {
			info[key] = fmt.Sprintf("Error: %v", err)
		} else {
			info[key] = strings.TrimSpace(output)
		}
	}
	
	return info, nil
}
