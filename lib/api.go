package lib

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

// QueryResponse 查询响应结构体
type QueryResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Total   int64       `json:"total"`
}

type queryOptionTradesResp struct {
	Time         string
	Symbol       string
	CreationDate string
	Exp          string
	Strike       string
	CP           string
	Spot         string
	Details      string
	Type         string
	Value        string
	Iv           string
	Color        string
}

// QueryRequest 简单查询请求结构体
type QueryRequest struct {
	OptionType         []string `json:"optionType"`         // 期权类型：CALL/PUT，可选
	Symbol             string   `json:"symbol"`             // 股票代码，可选
	FlowColor          []string `json:"flowColor"`          // 颜色，可选
	SecurityType       []string `json:"securityType"`       // 证券类型：ETF/Stock，可选
	Sector             []string `json:"sector"`             // 行业分类，可选
	Limit              int      `json:"limit"`              // 限制返回数量，默认100
	Offset             int      `json:"offset"`             // 偏移量，默认0
	StartDate          string   `json:"startDate"`          // 开始日期，格式：2025-01-01
	EndDate            string   `json:"endDate"`            // 结束日期，格式：2025-01-01
	BidAsk             []string `json:"bidAsk"`             // 买卖方向，可选
	MarketCapAbove750B bool     `json:"marketCapAbove750B"` // 市场资本大于750B
	PreValue           []string `json:"preValue"`

	WeepOnly        bool   `json:"weepOnly"`
	WeeklyOnly      bool   `json:"weeklyOnly"`
	Earnings        bool   `json:"earnings"`
	Unusual         bool   `json:"unusual"`
	ShowExDiv       bool   `json:"showExDiv"`
	LastId          int    `json:"lastId"`
	LastCreationate string `json:"lastCreationate"`
}

// queryOptionTrades 简单查询期权交易数据
func queryOptionTrades(c *gin.Context) {
	var req QueryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, QueryResponse{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Limit <= 0 {
		req.Limit = 100
	}
	if req.Offset < 0 {
		req.Offset = 0
	}

	query := DB.Model(&OptionTrade{})

	if len(req.PreValue) > 0 {
		// 获取req.PreValue最小的
		minPreValue, _ := strconv.Atoi(req.PreValue[0])
		for _, v := range req.PreValue {
			preValue, _ := strconv.Atoi(v)
			if preValue < minPreValue {
				minPreValue = preValue
			}
		}
		query = query.Where("premium >= ?", minPreValue)
	}

	// 构建查询

	// 添加过滤条件
	if req.Symbol != "" {
		query = query.Where("symbol LIKE ?", "%"+req.Symbol+"%")
	}

	if len(req.OptionType) > 0 {
		query = query.Where("option_type IN ?", req.OptionType)
	}
	if len(req.FlowColor) > 0 {
		query = query.Where("color IN ?", req.FlowColor)
	}

	if len(req.SecurityType) > 0 {
		query = query.Where("security_type IN ?", req.SecurityType)
	}

	if req.WeepOnly {
		query = query.Where("weekly_option = ?", "T")
	}

	if req.Earnings {
		query = query.Where("earnings_report = ?", "T")
	}
	if req.Unusual {
		query = query.Where("unusual_activity = ?", "T")
	}
	if req.ShowExDiv {
		query = query.Where("ex_div = ?", "1")
	}
	if req.LastCreationate != "" {
		query = query.Where("creation_date > ?", req.LastCreationate)
	}

	if len(req.Sector) > 0 && len(req.Sector) != 10 {
		query = query.Where("sector IN ?", append(req.Sector, "None"))
	}
	if len(req.BidAsk) > 0 {
		query = query.Where("bid_ask IN ?", req.BidAsk)
	}
	if req.MarketCapAbove750B {
		query = query.Where("market_cap < ?", 750000000000)
	}
	if req.LastId > 0 {
		query = query.Where("id > ?", req.LastId)
	}

	// if req.StartDate != "" {
	// 	query = query.Where("created_at >= ?", req.StartDate+" 00:00:00")
	// }

	// if req.EndDate != "" {
	// 	query = query.Where("created_at <= ?", req.EndDate+" 23:59:59")
	// }

	// 获取总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		Logger.Error("查询总数失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, QueryResponse{
			Code:    500,
			Message: "查询总数失败",
		})
		return
	}

	// 分页查询
	var trades []OptionTrade
	if err := query.Order("creation_date DESC").
		Offset(req.Offset).
		Limit(req.Limit).
		Find(&trades).Error; err != nil {
		Logger.Error("查询数据失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, QueryResponse{
			Code:    500,
			Message: "查询数据失败",
		})
		return
	}

	result := make([]queryOptionTradesResp, len(trades))
	for k, v := range trades {
		// 将 int64 时间戳转换为 time.Time
		creationTime := time.Unix(v.CreationDate, 0)
		expirationTime := time.Unix(v.Expiration, 0)

		result[k] = queryOptionTradesResp{
			Time:         creationTime.UTC().Format("15:04:05"),
			Symbol:       v.Symbol,
			Exp:          expirationTime.UTC().Format("01/02/06"),
			Strike:       v.Strike.String(),
			CP:           v.OptionType,
			Spot:         v.Spot.String(),
			Details:      v.Details,
			CreationDate: fmt.Sprintf("%d", v.CreationDate),
			Type:         v.TradeType,
			Value:        v.Premium.String(),
			Iv:           v.ImpliedVolatility.Mul(decimal.NewFromInt(100)).Round(2).String(),
			Color:        v.Color,
		}
	}

	// 返回结果
	c.JSON(http.StatusOK, QueryResponse{
		Code:    200,
		Message: "查询成功",
		Data:    result,
		Total:   total,
	})
}

// GetStats 获取统计信息
func GetStats(c *gin.Context) {
	var total int64
	var todayCount int64
	var callCount int64
	var putCount int64

	// 获取总数
	if err := DB.Model(&OptionTrade{}).Count(&total).Error; err != nil {
		Logger.Error("获取总数失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取统计信息失败"})
		return
	}

	// 获取今日数量
	today := time.Now().Format("2006-01-02")
	if err := DB.Model(&OptionTrade{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayCount).Error; err != nil {
		Logger.Error("获取今日数量失败", zap.Error(err))
	}

	// 获取CALL数量
	if err := DB.Model(&OptionTrade{}).
		Where("option_type = ?", "CALL").
		Count(&callCount).Error; err != nil {
		Logger.Error("获取CALL数量失败", zap.Error(err))
	}

	// 获取PUT数量
	if err := DB.Model(&OptionTrade{}).
		Where("option_type = ?", "PUT").
		Count(&putCount).Error; err != nil {
		Logger.Error("获取PUT数量失败", zap.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"todayCount": todayCount,
		"callCount":  callCount,
		"putCount":   putCount,
	})
}

// GetSymbols 获取所有股票代码
func GetSymbols(c *gin.Context) {
	var symbols []string

	if err := DB.Model(&OptionTrade{}).
		Distinct().
		Pluck("symbol", &symbols).Error; err != nil {
		Logger.Error("获取股票代码失败", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取股票代码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"symbols": symbols,
		"count":   len(symbols),
	})
}

// ToUTCTime 将时间转换为指定格式的UTC时间字符串
func ToUTCTime(t *time.Time, format string) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(format)
}
