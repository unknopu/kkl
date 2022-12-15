package token

type RefreshTokenForm struct {
	RefreshToken *string `json:"refresh_token" form:"refresh_token" validate:"required"`
}
