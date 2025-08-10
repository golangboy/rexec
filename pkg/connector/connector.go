package connector

import (
	"rexec/internal/executor"
	"rexec/internal/server"
)

// TestConnectionByOS 根据操作系统类型测试连接
func TestConnectionByOS(serverName string) error {
	// 现在所有系统都使用SSH连接
	return server.TestConnection(serverName)
}

// ExecuteRemoteCommandByOS 根据操作系统类型执行远程命令
func ExecuteRemoteCommandByOS(serverName string, command string, args []string) error {
	// 现在所有系统都使用SSH连接
	return executor.ExecuteRemoteCommand(serverName, command, args)
}
