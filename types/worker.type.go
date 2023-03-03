package types

type SocialsResponse map[string]string
type ListedContractsResponse map[string][]string
type ListedContracts []string

type ActualPrice struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

type ActualMarket struct {
	Contract string  `json:"contract"`
	Market   string  `json:"market"`
	Ticker   string  `json:"ticker"`
	Price    float64 `json:"price"`
	Volume   float64 `json:"volume"`
}

type ActualResponse struct {
	Price   ActualPrice    `json:"actual"`
	Markets []ActualMarket `json:"markets"`
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
