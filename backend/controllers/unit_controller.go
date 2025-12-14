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
// @Param unit body dto.UnitParamRequest true "Unit data"
// @Router /units [post]
func PostUnits(c *gin.Context) {
	var input dto.UnitParamRequest

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
		Name:   input.Name,
		Symbol: input.Symbol,
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

// UpdateUnit godoc
// @Summary Update Unit
// @Description Update an existing unit by ID
// @Tags Units
// @Param id path int true "Unit ID"
// @Param unit body dto.UnitParamRequest true "Updated unit data"
// @Router /units/{id} [put]
func UpdateUnit(c *gin.Context) {
	id := c.Param("id")

	var input dto.UnitParamRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	var unit models.Unit
	if err := config.DB.First(&unit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Unit not found",
		})
		return
	}

	unit.Name = input.Name
	unit.Symbol = input.Symbol

	if err := config.DB.Save(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update unit",
			"error":   err.Error(),
		})
		return
	}

	var result dto.Unit
	copier.Copy(&result, &unit)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Unit updated successfully",
		"data":    result,
	})
}

// DeleteUnit godoc
// @Summary Delete Unit
// @Description Delete a unit by ID
// @Tags Units
// @Param id path string true "Unit ID"
// @Router /units/{id} [delete]
func DeleteUnit(c *gin.Context) {
	id := c.Param("id")

	var unit models.Unit

	// Cek apakah unit dengan ID tersebut ada
	if err := config.DB.First(&unit, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Unit not found",
			"error":   err.Error(),
		})
		return
	}

	// Hapus data unit
	if err := config.DB.Delete(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete unit",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Unit deleted successfully",
	})
}
