package passwordHandler

import (
	"fmt"

	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/checkparameter"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login"
	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
)

func DeletePassword(c *fiber.Ctx) error {
	userInfo := login.GetUserInfo(c)
	type Req struct {
		// 密码记录id
		Id string
	}

	req := new(Req)

	//解析路径参数
	err := c.ParamsParser(req)
	if err != nil {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
	}

	// 参数检查
	if err := checkparameter.CheckStringsIsNil(c, userInfo.Id); err != nil {
		return err
	}
	if err := checkparameter.CheckStringsIsNil(c, req.Id); err != nil {
		return err
	}

	db := database.DB

	tpi := new(model.TPasswordInfo)
	// 检查记账本是否已经存在
	res := db.Model(&model.TPasswordInfo{}).Where(&model.TPasswordInfo{
		ID:     req.Id,
		UserId: userInfo.Id,
	}).Find(tpi)
	if res.Error != nil {
		return fiberresp.RespError(c, fiberresp.DELETE_ERROR, res.Error.Error())
	}
	if res.RowsAffected <= 0 {
		// 不存在返回
		return fiberresp.RespError(c, fiberresp.RECORD_NOT_EXISTS)
	}

	// 删除浏览记录
	if err := db.Unscoped().Where(&model.TGetPasswordRecords{PasswordId: req.Id, UserId: userInfo.Id}).Delete(&model.TGetPasswordRecords{}).Error; err != nil {
		return fiberresp.RespError(c, fiberresp.DELETE_ERROR, err.Error())
	}

	// 删除账本
	err = db.Unscoped().Delete(tpi).Error
	if err != nil {
		return fiberresp.RespError(c, fiberresp.DELETE_ERROR, err.Error())
	}

	return fiberresp.RespSuccess(c, fmt.Sprintf("%v delete ok.", req.Id))
}
