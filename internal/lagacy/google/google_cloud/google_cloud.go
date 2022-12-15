package googleCloud

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	cloud "cloud.google.com/go/storage"
	"context"
	"errors"
	firebase "firebase.google.com/go"
	"fmt"
	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type GoogleCloudInterface interface {
	UploadFile(file multipart.File, path string) (string, error)
	GetBucketName() string
	UploadFileUsers(request *UploadForm) (*ImageStructure, error)
	SignedURL(object string) (string, error)
	FindAllBooks() ([]*Books, error)
	FindOneBooks(id *string) ([]*Books, error)
	CreateBooks(i *Books) error
	UpdateBooks(i *Books) error
	DeleteBooks(request *DeleteUsersForm) error
	UploadImage(request *UploadForm) (*string, error)
}

type GoogleCloudStorage struct {
	cl         *storage.Client
	storage    *cloud.Client
	client     *firestore.Client
	projectID  string
	bucketName string
	basePath   string
}

func NewGoogleCloudStorage(db *mongo.Database) GoogleCloudInterface {
	key := option.WithCredentialsFile("internal/env/firebase_secret_key.json")
	app, err := firebase.NewApp(context.Background(), nil, key)
	if err != nil {
		return nil
	}
	client, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	storage, err := cloud.NewClient(context.Background(), key)
	if err != nil {
		log.Fatalln(err)
	}

	return &GoogleCloudStorage{
		client:  client,
		storage: storage,
	}
}

func (g *GoogleCloudStorage) UploadFile(file multipart.File, path string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	obj := g.basePath + path

	wc := g.cl.Bucket(g.bucketName).Object(obj).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %v", err)
	}

	err := g.cl.Bucket(g.bucketName).Object(obj).ACL().Set(ctx, storage.AllUsers, storage.RoleReader)
	if err != nil {
		return "", err
	}

	rObj := g.cl.Bucket(g.bucketName).Object(obj)
	return fmt.Sprintf("%s/%s/%s", viper.GetString("gcs.baseURL"), rObj.BucketName(), rObj.ObjectName()), nil
}

func (g *GoogleCloudStorage) GetBucketName() string {
	return g.bucketName
}

func (g *GoogleCloudStorage) SignedURL(object string) (string, error) {
	bucket := Bucket
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(30 * time.Minute),
	}
	url, err := g.storage.Bucket(bucket).SignedURL(object, opts)
	if err != nil {
		return "", err
	}
	return url, nil
}

func (g *GoogleCloudStorage) FindAllBooks() ([]*Books, error) {
	BooksData := []*Books{}
	iter := g.client.Collection("books").Documents(context.Background())
	for {
		BookData := &Books{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, errors.New("something wrong, please try again")
		}
		mapstructure.Decode(doc.Data(), &BookData)
		BooksData = append(BooksData, BookData)
	}
	return BooksData, nil
}

func (g *GoogleCloudStorage) FindOneBooks(id *string) ([]*Books, error) {
	BooksData := []*Books{}
	iter := g.client.Collection("books").Where("id", "==", id).Documents(context.Background())
	for {
		BookData := &Books{}
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		mapstructure.Decode(doc.Data(), &BookData)
		BooksData = append(BooksData, BookData)
	}
	return BooksData, nil
}
func (g *GoogleCloudStorage) CreateBooks(i *Books) error {
	uid := uuid.New()
	log.Println("uid:", uid)
	splitID := strings.Split(uid.String(), "-")
	log.Println("splitID:", splitID)
	id := splitID[0] + splitID[1] + splitID[2] + splitID[3] + splitID[4]
	log.Println("id:", id)
	i.ID = id
	_, _, err := g.client.Collection("books").Add(context.Background(), i)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil
}

func (g *GoogleCloudStorage) UpdateBooks(i *Books) error {
	var docID string

	iter := g.client.Collection("books").Where("id", "==", i.ID).Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		docID = doc.Ref.ID
		log.Println("docID:", docID)
	}

	_, err := g.client.Collection("books").Doc(docID).Set(context.Background(), i)
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil
}

func (g *GoogleCloudStorage) DeleteBooks(request *DeleteUsersForm) error {
	var docID string

	iter := g.client.Collection("books").Where("id", "==", request.ID).Documents(context.Background())
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		docID = doc.Ref.ID
		log.Println("docID:", docID)
	}

	_, err := g.client.Collection("books").Doc(docID).Delete(context.Background())
	if err != nil {
		log.Printf("An error has occurred: %s", err)
	}
	return nil

}

func CreateImageUrl(imagePath string, bucket string, ctx context.Context, client *firestore.Client) (*string, error) {
	imageStructure := ImageStructure{
		ImageName: imagePath,
		URL:       "https://storage.cloud.google.com/" + bucket + "/" + imagePath,
	}

	_, _, err := client.Collection("images").Add(ctx, imageStructure)
	if err != nil {
		return nil, err
	}
	log.Println("imageStructure.ImageName:", imageStructure.ImageName)
	log.Println("imageStructure.URL:", imageStructure.URL)
	return &imageStructure.URL, nil
}
func CreateImageUrlUsers(imagePath string, bucket string, ctx context.Context, client *firestore.Client) (*ImageStructure, error) {
	imageStructure := ImageStructure{
		ImageName: imagePath,
		URL:       "https://storage.cloud.google.com/" + bucket + "/" + imagePath,
	}

	// _, _, err := client.Collection("images").Add(ctx, imageStructure)
	// if err != nil {
	// 	return nil, err
	// }
	log.Println("imageStructure.ImageName:", imageStructure.ImageName)
	log.Println("imageStructure.URL:", imageStructure.URL)
	return &imageStructure, nil
}
func (g *GoogleCloudStorage) UploadImage(request *UploadForm) (*string, error) {
	imagePath := fmt.Sprintf("tests/%s-%s.jpeg", uuid.New().String(), strconv.FormatInt(time.Now().UnixNano(), 10))
	bucket := Bucket
	src, err := request.File.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	log.Println("imagePath:", imagePath)
	wc := g.storage.Bucket(bucket).Object(imagePath).NewWriter(context.Background())
	if _, err = io.Copy(wc, src); err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	_, err = CreateImageUrl(imagePath, bucket, context.Background(), g.client)
	if err != nil {
		return nil, err
	}
	objectHandle := g.storage.Bucket(bucket).Object(imagePath)
	objectName := objectHandle.ObjectName()
	return &objectName, nil
}

func (g *GoogleCloudStorage) UploadFileUsers(request *UploadForm) (*ImageStructure, error) {
	imagePath := fmt.Sprintf("users/%s-%s.jpeg", uuid.New().String(), strconv.FormatInt(time.Now().UnixNano(), 10))
	bucket := Bucket
	src, err := request.File.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	log.Println("imagePath:", imagePath)
	wc := g.storage.Bucket(bucket).Object(imagePath).NewWriter(context.Background())
	if _, err = io.Copy(wc, src); err != nil {
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	imageStructure, err := CreateImageUrlUsers(imagePath, bucket, context.Background(), g.client)
	if err != nil {
		return nil, err
	}
	// objectHandle := g.storage.Bucket(bucket).Object(imagePath)
	// objectName := objectHandle.ObjectName()
	return imageStructure, nil
}
