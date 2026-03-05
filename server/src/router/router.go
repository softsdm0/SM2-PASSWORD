package router

import (
	"fmt"

	"gitee.com/ouhaoqiang/passwordserver/server/src/handler"
	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/passwordHandler"
	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/passwordHandler/recordHandler"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login/casdoor"
	"gitee.com/ouhaoqiang/passwordserver/server/src/session"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	// 公共接口
	publicRoutes(app)
	// 私有接口
	privateRoutes(app)
}

// 代办
func todo(c *fiber.Ctx) error {
	msg := fmt.Sprintf("%s, 功能正在开发中...", c.Route().Name)
	return fiberresp.RespSuccess(c, msg)
}

// 公有接口
//
// 可以提供给公网使用的接口
func publicRoutes(app *fiber.App) {
	public := app.Group("/public", login.CheckLoginAndUserInfo).Name("public.")

	api := public.Group("/api").Name("api.")
	api.Get("", handler.Hello)

	// 用户
	user := api.Group("/userinfo").Name("user.")
	// 获取自己的用户信息
	user.Get("", handler.GetUserInfoHandler).Name("获取自己的用户信息(需要登录)")

	// 密码
	password := api.Group("/password", session.Sm4FiberHandler).Name("password.")
	// 获取密码记录
	record := password.Group("/record").Name("record.")
	// 查看获取密码记录
	record.Get("/", recordHandler.GetPasswordRecord).Name("获取密码记录")
	// 创建密码记录
	password.Post("", passwordHandler.CreatePassword).Name("创建密码记录(需要登录)")
	// 查看当前用户所有密码记录名称
	password.Get("", passwordHandler.GetSelfAllPasswordInfo).Name("查看当前用户所有密码记录信息, 不包含密码(需要登录)")
	// 通过密码记录id查看密码记录信息
	password.Get("/:id", passwordHandler.IdToGetPassword).Name("通过密码记录id获取密码信息(需要登录)")
	// 删除密码记录
	password.Delete("/:id", passwordHandler.DeletePassword).Name("删除密码记录(需要登录)")
	// 修改密码记录
	password.Post("/:id", passwordHandler.ChangePassword).Name("修改密码记录(需要登录)")

	// casdoor验证
	api.Post("/signin", casdoor.SigninHandler).Name("casdoor登录接口")

	// session
	sessionApi := api.Group("/session").Name("session.")
	sessionApi.Post("/id", session.GetSessionIdHanbler).Name("获取session_id接口")
	sessionApi.Get("/sm2PublicKey", session.GetSm2PublicKeyHandler).Name("获取sm2公钥接口")
	sessionApi.Get("/bootstrapConfig", session.GetBootstrapConfigHandler).Name("获取前端引导配置接口")
}

// 私有接口
//
// 不对外的接口，但是可以用于服务之间使用的接口。
func privateRoutes(app *fiber.App) {
	app.Group("/private", logger.New()).Name("private.")
}
