package conf

import (
	"encoding/json"
	"github.com/cat3306/gocommon/confutil"
)

var (
	ServerConfig *ServerConf
)

type ServerConf struct {
	Name            string                 `json:"name"`
	Host            string                 `json:"host"`     //ip
	Port            int                    `json:"port"`     //port
	MaxConn         int                    `json:"max_conn"` //最大连接数
	ConnWriteBuffer int                    `json:"conn_write_buffer"`
	ConnReadBuffer  int                    `json:"conn_read_buffer"`
	KV              map[string]interface{} `json:"kv"`
}
type MysqlConfig struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Pwd          string `json:"pwd"`
	ConnPoolSize int    `json:"conn_pool_size"`
	SetLog       bool   `json:"set_log"`
}
type RedisConfig struct {
	Dbs      []int  `json:"dbs"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

func Init(filePath string) error {
	c := confutil.Config{}
	cf := &ServerConf{}
	err := c.Load(filePath, cf)
	if err != nil {
		panic(err)
	}
	ServerConfig = cf
	return nil
}
func MapToStruct(v interface{}, m map[string]interface{}) error {
	raw, err := json.Marshal(m)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}
