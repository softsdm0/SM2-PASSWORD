package login

import (
	"gitee.com/ouhaoqiang/passwordserver/server/src/login/casdoor"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gofiber/fiber/v2"
)

// 检查登录，并且获取用户信息
func CheckLoginAndUserInfo(c *fiber.Ctx) error {
	// 检查是否白名单
	if filter(c) {
		return c.Next()
	}

	// 检测登录类型
	loginType, err := checkLoginTokenType(c)
	if err != nil {
		return err
	}

	// 检测登录状态
	switch loginType {
	case LOGIN_TYPE_CASDOOR:
		// casdoor登录
		if err := casdoor.CheckUserInfo(c, UserKeyName); err != nil {
			// 登录失败
			return err
		}
	default:
		return fiberresp.RespError(c, fiberresp.MISSING_OR_MALFORMED_JWT, "未登录")
	}

	return c.Next()
}

type LoginType string

// 登录类型
const (
	// casdoor
	LOGIN_TYPE_CASDOOR LoginType = "casdoor"
	// AuthServer, 废弃不用authserver进行登录了
	LOGIN_TYPE_AUTHSERVER LoginType = "authserver"
)

// 检测登录Token类型
func checkLoginTokenType(c *fiber.Ctx) (LoginType, error) {
	// 检测是否有casdoor的cookie key name
	token := c.Cookies(casdoor.CASDOOR_JWT_COOKIE_NAME)
	if token != "" {
		return LOGIN_TYPE_CASDOOR, nil
	}

	return "", fiberresp.RespError(c, fiberresp.MISSING_OR_MALFORMED_JWT, "未登录")
}

// user信息存放在Ctx中的key名称
var UserKeyName = "LoginUserKey"

// 从fiber上下文中获取用户信息
func GetUserInfo(c *fiber.Ctx) *casdoorsdk.User {
	user := c.Locals(UserKeyName).(*casdoorsdk.User)
	return user
}
