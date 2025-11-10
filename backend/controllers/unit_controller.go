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

// PostUnits godoc
// @Summary Post Units
// @Description Post Units
// @Tags Units
// @Param unit body dto.UnitCreateRequest true "Unit data"
// @Router /units [post]
func PostUnits(c *gin.Context) {
	var input dto.UnitCreateRequest

	// Bind JSON ke struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	// Mapping DTO ke model database
	unit := models.Unit{
		Name:        input.Name,
		Symbol: 	 input.Symbol,
	}

	// Simpan ke database
	if err := config.DB.Create(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create unit",
			"error":   err.Error(),
		})
		return
	}

	// Mapping ke DTO response
	var result dto.Unit
	copier.Copy(&result, &unit)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Unit created successfully",
		"data":    result,
	})
}
