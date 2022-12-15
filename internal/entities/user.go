package entities

import "github.com/khotchapan/KonLakRod-api/internal/core/mongodb"

type Users struct {
	mongodb.Model  `bson:",inline"`
	FirstName      string   `json:"firstName" bson:"first_name,omitempty"`
	LastName       string   `json:"lastName" bson:"last_name,omitempty"`
	Image          string   `json:"image" bson:"image,omitempty"`
	Email          string   `json:"email" bson:"email,omitempty"`
	PhoneNumber    string   `json:"phoneNumber" bson:"phone_number,omitempty"`
	Birthday       string   `json:"birthday" bson:"birthday,omitempty"`
	Username       string   `json:"username" bson:"username,omitempty"`
	PasswordHash   string   `json:"password_hash" bson:"password_hash,omitempty"`
	Roles          []string `json:"roles" bson:"roles,omitempty"`
	Activate       bool     `json:"activate" bson:"activate,omitempty"`
	FacebookID     string   `json:"-" bson:"facebook_id,omitempty"`
	FacebookActive bool     `json:"-" bson:"facebook_active,omitempty"`
	GoogleID       string   `json:"-" bson:"google_id,omitempty"`
	GoogleActive   bool     `json:"-" bson:"google_active,omitempty"`
	UserToken      string   `json:"-" bson:"user_token,omitempty"`
	UserSex        string   `json:"userSex" bson:"user_sex,omitempty"`
	//Address        []*primitive.ObjectID `json:"address" bson:"address,omitempty"`
	// AcceptConsent []*primitive.ObjectID `json:"acceptConsent" bson:"accept_consent,omitempty"`
	// HealthInfo    *HealthInfo            `json:"healthInfo" bson:"health_info,omitempty"`
}
