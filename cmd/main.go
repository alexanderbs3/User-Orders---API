package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/alexanderbs3/user-orders-api/internal/config"
	"github.com/alexanderbs3/user-orders-api/internal/handler"
	"github.com/alexanderbs3/user-orders-api/internal/repository"
	"github.com/alexanderbs3/user-orders-api/internal/service"
	"github.com/alexanderbs3/user-orders-api/pkg/middleware"
)

func main() {
	config.LoadEnv()

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco: %v", err)
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	userService := service.NewUserService(userRepo)
	orderService := service.NewOrderService(orderRepo, userRepo)

	userHandler := handler.NewUserHandler(userService)
	orderHandler := handler.NewOrderHandler(orderService)

	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(middleware.Logger(logger))
	router.Use(gin.Recovery())

	api := router.Group("/api/v1")
	userHandler.RegisterRoutes(api, orderHandler)
	orderHandler.RegisterRoutes(api)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciado na porta %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
