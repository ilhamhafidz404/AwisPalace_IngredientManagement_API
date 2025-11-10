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

	// route units
	unitRoutes := router.Group("/units")
	{
		unitRoutes.GET("/", controllers.GetUnits)
		unitRoutes.POST("/", controllers.PostUnits)
		unitRoutes.PUT("/:id", controllers.UpdateUnit)
		unitRoutes.DELETE("/:id", controllers.DeleteUnit)
	}

	// route ingredients
	ingredientRoutes := router.Group("/ingredients")
	{
		ingredientRoutes.GET("/", controllers.GetIngredients)
		ingredientRoutes.POST("/", controllers.PostIngredients)
		ingredientRoutes.PUT("/:id", controllers.UpdateIngredients)
	}
}
