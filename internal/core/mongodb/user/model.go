package user

import (
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
)


type UsersResponse struct {
	mongodb.Model  `bson:",inline"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Image          string `json:"profile_image"`
	Email          string `json:"email,omitempty"`
	PhoneNumber    string `json:"phone_number,omitempty"`
	Birthday       string `json:"birthday"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Activate       bool   `json:"activate"`
	FacebookID     string `json:"facebook_id"`
	FacebookActive bool   `json:"facebook_active"`
	GoogleID       string `json:"google_id"`
	GoogleActive   bool   `json:"google_active"`
	UserToken      string `json:"user_token"`
	UserSex        string `json:"user_sex"`
}
type GetAllUsersForm struct {
	mongodb.PageQuery
	Name *string `query:"name"`
}
