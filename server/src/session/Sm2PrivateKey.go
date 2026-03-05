package session

import (
	"errors"
	"fmt"

	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"gorm.io/gorm"
)

var (
	// sm2私钥
	Sm2PrivateKey *sm2.PrivateKey
	// sm2公钥
	Sm2PublicKeyStr string
)

func normalizeHexString(s string) string {
	if len(s)%2 != 0 {
		return "0" + s
	}
	return s
}

// 从数据库中获取sm2私钥加载到变量
func getSm2PrivateKey() error {
	db := database.DB

	// 读数据库
	tsi := new(model.TSm2Info)
	res := db.First(tsi)
	if res.Error != nil {
		// 没有记录
		if res.Error == gorm.ErrRecordNotFound {
			// 加载sm2私钥失败
			return errors.New("加载sm2私钥失败，在数据库里没有找到私钥")
		}
		fmt.Println("加载sm2私钥失败:", res.Error)
		return fmt.Errorf("加载sm2私钥失败: %s", res.Error)
	}

	// 读取私钥
	priv, err := x509.ReadPrivateKeyFromHex(normalizeHexString(tsi.PrivateKey))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return err
	}
	Sm2PrivateKey = priv
	fmt.Println("sm2 私钥加载成功.")

	Sm2PublicKeyStr = tsi.PublicKey
	fmt.Println("sm2 公钥加载成功.")

	return nil
}
