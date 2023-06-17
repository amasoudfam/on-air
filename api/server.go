package api

import (
	"net/http"
	"on-air/config"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Server struct {
	db     *gorm.DB
	router *echo.Echo
	cfg    *config.Config
}

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewServer(cfg *config.Config, db *gorm.DB) (*Server, error) {

	server := &Server{
		cfg: cfg,
		db:  db,
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := echo.New()
	router.Validator = &CustomValidator{validator: validator.New()}
	router.POST("/login", server.GetAuthToken)

	// authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// authRoutes.POST("/accounts", server.createAccount)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Start(address)
}
