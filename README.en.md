# rexec - Remote Execution Tool

[中文版](README.md) | **English**

rexec is a simple remote execution tool that supports remote connection and basic system information retrieval for both Linux and Windows servers.

## Features

- Add and manage remote Linux and Windows servers
- Unified SSH protocol connection for all servers (both Linux and Windows)
- Support for both password and SSH key authentication methods
- Secure storage of server configuration information in user home directory
- Connection testing functionality with remote system information retrieval
- Remote command execution with output redirection to local terminal

## Installation

1. Ensure Go 1.21 or higher is installed
2. Clone or download the project code
3. Run in the project directory:

```bash
go mod tidy
go build -o rexec
```

## Usage

### Adding Servers

**Add Linux server:**
```bash
rexec add linux <IP_Address> <Port> <Server_Name>
```

**Add Windows server:**
```bash
rexec add windows <IP_Address> <Port> <Server_Name>
```

Examples:
```bash
rexec add linux 192.168.1.100 22 mylinux
rexec add windows 192.168.1.101 22 mywindows
```

After running this command, the system will prompt you to enter:
- Username
- Authentication method: password or SSH key authentication
- Corresponding authentication information (password or key file path)

**Note:** Windows servers require OpenSSH Server service to be installed and enabled.

### Testing Connection

```bash
rexec <Server_Name> test
```

Examples:
```bash
rexec mylinux test
rexec mywindows test
```

This command will connect to the specified server and retrieve system information, including:

**Linux System Information:**
- Hostname, OS version, kernel version
- System uptime, CPU information, memory information
- Disk usage, system load

**Windows System Information:**
- Computer name, OS version, system architecture
- System uptime, CPU information, memory information
- C drive usage

### Remote Command Execution

```bash
rexec exec <Server_Name> <Command> [Arguments...]
```

**Linux Server Examples:**
```bash
# List remote directory
rexec exec mylinux ls -la /home

# Execute Python script
rexec exec mylinux python3 /path/to/script.py arg1 arg2

# View system processes
rexec exec mylinux ps aux

# Execute complex commands (automatic argument escaping)
rexec exec mylinux find /var/log -name "*.log" -mtime -1
```

**Windows Server Examples:**
```bash
# List directory contents
rexec exec mywindows dir C:\Users

# View system processes
rexec exec mywindows Get-Process

# Execute PowerShell commands
rexec exec mywindows Get-Service | Where-Object {$_.Status -eq 'Running'}

# View system information
rexec exec mywindows systeminfo
```

Key Features:
- Real-time output redirection to local terminal
- Support for command argument passing
- Automatic special character escaping to prevent injection attacks
- Support for standard input, output, and error streams

## Configuration File

Server configuration information is stored in the `.rexec/config.yaml` file in the user's home directory. This file contains all added server information (passwords are securely handled).

## Important Notes

**General:**
- Sensitive information (such as passwords) in the configuration file is stored in plain text, please ensure system security
- SSH connections currently use `InsecureIgnoreHostKey`, host key verification is recommended for production environments
- Default SSH port is 22

**Windows Server Requirements:**
- Requires OpenSSH Server service to be installed and enabled
- Windows 10 version 1809 and above have built-in OpenSSH support
- Older Windows versions require manual OpenSSH Server installation

## Future Plans

- Support for macOS systems
- Enhanced security (encrypted password storage, host key verification)
- Support for batch operations
- Add logging functionality
- Support for file transfer functionality
- Add session management functionality
