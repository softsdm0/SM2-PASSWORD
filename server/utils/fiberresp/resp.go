package fiberresp

import "github.com/gofiber/fiber/v2"

// 响应成功
func RespSuccess(c *fiber.Ctx, data interface{}) error {
	s, ok := sCodeMap[SUCCESS]
	if !ok {
		return c.JSON(fiber.Map{"status": SUCCESS, "message": "", "data": data})
	}
	return c.JSON(fiber.Map{"status": s.statusCode, "message": s.message, "data": data})
}

// 响应错误
func RespError(c *fiber.Ctx, sCode statusCode, data ...interface{}) error {
	var d interface{}
	if len(data) > 0 {
		d = data[0]
	}
	s, ok := sCodeMap[sCode]
	if !ok {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": sCode, "message": "未定义的错误码", "data": d})
		return fiber.NewError(sCode, "未定义的错误吗")
	}
	c.Status(s.httpCode).JSON(fiber.Map{"status": s.statusCode, "message": s.message, "data": d})
	return fiber.NewError(s.statusCode, s.message)
}
