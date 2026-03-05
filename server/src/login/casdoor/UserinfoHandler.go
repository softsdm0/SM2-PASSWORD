package casdoor

import (
	"context"
	"strings"

	"gitee.com/ouhaoqiang/passwordserver/server/utils/cache"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

// 获取用户信息
func UserinfoHandler(c *fiber.Ctx) error {
	// 解析token
	claims, err := ParseJwtToken(c)
	if err != nil {
		return err
	}

	return c.JSON(map[string]interface{}{
		"status": "ok",
		"data":   claims.User,
	})
}

// 检测是否有登录，并且解析token
func ParseJwtToken(c *fiber.Ctx) (*casdoorsdk.Claims, error) {
	// 检查是否有登录
	casdoorAuthToken, err := CheckLoginToken(c)
	if err != nil {
		return nil, err
	}

	// 在cache里面获取casdoorAuthToken对应的casdoor token
	token, err := cache.Rdb.Get(context.Background(), cache.AddPrefix(cachePrefixName, casdoorAuthToken)).Result()
	if err == redis.Nil {
		return nil, fiberresp.RespError(c, fiberresp.INVALID_TOKEN_ID, "token not fund")
	} else if err != nil {
		return nil, fiberresp.RespError(c, fiberresp.INTERNAL_SERVER_ERROR, err.Error())
	}

	// 解析token
	claims, err := casdoorsdk.ParseJwtToken(token)
	if err != nil {
		return nil, fiberresp.RespError(c, fiberresp.MISSING_OR_MALFORMED_JWT, "2")
	}

	return claims, nil
}

// 检查是否有登录
func CheckLoginToken(c *fiber.Ctx) (string, error) {
	var token string
	token, ok := CheckCookie(c)
	if ok {
		return token, nil
	}

	token, ok = CheckAuthorizationBearer(c)
	if ok {
		return token, nil
	}

	return "", fiberresp.RespError(c, fiberresp.MISSING_OR_MALFORMED_JWT, "1")
}

// 检查cookie里面是否有CASDOOR_JWT_COOKIE_NAME
func CheckCookie(c *fiber.Ctx) (string, bool) {
	token := c.Cookies(CASDOOR_JWT_COOKIE_NAME)
	if token == "" {
		return "", false
	}
	return token, true
}

// 检查请求头里面是否有Authorization Bearer
func CheckAuthorizationBearer(c *fiber.Ctx) (string, bool) {
	authHeader := c.Request().Header.Peek("Authorization")
	if string(authHeader) == "" {
		return "", false
	}

	token := strings.Split(string(authHeader), "Bearer ")
	if len(token) != 2 {
		return "", false
	}

	return token[1], true
}
