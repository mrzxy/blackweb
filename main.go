package main

import (
	"gorm.io/gorm"

	"blackweb/lib"
)

var db *gorm.DB

func main() {
	lib.InitLogger()
	lib.InitConfig()
	// 初始化数据库连接
	lib.InitDB()

	// 启动Gin服务器
	lib.InitServer()

	// 启动爬虫
	go lib.RunSpider()

	// 保持程序运行
	select {}
}
