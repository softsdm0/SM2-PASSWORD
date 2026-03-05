package session

import (
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
)

func GetSm2PublicKeyHandler(c *fiber.Ctx) error {
	return fiberresp.RespSuccess(c, Sm2PublicKeyStr)
}
