package types

type Price struct {
	Id       uint    `json:"id" gorm:"primaryKey,autoIncrement"`
	Contract string  `json:"contract"`
	Ticker   string  `json:"ticker"`
	Market   string  `json:"market"`
	Price    float64 `json:"price"`
	Volume   float64 `json:"volume"`
	Date     int64   `json:"date"`
}

type CreatePriceSchema struct {
	Contract string  `json:"contract" validate:"required"`
	Ticker   string  `json:"ticker" validate:"required"`
	Market   string  `json:"market" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Volume   float64 `json:"volume" validate:"required"`
	Date     int64   `json:"date" validate:"required"`
}
