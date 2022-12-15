package entities

import (
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PostReply struct {
	mongodb.Model `bson:",inline"`
	TopicRefId    *primitive.ObjectID `json:"topic_ref_id" bson:"topic_ref_id,omitempty"`
	Answer        *string             `json:"answer" bson:"answer,omitempty"`
}
