package post_topic

import (
	"log"

	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePostTopicForm struct {
	Title *string   `json:"title" validate:"required"`
	Tag   []*string `json:"tag" validate:"required"`
}

func (f *CreatePostTopicForm) fill(data *entities.PostTopic) *entities.PostTopic {
	if f.Title != nil {
		data.Title = f.Title
	}
	if f.Tag != nil {
		data.Tag = f.Tag
	}
	log.Println("data:", *data)
	return data
}

type UpdatePostTopicForm struct {
	ID    *primitive.ObjectID `json:"id" validate:"required"`
	Title *string             `json:"title"`
	Tag   []*string           `json:"tag"`
}

func (f *UpdatePostTopicForm) Fill(data *entities.PostTopic) *entities.PostTopic {
	if f.ID != nil {
		data.ID = *f.ID
	}
	if f.Title != nil {
		data.Title = f.Title
	}
	if f.Tag != nil {
		data.Tag = f.Tag
	}
	return data
}

type DeletePostTopicForm struct {
	ID *primitive.ObjectID `json:"id" validate:"required"`
}

type GetOneTopicForm struct {
	ID *primitive.ObjectID `param:"id"`
}
