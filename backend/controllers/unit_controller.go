package controllers

import (
	"net/http"

	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/dto"
	"AwisPalace_IngredientManagement/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// GetUnits godoc
// @Summary Get Units
// @Description Get Units
// @Tags Units
// @Router /units [get]
func GetUnits(c *gin.Context) {
	var units []models.Unit

	if err := config.DB.Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []dto.Unit
	copier.Copy(&result, &units)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Get Data Unit Succcess",
		"data":    result,
	})
}
