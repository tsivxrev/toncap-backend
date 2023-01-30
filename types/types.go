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
	Text     string `json:"text"`                                                       // ad caption
	Link     string `json:"link" gorm:"not null" validate:"required,url"`               // ad link
}
