package main

import (
	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/databases/migrations"
	"AwisPalace_IngredientManagement/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
 // Load file .env
  if err := godotenv.Load(); err != nil {
    log.Fatal("❌ Error loading .env file")
  }

	// Connect DB
	config.ConnectDB()
  migrations.Migrate()

	// Router init
	r := gin.Default()
	routes.SetupRoutes(r)

	// 
	r.Run(":8080")
}
