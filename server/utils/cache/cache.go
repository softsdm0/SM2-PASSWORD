package cache

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

// 写入cache时服务的前缀
const serverCachePrefixName = "password"

// 添加缓存名称前序
func AddPrefix(prefix string, key string) string {
	return fmt.Sprintf("%s:%s:%s", serverCachePrefixName, prefix, key)
}
