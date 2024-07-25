// 该代码文件初始化了一个 Redis 客户端，并提供了一个获取 Redis 客户端单例的函数。

package conf

import (
	"log"
	"sync"

	"github.com/redis/rueidis"
)

var (
	_redisClient rueidis.Client // 全局变量，存储 Redis 客户端实例
	_onceRedis   sync.Once      // 确保初始化仅执行一次的同步控制器
)

// MustRedisClient 返回 Redis 客户端的单例实例。
func MustRedisClient() rueidis.Client {
	_onceRedis.Do(func() {
		client, err := rueidis.NewClient(rueidis.ClientOption{
			InitAddress:      redisSetting.InitAddress,      // 初始化地址
			Username:         redisSetting.Username,         // 用户名
			Password:         redisSetting.Password,         // 密码
			SelectDB:         redisSetting.SelectDB,         // 选择的数据库
			ConnWriteTimeout: redisSetting.ConnWriteTimeout, // 连接写入超时时间
		})
		if err != nil {
			log.Fatalf("create a redis client failed: %s", err) // 如果初始化失败，则输出错误信息
		}
		_redisClient = client // 将初始化后的客户端赋值给全局变量
		// 同时初始化 CacheKeyPool
		initCacheKeyPool()
	})
	return _redisClient // 返回 Redis 客户端实例
}
