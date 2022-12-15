package guest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/khotchapan/KonLakRod-api/internal/core/bcrypt"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	"github.com/khotchapan/KonLakRod-api/internal/handlers/token"
	"github.com/labstack/echo/v4"
)

type ServiceInterface interface {
	LoginUsers(c echo.Context, request *LoginUsersForm) (*entities.TokenResponse, error)
	HealthCheck(c echo.Context) (interface{}, error)
}

type Service struct {
	con          *connection.Connection
	collection   *connection.Collection
	tokenService token.ServiceInterface
}

func NewService(app, collection context.Context) *Service {
	return &Service{
		con:          connection.GetConnect(app, connection.ConnectionInit),
		collection:   connection.GetCollection(collection, connection.CollectionInit),
		tokenService: token.NewService(app, collection),
	}
}

func (s *Service) LoginUsers(c echo.Context, request *LoginUsersForm) (*entities.TokenResponse, error) {
	us := &entities.Users{}
	err := s.collection.Users.FindOneByUserName(request.Username, us)
	if err != nil {
		return nil, err
	}
	if !bcrypt.ComparePassword(*request.Password, us.PasswordHash) {
		return nil, errors.New("password is incorrect")
	}
	tokenDetails, err := s.tokenService.GenerateTokensAndSetDatabase(c, us)
	if err != nil {
		return nil, err
	}

	//s.con.GCS.CreateBooks()
	token := &entities.TokenResponse{
		AccessToken:  &tokenDetails.AccessToken,
		RefreshToken: &tokenDetails.RefreshToken,
	}
	// json, err := json.Marshal(map[string]string{"some": "value"})
	// if err != nil {
	// 	return nil, err
	// }

	//s.con.Redis.Test("name", json, (60)*time.Second)
	// s.con.Redis.CheckPingPong()
	// s.con.Redis.SetKey("name", json, (10)*time.Second)
	// s.con.Redis.SetKey("name3", json, (10)*time.Second)
	//log.Println("=========================================")
	return token, nil
}

type TEST struct {
	TimeStamp string
}

func (s *Service) HealthCheck(c echo.Context) (interface{}, error) {
	now := time.Now().Format("Jan _2 15:04:05.000000")
	data := &TEST{TimeStamp: "last check: " + now}

	cachedata, err := s.con.Redis.GetKey("HealthCheck")
	if cachedata != nil {
		fmt.Println("=========================")
		fmt.Println("cache hit", *cachedata, err)
		fmt.Println("=========================")

		m := &TEST{}
		json.Unmarshal([]byte(*cachedata), m)
		return m, nil

	}
	fmt.Println("=========================")
	fmt.Println("cache not hit", cachedata, err)
	fmt.Println("=========================")

	s.collection.Tokens.Create(data)
	json, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	s.con.Redis.SetKey("HealthCheck", []byte(json), (10)*time.Second)
	return data, nil
}
