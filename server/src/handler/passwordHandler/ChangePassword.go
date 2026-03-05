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
	"gorm.io/gorm"
)

func ChangePassword(c *fiber.Ctx) error {
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
		// Url
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
		// Url
		Url string
		// 备注
		Notes string
		// 创建时间
		CreatedAt time.Time
		// 更新时间
		UpdatedAt time.Time
	}

	//解析路径参数
	id := c.Params("id")

	req := new(Req)
	// 映射参数
	err := c.BodyParser(req)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, "2")
	}

	// 检查参数
	if err := checkparameter.CheckStringsIsNil(c, userInfo.Id, req.AppName, req.AccountType, req.Account, req.Password); err != nil {
		return err
	}
	if err := checkparameter.CheckStringsIsNil(c, id); err != nil {
		return err
	}

	db := database.DB

	tpi := new(model.TPasswordInfo)
	// 检查记账本是否已经存在
	res := db.Model(&model.TPasswordInfo{}).Where(&model.TPasswordInfo{
		ID:     id,
		UserId: userInfo.Id, // 检查一下是不是这个用户的记录
	}).First(tpi)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return fiberresp.RespError(c, fiberresp.RECORD_NOT_EXISTS)
		}
		return fiberresp.RespError(c, fiberresp.SELECT_ERROR, res.Error.Error())
	}

	// 密码信息编码
	accountEncoded, passwordEncoded, err := EncodePasswordInfo(c, req.Account, req.Password)
	if err != nil {
		return err
	}

	// 更新数据
	tpi.AppName = req.AppName
	tpi.AccountType = req.AccountType
	tpi.Account = accountEncoded
	tpi.Password = passwordEncoded
	tpi.Url = req.Url
	tpi.Notes = req.Notes
	tpi.PasswordStrength = checkpassword.CheckPasswordStrength(req.Password) // 检测密码强度

	// 更新密码记录
	res = db.Model(&model.TPasswordInfo{}).
		Where(&model.TPasswordInfo{
			ID:     id,
			UserId: userInfo.Id, // 检查一下是不是这个用户的记录
		}).Updates(tpi)
	if res.Error != nil {
		return fiberresp.RespError(c, fiberresp.CREATE_ERROR, res.Error.Error())
	}
	if res.RowsAffected <= 0 {
		// 如果没有更新
		return fiberresp.RespError(c, fiberresp.RECORD_NOT_EXISTS, res.RowsAffected)
	}

	// 账户信息解码处理
	accountDecoded, passwordDecoded, err := DecodePasswordInfo(c, tpi.Account, tpi.Password)
	if err != nil {
		return err
	}
	tpi.Account = accountDecoded
	tpi.Password = passwordDecoded

	resp := &Resp{
		ID:          tpi.ID,
		UserId:      tpi.UserId,
		AppName:     tpi.AppName,
		AccountType: tpi.AccountType,
		Account:     tpi.Account,
		Url:         tpi.Url,
		Notes:       tpi.Notes,
		CreatedAt:   tpi.CreatedAt,
		UpdatedAt:   tpi.UpdatedAt,
	}
	return fiberresp.RespSuccess(c, resp)
}
