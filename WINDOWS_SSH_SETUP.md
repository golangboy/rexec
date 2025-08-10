# Windows OpenSSH Server 配置指南

要使用 rexec 通过 SSH 连接到 Windows 服务器，需要在目标 Windows 机器上安装并配置 OpenSSH Server。

## Windows 10 / Windows Server 2019 及更新版本

### 1. 检查 OpenSSH 可用性

以管理员身份打开 PowerShell，检查 OpenSSH 功能：

```powershell
# 查看可用的 OpenSSH 功能
Get-WindowsCapability -Online | Where-Object Name -like 'OpenSSH*'
```

### 2. 安装 OpenSSH Server

```powershell
# 安装 OpenSSH Server
Add-WindowsCapability -Online -Name OpenSSH.Server~~~~0.0.1.0

# 安装 OpenSSH Client（可选，通常已安装）
Add-WindowsCapability -Online -Name OpenSSH.Client~~~~0.0.1.0
```

### 3. 启动并配置 SSH 服务

```powershell
# 启动 SSH 服务
Start-Service sshd

# 设置 SSH 服务为自动启动
Set-Service -Name sshd -StartupType 'Automatic'

# 确认防火墙规则已创建（通常自动创建）
if (!(Get-NetFirewallRule -Name "OpenSSH-Server-In-TCP" -ErrorAction SilentlyContinue | Select-Object Name, Enabled)) {
    Write-Output "Firewall Rule 'OpenSSH-Server-In-TCP' does not exist, creating it..."
    New-NetFirewallRule -Name 'OpenSSH-Server-In-TCP' -DisplayName 'OpenSSH Server (sshd)' -Enabled True -Direction Inbound -Protocol TCP -Action Allow -LocalPort 22
} else {
    Write-Output "Firewall rule 'OpenSSH-Server-In-TCP' has been created and exists."
}
```

## 较旧版本的 Windows

### 使用 Chocolatey 安装

```powershell
# 安装 Chocolatey（如果未安装）
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# 安装 OpenSSH
choco install openssh -y
```

### 手动下载安装

1. 从 [GitHub 发布页面](https://github.com/PowerShell/Win32-OpenSSH/releases) 下载最新版本
2. 解压到 `C:\Program Files\OpenSSH`
3. 以管理员身份运行安装脚本：

```powershell
cd "C:\Program Files\OpenSSH"
.\install-sshd.ps1
```

## 配置说明

### 1. SSH 配置文件

默认配置文件位置：`C:\ProgramData\ssh\sshd_config`

常用配置选项：

```bash
# 端口设置（默认22）
Port 22

# 允许密码认证
PasswordAuthentication yes

# 允许公钥认证
PubkeyAuthentication yes

# 指定授权密钥文件位置
AuthorizedKeysFile .ssh/authorized_keys

# 日志级别
LogLevel INFO
```

### 2. 用户配置

#### 密码认证
确保用户账户有正确的密码，并且账户未被禁用。

#### 密钥认证
1. 在用户主目录创建 `.ssh` 文件夹：
   ```cmd
   mkdir C:\Users\%USERNAME%\.ssh
   ```

2. 将公钥内容添加到 `authorized_keys` 文件：
   ```cmd
   echo "ssh-rsa AAAAB3NzaC1yc2EAAAA..." >> C:\Users\%USERNAME%\.ssh\authorized_keys
   ```

3. 设置正确的权限：
   ```powershell
   icacls C:\Users\%USERNAME%\.ssh /inheritance:r
   icacls C:\Users\%USERNAME%\.ssh /grant:r "%USERNAME%:(OI)(CI)F"
   icacls C:\Users\%USERNAME%\.ssh\authorized_keys /inheritance:r
   icacls C:\Users\%USERNAME%\.ssh\authorized_keys /grant:r "%USERNAME%:F"
   ```

## PowerShell 作为默认 Shell

要将 PowerShell 设置为 SSH 的默认 shell：

```powershell
# 设置 PowerShell 为默认 shell
New-ItemProperty -Path "HKLM:\SOFTWARE\OpenSSH" -Name DefaultShell -Value "C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe" -PropertyType String -Force
```

## 测试连接

### 本地测试
```powershell
ssh localhost
```

### 远程测试
从另一台机器测试：
```bash
ssh username@windows-server-ip
```

## 故障排除

### 1. 检查服务状态
```powershell
Get-Service sshd
```

### 2. 查看 SSH 日志
日志位置：`C:\ProgramData\ssh\logs\sshd.log`

### 3. 检查防火墙
```powershell
Get-NetFirewallRule -Name "OpenSSH*"
```

### 4. 重启 SSH 服务
```powershell
Restart-Service sshd
```

### 5. 测试端口连通性
```powershell
Test-NetConnection -ComputerName localhost -Port 22
```

## 安全建议

1. **更改默认端口**：考虑将 SSH 端口从 22 改为其他端口
2. **禁用 root 登录**：Windows 中对应禁用 Administrator 直接 SSH 登录
3. **使用密钥认证**：比密码认证更安全
4. **配置防火墙**：只允许必要的 IP 地址访问
5. **定期更新**：保持 OpenSSH 版本为最新

## rexec 使用示例

配置完成后，使用 rexec 添加 Windows 服务器：

```bash
# 添加 Windows 服务器
rexec add windows 192.168.1.100 22 mywindows

# 测试连接
rexec mywindows test

# 执行 PowerShell 命令
rexec exec mywindows Get-Process
rexec exec mywindows dir C:\
```
