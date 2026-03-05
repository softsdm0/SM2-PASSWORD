package login

import (
	"github.com/gofiber/fiber/v2"
)

// 路由白名单
type routerWhiteList []routerInfo

// 路由信息
type routerInfo struct {
	// 路由名称
	Name string
	// 路由方法
	Method string
	// 路由路径
	Path string
}

var (
	// TODO: 临时的路由白名单
	rWhiteList routerWhiteList = routerWhiteList{
		{
			Name:   "public.api.signin",
			Method: fiber.MethodPost,
			Path:   "/public/api/signin",
		},
		{
			Name:   "public.api.session.id",
			Method: fiber.MethodPost,
			Path:   "/public/api/session/id",
		},
		{
			Name:   "public.api.session.sm2PublicKey",
			Method: fiber.MethodGet,
			Path:   "/public/api/session/sm2PublicKey",
		},
		{
			Name:   "public.api.session.bootstrapConfig",
			Method: fiber.MethodGet,
			Path:   "/public/api/session/bootstrapConfig",
		},
	}
)

// 鉴权路由白名单
func filter(c *fiber.Ctx) bool {
	for _, v := range rWhiteList {
		if v.Path == c.Path() && v.Method == c.Method() {
			return true
		}
	}
	return false
}
