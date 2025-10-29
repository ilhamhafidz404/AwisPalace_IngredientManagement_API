package routes

import (
	"AwisPalace_IngredientManagement/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	// route root
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Server running successfully!"})
	})

	// route users
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/", controllers.GetUsers)
	}
}
