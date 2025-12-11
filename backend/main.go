package main

import (
	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/databases/migrations"
	"AwisPalace_IngredientManagement/databases/seeders"
	"AwisPalace_IngredientManagement/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "AwisPalace_IngredientManagement/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Awis Palace Ingredient Management API
// @version 1.0
// @description API documentation for Awis Palace Ingredient Management built with Gin and GORM.
// @host localhost:8080
// @BasePath /
func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Error loading .env file")
	}

	// Connect to DB
	config.ConnectDB()
	migrations.Migrate()
	seeders.DatabaseSeeder(config.DB)

	// init routes
	r := gin.Default()
	routes.SetupRoutes(r)

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//
	r.Run(":8080")
}
