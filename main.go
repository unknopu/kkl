package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/khotchapan/KonLakRod-api/internal/core/connection"
	"github.com/khotchapan/KonLakRod-api/internal/core/memory"
	postReply "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_reply"
	postTopic "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_topic"
	tokens "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/token"
	users "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/user"
	coreValidator "github.com/khotchapan/KonLakRod-api/internal/core/validator"
	googleCloud "github.com/khotchapan/KonLakRod-api/internal/lagacy/google/google_cloud"
	coreMiddleware "github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/khotchapan/KonLakRod-api/internal/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx context.Context = context.Background()

func initViper() {

	viper.AddConfigPath("configs")                         // ระบุ path ของ config file
	viper.SetConfigName("config")                          // ชื่อ config file
	viper.AutomaticEnv()                                   // อ่าน value จาก ENV variable
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // แปลง _ underscore ใน env เป็น . dot notation ใน viper
	// read config
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read in viper config:%s", err)
	}

}
func initGoDotEnv() {
	if _, err := os.Stat(".env"); err != nil {
		log.Println("[*] initiate .env file")
		_, err := os.Create(".env")
		if err != nil {
			log.Fatal(err)
		}
	}
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func init() {
	log.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	initGoDotEnv()
	initViper()
	log.Println(viper.Get("app.env"))
	log.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
}
func main() {
	var (
		e             = initEcho()
		dbMonggo, _   = newMongoDB()
		redisDatabase = newRedisPool()
		gcs           = googleCloud.NewGoogleCloudStorage(dbMonggo)
	)
	app := context.WithValue(ctx, connection.ConnectionInit,
		connection.Connection{
			Mongo: dbMonggo,
			GCS:   gcs,
			Redis: memory.New(redisDatabase),
		})
	collection := context.WithValue(ctx, connection.CollectionInit,
		connection.Collection{
			Users:     users.NewRepo(dbMonggo),
			Tokens:    tokens.NewRepo(dbMonggo),
			PostTopic: postTopic.NewRepo(dbMonggo),
			PostReply: postReply.NewRepo(dbMonggo),
		})
	options := &router.Options{
		App:        app,
		Collection: collection,
		Echo:       e,
	}
	router.Router(options)

	//godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = viper.GetString("app.port")
		//port = "80" // Default port if not specified
	}
	address := fmt.Sprintf("%s:%s", "0.0.0.0", port)
	log.Println("address:", address)
	e.Logger.Fatal(e.Start(address))
}

func initEcho() *echo.Echo {
	e := echo.New()
	// e.HideBanner = false
	// e.HidePort = false
	// e.Debug = false
	// e.HideBanner = true
	//Validator
	e.Validator = coreValidator.NewValidator(validator.New())
	// Middleware
	e.Use(coreMiddleware.SetCustomContext)
	e.Use(middleware.Logger())    // Log everything to stdout
	e.Use(middleware.Recover())   // Recover from all panics to always have your server up
	e.Use(middleware.RequestID()) // Generate a request id on the HTTP response headers for identification
	return e
}

func newMongoDB() (*mongo.Database, context.Context) {
	EnvMongoURI := os.Getenv("MONGOURI")
	log.Println("Connecting to:", EnvMongoURI)
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI))
	if err != nil {
		log.Fatal(err)
	}

	contextDatabase, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(contextDatabase)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(contextDatabase, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB:", EnvMongoURI)
	return client.Database("konlakrod"), contextDatabase
}

func newRedisPool() *redis.Client {
	// return pool
	host := getenv("REDIS_URI", "localhost")
	log.Println("REDIS HOST::::::::::", host)
	val := fmt.Sprintf("%s:%s", host, "6379")
	rdb := redis.NewClient(&redis.Options{
		Addr:     val,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("redis error:", err)
	}
	log.Println("redis:", pong)
	return rdb
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
