package controllers

import (
	"net/http"

	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/dto"
	"AwisPalace_IngredientManagement/models"
	"AwisPalace_IngredientManagement/utils"

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

// PostIngredients godoc
// @Summary Post Ingredients
// @Description Post Ingredients
// @Tags Ingredients
// @Param unit body dto.IngredientParamRequest true "Create ingredient"
// @Router /ingredients [post]
func PostIngredients(c *gin.Context) {
	var input dto.IngredientParamRequest

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
	ingredient := models.Ingredient{
		Name:        input.Name,
		Slug: 		 utils.GenerateSlug(input.Name),
		Stock: 		 input.Stock,
		UnitID: 	 input.UnitID,
	}

	// Simpan ke database
	if err := config.DB.Create(&ingredient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create ingredient",
			"error":   err.Error(),
		})
		return
	}

	// Mapping ke DTO response
	var result dto.Unit
	copier.Copy(&result, &ingredient)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Ingredient created successfully",
		"data":    result,
	})
}
