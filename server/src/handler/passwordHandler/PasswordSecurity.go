package passwordHandler

import (
	"encoding/hex"
	"errors"

	"gitee.com/ouhaoqiang/passwordserver/server/utils/fiberresp"
	"github.com/gofiber/fiber/v2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/sm4"
)

// 对PasswordInfo结构体需要的内容进行base64编码
func EncodePasswordInfo(c *fiber.Ctx, account string, password string) (string, string, error) {
	accountEncode, err := encode(account)
	if err != nil {
		return "", "", nil
	}

	passwordEncode, err := encode(password)
	if err != nil {
		return "", "", err
	}

	return accountEncode, passwordEncode, nil
}

// 对PasswordInfo结构体需要的内容进行base64解码
func DecodePasswordInfo(c *fiber.Ctx, account string, password string) (string, string, error) {
	// 账户
	accountDecode, err := Decode(account)
	if err != nil {
		return "", "", fiberresp.RespError(c, fiberresp.ACCOUNT_DECRYPT, err.Error())
	}

	// 密码
	passwordDecode, err := Decode(password)
	if err != nil {
		return "", "", fiberresp.RespError(c, fiberresp.PASSWORD_DECRYPT, err.Error())
	}

	return accountDecode, passwordDecode, nil
}

// 加密算法
func encode(s string) (string, error) {
	if s == "" {
		return "", errors.New("加密失败: 字符串为空")
	}

	sByte := []byte(s)
	sm3EncodeByte := sm3.Sm3Sum(sByte)
	sm4EncodeByte, err := sm4.Sm4Ecb(getSm4Key(sm3EncodeByte), sByte, true)
	if err != nil {
		return "", err
	}

	saveEncodeByte := getSaveEncodeByte(sm4EncodeByte, sm3EncodeByte)

	rHex := hex.EncodeToString(saveEncodeByte)

	return rHex, nil
}

// 解密算法
func Decode(s string) (string, error) {
	sByte, err := hex.DecodeString(s)
	if err != nil {
		return "", err
	}

	if len(sByte) < 32+1 {
		// 32是sm3sum后的长度，1代表有东西
		return "", errors.New("解密失败: 字符串长度小于32 + 1 byte")
	}

	sm4EncodeByte, sm3EncodeByte := getSm4EncodeByteAndSm3EncodeByte(sByte)
	sm4Key := getSm4Key(sm3EncodeByte)
	sm4decodebyte, err := sm4.Sm4Ecb(sm4Key, sm4EncodeByte, false)
	if err != nil {
		return "", err
	}

	return string(sm4decodebyte), nil
}

// 在sm3sum中切片16位出来作为sm4的key
func getSm4Key(sm3Sumed []byte) []byte {
	sm4Key := make([]byte, 0, 16)
	sm4Key = append(sm4Key, sm3Sumed[16:18]...)
	sm4Key = append(sm4Key, sm3Sumed[2:4]...)
	sm4Key = append(sm4Key, sm3Sumed[4:6]...)
	sm4Key = append(sm4Key, sm3Sumed[20:22]...)
	sm4Key = append(sm4Key, sm3Sumed[30:32]...)
	sm4Key = append(sm4Key, sm3Sumed[10:12]...)
	sm4Key = append(sm4Key, sm3Sumed[9:11]...)
	sm4Key = append(sm4Key, sm3Sumed[15:17]...)
	return sm4Key
}

// 获取需要保存内容排序好的byte
func getSaveEncodeByte(sm4EncodeByte, sm3EncodeByte []byte) []byte {
	saveEncodeByte := make([]byte, 0, len(sm4EncodeByte)+len(sm3EncodeByte))
	saveEncodeByte = append(saveEncodeByte, sm3EncodeByte[:8]...)
	saveEncodeByte = append(saveEncodeByte, sm3EncodeByte[16:24]...)
	saveEncodeByte = append(saveEncodeByte, sm4EncodeByte...)
	saveEncodeByte = append(saveEncodeByte, sm3EncodeByte[24:]...)
	saveEncodeByte = append(saveEncodeByte, sm3EncodeByte[8:16]...)

	saveEncodeByte = append(saveEncodeByte[len(saveEncodeByte)/2:], saveEncodeByte[:len(saveEncodeByte)/2]...)
	return saveEncodeByte
}

// 获取sm4加密后的内容和sm3sum
func getSm4EncodeByteAndSm3EncodeByte(saveEncodeByte []byte) (sm4EncodeByte, sm3EncodeByte []byte) {
	saveEncodeByte = append(saveEncodeByte[len(saveEncodeByte)/2:], saveEncodeByte[:len(saveEncodeByte)/2]...)

	sm4EncodeByte = append(sm4EncodeByte, saveEncodeByte[16:len(saveEncodeByte)-16]...)

	sm3EncodeByte = append(sm3EncodeByte, saveEncodeByte[:8]...)
	sm3EncodeByte = append(sm3EncodeByte, saveEncodeByte[len(saveEncodeByte)-8:]...)
	sm3EncodeByte = append(sm3EncodeByte, saveEncodeByte[8:16]...)
	sm3EncodeByte = append(sm3EncodeByte, saveEncodeByte[len(saveEncodeByte)-16:len(saveEncodeByte)-8]...)
	return
}
