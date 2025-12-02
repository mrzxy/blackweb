package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/shopspring/decimal"
)

// API响应结构
type APIResponse []TradeData

type TradeData struct {
	ID struct {
		Timestamp    int64  `json:"timestamp"`
		CreationTime string `json:"creationTime"`
	} `json:"id"`
	Order             int     `json:"order"`
	CreatedDate       string  `json:"createdDate"`
	Symbol            string  `json:"symbol"`
	Type              string  `json:"type"`
	Details           string  `json:"details"`
	BidAsk            string  `json:"bidAsk"`
	ContractPrice     float64 `json:"contractPrice"`
	Volume            int     `json:"volume"`
	CallPut           string  `json:"callPut"`
	Strike            float64 `json:"strike"`
	Spot              float64 `json:"spot"`
	Premium           float64 `json:"premium"`
	Expiration        string  `json:"expiration"`
	Color             string  `json:"color"`
	ImpliedVolatility float64 `json:"impliedVolatility"`
	Dte               int     `json:"dte"`
	Er                string  `json:"er"`
	StockEtf          string  `json:"stockEtf"`
	Sector            string  `json:"sector"`
	Uoa               string  `json:"uoa"`
	Weekly            string  `json:"weekly"`
	MktCap            int64   `json:"mktCap"`
	Oi                int     `json:"oi"`
	Itm               int     `json:"itm"`
	Ex                int     `json:"ex"`
}

// 请求体结构
type RequestBody struct {
	Historical bool   `json:"historical"`
	Symbol     string `json:"symbol"`
	Strike     int    `json:"strike"`
	Count      int    `json:"count"`
	Filter     int64  `json:"filter"`
	Filters    struct {
		OptionsDate struct {
			End   string `json:"end"`
			Start string `json:"start"`
		} `json:"optionsDate"`
		ExpireOptionsDate struct {
			End   string `json:"end"`
			Start string `json:"start"`
		} `json:"expireOptionsDate"`
		OptionsFlowPuts                  bool `json:"optionsFlowPuts"`
		OptionsFlowCalls                 bool `json:"optionsFlowCalls"`
		OptionsFlowYellow                bool `json:"optionsFlowYellow"`
		OptionsFlowWhite                 bool `json:"optionsFlowWhite"`
		OptionsFlowMagenta               bool `json:"optionsFlowMagenta"`
		OptionsFlowAboveAskOnly          bool `json:"optionsFlowAboveAskOnly"`
		OptionsFlowBelowBidOnly          bool `json:"optionsFlowBelowBidOnly"`
		OptionsFlowAtOrAboveAsk          bool `json:"optionsFlowAtOrAboveAsk"`
		OptionsFlowAtOrBelowBid          bool `json:"optionsFlowAtOrBelowBid"`
		OptionsFlowMultileg              bool `json:"optionsFlowMultileg"`
		OptionsFlowOnlyMultiLeg          bool `json:"optionsFlowOnlyMultiLeg"`
		OptionsFlowBelowPoint5           bool `json:"optionsFlowBelowPoint5"`
		OptionsFlowBelow5                bool `json:"optionsFlowBelow5"`
		OptionsFlow100Contracts          bool `json:"optionsFlow100Contracts"`
		OptionsFlow500Contracts          bool `json:"optionsFlow500Contracts"`
		OptionsFlow5000Contracts         bool `json:"optionsFlow5000Contracts"`
		OptionsFlowStock                 bool `json:"optionsFlowStock"`
		OptionsFlowEtf                   bool `json:"optionsFlowEtf"`
		OptionsFlowAbove50k              bool `json:"optionsFlowAbove50k"`
		OptionsFlowAbove100k             bool `json:"optionsFlowAbove100k"`
		OptionsFlowAbove200k             bool `json:"optionsFlowAbove200k"`
		OptionsFlowAbove500k             bool `json:"optionsFlowAbove500k"`
		OptionsFlowAbove1m               bool `json:"optionsFlowAbove1m"`
		MarketCapAbove750B               bool `json:"marketCapAbove750B"`
		OptionsFlowInTheMoney            bool `json:"optionsFlowInTheMoney"`
		OptionsFlowOutOfTheMoney         bool `json:"optionsFlowOutOfTheMoney"`
		OptionsFlowSweepOnly             bool `json:"optionsFlowSweepOnly"`
		OptionsFlowWeeklyOnly            bool `json:"optionsFlowWeeklyOnly"`
		OptionsFlowEarningsReportOnly    bool `json:"optionsFlowEarningsReportOnly"`
		OptionsFlowUnusualOnly           bool `json:"optionsFlowUnusualOnly"`
		OptionsFlowExDiv                 bool `json:"optionsFlowExDiv"`
		OptionsFlowConsumerDiscretionary bool `json:"optionsFlowConsumerDiscretionary"`
		OptionsFlowIndustrials           bool `json:"optionsFlowIndustrials"`
		OptionsFlowInformationTechnology bool `json:"optionsFlowInformationTechnology"`
		OptionsFlowRealEstate            bool `json:"optionsFlowRealEstate"`
		OptionsFlowHealthCare            bool `json:"optionsFlowHealthCare"`
		OptionsFlowEnergy                bool `json:"optionsFlowEnergy"`
		OptionsFlowFinancials            bool `json:"optionsFlowFinancials"`
		OptionsFlowMaterials             bool `json:"optionsFlowMaterials"`
		OptionsFlowConsumerStaples       bool `json:"optionsFlowConsumerStaples"`
		OptionsFlowCommunicationServices bool `json:"optionsFlowCommunicationServices"`
		OptionsFlowUtilities             bool `json:"optionsFlowUtilities"`
		OptionsExpirationRange           bool `json:"optionsExpirationRange"`
		OptionsFlowSectorNone            bool `json:"optionsFlowSectorNone"`
	} `json:"filters"`
	FromDate string `json:"fromDate"`
	ToDate   string `json:"toDate"`
}

func fetchAndSaveData() {
	log.Println("开始抓取数据...")

	// 构建请求体
	now := time.Now()
	requestBody := buildRequestBody(now)

	// 发送HTTP请求
	data, err := sendHTTPRequest(requestBody)
	if err != nil {
		Logger.Infof("HTTP请求失败: %v", err)
		return
	}

	// 解析响应
	var response APIResponse
	if err := json.Unmarshal(data, &response); err != nil {
		Logger.Infof("解析响应失败: %v", err)
		return
	}

	// 保存数据到数据库
	saveDataToDB(response)

	Logger.Infof("本次抓取完成，获取到 %d 条数据", len(response))
}

func buildRequestBody(now time.Time) RequestBody {
	timeStr := now.Format("2006-01-02T15:04:05.000Z")

	return RequestBody{
		Historical: false,
		Symbol:     "",
		Strike:     0,
		Count:      300,
		Filter:     2198487171391,
		FromDate:   timeStr,
		ToDate:     timeStr,
		Filters: struct {
			OptionsDate struct {
				End   string `json:"end"`
				Start string `json:"start"`
			} `json:"optionsDate"`
			ExpireOptionsDate struct {
				End   string `json:"end"`
				Start string `json:"start"`
			} `json:"expireOptionsDate"`
			OptionsFlowPuts                  bool `json:"optionsFlowPuts"`
			OptionsFlowCalls                 bool `json:"optionsFlowCalls"`
			OptionsFlowYellow                bool `json:"optionsFlowYellow"`
			OptionsFlowWhite                 bool `json:"optionsFlowWhite"`
			OptionsFlowMagenta               bool `json:"optionsFlowMagenta"`
			OptionsFlowAboveAskOnly          bool `json:"optionsFlowAboveAskOnly"`
			OptionsFlowBelowBidOnly          bool `json:"optionsFlowBelowBidOnly"`
			OptionsFlowAtOrAboveAsk          bool `json:"optionsFlowAtOrAboveAsk"`
			OptionsFlowAtOrBelowBid          bool `json:"optionsFlowAtOrBelowBid"`
			OptionsFlowMultileg              bool `json:"optionsFlowMultileg"`
			OptionsFlowOnlyMultiLeg          bool `json:"optionsFlowOnlyMultiLeg"`
			OptionsFlowBelowPoint5           bool `json:"optionsFlowBelowPoint5"`
			OptionsFlowBelow5                bool `json:"optionsFlowBelow5"`
			OptionsFlow100Contracts          bool `json:"optionsFlow100Contracts"`
			OptionsFlow500Contracts          bool `json:"optionsFlow500Contracts"`
			OptionsFlow5000Contracts         bool `json:"optionsFlow5000Contracts"`
			OptionsFlowStock                 bool `json:"optionsFlowStock"`
			OptionsFlowEtf                   bool `json:"optionsFlowEtf"`
			OptionsFlowAbove50k              bool `json:"optionsFlowAbove50k"`
			OptionsFlowAbove100k             bool `json:"optionsFlowAbove100k"`
			OptionsFlowAbove200k             bool `json:"optionsFlowAbove200k"`
			OptionsFlowAbove500k             bool `json:"optionsFlowAbove500k"`
			OptionsFlowAbove1m               bool `json:"optionsFlowAbove1m"`
			MarketCapAbove750B               bool `json:"marketCapAbove750B"`
			OptionsFlowInTheMoney            bool `json:"optionsFlowInTheMoney"`
			OptionsFlowOutOfTheMoney         bool `json:"optionsFlowOutOfTheMoney"`
			OptionsFlowSweepOnly             bool `json:"optionsFlowSweepOnly"`
			OptionsFlowWeeklyOnly            bool `json:"optionsFlowWeeklyOnly"`
			OptionsFlowEarningsReportOnly    bool `json:"optionsFlowEarningsReportOnly"`
			OptionsFlowUnusualOnly           bool `json:"optionsFlowUnusualOnly"`
			OptionsFlowExDiv                 bool `json:"optionsFlowExDiv"`
			OptionsFlowConsumerDiscretionary bool `json:"optionsFlowConsumerDiscretionary"`
			OptionsFlowIndustrials           bool `json:"optionsFlowIndustrials"`
			OptionsFlowInformationTechnology bool `json:"optionsFlowInformationTechnology"`
			OptionsFlowRealEstate            bool `json:"optionsFlowRealEstate"`
			OptionsFlowHealthCare            bool `json:"optionsFlowHealthCare"`
			OptionsFlowEnergy                bool `json:"optionsFlowEnergy"`
			OptionsFlowFinancials            bool `json:"optionsFlowFinancials"`
			OptionsFlowMaterials             bool `json:"optionsFlowMaterials"`
			OptionsFlowConsumerStaples       bool `json:"optionsFlowConsumerStaples"`
			OptionsFlowCommunicationServices bool `json:"optionsFlowCommunicationServices"`
			OptionsFlowUtilities             bool `json:"optionsFlowUtilities"`
			OptionsExpirationRange           bool `json:"optionsExpirationRange"`
			OptionsFlowSectorNone            bool `json:"optionsFlowSectorNone"`
		}{
			OptionsDate: struct {
				End   string `json:"end"`
				Start string `json:"start"`
			}{
				End:   timeStr,
				Start: timeStr,
			},
			ExpireOptionsDate: struct {
				End   string `json:"end"`
				Start string `json:"start"`
			}{
				End:   timeStr,
				Start: timeStr,
			},
			OptionsFlowPuts:                  true,
			OptionsFlowCalls:                 true,
			OptionsFlowYellow:                true,
			OptionsFlowWhite:                 true,
			OptionsFlowMagenta:               true,
			OptionsFlowAboveAskOnly:          true,
			OptionsFlowBelowBidOnly:          false,
			OptionsFlowAtOrAboveAsk:          true,
			OptionsFlowAtOrBelowBid:          false,
			OptionsFlowMultileg:              false,
			OptionsFlowOnlyMultiLeg:          false,
			OptionsFlowBelowPoint5:           false,
			OptionsFlowBelow5:                false,
			OptionsFlow100Contracts:          false,
			OptionsFlow500Contracts:          false,
			OptionsFlow5000Contracts:         false,
			OptionsFlowStock:                 true,
			OptionsFlowEtf:                   true,
			OptionsFlowAbove50k:              false,
			OptionsFlowAbove100k:             false,
			OptionsFlowAbove200k:             false,
			OptionsFlowAbove500k:             false,
			OptionsFlowAbove1m:               false,
			MarketCapAbove750B:               false,
			OptionsFlowInTheMoney:            false,
			OptionsFlowOutOfTheMoney:         false,
			OptionsFlowSweepOnly:             false,
			OptionsFlowWeeklyOnly:            false,
			OptionsFlowEarningsReportOnly:    false,
			OptionsFlowUnusualOnly:           false,
			OptionsFlowExDiv:                 false,
			OptionsFlowConsumerDiscretionary: true,
			OptionsFlowIndustrials:           true,
			OptionsFlowInformationTechnology: true,
			OptionsFlowRealEstate:            true,
			OptionsFlowHealthCare:            true,
			OptionsFlowEnergy:                true,
			OptionsFlowFinancials:            true,
			OptionsFlowMaterials:             true,
			OptionsFlowConsumerStaples:       true,
			OptionsFlowCommunicationServices: true,
			OptionsFlowUtilities:             true,
			OptionsExpirationRange:           false,
			OptionsFlowSectorNone:            true,
		},
	}
}

func sendHTTPRequest(requestBody RequestBody) ([]byte, error) {
	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.blackboxstocks.com/api/v2/options/getFlowMobile", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	// 设置请求头
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("authorization", "Bearer "+Conf.Token)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("origin", "https://members.blackboxstocks.com")
	req.Header.Set("referer", "https://members.blackboxstocks.com/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP请求失败，状态码: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func saveDataToDB(trades []TradeData) {
	var savedCount, skippedCount int

	for _, trade := range trades {
		// 生成TradeID：timestamp + creationTime
		creationDate, err := time.Parse("2006-01-02T15:04:05Z", trade.CreatedDate)
		if err != nil {
			Logger.Infof("解析创建时间失败: %v", err)
			continue
		}

		tradeID := fmt.Sprintf("%d_%s", trade.ID.Timestamp, trade.ID.CreationTime)

		// 检查TradeID是否已存在
		var existingTrade OptionTrade
		result := DB.Where("trade_id = ?", tradeID).First(&existingTrade)
		if result.Error == nil {
			// TradeID已存在，跳过
			skippedCount++
			continue
		}

		// 解析其他时间字段
		expiration, _ := time.Parse("2006-01-02T15:04:05Z", trade.Expiration)

		// 创建OptionTrade记录
		optionTrade := OptionTrade{
			TradeID:           tradeID,
			Timestamp:         trade.ID.Timestamp,
			CreationDate:      creationDate.UnixMilli(),
			OrderID:           int64(trade.Order),
			Symbol:            trade.Symbol,
			TradeType:         trade.Type,
			Details:           trade.Details,
			BidAsk:            trade.BidAsk,
			ContractPrice:     decimal.NewFromFloat(trade.ContractPrice),
			Volume:            trade.Volume,
			OptionType:        trade.CallPut,
			Strike:            decimal.NewFromFloat(trade.Strike),
			Spot:              decimal.NewFromFloat(trade.Spot),
			Premium:           decimal.NewFromFloat(trade.Premium),
			Expiration:        expiration.UnixMilli(),
			Color:             trade.Color,
			ImpliedVolatility: decimal.NewFromFloat(trade.ImpliedVolatility),
			Dte:               trade.Dte,
			EarningsReport:    trade.Er,
			SecurityType:      trade.StockEtf,
			Sector:            trade.Sector,
			UnusualActivity:   trade.Uoa,
			WeeklyOption:      trade.Weekly,
			MarketCap:         trade.MktCap,
			OpenInterest:      trade.Oi,
			Itm:               strconv.Itoa(trade.Itm),
			ExDiv:             strconv.Itoa(trade.Ex),
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		// 保存到数据库
		if err := DB.Create(&optionTrade).Error; err != nil {
			Logger.Infof("保存数据失败: %v", err)
			continue
		}

		savedCount++
	}

	Logger.Infof("数据保存完成: 新增 %d 条，跳过 %d 条", savedCount, skippedCount)
}

func RunSpider() {
	// 启动定时任务
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	Logger.Infof("开始定时抓取数据，间隔10秒...")

	// 立即执行一次
	fetchAndSaveData()

	// 定时执行
	for range ticker.C {
		fetchAndSaveData()
	}
}
