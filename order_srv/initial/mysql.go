package initial

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/order_srv/global"
	"os"
	"time"
)

const (
	MysqlDsnFormat = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local"
	tablePrefix    = "order_srv_"
)

func InitMysql() {
	var err error

	myLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second * time.Duration(global.Config.Mysql.SlowThreshold),
			Colorful:      false,
			LogLevel:      logger.Info,
		})

	dsn := fmt.Sprintf(MysqlDsnFormat,
		global.Config.Mysql.User, global.Config.Mysql.Password,
		global.Config.Mysql.Addr, global.Config.Mysql.Db)

	global.Mysql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: myLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix,
			SingularTable: false,
		},
	})
	if err != nil {
		panic(err)
	}

	zap.S().Infof("InitMysql ok.")
}
