package tdameritrade

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/go-querystring/query"
)

var (
	validContractTypes = []string{"CALL", "PUT", "ALL"}
	validStrategies    = []string{"SINGLE", "ANALYTICAL", "COVERED", "VERTICAL", "CALENDAR", "STRANGLE", "STRADDLE", "BUTTERFLY", "CONDOR", "DIAGONAL", "COLLAR", "ROLL"}
	validRanges        = []string{"ITM", "NTM", "OTM", "SAK", "SBK", "SNK", "ALL"}
	validExpMonths     = []string{"JAN", "FEB", "MAR", "APR", "MAY", "JUN", "JUL", "AUG", "SEP", "OCT", "NOV", "DEC", "ALL"}
	validOptionTypes   = []string{"S", "NS", "ALL"}
)

const (
	defaultContractType = "ALL"
	defaultStrategy     = "SINGLE"
	defaultRange        = "ALL"
	defaultExpMonth     = "ALL"
	defaultOptionType   = "ALL"
)

// OptionChainService handles communication with the optionChain related methods of
// the TDAmeritrade API.
//
// TDAmeritrade API docs: https://developer.tdameritrade.com/option-chains/apis
type OptionChainService struct {
	client *Client
}

// OptionChainOptions is parsed and translated to query options in the https request
type OptionChainOptions struct {
	ContractType     string    `url:"contractType,omitempty"`
	StrikeCount      int       `url:"strikeCount,omitempty"`
	IncludeQuotes    *bool     `url:"includeQuotes,omitempty"`
	Strategy         string    `url:"strategy,omitempty"`
	Interval         int       `url:"interval,omitempty"`
	Strike           float64   `url:"strike,omitempty"`
	Range            string    `url:"range,omitempty"`
	FromDate         time.Time `url:"fromDate,omitempty"`
	ToDate           time.Time `url:"toDate,omitempty"`
	Volatility       float64   `url:"volatility,omitempty"`
	UnderlyingPrice  float64   `url:"underlyingPrice,omitempty"`
	InterestRate     float64   `url:"interestRate,omitempty"`
	DaysToExpiration float64   `url:"daysToExpiration,omitempty"`
	ExpMonth         string    `url:"expMonth,omitempty"`
	OptionType       string    `url:"optionType,omitempty"`
}

type OptionData struct {
	PutCall                string  `json:"putCall"`
	Symbol                 string  `json:"symbol"`
	Description            string  `json:"description"`
	ExchangeName           string  `json:"exchangeName"`
	BidPrice               float64 `json:"bidPrice"`
	AskPrice               float64 `json:"askPrice"`
	MarkPrice              float64 `json:"markPrice"`
	BidSize                int     `json:"bidSize"`
	AskSize                int     `json:"askSize"`
	LastSize               int     `json:"lastSize"`
	HighPrice              float64 `json:"highPrice"`
	LowPrice               float64 `json:"lowPrice"`
	OpenPrice              float64 `json:"openPrice"`
	ClosePrice             float64 `json:"closePrice"`
	TotalVolume            int     `json:"totalVolume"`
	QuoteTimeInLong        int     `json:"quoteTimeInLong"`
	TradeTimeInLong        int     `json:"tradeTimeInLong"`
	NetChange              float64 `json:"netChange"`
	Volatility             float64 `json:"volatility"`
	Delta                  float64 `json:"delta"`
	Gamma                  float64 `json:"gamma"`
	Theta                  float64 `json:"theta"`
	Vega                   float64 `json:"vega"`
	Rho                    float64 `json:"rho"`
	TimeValue              float64 `json:"timeValue"`
	OpenInterest           float64 `json:"openInterest"`
	IsInTheMoney           bool    `json:"isInTheMoney"`
	TheoreticalOptionValue float64 `json:"theoreticalOptionValue"`
	TheoreticalVolatility  float64 `json:"theoreticalVolatility"`
	IsMini                 bool    `json:"isMini"`
	IsNonStandard          bool    `json:"isNonStandard"`
	OptionDeliverablesList []struct {
		Symbol           string `json:"string"`
		AssetType        string `json:"assetType"`
		DeliverableUnits string `json:"deliverableUnits"`
		CurrencyType     string `json:"currencyType"`
	} `json:"optionDeliverablesList"`
	StrikePrice       float64 `json:"strikePrice"`
	ExpirationDate    int64   `json:"expirationDate"`
	ExpirationType    string  `json:"expirationType"`
	Multiplier        float64 `json:"multiplier"`
	SettlementType    string  `json:"settlementType"`
	DeliverableNote   string  `json:"deliverableNote"`
	IsIndexOption     bool    `json:"isIndexOption"`
	PercentChange     float64 `json:"percentChange"`
	MarkChange        float64 `json:"markChange"`
	MarkPercentChange float64 `json:"markPercentChange"`
}
type OptionChain struct {
	Symbol     string
	Status     string
	Underlying struct {
		Ask               float64
		AskSize           int
		Bid               float64
		BidSize           int
		Change            float64
		Close             float64
		Delayed           bool
		Description       string
		ExchangeName      string
		FiftyTwoWeekHigh  float64
		FiftyTwoWeekLow   float64
		HighPrice         float64
		Last              float64
		LowPrice          float64
		Mark              float64
		MarkChange        float64
		MarkPercentChange float64
		OpenPrice         float64
		PercentChange     float64
		QuoteTime         int64
		Symbol            string
		TotalVolume       int64
		TradeTime         int64
	}
	Strategy         string
	Interval         float64
	IsDelayed        bool
	IsIndex          bool
	DaysToExpiration float64
	InterestRate     float64
	UnderlyingPrice  float64
	Volatility       float64
	Calls            []struct {
		ExpDate    time.Time
		DaysTilExp int
		Strikes    []OptionData
	}
	Puts []struct {
		ExpDate    time.Time
		DaysTilExp int
		Strikes    []OptionData
	}
}

func (c *OptionChain) UnmarshalJSON(b []byte) error {
	var raw struct {
		Symbol     string `json:"symbol"`
		Status     string `json:"status"`
		Underlying struct {
			Ask               float64 `json:"ask"`
			AskSize           int     `json:"askSize"`
			Bid               float64 `json:"bid"`
			BidSize           int     `json:"bidSize"`
			Change            float64 `json:"change"`
			Close             float64 `json:"close"`
			Delayed           bool    `json:"delayed"`
			Description       string  `json:"description"`
			ExchangeName      string  `json:"exchangeName"`
			FiftyTwoWeekHigh  float64 `json:"fiftyTwoWeekHigh"`
			FiftyTwoWeekLow   float64 `json:"fiftyTwoWeekLow"`
			HighPrice         float64 `json:"highPrice"`
			Last              float64 `json:"last"`
			LowPrice          float64 `json:"lowPrice"`
			Mark              float64 `json:"mark"`
			MarkChange        float64 `json:"markChange"`
			MarkPercentChange float64 `json:"markPercentChange"`
			OpenPrice         float64 `json:"openPrice"`
			PercentChange     float64 `json:"percentChange"`
			QuoteTime         int64   `json:"quoteTime"`
			Symbol            string  `json:"symbol"`
			TotalVolume       int64   `json:"totalVolume"`
			TradeTime         int64   `json:"tradeTime"`
		} `json:"underlying"`
		Strategy         string                             `json:"strategy"`
		Interval         float64                            `json:"interval"`
		IsDelayed        bool                               `json:"isDelayed"`
		IsIndex          bool                               `json:"isIndex"`
		DaysToExpiration float64                            `json:"daysToExpiration"`
		InterestRate     float64                            `json:"interestRate"`
		UnderlyingPrice  float64                            `json:"underlyingPrice"`
		Volatility       float64                            `json:"volatility"`
		CallExpDateMap   map[string]map[string][]OptionData `json:"callExpDateMap"`
		PutExpDateMap    map[string]map[string][]OptionData `json:"putExpDateMap"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	spew.Dump(stirng(b), raw)
	c.Symbol = raw.Symbol
	c.Status = raw.Status
	c.Underlying.Ask = raw.Underlying.Ask
	c.Underlying.AskSize = raw.Underlying.AskSize
	c.Underlying.Bid = raw.Underlying.Bid
	c.Underlying.BidSize = raw.Underlying.BidSize
	c.Underlying.Change = raw.Underlying.Change
	c.Underlying.Close = raw.Underlying.Close
	c.Underlying.Delayed = raw.Underlying.Delayed
	c.Underlying.Description = raw.Underlying.Description
	c.Underlying.ExchangeName = raw.Underlying.ExchangeName
	c.Underlying.FiftyTwoWeekHigh = raw.Underlying.FiftyTwoWeekHigh
	c.Underlying.FiftyTwoWeekLow = raw.Underlying.FiftyTwoWeekLow
	c.Underlying.HighPrice = raw.Underlying.HighPrice
	c.Underlying.Last = raw.Underlying.Last
	c.Underlying.LowPrice = raw.Underlying.LowPrice
	c.Underlying.Mark = raw.Underlying.Mark
	c.Underlying.MarkChange = raw.Underlying.MarkChange
	c.Underlying.MarkPercentChange = raw.Underlying.MarkPercentChange
	c.Underlying.OpenPrice = raw.Underlying.OpenPrice
	c.Underlying.PercentChange = raw.Underlying.PercentChange
	c.Underlying.QuoteTime = raw.Underlying.QuoteTime
	c.Underlying.Symbol = raw.Underlying.Symbol
	c.Underlying.TotalVolume = raw.Underlying.TotalVolume
	c.Underlying.TradeTime = raw.Underlying.TradeTime
	c.Strategy = raw.Strategy
	c.Interval = raw.Interval
	c.IsDelayed = raw.IsDelayed
	c.IsIndex = raw.IsIndex
	c.DaysToExpiration = raw.DaysToExpiration
	c.InterestRate = raw.InterestRate
	c.UnderlyingPrice = raw.UnderlyingPrice
	c.Volatility = raw.Volatility
	c.Calls = make([]struct {
		ExpDate    time.Time
		DaysTilExp int
		Strikes    []OptionData
	}, len(raw.CallExpDateMap))
	c.Puts = make([]struct {
		ExpDate    time.Time
		DaysTilExp int
		Strikes    []OptionData
	}, len(raw.PutExpDateMap))
	i := 0
	var err error
	for dateStr, v := range raw.CallExpDateMap {
		dateParts := strings.Split(dateStr, ":")
		if c.Calls[i].ExpDate, err = time.Parse("2006-01-02", dateParts[0]); err != nil {
			return err
		}
		if c.Calls[i].DaysTilExp, err = strconv.Atoi(dateParts[1]); err != nil {
			return err
		}
		j := 0
		strikes := make([]OptionData, len(v))
		for _, optionData := range v {
			strikes[j] = optionData[0]
			j++
		}
		c.Calls[i].Strikes = strikes
		i++
	}
	i = 0
	for dateStr, v := range raw.PutExpDateMap {
		dateParts := strings.Split(dateStr, ":")
		if c.Puts[i].ExpDate, err = time.Parse("2006-01-02", dateParts[0]); err != nil {
			return err
		}
		if c.Puts[i].DaysTilExp, err = strconv.Atoi(dateParts[1]); err != nil {
			return err
		}
		j := 0
		strikes := make([]OptionData, len(v))
		for _, optionData := range v {
			strikes[j] = optionData[0]
			j++
		}
		c.Puts[i].Strikes = strikes
		i++
	}
	return nil
}

// OptionChange get the price history for a symbol
// TDAmeritrade API Docs: https://developer.tdameritrade.com/option-chains/apis/get/marketdata/chains
func (s *OptionChainService) OptionChain(ctx context.Context, symbol string, opts *OptionChainOptions) (*OptionChain, *Response, error) {
	u := "marketdata/chains"
	if opts != nil {
		if err := opts.validate(); err != nil {
			return nil, nil, err
		}
		q, err := query.Values(opts)
		if err != nil {
			return nil, nil, err
		}
		q.Add("symbol", symbol)
		u = fmt.Sprintf("%s?%s", u, q.Encode())
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	optionChain := new(OptionChain)
	resp, err := s.client.Do(ctx, req, optionChain)
	if err != nil {
		return nil, resp, err
	}
	if optionChain.Status != "SUCCESS" {
		return optionChain, resp, fmt.Errorf("error: %s", optionChain.Status)
	}
	return optionChain, resp, nil
}

func (opts *OptionChainOptions) validate() error {
	if opts.ContractType != "" {
		if !contains(opts.ContractType, validContractTypes) {
			return fmt.Errorf("invalid contractType, must have the value of one of the following %v", validContractTypes)
		}
	} else {
		opts.ContractType = defaultContractType
	}

	if opts.Strategy != "" {
		if !contains(opts.Strategy, validStrategies) {
			return fmt.Errorf("invalid strategy, must have the value of one of the following %v", validStrategies)
		}
	} else {
		opts.Strategy = defaultStrategy
	}

	if opts.ExpMonth != "" {
		if !contains(opts.ExpMonth, validExpMonths) {
			return fmt.Errorf("invalid expMonth, must have the value of one of the following %v", validExpMonths)
		}
	} else {
		opts.ExpMonth = defaultExpMonth
	}

	if opts.OptionType != "" {
		if !contains(opts.OptionType, validOptionTypes) {
			return fmt.Errorf("invalid optionType, must have the value of one of the following %v", validOptionTypes)
		}
	} else {
		opts.OptionType = defaultOptionType
	}

	return nil
}
