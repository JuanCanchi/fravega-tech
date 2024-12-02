package main

import (
	"fravega-tech/config"
	"fravega-tech/internal/handler"
	"fravega-tech/internal/repository/mongo"
	"fravega-tech/internal/usecase"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	db, err := config.InitMongoDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	productRepo := mongo.NewProductRepository(db)
	productUsecase := usecase.NewProductUsecase(productRepo)

	router := gin.Default()
	handler.NewProductHandler(router, productUsecase)

	log.Println("Starting server on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
