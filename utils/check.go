package utils

import (
	"regexp"
	"strings"
)

var (
	emailRegex = regexp.MustCompile(
		`^(?i)[a-z0-9!#$%&'*+/=?^_{|}~-]+` +
			`(?:\.[a-z0-9!#$%&'*+/=?^_{|}~-]+)*` +
			`@` +
			`(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+` +
			`[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$`,
	)
	// 限制域名的最大长度（RFC标准规定253字符）
	domainRegex = regexp.MustCompile(`^.{1,253}$`)
	// 限制本地部分最大长度（64字符）
	localRegex = regexp.MustCompile(`^.{1,64}$`)
)

// CheckEmail 执行多层级邮箱验证
func CheckEmail(email string) bool {
	// 去除首尾空格
	email = strings.TrimSpace(email)

	// 基础结构校验
	if len(email) < 6 || len(email) > 254 {
		return false
	}

	// 拆分本地部分和域名
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	local, domain := parts[0], parts[1]

	// 验证本地部分长度
	if !localRegex.MatchString(local) {
		return false
	}

	// 验证域名部分
	if !validateDomain(domain) {
		return false
	}

	// 正则表达式校验
	return emailRegex.MatchString(email)
}

// 辅助方法：验证域名部分
func validateDomain(domain string) bool {
	// 域名总长度校验
	if !domainRegex.MatchString(domain) {
		return false
	}

	// 检查点分隔的各个段
	segments := strings.Split(domain, ".")
	for _, seg := range segments {
		if len(seg) == 0 || len(seg) > 63 {
			return false
		}
		// 每段只能包含字母、数字和连字符
		if !regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(seg) {
			return false
		}
		// 连字符不能在开头或结尾
		if strings.HasPrefix(seg, "-") || strings.HasSuffix(seg, "-") {
			return false
		}
	}
	return true
}
