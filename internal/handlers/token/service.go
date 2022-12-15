package token

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceInterface interface {
	GenerateTokensAndSetDatabase(c echo.Context, u *entities.Users) (*entities.TokenDetails, error)
	//createJWTToken(c echo.Context, u *entities.Users) (*entities.TokenDetails, error)
	RefreshTokenAndSetDatabase(c echo.Context, request *RefreshTokenForm) (*entities.TokenResponse, error)
	//CreateAuthentication(c echo.Context,token *entities.TokenDetails)( *entities.TokenDetails,error)
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

func (s *Service) GenerateTokensAndSetDatabase(c echo.Context, u *entities.Users) (*entities.TokenDetails, error) {
	token, err := s.GenerateTokens(u)
	if err != nil {
		return nil, err
	}

	return token, nil
}
func (s *Service) CreateAuthentication(token *entities.TokenDetails) error {

	at := time.Unix(token.AccessTokenExpiresAt, 0)
	rt := time.Unix(token.RefreshTokenExpiresAt, 0)
	now := time.Now()
	// json, err := json.Marshal(map[string]string{"some2": "value2"})
	// if err != nil {
	// 	return err
	// }
	//s.con.Redis.SetKey("name2", json, (10)*time.Second)
	s.con.Redis.SetKey(token.AccessTokenId, token.UserID.Hex(), time.Duration(at.Sub(now)))
	s.con.Redis.SetKey(token.RefreshTokenId, token.UserID.Hex(), time.Duration(rt.Sub(now)))
	return nil
}

func (s *Service) GenerateTokens(u *entities.Users) (*entities.TokenDetails, error) {
	now := time.Now()
	tokenDetails := &entities.TokenDetails{}
	tokenDetails.IssuedAt = now.Unix()
	tokenDetails.AccessTokenExpiresAt = now.Add(time.Hour * 1).Unix()
	tokenDetails.RefreshTokenExpiresAt = time.Now().Add(time.Hour * 24 * 14).Unix()
	tokenDetails.AccessTokenId = uuid.New().String()
	tokenDetails.RefreshTokenId = uuid.New().String()

	claims := &coreMiddleware.Claims{}
	claims.Subject = "access_token"
	claims.Issuer = "KonLakRod"
	claims.IssuedAt = tokenDetails.IssuedAt
	claims.ExpiresAt = tokenDetails.AccessTokenExpiresAt
	claims.Id = tokenDetails.AccessTokenId
	claims.UserID = &u.ID
	claims.Roles = u.Roles
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	accessToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	tokenDetails.AccessToken = accessToken

	//====================================================================
	rtToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := rtToken.Claims.(jwt.MapClaims)
	rtClaims["sub"] = "refresh_token"
	rtClaims["iss"] = "KonLakRod"
	rtClaims["iat"] = tokenDetails.IssuedAt
	rtClaims["exp"] = tokenDetails.RefreshTokenExpiresAt
	rtClaims["jti"] = tokenDetails.RefreshTokenId
	rtClaims["user_id"] = &u.ID
	//rtClaims["user_id"] = &u.ID
	//rtClaims["roles"] = u.Roles
	refreshToken, err := rtToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	tokenDetails.RefreshToken = refreshToken
	tokenDetails.UserID = &u.ID
	tokenDetails.Roles = u.Roles
	//====================================================================
	// tk := &entities.Tokens{RefreshToken: tokenDetails.RefreshToken,
	// 	UserRefId: &u.ID}

	// err = s.collection.Tokens.Create(tk)
	// if err != nil {
	// 	return nil, err
	// }
	err = s.CreateAuthentication(tokenDetails)
	if err != nil {
		return nil, err
	}
	return tokenDetails, nil
}

func (s *Service) RefreshTokenAndSetDatabase(c echo.Context, request *RefreshTokenForm) (*entities.TokenResponse, error) {
	token, err := s.verifyToken(*request.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid token or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token or expired token")
		//return nil, echo.ErrUnauthorized
	}
	log.Println("claims:", claims)

	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, errors.New("invalid JWT Payload")
	}
	log.Println("userID:", userID)
	objUserID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.New("invalid JWT Payload")
	}
	RefreshTokenId, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("invalid JWT Payload")
	}
	//============================================================================
	//tk := &entities.Tokens{}
	deleted, err := s.con.Redis.Delete(RefreshTokenId)
	if err != nil || deleted == 0 {
		return nil, errors.New("Invalid token or expired token")
	}

	//============================================================================
	/*log.Println("request.RefreshToken:", *request.RefreshToken)
	err = s.collection.Tokens.FindOneByRefreshToken(request.RefreshToken, tk)
	log.Println("tk:", tk)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("error no documents")
		}
		return nil, err
	}
	err = s.collection.Tokens.Delete(tk)
	if err != nil {
		return nil, err
	}*/
	us := &entities.Users{}
	//err = s.collection.Users.FindOneByObjectID(tk.UserRefId, us)
	err = s.collection.Users.FindOneByObjectID(&objUserID, us)
	if err != nil {
		return nil, err
	}
	log.Println("us:", us)
	tokenDetails, err := s.GenerateTokens(us)
	if err != nil {
		return nil, err
	}
	tokenResponse := &entities.TokenResponse{
		AccessToken:  &tokenDetails.AccessToken,
		RefreshToken: &tokenDetails.RefreshToken,
	}
	return tokenResponse, nil

}

func (s *Service) verifyToken(tokenStr string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key.
	// The latter is especially useful if you use multiple keys for your application.
	// The standard is to use 'kid' in the head of the token to identify
	// which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte("secret"), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}
func (s *Service) createJWTTokenTest(c echo.Context, u *entities.Users) (*entities.TokenResponse, error) {
	rto, err := s.createRefreshToken(u)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	claims := &coreMiddleware.Claims{}
	claims.Subject = "access_token"
	claims.Issuer = "KonLakRod"
	claims.IssuedAt = now.Unix()
	claims.ExpiresAt = now.Add(time.Hour * 24).Unix()
	claims.Id = uuid.New().String()
	claims.UserID = &u.ID
	claims.Roles = u.Roles
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	td := &entities.TokenDetailsTest{}
	td.AtExpires = time.Now().Add(time.Hour * 2).Unix()
	td.AccessUuid = uuid.New().String()
	td.RtExpires = time.Now().Add(time.Minute * 24 * 7).Unix()
	td.RefreshUuid = uuid.New().String()
	td.AccessToken, err = token.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}
	//====================================================================
	rtToken := jwt.New(jwt.SigningMethodHS256)
	rtClaims := rtToken.Claims.(jwt.MapClaims)
	rtClaims["id"] = &u.ID
	rtClaims["sub"] = "refresh_token"
	rtClaims["exp"] = td.RtExpires
	rtClaims["jti"] = td.RefreshUuid
	td.RefreshToken, err = rtToken.SignedString([]byte("secret"))
	if err != nil {
		return nil, err
	}

	//====================================================================
	tk := &entities.Tokens{RefreshToken: td.RefreshToken,
		UserRefId: &u.ID}

	err = s.collection.Tokens.Create(tk)
	if err != nil {
		return nil, err
	}

	tkr := &entities.TokenResponse{
		AccessToken:  &t,
		RefreshToken: &rto,
		// AccessTokenTest:  &td.AccessToken,
		// RefreshTokenTest: &td.RefreshToken,
	}

	return tkr, nil
}
func (s *Service) createRefreshToken(u *entities.Users) (string, error) {
	rts := fmt.Sprintf("%d%s", u.ID, time.Now().String())
	h := sha1.New()
	_, err := h.Write([]byte(rts))
	if err != nil {
		return "", err
	}
	res := hex.EncodeToString(h.Sum(nil))
	//log.Println("EncodeToString:", res)
	return res, nil
}
