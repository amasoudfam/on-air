package middleware

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func DbMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("_db", db)
			return next(c)
		}
	}
}
