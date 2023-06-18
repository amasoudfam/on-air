package handlers

import (
	"net/http"
	"on-air/config"
	"on-air/repository"
	"on-air/utils"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

const Bearer = "Bearer"

type Auth struct {
	DB  *gorm.DB
	JWT *config.JWT
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type" binding:"required"`
}

func (a *Auth) Login(ctx echo.Context) error {
	var req LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}

	if err := ctx.Validate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	// TODO repository
	// TODO error package
	dbUser, err := repository.GetUserByEmail(a.DB, req.Email)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "Invalid credentials")
	}

	err = utils.CheckPassword(req.Password, dbUser.Password)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, "Invalid credentials")
	}

	accessToken, err := repository.CreateToken(a.JWT, int(dbUser.ID))
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, "Internal server error")
	}

	return ctx.JSON(http.StatusOK, LoginResponse{
		AccessToken: accessToken,
		TokenType:   Bearer,
	})
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type" binding:"required"`
}

func (a *Auth) Register(ctx echo.Context) error {
	user := new(RegisterRequest)
	if err := ctx.Bind(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, "")
	}
	if err := ctx.Validate(user); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	dbUser, _ := repository.GetUserByEmail(a.DB, user.Email)
	if dbUser != nil {
		return ctx.JSON(http.StatusBadRequest, "User exist")
	}
	_, err := repository.RegisterUser(a.DB, user.Email, user.Password)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	return ctx.JSON(http.StatusCreated, nil)
}
