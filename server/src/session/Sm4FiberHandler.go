package session

import (
	"context"
	"encoding/hex"
	"fmt"
	"strings"

	"gitee.com/ouhaoqiang/passwordserver/server/utils/cache"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tjfoc/gmsm/sm4"
)

func Sm4FiberHandler(c *fiber.Ctx) error {
	ctx := context.Background()
	sessionId := string(c.Request().Header.Peek("Password-Session-Id"))
	if sessionId == "" {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, "session id为空")
	}
	sm4Key, err := cache.Rdb.Get(ctx, cache.AddPrefix(cachePrefixName, sessionId)).Result()
	if err == redis.Nil {
		return fiberresp.RespError(c, fiberresp.SESSION_ID_EXPIRE, err.Error())
	} else if err != nil {
		return fiberresp.RespError(c, fiberresp.INTERNAL_SERVER_ERROR, err.Error())
	}

	sm4KeyByte, err := hex.DecodeString(sm4Key)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.INTERNAL_SERVER_ERROR, err.Error())
	}

	// 有请求body
	if len(c.Body()) != 0 {
		hexToByte, err := hex.DecodeString(strings.ReplaceAll(string(c.Body()), "\"", ""))
		if err != nil {
			return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
		}

		// 使用sm4解密内容
		sm4Msg, err := sm4.Sm4Ecb(sm4KeyByte, hexToByte, false)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
		}
		// 解密后重新写入body
		c.Request().SetBody(sm4Msg)
	}

	// 下一步
	err = c.Next()

	// 有响应body
	if len(c.Response().Body()) != 0 {
		// 使用sm4加密内容
		sm4respMsg, err := sm4.Sm4Ecb(sm4KeyByte, c.Response().Body(), true)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
		}
		sm4respMsgHex := hex.EncodeToString(sm4respMsg)
		// 解密后重新写入body
		c.Response().SetBody([]byte(sm4respMsgHex))
	}

	return err
}
