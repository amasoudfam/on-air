package api

import (
	"net/http"
	"on-air/schemas"
	"on-air/services"
	"on-air/utils"

	"github.com/labstack/echo/v4"
)

func (server *Server) GetAuthToken(ctx echo.Context) error {
	user := new(schemas.LoginUserRequest)
	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}
	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	dbUser, err := services.GetUserByEmail(server.db, user.Email)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err.Error())
	}

	err = utils.CheckPassword(user.Password, dbUser.Password)
	if err != nil {
		return ctx.JSON(http.StatusForbidden, "wrong password")
	}
	accessToken, _ := services.CreateToken(server.cfg, int(dbUser.ID))

	response := schemas.LoginUserResponse{
		AccessToken: accessToken,
		TokenType:   dbUser.Email,
	}
	return ctx.JSON(http.StatusOK, response)
}

func (server *Server) RegisterUser(ctx echo.Context) error {
	user := new(schemas.RegisterUserRequest)
	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}
	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	dbUser, _ := services.GetUserByEmail(server.db, user.Email)
	if dbUser != nil {
		return ctx.JSON(http.StatusBadRequest, "User exist")
	}
	_, err := services.RegisterUser(server.db, user.Email, user.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusOK, nil)
}
