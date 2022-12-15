package post_topic

import (
	"net/http"

	"github.com/khotchapan/KonLakRod-api/internal/core/mongodb"
	postTopic "github.com/khotchapan/KonLakRod-api/internal/core/mongodb/post_topic"
	"github.com/khotchapan/KonLakRod-api/internal/middleware"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	service ServiceInterface
}

func NewHandler(service ServiceInterface) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) CreatePostTopic(c echo.Context) error {
	request := &CreatePostTopicForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.service.CreatePostTopic(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulCreated())
}

func (h *Handler) UpdatePostTopic(c echo.Context) error {
	request := &UpdatePostTopicForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.service.UpdatePostTopic(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulOK())
}

func (h *Handler) DeletePostTopic(c echo.Context) error {
	request := &DeletePostTopicForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err := h.service.DeletePostTopic(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response := &mongodb.Response{}
	return c.JSON(http.StatusOK, response.SuccessfulOK())
}

func (h *Handler) GetOnePostTopic(c echo.Context) error {
	request := &GetOneTopicForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response, err := h.service.FindOnePostTopic(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}

func (h *Handler) GetAllPostTopic(c echo.Context) error {
	request := &postTopic.GetAllPostTopicForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	// uid := c.Request().Header.Get("UserID")
	// log.Println("uid:",uid)
	response, err := h.service.FindAllPostTopic(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}
