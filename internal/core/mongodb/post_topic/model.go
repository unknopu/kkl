package post_topic

import "github.com/khotchapan/KonLakRod-api/internal/core/mongodb"

type GetAllPostTopicForm struct {
	mongodb.PageQuery
}
type PostTopicResponse struct {
	mongodb.Model `bson:",inline"`
	Title         string    `json:"title" bson:"title,omitempty"`
	Tag           []*string `json:"tag" bson:"tag,omitempty"`
}
