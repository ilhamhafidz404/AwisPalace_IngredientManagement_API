package controllers

import (
	"net/http"

	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/dto"
	"AwisPalace_IngredientManagement/models"
	"AwisPalace_IngredientManagement/utils"

	"github.com/gin-gonic/gin"
)

// GetMenus godoc
// @Summary Get Menus
// @Description Get Menus
// @Tags Menus
// @Router /menus [get]
// GetMenus godoc
// @Summary Get Menus
// @Tags Menus
// @Router /menus [get]
func GetMenus(c *gin.Context) {
	var menus []models.Menu

	if err := config.DB.
		Preload("MenuIngredients").
		Preload("MenuIngredients.Ingredient").
		Preload("MenuIngredients.Unit").
		Find(&menus).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var response []dto.Menu

	for _, menu := range menus {
		menuDTO := dto.Menu{
			ID:          menu.ID,
			Name:        menu.Name,
			Slug:        menu.Slug,
			Image:       menu.Image,
			Price:       menu.Price,
			Description: menu.Description,
			CreatedAt:   menu.CreatedAt,
			UpdatedAt:   menu.UpdatedAt,
		}

		for _, mi := range menu.MenuIngredients {
			menuDTO.Ingredients = append(menuDTO.Ingredients, dto.MenuIngredient{
				ID: mi.ID,
				Ingredient: dto.MenuIngredientIngredient{
					ID:   mi.Ingredient.ID,
					Name: mi.Ingredient.Name,
					Slug: mi.Ingredient.Slug,
				},
				Quantity: mi.Quantity,
				Unit: dto.MenuIngredientUnit{
					ID:   mi.Unit.ID,
					Name: mi.Unit.Name,
				},
			})
		}

		response = append(response, menuDTO)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// PostMenus godoc
// @Summary Post Menu
// @Description Create Menu
// @Tags Menus
// @Param menu body dto.MenuCreateRequest true "Create menu"
// @Router /menus [post]
// PostMenu godoc
// @Summary Create Menu
// @Tags Menus
// @Router /menus [post]
func PostMenu(c *gin.Context) {
	var input dto.MenuCreateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	tx := config.DB.Begin()

	menu := models.Menu{
		Name:        input.Name,
		Slug:        utils.GenerateSlug(input.Name),
		Description: input.Description,
		Image:       input.Image,
		Price:       input.Price,
	}

	if err := tx.Create(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	for _, item := range input.Ingredients {
		menuIngredient := models.MenuIngredient{
			MenuID:       menu.ID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitID:       item.UnitID,
		}

		if err := tx.Create(&menuIngredient).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Menu created successfully",
	})
}

// UpdateMenus godoc
// @Summary Update Menu
// @Description Update Menu by ID
// @Tags Menus
// @Param id path int true "Menu ID"
// @Param menu body dto.MenuUpdateRequest true "Update menu"
// @Router /menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	id := c.Param("id")
	var input dto.MenuUpdateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	tx := config.DB.Begin()

	var menu models.Menu
	if err := tx.First(&menu, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Menu not found",
		})
		return
	}

	menu.Name = input.Name
	menu.Slug = utils.GenerateSlug(input.Name)
	menu.Description = input.Description
	menu.Image = input.Image
	menu.Price = input.Price

	if err := tx.Save(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Hapus relasi lama
	if err := tx.Where("menu_id = ?", menu.ID).
		Delete(&models.MenuIngredient{}).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Insert relasi baru
	for _, item := range input.Ingredients {
		if err := tx.Create(&models.MenuIngredient{
			MenuID:       menu.ID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitID:       item.UnitID,
		}).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Menu updated successfully",
	})
}

// DeleteMenus godoc
// @Summary Delete Menu
// @Description Delete Menu by ID
// @Tags Menus
// @Param id path int true "Menu ID"
// @Router /menus/{id} [delete]
func DeleteMenu(c *gin.Context) {
	id := c.Param("id")

	tx := config.DB.Begin()

	var menu models.Menu
	if err := tx.First(&menu, id).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Menu not found",
		})
		return
	}

	// Hapus pivot table dulu
	if err := tx.Where("menu_id = ?", menu.ID).
		Delete(&models.MenuIngredient{}).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Hapus menu
	if err := tx.Delete(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Menu deleted successfully",
	})
}
