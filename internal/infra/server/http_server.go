package server

import (
	"fmt"
	"time"

	"github.com/goccy/go-json"
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"

	_ "github.com/rescuebyte/rescuebyte-be/docs"
	authCtr "github.com/rescuebyte/rescuebyte-be/internal/app/auth/controller"
	authSvc "github.com/rescuebyte/rescuebyte-be/internal/app/auth/service"
	cartCtr "github.com/rescuebyte/rescuebyte-be/internal/app/cart/controller"
	cartRepo "github.com/rescuebyte/rescuebyte-be/internal/app/cart/repository"
	cartSvc "github.com/rescuebyte/rescuebyte-be/internal/app/cart/service"
	enumCtr "github.com/rescuebyte/rescuebyte-be/internal/app/enum/controller"
	enumSvc "github.com/rescuebyte/rescuebyte-be/internal/app/enum/service"
	orderCtr "github.com/rescuebyte/rescuebyte-be/internal/app/order/controller"
	orderRepo "github.com/rescuebyte/rescuebyte-be/internal/app/order/repository"
	orderSvc "github.com/rescuebyte/rescuebyte-be/internal/app/order/service"
	paymentCtr "github.com/rescuebyte/rescuebyte-be/internal/app/payment/controller"
	paymentRepo "github.com/rescuebyte/rescuebyte-be/internal/app/payment/repository"
	paymentSvc "github.com/rescuebyte/rescuebyte-be/internal/app/payment/service"
	productCtr "github.com/rescuebyte/rescuebyte-be/internal/app/product/controller"
	productRepo "github.com/rescuebyte/rescuebyte-be/internal/app/product/repository"
	productSvc "github.com/rescuebyte/rescuebyte-be/internal/app/product/service"
	storageCtr "github.com/rescuebyte/rescuebyte-be/internal/app/storage/controller"
	storageSvc "github.com/rescuebyte/rescuebyte-be/internal/app/storage/service"
	storeCtr "github.com/rescuebyte/rescuebyte-be/internal/app/store/controller"
	storeRepo "github.com/rescuebyte/rescuebyte-be/internal/app/store/repository"
	storeSvc "github.com/rescuebyte/rescuebyte-be/internal/app/store/service"
	userCtr "github.com/rescuebyte/rescuebyte-be/internal/app/user/controller"
	userRepo "github.com/rescuebyte/rescuebyte-be/internal/app/user/repository"
	userSvc "github.com/rescuebyte/rescuebyte-be/internal/app/user/service"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
	"github.com/rescuebyte/rescuebyte-be/internal/middlewares"
	errorhandler "github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/error_handler"
	"github.com/rescuebyte/rescuebyte-be/pkg/helpers/http/response"
	"github.com/rescuebyte/rescuebyte-be/pkg/jwt"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
	"github.com/rescuebyte/rescuebyte-be/pkg/midtrans"
	"github.com/rescuebyte/rescuebyte-be/pkg/oauth"
	"github.com/rescuebyte/rescuebyte-be/pkg/supabase"
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
	supabaseStorage := supabase.SupabaseStorage

	dbGetter := database.NewDBGetter(db)
	transactor := database.NewTransactor(db)

	// Repositories
	userRepository := userRepo.NewUserRepository(dbGetter)
	storeRepository := storeRepo.NewStoreRepository(dbGetter)
	productRepository := productRepo.NewProductRepository(dbGetter)
	cartRepository := cartRepo.NewCartRepository(dbGetter)
	orderRepository := orderRepo.NewOrderRepository(dbGetter)
	paymentRepository := paymentRepo.NewPaymentRepository(dbGetter)

	// Midtrans Client
	midtransService := midtrans.NewMidtransService(midtrans.Config{
		ServerKey:    env.AppEnv.MidtransServerKey,
		ClientKey:    env.AppEnv.MidtransClientKey,
		IsProduction: env.AppEnv.MidtransIsProduction,
	})

	// Services
	enumService := enumSvc.NewEnumService()
	userService := userSvc.NewUserService(transactor, userRepository, validatorService)
	authService := authSvc.NewAuthService(
		transactor,
		userRepository,
		validatorService,
		jwtService,
		googleOAuth,
		uuidService,
	)
	storeService := storeSvc.NewStoreService(
		transactor,
		storeRepository,
		userRepository,
		validatorService,
		uuidService,
	)
	productService := productSvc.NewProductService(
		transactor,
		productRepository,
		storeRepository,
		validatorService,
		uuidService,
	)
	storageService := storageSvc.NewStorageService(
		supabaseStorage,
		uuidService,
	)
	cartService := cartSvc.NewCartService(
		transactor,
		cartRepository,
		productRepository,
		validatorService,
		uuidService,
	)
	orderService := orderSvc.NewOrderService(
		transactor,
		orderRepository,
		cartRepository,
		productRepository,
		storeRepository,
		paymentRepository,
		midtransService,
		validatorService,
		uuidService,
	)
	paymentService := paymentSvc.NewPaymentService(
		transactor,
		paymentRepository,
		orderRepository,
		productRepository,
		midtransService,
	)

	// Middlewares
	middleware := middlewares.NewMiddleware(jwtService)

	// Swagger Documentation
	s.app.Get("/swagger", func(c fiber.Ctx) error {
		return c.Redirect().To("/swagger/index.html")
	})
	s.app.Get("/swagger/*", swaggo.New(swaggo.Config{
		Title: "RescueByte API Documentation",
	}))

	s.app.Get("/", func(c fiber.Ctx) error {
		return c.Redirect().Status(fiber.StatusMovedPermanently).To("/swagger/index.html")
	})

	s.app.Get("/health", s.healthCheck)

	api := s.app.Group("/api")
	v1 := api.Group("/v1")

	// Domain Controllers Initialization (mounting endpoints within their respective modules)
	enumCtr.InitEnumController(v1, enumService)
	authCtr.InitAuthController(v1, authService, middleware)
	userCtr.InitUserController(v1, userService, middleware)
	storeCtr.InitStoreController(v1, storeService, middleware)
	productCtr.InitProductController(v1, productService, middleware)
	storageCtr.InitStorageController(v1, storageService, middleware)
	cartCtr.InitCartController(v1, cartService, middleware)
	orderCtr.InitOrderController(v1, orderService, middleware)
	paymentCtr.InitPaymentController(v1, paymentService, middleware)
}

// healthCheck godoc
// @Summary Health check endpoint
// @Description Check API server status and liveness
// @Tags Health Check
// @Produce json
// @Success 200 {object} response.Response{payload=map[string]string} "API is healthy"
// @Router /health [get]
func (s *httpServer) healthCheck(c fiber.Ctx) error {
	return response.SendResponse(c, fiber.StatusOK, map[string]string{
		"message": "RescueByte API is up and running",
		"status":  "healthy",
	})
}
