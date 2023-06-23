package middlewares

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"on-air/config"
	"on-air/repository"
	"strings"
)

const (
	AuthHeader         = "Authorization"
	Bearer             = "bearer"
	UserIdContextField = "user_id"
)

type Auth struct {
	JWT *config.JWT
}

func (a *Auth) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		authHeader := ctx.Request().Header.Get(AuthHeader)
		if authHeader == "" {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		authParams := strings.Split(authHeader, "")
		if len(authParams) < 2 {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		authType := strings.ToLower(authParams[0])
		if authType != Bearer {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		accessToken := authParams[1]
		payload, err := repository.VerifyToken(a.JWT, accessToken)
		if err != nil {
			return ctx.NoContent(http.StatusUnauthorized)
		}

		ctx.Set(UserIdContextField, payload.UserID)
		return next(ctx)
	}
}
