package entities

import "github.com/khotchapan/KonLakRod-api/internal/core/mongodb"

type PostTopic struct {
	mongodb.Model `bson:",inline"`
	Title         *string   `json:"title" bson:"title,omitempty"`
	Tag           []*string `json:"tag" bson:"tag,omitempty"`
}
