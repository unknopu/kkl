package post_reply

import (
	"context"
	postReply "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_reply"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"github.com/labstack/echo/v4"
)

type ServiceInterface interface {
	CreatePostReply(c echo.Context, request *CreatePostReplyForm) error
	UpdatePostReply(c echo.Context, request *UpdatePostReplyForm) error
	DeletePostReply(c echo.Context, request *DeletePostReplyForm) error
	FindOnePostReply(c echo.Context, request *GetOneReplyForm) (*entities.PostReply, error)
	FindAllPostReply(c echo.Context, request *postReply.GetAllPostReplayForm) (*mongodb.Page, error)
}

type Service struct {
	con        *connection.Connection
	collection *connection.Collection
}

func NewService(app, collection context.Context) *Service {
	return &Service{
		con:        connection.GetConnect(app, connection.ConnectionInit),
		collection: connection.GetCollection(collection, connection.CollectionInit),
	}
}

func (s *Service) CreatePostReply(c echo.Context, request *CreatePostReplyForm) error {
	us := &entities.PostReply{}
	u := request.fill(us)
	err := s.collection.PostReply.Create(u)
	if err != nil {
		return err
	}
	return nil
}


func (s *Service) UpdatePostReply(c echo.Context, request *UpdatePostReplyForm) error {
	pr := &entities.PostReply{}
	err := s.collection.PostReply.FindOneByObjectID(request.ID, pr)
	if err != nil {
		return err
	}
	b := request.Fill(pr)
	err = s.collection.PostReply.Update(b)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeletePostReply(c echo.Context, request *DeletePostReplyForm) error {
	pr := &entities.PostReply{
		Model: mongodb.Model{ID: *request.ID},
	}
	err := s.collection.PostReply.Delete(pr)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) FindOnePostReply(c echo.Context, request *GetOneReplyForm) (*entities.PostReply, error) {
	response := &entities.PostReply{}
	err := s.collection.PostReply.FindOneByObjectID(request.ID, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) FindAllPostReply(c echo.Context, request *postReply.GetAllPostReplayForm) (*mongodb.Page, error) {
	// cc := c.(*middleware.CustomContext)
	// userID := cc.GetClaims().UserID
	// log.Println("userID:", userID)
	// user := c.Get("user").(*jwt.Token)
	// claims := user.Claims.(*coreMiddleware.Claims)
	// log.Println("claims.UserID:", claims.UserID)

	response, err := s.collection.PostReply.FindAllPostReply(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}