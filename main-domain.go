package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"os/exec"

//	"github.com/likexian/whois"
)

// 获取域名到期时间
func getDomainExpiry(domain string) (time.Time, error) {
	// 查询域名的 WHOIS 信息
//	result, err := whois.Whois(domain)
        cmd := exec.Command("whois", domain)


	output, err := cmd.CombinedOutput() // 同时捕获标准输出和错误输出
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 将 WHOIS 查询结果按行分割
	lines := strings.Split(string(output), "\n")

	// 使用正则表达式查找包含 Registry Expiry Date 的行
	reg := regexp.MustCompile(`(?i)Registry Expiry Date:\s*(\S+)`)
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// 正则匹配 Registry Expiry Date 行
		matches := reg.FindStringSubmatch(line)
		if len(matches) > 1 {
			// 提取日期并解析
			expiryDate, err := time.Parse(time.RFC3339, matches[1])
			if err != nil {
				return time.Time{}, fmt.Errorf("解析日期失败: %v", err)
			}
			return expiryDate, nil
		}
	}

	return time.Time{}, fmt.Errorf("未找到域名到期时间")
}

// 处理文件中的每个域名
func processDomainsFromFile(filePath string) error {
	// 打开包含域名的文件
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("无法打开文件 %s: %v", filePath, err)
	}
	defer file.Close()

	// 创建输出文件（以当前时间戳为文件名）
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	outFile, err := os.Create(fmt.Sprintf("%s.txt", timestamp))
	if err != nil {
		return fmt.Errorf("无法创建输出文件: %v", err)
	}
	defer outFile.Close()

	// 逐行读取文件中的域名
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domain := strings.TrimSpace(scanner.Text())
		if domain == "" {
			continue
		}

		// 获取域名的到期时间
		expiryDate, err := getDomainExpiry(domain)
		if err != nil {
			log.Printf("查询域名 %s 失败: %v\n", domain, err)
			continue
		}

		log.Printf("域名 %s 的到期时间是: %s\n", domain, expiryDate.Format("2006-01-02 15:04:05"))

		// 将结果写入输出文件
		_, err = outFile.WriteString(fmt.Sprintf("域名 %s 的到期时间是: %s\n", domain, expiryDate.Format("2006-01-02 15:04:05")))
		if err != nil {
			log.Printf("写入文件失败: %v\n", err)
		}
	}

	// 检查扫描是否遇到错误
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件时发生错误: %v", err)
	}

	return nil
}

func main() {
	// 处理域名文件并输出到时间戳文件
	err := processDomainsFromFile("domains.txt")
	if err != nil {
		log.Fatalf("处理域名文件时出错: %v\n", err)
	} else {
		fmt.Println("所有域名的到期时间已保存到文件。")
	}
}

