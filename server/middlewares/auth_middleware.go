package middlewares

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"on-air/config"
	"on-air/repository"
	"strings"
)

const Bearer = "bearer"

type Auth struct {
	JWT *config.JWT
}

func (a *Auth) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {

		authorizationHeader := ctx.Request().Header.Get("Authorization")
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header not provided")
			return ctx.JSON(http.StatusUnauthorized, err.Error())
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			return ctx.JSON(http.StatusUnauthorized, err.Error())
		}
		authorizationType := strings.ToLower(fields[0])
		if authorizationType != Bearer {
			err := fmt.Errorf("authorization type %s not supported", authorizationType)
			return ctx.JSON(http.StatusUnauthorized, err.Error())

		}
		accessToken := fields[1]
		payload, err := repository.VerifyToken(a.JWT, accessToken)
		if err != nil {
			return ctx.JSON(http.StatusUnauthorized, err.Error())
		}
		ctx.Request().Header.Set("userId", string(payload.UserID))

		return next(ctx)
	}
}
