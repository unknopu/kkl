package user

import (
	"mime/multipart"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetOneUsersForm struct {
	ID *primitive.ObjectID `param:"id"`
}

type CreateUsersForm struct {
	FirstName   *string `json:"first_name" validate:"required"`
	LastName    *string `json:"last_name" validate:"required"`
	Username    *string `json:"username" validate:"required"`
	Password    *string `json:"password" validate:"required"`
	Image       *string `json:"image"`
	Email       *string `json:"email"`
	PhoneNumber *string `json:"phone_number"`
}

func (f *CreateUsersForm) fill(data *entities.Users) *entities.Users {
	if f.FirstName != nil {
		data.FirstName = *f.FirstName
	}
	if f.LastName != nil {
		data.LastName = *f.LastName
	}
	if f.Username != nil {
		data.Username = *f.Username
	}
	if f.Image != nil {
		data.Image = *f.Image
	}
	if f.Email != nil {
		data.Email = *f.Email
	}
	if f.PhoneNumber != nil {
		data.PhoneNumber = *f.PhoneNumber
	}
	return data
}

type UpdateUsersForm struct {
	ID          *primitive.ObjectID `json:"id" validate:"required"`
	FirstName   *string             `json:"first_name" validate:"required"`
	LastName    *string             `json:"last_name" validate:"required"`
	Image       *string             `json:"image"`
	Email       *string             `json:"email"`
	PhoneNumber *string             `json:"phone_number"`
}

func (f *UpdateUsersForm) fill(data *entities.Users) *entities.Users {
	if f.ID != nil {
		data.ID = *f.ID
	}
	if f.FirstName != nil {
		data.FirstName = *f.FirstName
	}
	if f.LastName != nil {
		data.LastName = *f.LastName
	}
	if f.Image != nil {
		data.Image = *f.Image
	}
	if f.Email != nil {
		data.Email = *f.Email
	}
	if f.PhoneNumber != nil {
		data.PhoneNumber = *f.PhoneNumber
	}
	return data
}

type DeleteUsersForm struct {
	ID *primitive.ObjectID `json:"id" validate:"required"`
}

type UploadForm struct {
	Path string                `form:"path"`
	Mime string                `form:"mime"`
	File *multipart.FileHeader `form:"file"`
}
