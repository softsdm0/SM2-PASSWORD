package handler

import (
	"gitee.com/ouhaoqiang/passwordserver/server/src/login"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
)

func GetUserInfoHandler(c *fiber.Ctx) error {
	uInfo := login.GetUserInfo(c)
	return fiberresp.RespSuccess(c, uInfo)
}
