package thirdmodule

import (
	"errors"
	"github.com/cat3306/gnetrpc/example/gameserver/conf"
	"github.com/cat3306/gocommon/goredisutil"
	"github.com/go-redis/redis/v8"
)

var (
	RedisClients *goredisutil.RedisClientPool
	NilErr       = redis.Nil
)

func IgnoreRedisNil(err error) error {
	if errors.Is(err, redis.Nil) {
		err = nil
	}
	return err
}
func InitCache() {
	redisConf := &conf.RedisConfig{}
	m := conf.ServerConfig.KV["redis"].(map[string]interface{})
	err := conf.MapToStruct(redisConf, m)
	if err != nil {
		panic(err)
	}
	RedisClients = goredisutil.NewRedisClients(&goredisutil.ClientConf{
		Options: &redis.Options{
			Addr:     redisConf.Addr,
			Password: redisConf.Password,
		},
		DB: redisConf.Dbs,
	})
}

func CacheSelect(idx int) *redis.Client {
	return RedisClients.Select(idx)
}
