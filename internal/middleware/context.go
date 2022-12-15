package middleware

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomContext struct {
	echo.Context
	Parameters interface{}
}

func SetCustomContext(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := &CustomContext{Context: c}
		return next(cc)
	}
}

//const pathKey = "path"

// BindAndValidate bind and validate form
func (c *CustomContext) BindAndValidate(i interface{}) error {
	if err := c.Bind(i); err != nil {
		return err
	}
	//c.parsePathParams(i)
	if err := c.Validate(i); err != nil {
		return err
	}
	//c.Parameters = i
	return nil
}

// func (c *CustomContext) parsePathParams(form interface{}) {
// 	formValue := reflect.ValueOf(form)
// 	if formValue.Kind() == reflect.Ptr {
// 		formValue = formValue.Elem()
// 	}
// 	t := reflect.TypeOf(formValue.Interface())
// 	for i := 0; i < t.NumField(); i++ {
// 		tag := t.Field(i).Tag.Get(pathKey)
// 		if tag != "" {
// 			fieldName := t.Field(i).Name
// 			paramValue := formValue.FieldByName(fieldName)
// 			if paramValue.IsValid() {
// 				paramValue.Set(reflect.ValueOf(c.Param(tag)))
// 			}
// 		}
// 	}
// }

const (
	pathKey     = "path"
	userContext = "user"
)

// Claims jwt claims
type Claims struct {
	jwt.StandardClaims
	Roles          []string            `json:"roles"`
	UserID         *primitive.ObjectID `json:"user_id"`
}

// GetClaims get user claims
func (c *CustomContext) GetClaims() *Claims {
	user := c.Get(userContext).(*jwt.Token)
	return user.Claims.(*Claims)
}
