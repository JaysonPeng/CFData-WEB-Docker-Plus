package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const countryIPSourceURL = "https://zip.cm.edu.kg/all.txt"

// countryIPRegionNames follows the region-code convention used by
// Cloudflare-Country-Specific-IP-Filter. Unknown codes are still returned as
// their two-letter code so a new region never breaks the UI.
var countryIPRegionNames = map[string]string{
	"JP": "日本", "KR": "韩国", "SG": "新加坡", "HK": "香港", "TW": "台湾",
	"MY": "马来西亚", "TH": "泰国", "VN": "越南", "PH": "菲律宾", "ID": "印尼",
	"IN": "印度", "AU": "澳大利亚", "NZ": "新西兰", "KH": "柬埔寨", "MO": "澳门",
	"US": "美国", "CA": "加拿大", "MX": "墨西哥", "GB": "英国", "UK": "英国",
	"DE": "德国", "FR": "法国", "NL": "荷兰", "IT": "意大利", "ES": "西班牙",
	"PT": "葡萄牙", "RU": "俄罗斯", "UA": "乌克兰", "PL": "波兰", "SE": "瑞典",
	"FI": "芬兰", "NO": "挪威", "DK": "丹麦", "IS": "冰岛", "IE": "爱尔兰",
	"BE": "比利时", "CH": "瑞士", "AT": "奥地利", "CZ": "捷克", "HU": "匈牙利",
	"RO": "罗马尼亚", "BG": "保加利亚", "GR": "希腊", "TR": "土耳其", "BR": "巴西",
	"AR": "阿根廷", "CL": "智利", "CO": "哥伦比亚", "PE": "秘鲁", "ZA": "南非",
	"AE": "阿联酋", "SA": "沙特", "IL": "以色列", "KZ": "哈萨克斯坦",
}

type countryIPRegion struct {
	Code  string `json:"code"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type countryIPFetchRequest struct {
	Regions       []string `json:"regions"`
	Limit         int      `json:"limit"`
	Port          int      `json:"port"`
	UseSourcePort bool     `json:"useSourcePort"`
	Format        string   `json:"format"`
	Carrier       string   `json:"carrier"`
}

type countryIPFetchResult struct {
	Success       bool   `json:"success"`
	Message       string `json:"message,omitempty"`
	Count         int    `json:"count"`
	SourceURL     string `json:"sourceURL"`
	FormattedText string `json:"formattedText"`
	ScanText      string `json:"scanText"`
}

type countryIPItem struct {
	host       string
	port       int
	code       string
	sourceLine string
}

func fetchCountryIPSource(ctx context.Context) ([]countryIPItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, countryIPSourceURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "CFData-WEB/1.0")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("获取地区 IP 数据失败: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("地区 IP 数据源返回 HTTP %d", resp.StatusCode)
	}

	// The source is a plain-text list. Keep a hard limit so a malformed or
	// unexpectedly large upstream response cannot consume unbounded memory.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 10<<20))
	if err != nil {
		return nil, fmt.Errorf("读取地区 IP 数据失败: %w", err)
	}

	items := make([]countryIPItem, 0, 4096)
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimPrefix(string(body), "\uFEFF")))
	scanner.Buffer(make([]byte, 4096), 64*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.Contains(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "#", 2)
		address := strings.TrimSpace(parts[0])
		code := strings.ToUpper(strings.TrimSpace(parts[1]))
		if address == "" || code == "" {
			continue
		}
		host, port, ok := splitHostPortLoose(address)
		if !ok || net.ParseIP(host) == nil || port <= 0 || port > 65535 {
			continue
		}
		items = append(items, countryIPItem{host: host, port: port, code: code, sourceLine: line})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("解析地区 IP 数据失败: %w", err)
	}
	return items, nil
}

func splitHostPortLoose(address string) (string, int, bool) {
	address = strings.TrimSpace(address)
	if strings.HasPrefix(address, "[") {
		host, portText, err := net.SplitHostPort(address)
		if err != nil {
			return "", 0, false
		}
		port, err := strconv.Atoi(portText)
		return host, port, err == nil
	}
	if host, portText, err := net.SplitHostPort(address); err == nil {
		port, err := strconv.Atoi(portText)
		return host, port, err == nil
	}
	idx := strings.LastIndex(address, ":")
	if idx <= 0 || idx == len(address)-1 {
		return "", 0, false
	}
	host := address[:idx]
	port, err := strconv.Atoi(address[idx+1:])
	return host, port, err == nil
}

func formatCountryIPHostPort(host string, port int) string {
	if strings.Contains(host, ":") {
		return fmt.Sprintf("[%s]:%d", host, port)
	}
	return fmt.Sprintf("%s:%d", host, port)
}

func countryIPRegionStats(items []countryIPItem) []countryIPRegion {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.code]++
	}
	regions := make([]countryIPRegion, 0, len(counts))
	for code, count := range counts {
		name := countryIPRegionNames[code]
		if name == "" {
			name = code
		}
		regions = append(regions, countryIPRegion{Code: code, Name: name, Count: count})
	}
	sortCountryIPRegions(regions)
	return regions
}

func sortCountryIPRegions(regions []countryIPRegion) {
	for i := 0; i < len(regions); i++ {
		for j := i + 1; j < len(regions); j++ {
			if regions[j].Count > regions[i].Count || (regions[j].Count == regions[i].Count && regions[j].Code < regions[i].Code) {
				regions[i], regions[j] = regions[j], regions[i]
			}
		}
	}
}

func countryIPFlag(code string) string {
	if code == "TW" {
		return "🇹🇼"
	}
	if code == "UK" {
		return "🇬🇧"
	}
	if len(code) != 2 {
		return "🌐"
	}
	r1 := rune(code[0])
	r2 := rune(code[1])
	return string([]rune{127397 + r1, 127397 + r2})
}

func buildCountryIPFetchResult(req countryIPFetchRequest, all []countryIPItem) countryIPFetchResult {
	selectedRegions := make(map[string]bool)
	for _, region := range req.Regions {
		region = strings.ToUpper(strings.TrimSpace(region))
		if region != "" {
			selectedRegions[region] = true
		}
	}

	pools := make(map[string][]countryIPItem)
	for _, item := range all {
		if selectedRegions[item.code] {
			pools[item.code] = append(pools[item.code], item)
		}
	}

	limit := req.Limit
	if limit < 0 {
		limit = 0
	}
	if limit > 5000 {
		limit = 5000
	}
	selected := make([]countryIPItem, 0)
	for _, region := range req.Regions {
		region = strings.ToUpper(strings.TrimSpace(region))
		pool := pools[region]
		if len(pool) == 0 {
			continue
		}
		if limit > 0 && len(pool) > limit {
			// 与上游 Worker 的行为保持一致：限制数量时随机抽取，避免每次都固定取数据源前 N 条。
			perm := rand.Perm(len(pool))
			for _, index := range perm[:limit] {
				selected = append(selected, pool[index])
			}
			continue
		}
		selected = append(selected, pool...)
	}

	format := strings.ToLower(strings.TrimSpace(req.Format))
	if format == "" {
		format = "ipport_region"
	}
	carrier := strings.TrimSpace(req.Carrier)
	formatted := make([]string, 0, len(selected))
	scan := make([]string, 0, len(selected))
	for _, item := range selected {
		port := item.port
		if !req.UseSourcePort && req.Port > 0 {
			port = req.Port
		}
		ipPort := formatCountryIPHostPort(item.host, port)
		name := countryIPRegionNames[item.code]
		if name == "" {
			name = item.code
		}
		flag := countryIPFlag(item.code)
		label := fmt.Sprintf("%s %s", flag, name)
		switch format {
		case "ip":
			formatted = append(formatted, item.host)
		case "ipport":
			formatted = append(formatted, ipPort)
		case "ipport_carrier_region":
			if carrier == "" {
				carrier = "未指定"
			}
			formatted = append(formatted, fmt.Sprintf("%s#%s+%s", ipPort, carrier, name))
		default:
			formatted = append(formatted, fmt.Sprintf("%s#%s", ipPort, label))
		}
		scan = append(scan, fmt.Sprintf("%s %d", item.host, port))
	}

	return countryIPFetchResult{
		Success:       true,
		Count:         len(selected),
		SourceURL:     countryIPSourceURL,
		FormattedText: strings.Join(formatted, "\n"),
		ScanText:      strings.Join(scan, "\n"),
	}
}
