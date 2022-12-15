package post_reply

import (
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreatePostReplyForm struct {
	TopicRefId *primitive.ObjectID `json:"topic_ref_id" validate:"required"`
	Answer     *string             `json:"answer" validate:"required"`
}

func (f *CreatePostReplyForm) fill(data *entities.PostReply) *entities.PostReply {
	if f.TopicRefId != nil {
		data.TopicRefId = f.TopicRefId
	}
	if f.Answer != nil {
		data.Answer = f.Answer
	}

	return data
}

type UpdatePostReplyForm struct {
	ID         *primitive.ObjectID `json:"id" validate:"required"`
	Answer     *string             `json:"answer"`
}

func (f *UpdatePostReplyForm) Fill(data *entities.PostReply) *entities.PostReply {
	if f.ID != nil {
		data.ID = *f.ID
	}
	if f.Answer != nil {
		data.Answer = f.Answer
	}
	return data
}

type DeletePostReplyForm struct {
	ID *primitive.ObjectID `json:"id" validate:"required"`
}
type GetOneReplyForm struct {
	ID *primitive.ObjectID `param:"id"`
}
