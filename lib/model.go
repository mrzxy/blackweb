package lib

import (
	"time"

	"github.com/shopspring/decimal"
)

type OptionTrade struct {
	ID                int             `gorm:"primaryKey;autoIncrement" json:"id"`
	TradeID           string          `gorm:"type:varchar(50);not null" json:"trade_id"`
	Timestamp         int64           `gorm:"not null" json:"timestamp"`
	CreationDate      int64           `gorm:"not null" json:"creation_date"`
	OrderID           int64           `gorm:"not null" json:"order_id"`
	Symbol            string          `gorm:"type:varchar(20);not null" json:"symbol"`
	TradeType         string          `gorm:"type:varchar(20);not null" json:"trade_type"`
	Details           string          `gorm:"type:varchar(255)" json:"details"`
	BidAsk            string          `gorm:"type:varchar(20)" json:"bid_ask"`
	ContractPrice     decimal.Decimal `gorm:"not null" json:"contract_price"`
	Volume            int             `gorm:"not null" json:"volume"`
	OptionType        string          `gorm:"type:varchar(10);not null" json:"option_type"`
	Strike            decimal.Decimal `gorm:"not null" json:"strike"`
	Spot              decimal.Decimal `gorm:"not null" json:"spot"`
	Premium           decimal.Decimal `gorm:"not null" json:"premium"`
	Expiration        int64           `gorm:"not null" json:"expiration"`
	Color             string          `gorm:"type:varchar(20);not null" json:"color"`
	ImpliedVolatility decimal.Decimal `gorm:"not null" json:"implied_volatility"`
	Dte               int             `gorm:"not null" json:"dte"`
	EarningsReport    string          `gorm:"type:varchar(1);not null" json:"earnings_report"`
	SecurityType      string          `gorm:"type:varchar(20);not null" json:"security_type"`
	Sector            string          `gorm:"type:varchar(100);not null" json:"sector"`
	UnusualActivity   string          `gorm:"type:varchar(1);not null" json:"unusual_activity"`
	WeeklyOption      string          `gorm:"type:varchar(1);not null" json:"weekly_option"`
	MarketCap         int64           `gorm:"not null" json:"market_cap"`
	OpenInterest      int             `gorm:"not null" json:"open_interest"`
	Itm               string          `gorm:"type:varchar(1);not null" json:"itm"`
	ExDiv             string          `gorm:"type:varchar(1);not null" json:"ex_div"`
	CreatedAt         time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"not null" json:"updated_at"`
}

func (OptionTrade) TableName() string {
	return "option_trades"
}
