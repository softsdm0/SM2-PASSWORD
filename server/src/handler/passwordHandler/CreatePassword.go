package passwordHandler

import (
	"time"

	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/checkparameter"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login"
	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/checkpassword"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func CreatePassword(c *fiber.Ctx) error {
	// 登录状态
	userInfo := login.GetUserInfo(c)
	type Req struct {
		// 应用名称
		AppName string
		// 账户类型,一个账户可以有多种类型
		AccountType string
		// 账户
		Account string
		// 密码
		Password string
		// url
		Url string
		// 备注
		Notes string
	}
	type Resp struct {
		// 密码记录id
		ID string
		// 创建者id
		UserId string
		// 应用名称
		AppName string
		// 账户类型,一个账户可以有多种类型
		AccountType string
		// 账户
		Account string
		// url
		Url string
		// 备注
		Notes string
		// 创建时间
		CreatedAt time.Time
	}
	req := new(Req)

	// 映射参数
	err := c.BodyParser(req)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT)
	}

	// 检查参数
	if err := checkparameter.CheckStringsIsNil(c, userInfo.Id, req.AppName, req.AccountType, req.Account, req.Password); err != nil {
		return err
	}

	// 账户和密码信息编码处理
	accountEncoded, passwordEncoded, err := EncodePasswordInfo(c, req.Account, req.Password)
	if err != nil {
		return err
	}

	db := database.DB

	// 检查记账本是否已经存在
	res := db.Model(&model.TPasswordInfo{}).Where(&model.TPasswordInfo{
		UserId:      userInfo.Id,
		AppName:     req.AppName,
		AccountType: req.AccountType,
		Account:     accountEncoded,
	}).Find(nil)
	if res.Error != nil {
		return fiberresp.RespError(c, fiberresp.SELECT_ERROR, res.Error.Error())
	}
	if res.RowsAffected > 0 {
		return fiberresp.RespError(c, fiberresp.RECORD_EXISTS)
	}

	tpi := &model.TPasswordInfo{
		ID:               uuid.NewString(),
		UserId:           userInfo.Id,
		AppName:          req.AppName,
		AccountType:      req.AccountType,
		Account:          accountEncoded,
		Password:         passwordEncoded,
		Url:              req.Url,
		Notes:            req.Notes,
		PasswordStrength: checkpassword.CheckPasswordStrength(req.Password), // 检测密码强度
	}
	// 创建密码记录
	res = db.Create(tpi)
	if res.Error != nil {
		return fiberresp.RespError(c, fiberresp.CREATE_ERROR, res.Error.Error())
	}

	resp := &Resp{
		ID:          tpi.ID,
		UserId:      tpi.UserId,
		AppName:     tpi.AppName,
		AccountType: tpi.AccountType,
		Account:     req.Account,
		Url:         tpi.Url,
		Notes:       tpi.Notes,
		CreatedAt:   tpi.CreatedAt,
	}
	return fiberresp.RespSuccess(c, resp)
}
