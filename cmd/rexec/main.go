package main

import (
	"fmt"
	"os"

	"rexec/internal/server"
	"rexec/pkg/connector"
)

func main() {
	if len(os.Args) < 2 {
		showUsage()
		return
	}

	command := os.Args[1]
	switch command {
	case "add":
		if len(os.Args) < 6 {
			fmt.Println("Usage: rexec add <linux|windows> <ip> <port> <name>")
			return
		}
		osType := os.Args[2]
		ip := os.Args[3]
		port := os.Args[4]
		name := os.Args[5]
		
		if osType != "linux" && osType != "windows" {
			fmt.Println("Supported OS types: linux, windows")
			return
		}
		
		err := server.AddServer(ip, port, name, osType)
		if err != nil {
			fmt.Printf("Failed to add server: %v\n", err)
			return
		}
		fmt.Printf("Server %s added successfully\n", name)
		
	case "exec":
		if len(os.Args) < 4 {
			fmt.Println("Usage: rexec exec <server_name> <command> [args...]")
			return
		}
		serverName := os.Args[2]
		execCommand := os.Args[3]
		args := []string{}
		if len(os.Args) > 4 {
			args = os.Args[4:]
		}
		
		err := connector.ExecuteRemoteCommandByOS(serverName, execCommand, args)
		if err != nil {
			fmt.Printf("Execution failed: %v\n", err)
			return
		}
		
	default:
		// Check if command is a server name for test
		if len(os.Args) == 3 && os.Args[2] == "test" {
			err := connector.TestConnectionByOS(command)
			if err != nil {
				fmt.Printf("Test failed: %v\n", err)
				return
			}
		} else {
			showUsage()
		}
	}
}

func showUsage() {
	fmt.Println("Remote Execution Tool (rexec)")
	fmt.Println("Usage:")
	fmt.Println("  rexec add <linux|windows> <ip> <port> <name>  - Add a server")
	fmt.Println("  rexec <name> test                             - Test connection to server")
	fmt.Println("  rexec exec <name> <command> [args...]         - Execute command on remote server")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  rexec add linux 192.168.1.100 22 mylinux")
	fmt.Println("  rexec add windows 192.168.1.101 5985 mywindows")
	fmt.Println("  rexec mylinux test")
	fmt.Println("  rexec mywindows test")
	fmt.Println("  rexec exec mylinux ls -la /home")
	fmt.Println("  rexec exec mywindows dir C:\\Users")
	fmt.Println("  rexec exec mywindows Get-Process")
}
