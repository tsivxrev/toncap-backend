package types

import (
	"time"
)

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Ad struct {
	Id       string `json:"id"`        // ad id (uuid_v4)
	Type     string `json:"type"`      // banner or text
	ImageURL string `json:"image_url"` // banner image
	Text     string `json:"text"`      // ad caption
	Link     string `json:"link"`      // ad link
}

type Currency struct {
	Byn float64 `json:"byn"`
	Cny float64 `json:"cny"`
	Eur float64 `json:"eur"`
	Gbp float64 `json:"gbp"`
	Kzt float64 `json:"kzt"`
	Rub float64 `json:"rub"`
	Uah float64 `json:"uah"`
	Usd float64 `json:"usd"`
}

type Price struct {
	Id     int64      `json:"id"`
	Ticker string     `json:"ticker"`
	Market string     `json:"market"`
	Price  float64    `json:"price"`
	Volume float64    `json:"volume"`
	Date   *time.Time `json:"date"`
}

type JettonMarket struct {
	Contract string `json:"contract"`
	Image    string `json:"image"`
	Name     string `json:"name"`
	Template string `json:"template"`
	Type     string `json:"type"`
	URL      string `json:"url"`
}

type Jetton struct {
	Contract string                  `json:"contract"`
	Decimals int                     `json:"decimals"`
	Supply   float64                 `json:"supply"`
	Image    string                  `json:"image"`
	Markets  map[string]JettonMarket `json:"markets"`
	Links    map[string]string       `json:"links"`
	Name     string                  `json:"name"`
	Ticker   string                  `json:"ticker"`
}
