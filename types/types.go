package types

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type ValidateErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

type TokenData struct {
	Id        string `json:"id"`
	UserId    int    `json:"user_id" validate:"required"`
	Type      string `json:"type" validate:"required,oneof=user service"`
	ExpiresIn int64  `json:"expires_in" validate:"required,min=0"`
}
