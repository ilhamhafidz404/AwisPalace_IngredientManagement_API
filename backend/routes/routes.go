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

	auth := router.Group("/auth")
	{
		auth.POST("/google", controllers.GoogleAuth)
		auth.GET("/verify", controllers.VerifyToken)
		auth.POST("/refresh", controllers.RefreshToken)
	}

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
		ingredientRoutes.DELETE("/:id", controllers.DeleteIngredients)
	}

	// route menus
	menuRoutes := router.Group("/menus")
	{
		menuRoutes.GET("", controllers.GetMenus)
		menuRoutes.GET("/:id", controllers.ShowMenu)
		menuRoutes.POST("", controllers.PostMenu)
		menuRoutes.PUT("/:id", controllers.UpdateMenu)
		menuRoutes.DELETE("/:id", controllers.DeleteMenu)
	}

	// route transactions
	transactionRoutes := router.Group("/transactions")
	{
		transactionRoutes.GET("/", controllers.GetTransactions)
		transactionRoutes.GET("/:id", controllers.GetTransaction)
		transactionRoutes.POST("/", controllers.PostTransaction)
		transactionRoutes.DELETE("/:id", controllers.DeleteTransaction)
	}

	// route exports
	exportRoutes := router.Group("/export")
	{
		exportRoutes.GET("/transactions", controllers.ExportTransactions)
	}

}
