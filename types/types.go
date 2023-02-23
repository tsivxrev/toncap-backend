package types

type HTTPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type TokenData struct {
	Id        string `json:"id"`
	UserId    int    `json:"user_id" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=default service"`
	ExpiresIn int64  `json:"expires_in" validate:"required,min=0"`
}

type Ad struct {
	Id       uint64 `json:"id" gorm:"primaryKey,index,autoIncrement"`                   // ad id
	Type     string `json:"type" gorm:"not null" validate:"required,oneof=banner text"` // banner or text
	ImageURL string `json:"image_url" gorm:"not null" validate:"required,url"`          // banner image
	Text     string `json:"text" gorm:"not null" validate:"required"`                   // ad caption
	Link     string `json:"link" gorm:"not null" validate:"required,url"`               // ad link
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
	Id       uint    `json:"id" gorm:"primaryKey,autoIncrement"`
	Contract string  `json:"contract" validate:"required"`
	Ticker   string  `json:"ticker" validate:"required"`
	Market   string  `json:"market" validate:"required"`
	Price    float64 `json:"price" validate:"required"`
	Volume   float64 `json:"volume" validate:"required"`
	Day      int     `json:"day" validate:"required"`
	Month    int     `json:"month" validate:"required"`
	Year     int     `json:"year" validate:"required"`
}

type ActualResponsePriceVolume struct {
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}

type ActualResponseMarket struct {
	Contract string  `json:"contract"`
	Market   string  `json:"market"`
	Ticker   string  `json:"ticker"`
	Price    float64 `json:"price"`
	Volume   float64 `json:"volume"`
}

type ActualResponse struct {
	Actual  ActualResponsePriceVolume `json:"actual"`
	Markets []ActualResponseMarket    `json:"markets"`
}

type Graph struct {
	Date   int64   `json:"date"`
	Price  float64 `json:"price"`
	Volume float64 `json:"volume"`
}
