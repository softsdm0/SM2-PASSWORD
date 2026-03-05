package recordHandler

import (
	"fmt"

	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/checkparameter"
	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/passwordHandler"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login"
	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
)

func GetPasswordRecord(c *fiber.Ctx) error {
	user := login.GetUserInfo(c)

	// 检查参数
	if err := checkparameter.CheckStringsIsNil(c, user.Id); err != nil {
		return err
	}

	tgpr := []*model.TGetPasswordRecords{}

	db := database.DB.Debug()

	sqlResp := db.Model(&model.TGetPasswordRecords{}).Where(&model.TGetPasswordRecords{UserId: user.Id}).Preload("PasswordInfo").Find(&tgpr)
	if sqlResp.Error != nil {
		return fiberresp.RespError(c, fiberresp.SELECT_ERROR, "查询记录失败")
	}

	// 清空一下密码
	for _, v := range tgpr {
		v.PasswordInfo.Password = ""
		accountDecode, err := passwordHandler.Decode(v.PasswordInfo.Account)
		if err != nil {
			fmt.Printf("id: %v get account decode error. %s\n", v.ID, err.Error())
			// 这里解析失败只报错就算了，不直接抛出错误给前端
			// 这里解释失败了代表数据库有藏数据了
			// 如果v.Account为空，就让前端提示用户更新account
		}
		v.PasswordInfo.Account = accountDecode
	}

	return fiberresp.RespSuccess(c, tgpr)
}
