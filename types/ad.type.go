package types

type Ad struct {
	Id       uint64 `json:"id" gorm:"primaryKey,index,autoIncrement"`
	Type     string `json:"type" gorm:"not null"`
	ImageURL string `json:"image_url" gorm:"not null"`
	Text     string `json:"text" gorm:"not null"`
	Link     string `json:"link" gorm:"not null"`
}

type AdCreateSchema struct {
	Type     string `json:"type" validate:"required,oneof=banner text"`
	ImageURL string `json:"image_url" validate:"required,url"`
	Text     string `json:"text" validate:"required"`
	Link     string `json:"link" validate:"required,url"`
}

type AdUpdateSchema struct {
	Type     string `json:"type" validate:"required,oneof=banner text"`
	ImageURL string `json:"image_url" validate:"required,url"`
	Text     string `json:"text" validate:"required"`
	Link     string `json:"link" validate:"required,url"`
}
