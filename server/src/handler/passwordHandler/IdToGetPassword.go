package passwordHandler

import (
	"fmt"

	"gitee.com/ouhaoqiang/passwordserver/server/src/handler/checkparameter"
	"gitee.com/ouhaoqiang/passwordserver/server/src/login"
	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// 通过id获取密码记录信息
func IdToGetPassword(c *fiber.Ctx) error {
	userInfo := login.GetUserInfo(c)
	var ips string
	if len(c.IPs()) > 0 {
		ips = c.IPs()[0]
	} else {
		ips = c.IP()
	}
	type Req struct {
		Id string
	}
	type Resp struct {
		// 密码记录id
		ID string
		// 密码
		Password string
	}
	req := new(Req)

	if err := c.ParamsParser(req); err != nil {
		return fiberresp.RespError(c, fiberresp.REVIEW_YOUR_INPUT, err.Error())
	}

	// 检查参数
	if err := checkparameter.CheckStringsIsNil(c, userInfo.Id); err != nil {
		return err
	}
	if err := checkparameter.CheckStringsIsNil(c, req.Id); err != nil {
		return err
	}

	db := database.DB
	resp := new(Resp)
	// 查询记录
	res := db.Model(&model.TPasswordInfo{}).
		Where(&model.TPasswordInfo{ID: req.Id}).
		First(resp)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			recordRequest(ips, userInfo.Id, req.Id, model.NotFound)
			return fiberresp.RespError(c, fiberresp.RECORD_NOT_EXISTS)
		}
		return fiberresp.RespError(c, fiberresp.SELECT_ERROR, res.Error.Error())
	}
	// 密码
	passwordDecoded, err := Decode(resp.Password)
	if err != nil {
		fmt.Printf("id: %v get password decode error. %s\n", req.Id, err.Error())
		// 这里解析失败只报错就算了，不直接抛出错误给前端
		// 这里解释失败了代表数据库有脏数据了
		// 如果resp.Password为空，就让前端提示用户更新password
	}
	resp.Password = passwordDecoded

	recordRequest(ips, userInfo.Id, req.Id, model.Success)
	return fiberresp.RespSuccess(c, resp)
}

// 记录密码请求
func recordRequest(ip string, userId string, passwordId string, status model.GetPasswordRecordsStatus) {
	db := database.DB
	db.Create(&model.TGetPasswordRecords{
		UserId:     userId,
		PasswordId: passwordId,
		Status:     status,
		IP:         ip,
	})
}
