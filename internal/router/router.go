package router

import (
	"context"
	"fmt"
	"net/http"
	"path"

	"github.com/golang-jwt/jwt"
	"github.com/khotchapan/KonLakRod-api/internal/entities"
	guestHandler "github.com/khotchapan/KonLakRod-api/internal/handlers/guest"
	postReplyHandler "github.com/khotchapan/KonLakRod-api/internal/handlers/post_reply"
	postTopicHandler "github.com/khotchapan/KonLakRod-api/internal/handlers/post_topic"
	testHandler "github.com/khotchapan/KonLakRod-api/internal/handlers/test"
	tokenHandler "github.com/khotchapan/KonLakRod-api/internal/handlers/token"
	userHandler "github.com/khotchapan/KonLakRod-api/internal/handlers/user"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Options struct {
	App        context.Context
	Collection context.Context
	Echo       *echo.Echo
}

func Router(options *Options) {
	app := options.App
	collection := options.Collection
	e := options.Echo
	//===============================================================================
	// Configure middleware with the custom claims type
	config := middleware.JWTConfig{
		Claims:        &coreMiddleware.Claims{},
		SigningKey:    []byte("secret"),
		SigningMethod: jwt.SigningMethodHS256.Name,
	}
	checkSessionMiddleware := middleware.JWTWithConfig(config)
	requiredUser := coreMiddleware.RequiredRoles(entities.UserRole)
	requiredAdmin := coreMiddleware.RequiredRoles(entities.AdminRole)
	//requiredGarage := coreMiddleware.RequiredRoles(entities.GarageRole)
	//===============================================================================
	//home
	e.GET(path.Join("/"), Version)
	api := e.Group("/v1/api")
	api.GET("", HelloWorld)

	// Unauthenticated route
	api.GET("/", accessible)
	// Restricted group
	r := api.Group("/restricted", checkSessionMiddleware, requiredAdmin)
	{
		//api.Use(checkSessionMiddleware)
		r.GET("", restricted)
	}
	//guest
	guestEndpoint := guestHandler.NewHandler(guestHandler.NewService(app, collection))
	guestGroup := api.Group("/guest")
	{
		guestGroup.POST("/login", guestEndpoint.LoginUsers)
		guestGroup.GET("/healthcheck", guestEndpoint.HealthCheck)
	}
	//token
	tokenEndpoint := tokenHandler.NewHandler(tokenHandler.NewService(app, collection))
	tokensGroup := api.Group("/tokens")
	{
		tokensGroup.POST("/refreshToken", tokenEndpoint.RefreshToken)
	}

	//user
	usersEndpoint := userHandler.NewHandler(userHandler.NewService(app, collection))
	usersGroup := api.Group("/users")
	usersGroup.GET("/me", usersEndpoint.GetMe, checkSessionMiddleware, requiredUser)
	usersGroup.GET("", usersEndpoint.GetAllUsers)
	usersGroup.GET("/:id", usersEndpoint.GetOneUsers)
	usersGroup.POST("", usersEndpoint.CreateUsers)
	usersGroup.PUT("", usersEndpoint.UpdateUsers)
	usersGroup.DELETE("", usersEndpoint.DeleteUsers)
	usersGroup.POST("/image/upload", usersEndpoint.UploadFileUsers)

	//post topic
	postTopicEndpoint := postTopicHandler.NewHandler(postTopicHandler.NewService(app, collection))
	postTopicGroup := api.Group("/post-topic")
	postTopicGroup.GET("/:id", postTopicEndpoint.GetOnePostTopic)
	postTopicGroup.GET("", postTopicEndpoint.GetAllPostTopic)
	postTopicGroup.POST("", postTopicEndpoint.CreatePostTopic)
	postTopicGroup.PUT("", postTopicEndpoint.UpdatePostTopic)
	postTopicGroup.DELETE("", postTopicEndpoint.DeletePostTopic)

	//post reply
	postReplyEndpoint := postReplyHandler.NewHandler(postReplyHandler.NewService(app, collection))
	postReplyGroup := api.Group("/post-reply")
	postReplyGroup.GET("/:id", postReplyEndpoint.GetOnePostReply)
	postReplyGroup.GET("", postReplyEndpoint.GetAllPostReply)
	postReplyGroup.POST("", postReplyEndpoint.CreatePostReply)
	postReplyGroup.PUT("", postReplyEndpoint.UpdatePostReply)
	postReplyGroup.DELETE("", postReplyEndpoint.DeletePostReply)

	// test zone
	testEndpoint := testHandler.NewHandler(testHandler.NewService(app, collection))
	testGroup := api.Group("/tests")
	testGroup.GET("/google-cloud/books", testEndpoint.GetFile)
	testGroup.GET("/google-cloud/books/:id", testEndpoint.GetOneGoogleCloudBooks)
	testGroup.POST("/google-cloud/books", testEndpoint.CreateGoogleCloudBooks)
	testGroup.PUT("/google-cloud/books", testEndpoint.UpdateBooks)
	testGroup.DELETE("/google-cloud/books", testEndpoint.DeleteBooks)
	testGroup.POST("/google-cloud/image/upload", testEndpoint.UploadImage)
}
func Version(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{"version": 4})
}
func HelloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! KonLakRod&")
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*coreMiddleware.Claims)
	name := claims.Subject
	userID := claims.UserID
	text := fmt.Sprintf("Welcome %s : %v \n", name, userID)
	return c.String(http.StatusOK, text)
}
