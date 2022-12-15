package post_topic

import (
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "post_topic"
)

type RepoInterface interface {
	Create(i interface{}) error
	Update(i interface{}) error
	Delete(i interface{}) error
	FindOneByObjectID(oid *primitive.ObjectID, i interface{}) error
	FindOneByID(id string, i interface{}) error
	FindAllPostTopic(f *GetAllPostTopicForm) (*mongodb.Page, error)
}

type Repo struct {
	mongodb.Repo
}

func NewRepo(db *mongo.Database) *Repo {
	return &Repo{
		Repo: mongodb.Repo{
			Collection: mongodb.
				DB(db).
				Collection(collectionName),
		},
	}
}

func (r *Repo) FindAllPostTopic(f *GetAllPostTopicForm) (*mongodb.Page, error) {
	var filterElements primitive.D
	filterElements = append(filterElements, primitive.E{})
	filterElements = append(filterElements, primitive.E{
		Key: "deleted_at", Value: primitive.M{
			"$exists": false,
		},
	})
	//if f.Name != nil {
	//filterElements = append(filterElements, primitive.E{Key: "$or", Value: primitive.A{primitive.M{"detail.name": primitive.Regex{Pattern: *f.Query, Options: "ig"}}}})
	//}
	pipeline := []primitive.M{
		{"$match": filterElements},
		// {"$lookup": primitive.M{
		// 	"from":         "drug_categories",
		// 	"localField":   "category",
		// 	"foreignField": "_id",
		// 	"as":           "drug_category_docs",
		// }},
	}
	postTopicResponse := []*PostTopicResponse{}
	response, err := r.Aggregate(pipeline, &postTopicResponse, &f.PageQuery)
	if err != nil {
		return nil, err
	}
	return response, err
}
