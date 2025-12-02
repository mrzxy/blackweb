package lib

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var router *gin.Engine

// setupStaticFiles 设置静态文件服务
func setupStaticFiles() {
	// 设置根路径重定向到web目录

	// 设置静态文件目录（放在最后，避免路由冲突）
	router.Static("/web", "./web")

	Logger.Info("静态文件服务已配置，web目录: ./web")
}

// InitServer 初始化Gin服务器
func InitServer() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()

	// 设置静态文件服务
	setupStaticFiles()

	// 设置路由
	setupRoutes()

	// 启动服务器
	go func() {
		if err := router.Run(fmt.Sprintf(":%d", Conf.Port)); err != nil {
			Logger.Fatal("启动Gin服务器失败", zap.Error(err))
		}
	}()

	Logger.Info("Gin服务器已启动，监听端口: 8080")
}

// setupRoutes 设置路由
func setupRoutes() {
	// 健康检查接口
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "服务器运行正常",
		})
	})

	// 查询期权交易数据接口 - 使用RequestBody结构体

	// 简单查询接口
	router.POST("/api/option-trades", queryOptionTrades)

	// 统计信息接口
	router.GET("/api/stats", GetStats)

	// 获取所有股票代码
	router.GET("/api/symbols", GetSymbols)
}
