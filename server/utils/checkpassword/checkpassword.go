package checkpassword

import (
	"regexp"
)

// PasswordStrength 定义密码强度类型
type PasswordStrength int

// 定义密码强度等级常量
const (
	// 非常弱：长度小于 6
	VeryWeak PasswordStrength = iota + 1
	// 弱：长度小于 8 或字符类型少于 2 种
	Weak
	// 中等：长度小于 10 或字符类型少于 3 种
	Medium
	// 强：长度大于等于 10 且字符类型至少 3 种
	Strong
)

// CheckPasswordStrength 检查密码强度，返回 1 - 4 的等级
func CheckPasswordStrength(password string) PasswordStrength {
	length := len(password)
	hasNumber := regexp.MustCompile(`\d`).MatchString(password)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*(),.?":{}|<>]`).MatchString(password)

	// 计算密码包含的字符类型数量
	charTypes := 0
	if hasNumber {
		charTypes++
	}
	if hasLower {
		charTypes++
	}
	if hasUpper {
		charTypes++
	}
	if hasSpecial {
		charTypes++
	}

	// 根据长度和字符类型数量确定密码强度等级
	if length < 6 {
		return VeryWeak // 非常弱：长度小于 6
	} else if length < 8 || charTypes < 2 {
		return Weak // 弱：长度小于 8 或字符类型少于 2 种
	} else if length < 10 || charTypes < 3 {
		return Medium // 中等：长度小于 10 或字符类型少于 3 种
	}
	return Strong // 强：长度大于等于 10 且字符类型至少 3 种
}
