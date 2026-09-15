# CFData-WEB Plus

> 基于 [PoemMisty/CFData-WEB](https://github.com/PoemMisty/CFData-WEB) 的 Docker / GHCR 增强版本。

CFData-WEB Plus 是一个基于 Go + WebSocket 的 Cloudflare IP 扫描、延迟测试、HTTPing、测速、筛选与结果导出工具。Plus 版本在尽量保持原项目功能和使用方式的基础上，增加了 **Docker / GHCR 多架构部署、服务器端配置持久化、Cloudflare 地区 IP 数据源、地区筛选与一键导入非标优选测速** 等功能。

## 目录

- [主要功能](#主要功能)
- [Plus 增强功能](#plus-增强功能)
- [Docker 部署](#docker-部署)
- [NAS 部署](#nas-部署)
- [从源码构建](#从源码构建)
- [Web 使用](#web-使用)
- [官方优选](#官方优选)
- [非标优选](#非标优选)
- [Cloudflare 地区 IP](#cloudflare-地区-ip)
- [服务器配置持久化](#服务器配置持久化)
- [数据目录与持久化](#数据目录与持久化)
- [GHCR 镜像](#ghcr-镜像)
- [配置说明](#配置说明)
- [项目结构](#项目结构)
- [第三方项目与数据源](#第三方项目与数据源)
- [免责声明](#免责声明)
- [License](#license)

---

## 主要功能

### 官方优选

- Cloudflare IPv4 / IPv6 地址扫描
- TCPing 延迟测试
- HTTPing 延迟测试
- 数据中心（DC）结果整理
- 延迟阈值筛选
- 扫描并发控制
- 详细测试与测速
- 测速地址自动选择或手动指定

### 非标优选

- 支持本地 TXT / CSV 文件
- 支持网络 URL 输入
- 支持 `IP`、`IP PORT`、`IP:PORT` 等常见输入形式
- IPv4 / IPv6 地址处理
- 缺少端口时自动使用备用端口
- TCPing / HTTPing
- TLS 开关
- 并发、延迟、结果数量和测速阈值控制
- 测试结果筛选、导出

### 结果处理

- CSV / TXT 导出
- 自定义导出字段
- IP 类型筛选
- 合格结果筛选
- GitHub 上传

---

## Plus 增强功能

### 1. Docker 化部署

项目提供正式 Dockerfile，可通过 Docker / Docker Compose 部署。

支持构建：

- `linux/amd64`
- `linux/arm64`

适合：

- x86 NAS
- ARM64 NAS
- Linux 服务器
- 家用服务器
- Docker 环境

### 2. GHCR 镜像

Plus 版本使用独立镜像，不覆盖原项目镜像：

```text
ghcr.io/jaysonpeng/cfdata-web-plus:latest
```

### 3. 服务器配置持久化

Web 页面支持：

- 自动读取服务器配置
- 保存当前 Web 配置到服务器
- 页面重新打开后自动同步
- WebSocket 重连后自动同步

Docker 中配置文件默认保存为：

```text
/data/cfdata-config.json
```

通过 Docker Volume 持久化后，容器删除、更新或重建不会丢失配置。

### 4. Cloudflare 地区 IP

在「非标优选」中增加 Cloudflare 地区 IP 功能：

- 自动加载地区列表
- 显示各地区 IP 数量
- 地区多选
- 每地区数量根据该地区实际 IP 数量自动设置
- 支持自定义端口
- 支持使用数据源原始端口
- 支持 IP / IP:PORT / 地区标签等输出格式
- 支持运营商标签：电信、移动、联通
- 一键复制
- 一键导入非标优选测速

地区 IP 获取完成后，可以直接进入原有的非标扫描、TCPing、HTTPing 和测速流程。

---

# Docker 部署

## 使用 Docker Compose（推荐）

项目已经提供 `docker-compose.yml`：

```yaml
services:
  cfdata-web:
    image: ghcr.io/jaysonpeng/cfdata-web-plus:latest
    container_name: cfdata-web
    restart: unless-stopped

    working_dir: /data

    ports:
      - "13335:13335"

    environment:
      TZ: Asia/Shanghai
      CFDATA_CONFIG_PATH: /data/cfdata-config.json

    volumes:
      - ./data:/data
```

进入项目目录后执行：

```bash
docker compose pull
docker compose up -d
```

查看容器：

```bash
docker ps
```

查看日志：

```bash
docker logs -f cfdata-web
```

停止：

```bash
docker compose down
```

更新：

```bash
docker compose pull
docker compose up -d
```

Web 默认端口：

```text
13335
```

浏览器访问：

```text
http://NAS_IP:13335
```

例如：

```text
http://192.168.1.100:13335
```

---

# NAS 部署

以常见 NAS Docker 环境为例，建议建立：

```text
cfdata-web/
├── docker-compose.yml
└── data/
```

然后：

```bash
cd /path/to/cfdata-web
docker compose pull
docker compose up -d
```

## 非常重要：不要挂载 `/app`

程序本体位于容器：

```text
/app/cfdata
```

数据目录位于：

```text
/data
```

因此应该使用：

```yaml
volumes:
  - ./data:/data
```

不要使用：

```yaml
volumes:
  - ./data:/app
```

否则宿主机目录会覆盖程序所在目录，导致容器无法正常启动。

---

# 从源码构建

正式 Web 源码位于：

```text
combined_refactor/
```

其 Go 模块要求：

```text
go 1.25.4
```

本地构建：

```bash
cd combined_refactor
go mod download
go build -o cfdata .
```

启动：

```bash
./cfdata -port 13335
```

然后访问：

```text
http://127.0.0.1:13335
```

---

# Web 使用

打开 Web 页面后，程序会建立 WebSocket 连接。

正常情况下顶部状态应显示：

```text
已连接服务器
```

Plus 版本启动后会自动同步服务器配置，不需要每次手动点击「从服务器同步配置」。

同时，Cloudflare 地区 IP 数据也会自动加载，不需要手动点击「加载地区」。

---

# 官方优选

进入「官方优选」后，可以设置：

- IP 类型：IPv4 / IPv6
- 测试端口
- 扫描线程
- 延迟阈值
- 测速地址
- 测速阈值
- 自动测速等

点击：

```text
开始扫描与测试
```

完成初步扫描后，可根据数据中心继续进行详细测试和测速。

---

# 非标优选

进入「非标优选」后，可以使用：

1. 本地文件
2. 网络 URL
3. Cloudflare 地区 IP

## 推荐输入格式

```text
1.2.3.4 443
5.6.7.8 8443
2606:4700::1111 443
1.1.1.1
```

未提供端口时，程序会使用备用端口。

常见规则：

- TLS 开启：默认备用端口 `443`
- TLS 关闭：默认备用端口 `80`

---

# Cloudflare 地区 IP

「非标优选」内置「Cloudflare 地区 IP」功能。

## 自动加载

进入页面后会自动获取地区统计信息，例如：

```text
日本 JP
香港 HK
新加坡 SG
台湾 TW
韩国 KR
```

每个地区同时显示数据源中的 IP 数量。

## 数量逻辑

Plus 版本的默认数量不是固定的 50，也不是将 50 平均分配到各地区。

**默认数量直接根据该地区实际可用 IP 总数自动设置。**

例如数据源返回：

```text
日本 JP：400
香港 HK：297
新加坡 SG：571
```

选择日本后：

```text
每地区数量 = 400
```

选择香港后：

```text
每地区数量 = 297
```

选择新加坡后：

```text
每地区数量 = 571
```

因此默认情况下，会尽可能获取所选地区的数据源全部 IP。

如果用户手动修改数量，则按照手动输入数量执行。

## 端口

支持：

- 使用数据源原始端口
- 自定义端口

例如：

```text
443
8443
2053
2083
2087
2096
```

## 输出格式

支持：

```text
IP:PORT#🇯🇵 日本
```

```text
IP:PORT#运营商+地区
```

```text
IP:PORT
```

```text
仅 IP
```

运营商可以选择：

```text
电信
移动
联通
```

## 导入非标测速

生成地区 IP 后点击：

```text
导入非标优选测速
```

程序会将生成的数据直接交给非标优选模块。

之后点击原有的：

```text
开始扫描与测试
```

即可继续 TCPing / HTTPing / 测速。

---

# 服务器配置持久化

Plus 版本将 Web 配置保存到服务器，而不是只保存到浏览器 `localStorage`。

默认配置文件：

```text
/data/cfdata-config.json
```

Docker Compose 中通过：

```yaml
environment:
  CFDATA_CONFIG_PATH: /data/cfdata-config.json
```

指定配置路径。

页面启动后：

```text
WebSocket 连接
    ↓
自动请求服务器配置
    ↓
读取 cfdata-config.json
    ↓
同步到 Web 表单
```

点击「保存当前配置到服务器」后，配置会写入持久化目录。

如果使用：

```yaml
volumes:
  - ./data:/data
```

那么宿主机中会出现：

```text
./data/cfdata-config.json
```

---

# 数据目录与持久化

推荐目录：

```text
cfdata-web/
├── docker-compose.yml
└── data/
    ├── cfdata-config.json
    ├── ips-v4.txt
    ├── ips-v6.txt
    ├── locations.json
    └── 其他运行数据
```

实际文件会根据运行功能产生，不要求提前创建。

**建议定期备份 `data/` 目录。**

---

# GHCR 镜像

GitHub Actions 自动构建并发布：

```text
ghcr.io/jaysonpeng/cfdata-web-plus:latest
```

当前 Workflow 支持：

```text
linux/amd64
linux/arm64
```

构建文件：

```text
.github/workflows/docker-ghcr.yml
```

触发方式：

- 推送到 `main`
- GitHub Actions 手动运行 `workflow_dispatch`

更新镜像：

```bash
docker compose pull
docker compose up -d
```

---

# 配置说明

## Docker 环境变量

| 变量 | 默认值 | 说明 |
|---|---|---|
| `TZ` | `Asia/Shanghai` | 容器时区 |
| `CFDATA_CONFIG_PATH` | 程序目录下 `cfdata-config.json` | Web 服务器配置文件路径 |

例如：

```yaml
environment:
  TZ: Asia/Shanghai
  CFDATA_CONFIG_PATH: /data/cfdata-config.json
```

## Web 默认端口

```text
13335
```

可以修改宿主机端口，例如：

```yaml
ports:
  - "18080:13335"
```

访问：

```text
http://NAS_IP:18080
```

---

# TCPing 与 HTTPing

## TCPing

TCPing 主要测量 TCP 建立连接所需的时间。

## HTTPing

HTTPing 主要关注 HTTP TTFB（Time To First Byte）。

HTTPing 与 TCPing 的测试对象不同，因此：

> 不同扫描模式得到的延迟数据不建议直接横向比较，应在同一种扫描模式下比较结果。

HTTPing 的实际延迟通常高于 TCPing，这是正常现象。

---

# 测速地址

默认测速地址：

```text
auto
```

表示由程序自动选择测速源。

Web 页面支持：

- 自动选择
- Cloudflare
- CM 提供
- 移动专属
- 手动输入

CLI 可通过 `-offurl` / `-nsburl` 指定测速地址。

示例：

```bash
./cfdata-linux-amd64 -cli -offurl auto
```

```bash
./cfdata-linux-amd64 -cli -offurl speed.cloudflare.com/__down?bytes=99999999
```

测速过程只读取响应数据流用于计算速度，不会将测速文件作为普通文件保存到本地。

---

# CLI

默认 CLI：

```bash
./cfdata-linux-amd64 -cli
```

官方模式示例：

```bash
./cfdata-linux-amd64 -cli \
  -mode official \
  -offiptype 4 \
  -offport 443 \
  -offurl auto
```

非标模式示例：

```bash
./cfdata-linux-amd64 -cli \
  -mode nsb \
  -nsbfile ip.txt \
  -nsbtls=true \
  -nsbspeedtest 5 \
  -offurl auto
```

查看完整参数：

```bash
./cfdata-linux-amd64 -h
```

常用参数包括：

```text
-cli
-mode
-scanmode
-offthreads
-nsbthreads
-offport
-offdelay
-nsbdelay
-offurl
-nsburl
-dns
-debug
-offout
-nsbout
-nsbfile
-nsbsourceurl
-nsbfallbackport
-nsbtls
-nsbspeedtest
-nsbresultlimit
-nsbspeedmin
-nsbspeedlimit
```

---

# 项目结构

核心结构：

```text
.
├── .github/
│   └── workflows/
│       └── docker-ghcr.yml
├── combined_refactor/
│   ├── main.go
│   ├── server.go
│   ├── country_ip.go
│   ├── config_path.go
│   ├── index.html
│   ├── login.html
│   ├── types.go
│   ├── ca.go
│   ├── go.mod
│   └── ...
├── Dockerfile
├── docker-compose.yml
├── PLUS_CHANGELOG.md
├── THIRD_PARTY_NOTICES.md
├── LICENSE
└── README.md
```

Docker 正式构建入口为：

```text
combined_refactor/
```

---

# 第三方项目与数据源

本项目基于：

[PoemMisty/CFData-WEB](https://github.com/PoemMisty/CFData-WEB)

Cloudflare 地区 IP 功能参考并集成：

[alienwaregf/Cloudflare-Country-Specific-IP-Filter](https://github.com/alienwaregf/Cloudflare-Country-Specific-IP-Filter)

地区 IP 数据源：

`https://zip.cm.edu.kg/all.txt`

本项目没有直接运行上游 Cloudflare Worker，而是由 CFData-WEB 后端获取数据并转换为现有非标优选模块能够使用的输入格式。

具体第三方说明请参阅：

```text
THIRD_PARTY_NOTICES.md
```

请在使用相关第三方代码、数据和服务前自行核对其当前许可证、项目声明和使用条件。

---

# 隐私与安全建议

本项目可以部署在局域网、NAS 或公网服务器上。

如果通过公网访问，建议：

- 使用 HTTPS
- 使用反向代理
- 限制访问来源
- 不要公开暴露不必要的管理接口
- 定期备份配置
- 不要在公开仓库提交敏感 Token、密码或 GitHub 凭据

如果启用 GitHub 上传功能，请妥善保护相关 Token。

---

# 常见问题

## 1. Docker 报 `container name ... is already in use`

说明已经存在同名容器，例如：

```text
cfdata-web
```

先查看：

```bash
docker ps -a --filter "name=cfdata-web"
```

确认旧容器不再需要后：

```bash
docker stop cfdata-web
docker rm cfdata-web
```

然后重新：

```bash
docker compose up -d
```

删除容器不会自动删除宿主机 `./data` 目录。

## 2. 页面一直显示“正在连接服务器”

检查：

```bash
docker ps
```

以及：

```bash
docker logs --tail 200 cfdata-web
```

同时确认浏览器访问的端口与 Docker 映射一致。

## 3. 更新镜像后配置还在吗？

只要使用：

```yaml
volumes:
  - ./data:/data
```

并且没有删除宿主机 `data/` 目录，配置文件会保留。

## 4. 地区 IP 获取失败

该功能依赖外部数据源：

```text
https://zip.cm.edu.kg/all.txt
```

如果数据源暂时不可访问，地区列表和 IP 获取会失败。此时可以直接使用原有的本地文件 / URL 非标输入功能。

---

# 更新记录

Plus 的具体修改记录请查看：

```text
PLUS_CHANGELOG.md
```

当前 Plus 版本重点包括：

- Docker / GHCR
- amd64 / arm64 多架构镜像
- `/data` 持久化
- Web 配置服务器保存与自动同步
- Cloudflare 地区 IP
- 地区自动加载
- 地区 IP 数量自动设置
- 地区 IP 一键导入非标优选测速
- 非标优选界面布局优化

---

# 免责声明

本项目仅用于合法的网络测试、学习、研究和运维场景。

使用者应自行确保其扫描、测速、数据获取、网络访问以及 GitHub 上传等行为符合所在地区法律法规、目标网络所有者的授权要求以及相关第三方服务条款。

对于因使用本项目产生的任何直接或间接损失、服务中断、数据问题或法律责任，由使用者自行承担。

---

# License

本项目的原始 CFData-WEB 部分遵循仓库中现有的 `LICENSE` 文件及其适用条款。

Plus 增强部分涉及第三方项目和数据源，请同时遵守 `THIRD_PARTY_NOTICES.md` 中列出的相关条件。

在重新发布、二次开发或商业使用前，请务必阅读并核对仓库内的完整许可证文本以及第三方项目当前许可证和声明。
