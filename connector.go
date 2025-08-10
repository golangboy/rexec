package main

import (

)

// testConnectionByOS 根据操作系统类型测试连接
func testConnectionByOS(serverName string) error {
	// 现在所有系统都使用SSH连接
	return testConnection(serverName)
}

// executeRemoteCommandByOS 根据操作系统类型执行远程命令
func executeRemoteCommandByOS(serverName string, command string, args []string) error {
	// 现在所有系统都使用SSH连接
	return executeRemoteCommand(serverName, command, args)
}
