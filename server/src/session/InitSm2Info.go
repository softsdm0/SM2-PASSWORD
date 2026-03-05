package session

import (
	"crypto/rand"
	"fmt"
	"os"

	"gitee.com/ouhaoqiang/passwordserver/server/src/model"
	"gitee.com/ouhaoqiang/passwordserver/server/utils/database"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/x509"
	"gorm.io/gorm"
)

// 初始化sm2密钥
func InitSm2Info() {
	db := database.DB

	tsi := new(model.TSm2Info)
	res := db.First(tsi)
	if res.Error != nil {
		if res.Error != gorm.ErrRecordNotFound {
			fmt.Println("查询sm2 info失败:", res.Error.Error())
			os.Exit(1000)
		}
		// 没有记录
		// 生成私钥
		privateKey, err := sm2.GenerateKey(rand.Reader)
		if err != nil {
			fmt.Printf("生成sm2密钥失败: %v\n", err.Error())
			os.Exit(10002)
		}
		// 生成私钥16进制字符串
		privHex := normalizeHexString(x509.WritePrivateKeyToHex(privateKey))

		// 公钥
		pubkey := privateKey.PublicKey
		// 生成公钥16进制字符串
		pubHex := normalizeHexString(x509.WritePublicKeyToHex(&pubkey))

		// 没有记录，需要插入
		tsi.PublicKey = pubHex
		tsi.PrivateKey = privHex
		res = db.Create(tsi)
		if res.Error != nil {
			fmt.Println("写入sm2 info失败:", res.Error.Error())
			os.Exit(10001)
		}

		fmt.Printf("sm2密钥新创建成功, 公钥: %s\n", tsi.PublicKey)
	}

	// 私钥到内存
	err := getSm2PrivateKey()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(10003)
	}

	// 如果存在就不创建了
}
