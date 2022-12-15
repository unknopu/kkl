package user

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/khotchapan/KonLakRod-api/internal/core/bcrypt"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

type ServiceInterface interface {
	CallGetMe(c echo.Context) (*entities.Users, error)
	FindAllUsers(c echo.Context, request *user.GetAllUsersForm) (*mongodb.Page, error)
	FindOneUsers(c echo.Context, request *GetOneUsersForm) (*entities.Users, error)
	CreateUsers(c echo.Context, request *CreateUsersForm) error
	UpdateUsers(c echo.Context, request *UpdateUsersForm) error
	DeleteUsers(c echo.Context, request *DeleteUsersForm) error
	UploadFileUsers(c echo.Context, req *googleCloud.UploadForm) (*googleCloud.ImageStructure, error)
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
func (s *Service) CallGetMe(c echo.Context) (*entities.Users, error) {
	response := &entities.Users{}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*coreMiddleware.Claims)
	log.Println("UserID:", claims.UserID)
	userID := claims.UserID
	err := s.collection.Users.FindOneByObjectID(userID, response)
	if err != nil {
		return nil, err
	}
	//err := s.collection.Users.FindOneByID(c.Get("user").(*jwt.User).ID.Hex(), response)
	// if err != nil {
	// 	return nil, errs.New(http.StatusConflict, "10003", "user not found")
	// }

	if response.Image != "" && !strings.Contains(response.Image, "http") {
		url, err := s.con.GCS.SignedURL(response.Image)
		if err != nil {
			return nil, errors.New("can not singed url")
		}
		response.Image = url
	}

	// if response.HealthInfo.Birthday != "" {
	// 	birthday, err := time.Parse("2006-01-02", response.HealthInfo.Birthday)
	// 	if err != nil {
	// 		return nil, errs.NewBadRequest("Can not convert to time", err.Error())
	// 	}

	// 	y, m, d := birthday.Date()
	// 	dob := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	// 	response.HealthInfo.Age = age.Age(dob)
	// }

	return response, nil
}
func (s *Service) FindAllUsers(c echo.Context, request *user.GetAllUsersForm) (*mongodb.Page, error) {
	//objectUserID := &c.Get("user").(*jwt.User).ID
	response, err := s.collection.Users.FindAllUsers(request)
	if err != nil {
		return nil, err
	}
	return response, nil
}
func (s *Service) FindOneUsers(c echo.Context, request *GetOneUsersForm) (*entities.Users, error) {
	response := &entities.Users{}
	err := s.collection.Users.FindOneByObjectID(request.ID, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (s *Service) CreateUsers(c echo.Context, request *CreateUsersForm) error {
	log.Println("TEST:1")
	us := &entities.Users{}
	response := &entities.Users{}
	err := s.collection.Users.FindOneByUserName(request.Username, response)
	log.Println("TEST:1.2")
	if err != nil && err != mongo.ErrNoDocuments {
		log.Println("err:", err)
		// ErrNoDocuments means that the filter did not match any documents in the collection
		//	if err == mongo.ErrNoDocuments {
		//return errors.New("error no documents")
		return err
		//}

	}
	log.Println("TEST:2")
	password, err := bcrypt.GeneratePassword(*request.Password)
	if err != nil {
		//c.Error(err)
		return err
	}
	us.PasswordHash = password
	us.Roles = []string{entities.UserRole}
	u := request.fill(us)
	err = s.collection.Users.Create(u)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateUsers(c echo.Context, request *UpdateUsersForm) error {
	us := &entities.Users{}
	err := s.collection.Users.FindOneByObjectID(request.ID, us)
	if err != nil {
		return err
	}
	u := request.fill(us)
	err = s.collection.Users.Update(u)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) DeleteUsers(c echo.Context, request *DeleteUsersForm) error {
	u := &entities.Users{
		Model: mongodb.Model{ID: *request.ID},
	}
	err := s.collection.Users.Delete(u)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UploadFileUsers(c echo.Context, request *googleCloud.UploadForm) (*googleCloud.ImageStructure, error) {
	imageStructure, err := s.con.GCS.UploadFileUsers(request)
	if err != nil {
		return nil, err
	}
	//log.Println("objectName:", *objectName)
	signedUrl, err := s.con.GCS.SignedURL(imageStructure.ImageName)
	if err != nil {
		return nil, err
	}
	imageStructure.URL = signedUrl
	return imageStructure, nil
}
