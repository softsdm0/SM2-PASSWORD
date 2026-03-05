package cache

import (
	"context"
	"fmt"
	"os"

	"gitee.com/ouhaoqiang/passwordserver/server/utils/config"
	"github.com/redis/go-redis/v9"
)

// ConnectDB connect to db
func ConnectCache() {
	c := config.Config.Server.Redis

	Rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", c.Host, c.Port),
		Password: c.Pass,
		DB:       c.DbName,
	})

	ctx := context.Background()
	val, err := Rdb.Ping(ctx).Result()
	if err != nil {
		fmt.Printf("连接redis失败: %s", err.Error())
		os.Exit(313)
	}
	fmt.Println("连接redis成功.", val)
}
