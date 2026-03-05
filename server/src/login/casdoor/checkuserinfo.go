package casdoor

import (
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gofiber/fiber/v2"
)

// 访问casdoor服务检查token是否过期，并且获取用户信息
func CheckUserInfo(c *fiber.Ctx, UserKeyName string) error {
	Claims, err := ParseJwtToken(c)
	if err != nil {
		return err
	}

	setUserInfo2Fiber(c, UserKeyName, &Claims.User)

	return nil
}

// 设置用户信息到fiber上下文
func setUserInfo2Fiber(c *fiber.Ctx, UserKeyName string, r *casdoorsdk.User) {
	// 设置到Ctx上下文中
	c.Locals(UserKeyName, r)
}
