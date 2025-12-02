package lib

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func InitDB() {
	var err error
	dsn := Conf.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:          NewDBLogger(Logger),
		CreateBatchSize: 1000,
		// 启用SQL打印
		// 需要确保NewDBLogger(Logger)支持Info级别SQL日志输出
		// 如果想用GORM自带的Logger，可以这样写：
		// logger := logger.Default.LogMode(logger.Info)
		// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		//     Logger: logger,
		//     CreateBatchSize: 1000,
		// })
		// 这里假设NewDBLogger(Logger)已实现SQL打印
	})

	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}
	DB = db

	log.Println("数据库连接成功")
}
