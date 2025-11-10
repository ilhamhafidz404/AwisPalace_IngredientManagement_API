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
// @Param ingredient body dto.IngredientParamRequest true "Create ingredient"
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
	var result dto.Ingredient
	copier.Copy(&result, &ingredient)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Ingredient created successfully",
		"data":    result,
	})
}

// UpdateIngredient godoc
// @Summary Update Ingredient
// @Description Update an existing ingredient by ID
// @Tags Ingredients
// @Param id path int true "Ingredient ID"
// @Param ingredient body dto.IngredientParamRequest true "Updated ingredient data"
// @Router /ingredients/{id} [put]
func UpdateIngredients(c *gin.Context) {
	id := c.Param("id")

	var input dto.IngredientParamRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid input data",
			"error":   err.Error(),
		})
		return
	}

	var ingredient models.Ingredient
	if err := config.DB.First(&ingredient, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Ingredient not found",
		})
		return
	}

	ingredient.Name 	= input.Name
	ingredient.Slug 	= utils.GenerateSlug(input.Name)
	ingredient.Stock 	= input.Stock
	ingredient.UnitID 	= input.UnitID

	if err := config.DB.Save(&ingredient).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update ingredient",
			"error":   err.Error(),
		})
		return
	}

	var result dto.Ingredient
	copier.Copy(&result, &ingredient)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Ingredient updated successfully",
		"data":    result,
	})
}