package controllers

import (
	"net/http"

	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/dto"
	"AwisPalace_IngredientManagement/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// GetIngredients godoc
// @Summary Get Ingredients
// @Description Get Ingredients
// @Tags Ingredients
// @Router /ingredients [get]
func GetIngredients(c *gin.Context) {
	var ingredients []models.Ingredient

	if err := config.DB.Preload("Unit").Find(&ingredients).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []dto.Ingredient
	copier.Copy(&result, &ingredients)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Get Data Ingredient Succcess",
		"data":    result,
	})
}
