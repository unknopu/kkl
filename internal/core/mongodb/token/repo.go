package token

import (
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	collectionName = "tokens"
)

type RepoInterface interface {
	FindOneByRefreshToken(refresh_token *string, i interface{}) error
	Delete(i interface{}) error
	Create(i interface{}) error
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

func (r *Repo) FindOneByRefreshToken(refresh_token *string, i interface{}) error {
	err := r.FindOneByPrimitiveM(primitive.M{
		"deleted_at": primitive.M{
			"$exists": false,
		},
		"refresh_token": refresh_token,
		//"user_id": userID,
	}, i)
	if err != nil {
		return err
	}
	return nil
}
