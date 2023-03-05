package types

type ContractMeta struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Symbol      string `json:"symbol"`
	Image       string `json:"image"`

	Socials     SocialsResponse `json:"socials"`
	TotalSupply float64         `json:"total_supply"`
	Decimals    int             `json:"decimals"`
}

type ContractGraph struct {
	Date   int64   `json:"date"`
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

type ContractResponse struct {
	Contract string          `json:"contract"`
	Graph    []ContractGraph `json:"graph"`
	Meta     ContractMeta    `json:"meta"`
	Actual   ActualPrice     `json:"actual"`
	Markets  []ActualMarket  `json:"markets"`
}
