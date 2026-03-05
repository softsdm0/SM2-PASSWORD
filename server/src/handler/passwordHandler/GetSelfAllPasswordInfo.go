package passwordHandler

import (
	"fmt"
	"time"

	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/checkparameter"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login"
	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/checkpassword"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
)

// 获取自己的所有账本
func GetSelfAllPasswordInfo(c *fiber.Ctx) error {
	// 登录状态
	userInfo := login.GetUserInfo(c)

	// 响应
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
		// 密码强度
		PasswordStrength checkpassword.PasswordStrength
		// Url
		Url string
		// 备注
		Notes string
		// 创建时间
		CreatedAt time.Time
		// 更新信息时间
		UpdatedAt time.Time
	}

	if err := checkparameter.CheckStringsIsNil(c, userInfo.Id); err != nil {
		return err
	}

	var resp = []*Resp{}

	db := database.DB

	// 查询用户的所有密码记录
	res := db.Model(&model.TPasswordInfo{}).
		Where(&model.TPasswordInfo{UserId: userInfo.Id}).Order("app_name, account, account_type, created_at").
		Find(&resp)
	if res.Error != nil {
		return fiberresp.RespError(c, fiberresp.SELECT_ERROR, res.Error.Error())
	}

	// 账户信息解码处理
	for _, v := range resp {
		accountDecode, err := Decode(v.Account)
		if err != nil {
			fmt.Printf("id: %v get account decode error. %s\n", v.ID, err.Error())
			// 这里解析失败只报错就算了，不直接抛出错误给前端
			// 这里解释失败了代表数据库有藏数据了
			// 如果v.Account为空，就让前端提示用户更新account
		}
		v.Account = accountDecode
	}

	return fiberresp.RespSuccess(c, resp)
}
