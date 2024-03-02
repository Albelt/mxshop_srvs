package main

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"mxshop_srvs/stock_srv/model"
	"os"
	"time"
)
import "gorm.io/driver/mysql"

const (
	mysqlDsn = "root:Albelt2017Go@tcp(ubuntu-learn:3306)/mxshop_srvs?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	myLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second * 2,
			Colorful:      false,
			LogLevel:      logger.Info,
		})

	db, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{
		Logger: myLogger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "stock_srv_",
			SingularTable: false,
		},
	})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(model.Stock{}, model.StockHistory{})
	if err != nil {
		log.Fatalf("AutoMigrate failed, err:%s", err.Error())
	}
}
