package session

import (
	"gitee.com/ouhaoqiang/passwordserver/server/utils/config"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
)

func GetBootstrapConfigHandler(c *fiber.Ctx) error {
	casdoorCfg := config.Config.Server.Casdoor
	return fiberresp.RespSuccess(c, fiber.Map{
		"sm2PublicKey":        Sm2PublicKeyStr,
		"casdoorClientId":     casdoorCfg.ClientID,
		"casdoorOrganization": casdoorCfg.Organization,
		"casdoorApplication":  casdoorCfg.Application,
	})
}
