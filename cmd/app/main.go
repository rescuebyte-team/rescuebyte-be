package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/env"
	redisConn "github.com/rescuebyte/rescuebyte-be/internal/infra/redis"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/server"
	"github.com/rescuebyte/rescuebyte-be/pkg/log"
)

// @title RescueByte Backend API
// @version 1.0
// @description RESTful API documentation for RescueByte Backend.
// @termsOfService http://swagger.io/terms/

// @contact.name RescueByte Support
// @contact.email support@rescuebyte.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT access token.
func main() {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		log.Fatal(log.Fields{"error": err.Error()}, "[Main] Failed to load timezone")
	}
	time.Local = loc

	enums.DefaultLanguage = env.AppEnv.AppLang

	server := server.NewHTTPServer()
	psqlDB := database.NewPgsqlConn()
	defer psqlDB.Close()

	redisClient := redisConn.NewRedisConnWithDB(env.AppEnv.RedisDB)
	defer redisClient.Close()

	app := server.GetApp()

	server.MountMiddlewares()
	server.MountRoutes(psqlDB, redisClient)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info(log.Fields{}, "[Main] Shutting down server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		if shutdownErr := app.ShutdownWithContext(shutdownCtx); shutdownErr != nil {
			log.Error(log.Fields{"error": shutdownErr.Error()}, "[Main] Server forced to shutdown")
		}
	}()

	server.Start(env.AppEnv.AppPort)
}
