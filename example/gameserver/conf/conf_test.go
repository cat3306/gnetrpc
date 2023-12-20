package conf

import (
	"github.com/cat3306/gocommon/confutil"
	"testing"
)

func TestGenConf(t *testing.T) {
	c := confutil.Config{}
	cf := ServerConf{
		Name:            "gameserver",
		Host:            "0.0.0.0",
		Port:            7898,
		MaxConn:         1000,
		ConnWriteBuffer: 1024,
		ConnReadBuffer:  1024,
		KV: map[string]interface{}{
			"mysql": &MysqlConfig{
				Host:         "0.0.0.0",
				Port:         3306,
				User:         "root",
				Pwd:          "12345678",
				ConnPoolSize: 20,
				SetLog:       true,
			},
			"redis": &RedisConfig{
				Dbs:      []int{0, 1},
				Addr:     "0.0.0.0:6379",
				Password: "redis-hahah@123",
			},
		},
	}
	c.Save("./conf.json", cf)
}
