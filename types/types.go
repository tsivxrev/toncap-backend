package types

import "time"

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Ad struct {
	Type     string `json:"type"`      // banner or text
	ImageURL string `json:"image_url"` // banner image
	Text     string `json:"text"`      // ad caption
	Link     string `json:"link"`      // ad link
}

type AppStorage struct {
	ServerAddress   string       `json:"server_address"`
	ContractAddress string       `json:"contract_address"`
	AdsFilePath     string       `json:"ads_file_path"`
	Ads             []Ad         `json:"ads"`
	ExchangeRate    ExchangeRate `json:"exchange_rate"`
}

type Task struct {
	Func     func(chan bool)
	Interval time.Duration
	Stop     chan bool
}

type ExchangeRate struct {
	Byn float64 `json:"byn"`
	Cny float64 `json:"cny"`
	Eur float64 `json:"eur"`
	Gbp float64 `json:"gbp"`
	Kzt float64 `json:"kzt"`
	Rub float64 `json:"rub"`
	Uah float64 `json:"uah"`
	Usd float64 `json:"usd"`
}
