package middleware

import (
	"github.com/labstack/echo/v4"
)

// CustomeMiddleware has the access of request object
// and we can decide if we want to go ahead with the req
func RequiredRoles(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			//do the things
			claims := c.(*CustomContext).GetClaims()
			tokenRole := claims.Roles
			for _, roleRequired := range roles {
				for _, userRole := range tokenRole {
					if roleRequired == userRole {
						return next(c)
					}
				}
			}
			return echo.ErrUnauthorized
		}

	}
}
