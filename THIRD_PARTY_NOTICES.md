# Third-Party Notices

## Cloudflare-Country-Specific-IP-Filter

本项目的“Cloudflare 地区 IP”功能参考并集成了以下开源项目的地区 IP 筛选思路及地区编码映射：

- 项目：https://github.com/alienwaregf/Cloudflare-Country-Specific-IP-Filter
- 数据源：https://zip.cm.edu.kg/all.txt

上游项目 README 声明支持二次创作，并要求保留鸣谢链接，同时声明禁止商用。使用本集成功能时请遵守上游项目的声明、许可条件及当地法律法规。

本集成功能未直接运行 Cloudflare Worker；数据由 CFData-WEB 后端获取并转换后交给现有非标优选测速模块处理。
