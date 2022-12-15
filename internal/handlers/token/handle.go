package token

import (
	"net/http"
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
func (h *Handler) RefreshToken(c echo.Context) error {
	request := &RefreshTokenForm{}
	cc := c.(*middleware.CustomContext)
	if err := cc.BindAndValidate(request); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	response, err := h.service.RefreshTokenAndSetDatabase(c, request)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
 
	return c.JSON(http.StatusOK, response)
}
