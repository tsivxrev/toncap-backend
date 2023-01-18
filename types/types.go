package types

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
	Id       uint    `json:"id" gorm:"primaryKey,autoIncrement"`
	Contract string  `json:"contract" binding:"required"`
	Ticker   string  `json:"ticker" binding:"required"`
	Market   string  `json:"market" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
	Volume   float64 `json:"volume" binding:"required"`
	Date     int64   `json:"date" gorm:"autoCreateTime"`
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

type TokenData struct {
	Id        string `json:"id"`
	UserId    int    `json:"user_id" binding:"required"`
	Type      string `json:"type" binding:"required,oneof=default service"`
	ExpiresIn int64  `json:"expires_in" binding:"required,gte=0"`
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
