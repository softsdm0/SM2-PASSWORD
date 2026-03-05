package casdoor

import (
	"context"
	"fmt"
	"time"

	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/checkparameter"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/cache"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const (
	CASDOOR_JWT_COOKIE_NAME = "casdoorAuthToken"
	cachePrefixName         = "casdoorAuthToken"
)

func SigninHandler(c *fiber.Ctx) error {
	type Req struct {
		Code  string `json:"code"`
		State string `json:"state"`
	}

	req := new(Req)
	err := c.QueryParser(req)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
	}

	// 检查参数
	if err := checkparameter.CheckStringsIsNil(c, req.Code, req.State); err != nil {
		return err
	}
	fmt.Println(req.Code, req.State)

	token, err := casdoorsdk.GetOAuthToken(req.Code, req.State)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.JWT_SIGN_ERROR, err.Error())
	}

	// 由于这个token.AccessToken太大了，会导致超出服务接受请求头大小，所以要转换一下存在cache里面
	casdoorAuthToken := uuid.New().String()
	// 写入到缓存
	ctx := context.Background()
	cache.Rdb.Set(ctx, cache.AddPrefix(cachePrefixName, casdoorAuthToken), token.AccessToken, time.Until(token.Expiry))

	c.Cookie(&fiber.Cookie{
		Name:    CASDOOR_JWT_COOKIE_NAME,
		Value:   casdoorAuthToken,
		Expires: token.Expiry,
	})

	return fiberresp.RespSuccess(c, token.AccessToken)
}
