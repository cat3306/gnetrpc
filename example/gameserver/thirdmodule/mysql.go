package thirdmodule

import (
	"fmt"
	"github.com/cat3306/gnetrpc/example/gameserver/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var (
	MysqlDb *gorm.DB
)

func InitDb() {
	mysqlConf := &conf.MysqlConfig{}
	m := conf.ServerConfig.KV["mysql"].(map[string]interface{})
	err := conf.MapToStruct(mysqlConf, m)
	if err != nil {
		panic(err)
	}
	conn := fmt.Sprintf("%s:%s@(%s:%d)/?charset=utf8mb4&parseTime=True&loc=Local", mysqlConf.User, mysqlConf.Pwd, mysqlConf.Host, mysqlConf.Port)
	db, err := gorm.Open(mysql.Open(conn))
	if err != nil {
		panic(err)
	}
	sqlDb, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDb.SetMaxOpenConns(mysqlConf.ConnPoolSize)
	sqlDb.SetMaxIdleConns(mysqlConf.ConnPoolSize / 2)
	err = sqlDb.Ping()
	if err != nil {
		panic(err)
	}
	if mysqlConf.SetLog {
		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             time.Second, // Slow SQL threshold
				LogLevel:                  logger.Info, // Log level
				IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,        // Disable color
			},
		)
		db.Logger = newLogger
	}
	MysqlDb = db
}
