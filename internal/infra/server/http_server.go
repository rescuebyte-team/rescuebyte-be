package server

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	authCtr "github.com/rescuebyte/rescuebyte-be/internal/app/auth/controller"
	authSvc "github.com/rescuebyte/rescuebyte-be/internal/app/auth/service"
	userCtr "github.com/rescuebyte/rescuebyte-be/internal/app/user/controller"
	userRepo "github.com/rescuebyte/rescuebyte-be/internal/app/user/repository"
	userSvc "github.com/rescuebyte/rescuebyte-be/internal/app/user/service"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
	errorhandler "github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/error_handler"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
	"github.com/rescuebyte/rescuebyte-be/pkg/oauth"
	"github.com/rescuebyte/rescuebyte-be/pkg/uuid"
	"github.com/rescuebyte/rescuebyte-be/pkg/validator"
)

type HTTPServer interface {
	Start(port string)
	MountMiddlewares()
	MountRoutes(db *sqlx.DB, redisClient *redis.Client)
	GetApp() *fiber.App
}

type httpServer struct {
	app *fiber.App
}

func NewHTTPServer() HTTPServer {
	config := fiber.Config{
		CaseSensitive:  true,
		AppName:        "RescueByte Backend",
		ServerHeader:   "RescueByte API",
		JSONEncoder:    json.Marshal,
		JSONDecoder:    json.Unmarshal,
		ErrorHandler:   errorhandler.ErrorHandler,
		ReadTimeout:    1 * time.Minute,
		WriteTimeout:   1 * time.Minute,
		IdleTimeout:    2 * time.Minute,
		ReadBufferSize: 32768,
	}

	app := fiber.New(config)

	return &httpServer{
		app: app,
	}
}

func (s *httpServer) GetApp() *fiber.App {
	return s.app
}

func (s *httpServer) Start(port string) {
	if len(port) > 0 && port[0] != ':' {
		port = fmt.Sprintf(":%s", port)
	}

	listenConfig := fiber.ListenConfig{
		EnablePrefork: false,
	}

	log.Info(log.Fields{"port": port}, "[Server] Starting HTTP server...")
	if err := s.app.Listen(port, listenConfig); err != nil {
		log.Fatal(log.Fields{"error": err.Error()}, "[Server] Failed to start server")
	}
}

func (s *httpServer) MountMiddlewares() {
	s.app.Use(requestid.New())
	s.app.Use(middlewares.RecoverConfig())
	s.app.Use(middlewares.RequestIDContext())
	s.app.Use(middlewares.LoggerConfig())
	s.app.Use(middlewares.Cors())
}

func (s *httpServer) MountRoutes(db *sqlx.DB, redisClient *redis.Client) {
	uuidService := uuid.UUID
	validatorService := validator.Validator
	jwtService := jwt.JWT
	googleOAuth := oauth.GoogleOAuth

	dbGetter := database.NewDBGetter(db)
	transactor := database.NewTransactor(db)

	// Repositories
	userRepository := userRepo.NewUserRepository(dbGetter)

	// Services
	userService := userSvc.NewUserService(transactor, userRepository, validatorService)
	authService := authSvc.NewAuthService(
		transactor,
		userRepository,
		validatorService,
		jwtService,
		googleOAuth,
		uuidService,
	)

	// Middlewares
	middleware := middlewares.NewMiddleware(jwtService)

	// Root health check endpoint
	s.app.Get("/", func(c fiber.Ctx) error {
		return response.SendResponse(c, fiber.StatusOK, map[string]string{
			"message": "RescueByte API is up and running",
			"status":  "healthy",
		})
	})

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	// Domain Controllers Initialization (mounting endpoints within their respective modules)
	authCtr.InitAuthController(v1, authService, middleware)
	userCtr.InitUserController(v1, userService, middleware)
}
