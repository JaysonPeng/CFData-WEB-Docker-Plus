# CFData-WEB Plus

## 当前版本

基于 CFData-WEB Docker 版本，新增 Cloudflare 地区 IP 数据源与非标优选联动。

## 本次变更

- 新增 `combined_refactor/country_ip.go`
- 新增 Web 端「Cloudflare 地区 IP」面板
- 后端通过 WebSocket 提供地区列表与 IP 获取接口
- 支持地区多选、全选、清空
- 支持每地区数量限制；限制数量时随机抽取
- 支持使用数据源端口或统一指定端口
- 支持 IP、IP:PORT、IP:PORT#地区、IP:PORT#运营商+地区
- 支持 IPv4 / IPv6 格式化
- 支持一键复制结果
- 支持一键导入「非标优选」测速
- 保留原有 WebSocket、服务器配置保存、配置持久化与 Docker 功能
- GHCR 使用独立镜像名：`ghcr.io/jaysonpeng/cfdata-web-plus`

## 上游致谢

地区 IP 功能参考：
https://github.com/alienwaregf/Cloudflare-Country-Specific-IP-Filter

数据源：
https://zip.cm.edu.kg/all.txt

请遵守上游项目声明及当地法律法规。
