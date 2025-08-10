# rexec - 远程执行工具

**中文** | [English](README.en.md)

rexec 是一个简单的远程执行工具，支持 Linux 和 Windows 服务器的远程连接和基本系统信息获取。

## 功能特性

- 支持添加和管理远程 Linux 和 Windows 服务器
- 统一使用 SSH 协议连接所有服务器（Linux 和 Windows）
- 支持密码和 SSH 密钥两种认证方式
- 安全存储服务器配置信息到用户家目录
- 提供连接测试功能，获取远程系统信息
- 支持远程命令执行，输出重定向到本地终端

## 安装

1. 确保已安装 Go 1.21 或更高版本
2. 克隆或下载项目代码
3. 在项目目录中运行：

```bash
go mod tidy
go build -o rexec
```

## 使用方法

### 添加服务器

**添加 Linux 服务器：**
```bash
rexec add linux <IP地址> <端口> <服务器名称>
```

**添加 Windows 服务器：**
```bash
rexec add windows <IP地址> <端口> <服务器名称>
```

例如：
```bash
rexec add linux 192.168.1.100 22 mylinux
rexec add windows 192.168.1.101 22 mywindows
```

运行此命令后，系统会提示您输入：
- 用户名
- 认证方式：密码或 SSH 密钥认证
- 相应的认证信息（密码或密钥文件路径）

**注意：** Windows 服务器需要安装并启用 OpenSSH Server 服务。

### 测试连接

```bash
rexec <服务器名称> test
```

例如：
```bash
rexec mylinux test
rexec mywindows test
```

此命令会连接到指定服务器并获取系统信息，包括：

**Linux 系统信息：**
- 主机名、操作系统版本、内核版本
- 系统运行时间、CPU 信息、内存信息
- 磁盘使用情况、系统负载

**Windows 系统信息：**
- 计算机名、操作系统版本、系统架构
- 系统运行时间、CPU 信息、内存信息
- C盘使用情况

### 执行远程命令

```bash
rexec exec <服务器名称> <命令> [参数...]
```

**Linux 服务器示例：**
```bash
# 列出远程目录
rexec exec mylinux ls -la /home

# 执行Python脚本
rexec exec mylinux python3 /path/to/script.py arg1 arg2

# 查看系统进程
rexec exec mylinux ps aux

# 执行复杂命令（会自动处理参数转义）
rexec exec mylinux find /var/log -name "*.log" -mtime -1
```

**Windows 服务器示例：**
```bash
# 列出目录内容
rexec exec mywindows dir C:\Users

# 查看系统进程
rexec exec mywindows Get-Process

# 执行PowerShell命令
rexec exec mywindows Get-Service | Where-Object {$_.Status -eq 'Running'}

# 查看系统信息
rexec exec mywindows systeminfo
```

该功能特点：
- 实时输出重定向到本地终端
- 支持命令参数传递
- 自动处理特殊字符转义，防止注入攻击
- 支持标准输入、输出、错误流

## 配置文件

服务器配置信息存储在用户家目录的 `.rexec/config.yaml` 文件中。该文件包含所有添加的服务器信息（密码经过安全处理）。

## 注意事项

**通用：**
- 配置文件中的敏感信息（如密码）会以明文形式存储，请确保系统安全
- SSH 连接暂时使用 `InsecureIgnoreHostKey`，生产环境中建议实现主机密钥验证
- 默认 SSH 端口为 22

**Windows 服务器要求：**
- 需要安装并启用 OpenSSH Server 服务
- Windows 10 版本 1809 及以上版本内置支持 OpenSSH
- 旧版本 Windows 需要手动安装 OpenSSH Server

## 未来计划

- 支持 macOS 系统
- 增强安全性（密码加密存储、主机密钥验证）
- 支持批量操作
- 添加日志记录功能
- 支持文件传输功能
- 添加会话管理功能
