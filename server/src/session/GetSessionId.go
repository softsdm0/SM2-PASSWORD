package session

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"gitee.com/ouhaoqiang/passwordserver/server/utils/cache"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/config"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm4"
)

func GetSessionIdHanbler(c *fiber.Ctx) error {
	type Req struct {
		// sm4密钥
		Sm4Key string
		// 默认就是为 ok.
		Msg string
	}
	type Resp struct {
		SessionId string
	}

	// 卧槽找了我几天，原来nodejs传过来的c1没有带04开头，需要自己加一个，sm2才能解析，好奇怪。。。
	// 还是在看nodejs的解密过程时发现要加04的 `const c1 = _.getGlobalCurve().decodePointHex('04' + encryptData.substr(0, 128))``
	sm2BodyAdd04 := []byte("04")
	sm2BodyAdd04 = append(sm2BodyAdd04, c.Body()...)
	// fmt.Printf("len(c.Body()): %v\n", len(c.Body()))

	// 将16进制string转为byte切片
	data, err := hex.DecodeString(string(sm2BodyAdd04))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fiberresp.RespError(c, fiberresp.SM2_DECRYPT_ERROR, err.Error())
	}

	// sm2解密
	decrypt, err := sm2.Decrypt(Sm2PrivateKey, data, sm2.C1C3C2)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fiberresp.RespError(c, fiberresp.SM2_DECRYPT_ERROR, err.Error())
	}

	// 将解密后的内容转到结构体
	req := new(Req)
	err = json.Unmarshal(decrypt, req)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
	}

	// fmt.Printf("req.Sm4Key: %v\n", req.Sm4Key)
	// fmt.Printf("req.Msg: %v\n", req.Msg)

	if req.Msg == "" {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, "req.Msg is empty.")
	}

	// sm4key16进制字符串传byte切片
	sm4KeyByte, err := hex.DecodeString(req.Sm4Key)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
	}
	// fmt.Printf("sm4KeyByte: %v\n", sm4KeyByte)

	// 内容16进制字符串传byte切片
	msgByte, err := hex.DecodeString(req.Msg)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
	}

	// 使用sm4解密req.Msg内容
	sm4Msg, err := sm4.Sm4Ecb(sm4KeyByte, msgByte, false)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
	}
	// fmt.Printf("sm4Msg: %v\n", string(sm4Msg))

	// 判断解密是否成功
	if !bytes.Equal(sm4Msg, []byte("ok.")) {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, "msg is not ok.")
	}

	// 生成sessionId
	sessionId := uuid.New().String()
	// fmt.Printf("sessionId: %v\n", sessionId)

	// 写入到缓存
	ctx := context.Background()
	// 这里存的是sm4密钥16进制字符串，不是byte切片
	cache.Rdb.Set(ctx, cache.AddPrefix(cachePrefixName, sessionId), req.Sm4Key, time.Second*config.Config.Server.Redis.SessionIdExpiration)

	// 解析为json
	resp := &Resp{
		SessionId: sessionId,
	}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return fiberresp.RespError(c, fiberresp.INTERNAL_SERVER_ERROR, err.Error())
	}

	// sm4加密
	respSm4, err := sm4.Sm4Ecb(sm4KeyByte, jsonResp, true)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.SM4_ENCRYPT_ERROR, err.Error())
	}
	// 十六进制字符串格式
	respSm4Hex := hex.EncodeToString(respSm4)

	// 返回sm4加密内容
	return fiberresp.RespSuccess(c, respSm4Hex)
}
