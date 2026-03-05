package checkparameter

import (
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
)

// 检查uint是否为0
//
// 如果是0就返回位置id和true
func CheckUintsIsZero[T uint | uint64](c *fiber.Ctx, u ...T) error {
	for i, v := range u {
		if v == 0 {
			return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, i+1)
		}
	}
	return nil
}

// 检查int是否为0
//
// 如果是0就返回位置id和true
func CheckIntsIsZero(c *fiber.Ctx, u ...int) error {
	for i, v := range u {
		if v == 0 {
			return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, i+1)
		}
	}
	return nil
}

// 检查string是否为空
//
// 如果是0就返回位置id和true
func CheckStringsIsNil(c *fiber.Ctx, u ...string) error {
	for i, v := range u {
		if v == "" {
			return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, i+1)
		}
	}
	return nil
}
