package mongodb

import (
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EmptyStruct struct {
}

type Model struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty" firestore:"-"`
	DeletedAt *time.Time         `json:"-" bson:"deleted_at,omitempty" firestore:"-"`
	DeletedBy primitive.ObjectID `json:"-" bson:"deleted_by,omitempty" firestore:"-"`
	UpdatedAt *time.Time         `json:"updated_at" bson:"updated_at,omitempty" firestore:"-"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at,omitempty" firestore:"-"`
}
type ModelInterface interface {
	GetID() primitive.ObjectID
	SetID(id primitive.ObjectID)
	Stamp()
	UpdateStamp()
	UpdateDeletedStamp()
}

func (model *Model) SetID(id primitive.ObjectID) {
	model.ID = id
}

func (model *Model) GetID() primitive.ObjectID {
	return model.ID
}

func (model *Model) Stamp() {
	timeNow := time.Now()
	model.UpdatedAt = &timeNow
	model.CreatedAt = timeNow
}

func (model *Model) UpdateStamp() {
	timeNow := time.Now()
	model.UpdatedAt = &timeNow
	if model.CreatedAt.IsZero() {
		model.CreatedAt = timeNow
	}
}

func (model *Model) UpdateDeletedStamp() {
	timeNow := time.Now()
	model.DeletedAt = &timeNow
}

type PageQuery struct {
	Page int64 `query:"page"`
	Size int64 `query:"size"`
}

// DefaultPageSize default page size
var DefaultPageSize int64 = 20
var DefaultPage int64 = 1

// GetPage get page
func (form *PageQuery) GetPage() int64 {
	if form.Page > 0 {
		return form.Page
	}
	return DefaultPage
}

// PageSize page size
func (form *PageQuery) GetSize() int64 {
	if form.Size > 0 {
		return form.Size
	}
	return DefaultPageSize
}

type Response struct {
	Status     string `json:"status"`
	StatusCode int    `json:"statusCode"`
}

func (form *Response) SuccessfulOK() *Response {
	form.Status = http.StatusText(http.StatusOK)
	form.StatusCode = http.StatusOK
	return form
}

func (form *Response) SuccessfulCreated() *Response {
	form.Status = http.StatusText(http.StatusCreated)
	form.StatusCode = http.StatusCreated
	return form
}

type CountDocument struct {
	Count int64 `json:"count,omitempty" bson:"count,omitempty"`
}

func ParseStartDateEndDate(tempDate *time.Time) (*time.Time, *time.Time) {
	startDate := time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day(), 0, 0, 0, 0, time.UTC)
	endDate := time.Date(tempDate.Year(), tempDate.Month(), tempDate.Day(), 23, 59, 59, 0, time.UTC)
	return &startDate, &endDate
}
