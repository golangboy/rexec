package executor

import (
	"fmt"
	"io"
	"os"
	"strings"

	"rexec/internal/config"
	"rexec/internal/ssh"
)

// ExecuteRemoteCommand 在远程服务器上执行命令并实时输出
func ExecuteRemoteCommand(serverName string, command string, args []string) error {
	// 查找服务器配置
	server, err := config.FindServer(serverName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Executing on server '%s' (%s:%s): %s\n", serverName, server.IP, server.Port, command)
	if len(args) > 0 {
		fmt.Printf("Arguments: %s\n", strings.Join(args, " "))
	}
	fmt.Println("----------------------------------------")
	
	// 创建SSH连接
	client, err := ssh.NewClient(server)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer client.Close()
	
	// 创建SSH会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	
	// 构建完整的命令
	fullCommand := command
	if len(args) > 0 {
		// 对参数进行转义以防止shell注入
		escapedArgs := make([]string, len(args))
		for i, arg := range args {
			escapedArgs[i] = shellEscape(arg)
		}
		fullCommand = command + " " + strings.Join(escapedArgs, " ")
	}
	
	// 设置输出重定向到当前终端
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	
	// 执行命令
	err = session.Run(fullCommand)
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}
	
	return nil
}

// ExecuteRemoteCommandInteractive 交互式执行远程命令（支持实时输入输出）
func ExecuteRemoteCommandInteractive(serverName string, command string, args []string) error {
	// 查找服务器配置
	server, err := config.FindServer(serverName)
	if err != nil {
		return err
	}
	
	fmt.Printf("Executing on server '%s' (%s:%s): %s\n", serverName, server.IP, server.Port, command)
	if len(args) > 0 {
		fmt.Printf("Arguments: %s\n", strings.Join(args, " "))
	}
	fmt.Println("----------------------------------------")
	
	// 创建SSH连接
	client, err := ssh.NewClient(server)
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer client.Close()
	
	// 创建SSH会话
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()
	
	// 构建完整的命令
	fullCommand := command
	if len(args) > 0 {
		escapedArgs := make([]string, len(args))
		for i, arg := range args {
			escapedArgs[i] = shellEscape(arg)
		}
		fullCommand = command + " " + strings.Join(escapedArgs, " ")
	}
	
	// 创建管道连接
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}
	
	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}
	
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %v", err)
	}
	
	// 启动命令
	err = session.Start(fullCommand)
	if err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}
	
	// 启动goroutine来处理输入输出
	go func() {
		io.Copy(os.Stdout, stdout)
	}()
	
	go func() {
		io.Copy(os.Stderr, stderr)
	}()
	
	go func() {
		io.Copy(stdin, os.Stdin)
	}()
	
	// 等待命令完成
	err = session.Wait()
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}
	
	return nil
}

// shellEscape 转义shell参数，防止注入攻击
func shellEscape(arg string) string {
	// 如果参数包含特殊字符，用单引号包围
	if strings.ContainsAny(arg, " \t\n\r;|&<>(){}[]$`\"'\\*?") {
		// 替换单引号为 '\''
		escaped := strings.ReplaceAll(arg, "'", "'\"'\"'")
		return "'" + escaped + "'"
	}
	return arg
}
