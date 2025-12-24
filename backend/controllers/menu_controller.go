package controllers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/dto"
	"AwisPalace_IngredientManagement/models"
	"AwisPalace_IngredientManagement/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ==================== GET MENUS ====================

// GetMenus godoc
// @Summary Get Menus
// @Description Get all menus with ingredients
// @Tags Menus
// @Produce json
// @Success 200 {object} map[string]interface{}
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

// ==================== SHOW MENU ====================

// GetMenu godoc
// @Summary Get Menu
// @Description Get Menu by ID
// @Tags Menus
// @Param id path int true "Menu ID"
// @Router /menus/{id} [get]
func ShowMenu(c *gin.Context) {
	id := c.Param("id")
	var menu models.Menu

	if err := config.DB.
		Preload("MenuIngredients").
		Preload("MenuIngredients.Ingredient").
		Preload("MenuIngredients.Unit").
		First(&menu, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Menu not found",
		})
		return
	}

	menuDTO := dto.Menu{
		ID:          menu.ID,
		Name:        menu.Name,
		Slug:        menu.Slug,
		Image:       menu.Image,
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

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   menuDTO,
	})
}

// ==================== CREATE MENU ====================

// PostMenu godoc
// @Summary Create Menu
// @Description Create menu with image upload
// @Tags Menus
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Menu Name"
// @Param price formData number true "Menu Price"
// @Param description formData string false "Menu Description"
// @Param ingredients formData string true "Ingredients JSON array"
// @Param image formData file true "Menu Image"
// @Success 201 {object} map[string]interface{}
// @Router /menus [post]
func PostMenu(c *gin.Context) {
	// Parse form data
	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	ingredientsStr := c.PostForm("ingredients")

	// Validate required fields
	if name == "" || priceStr == "" || ingredientsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "name, price, and ingredients are required",
		})
		return
	}

	// Parse price
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid price format",
		})
		return
	}

	// Parse ingredients JSON
	var ingredients []dto.MenuIngredientRequest
	if err := json.Unmarshal([]byte(ingredientsStr), &ingredients); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ingredients format: " + err.Error(),
		})
		return
	}

	// Validate ingredients
	if len(ingredients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "At least one ingredient is required",
		})
		return
	}

	// Handle file upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Image is required",
		})
		return
	}

	// Validate file extension
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedExts[ext] {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Only image files (jpg, jpeg, png, gif) are allowed",
		})
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	uploadPath := "./uploads/" + filename

	// Create uploads directory if not exists
	if err := os.MkdirAll("./uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create upload directory",
		})
		return
	}

	// Save file
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to upload image: " + err.Error(),
		})
		return
	}

	// Start transaction
	tx := config.DB.Begin()

	menu := models.Menu{
		Name:        name,
		Slug:        utils.GenerateSlug(name),
		Description: description,
		Image:       filename,
		Price:       price,
	}

	if err := tx.Create(&menu).Error; err != nil {
		tx.Rollback()
		os.Remove(uploadPath)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Insert ingredients
	for _, item := range ingredients {
		menuIngredient := models.MenuIngredient{
			MenuID:       menu.ID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitID:       item.UnitID,
		}

		if err := tx.Create(&menuIngredient).Error; err != nil {
			tx.Rollback()
			os.Remove(uploadPath)
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	tx.Commit()

	// Generate full image URL
	baseURL := getBaseURL(c)
	imageURL := baseURL + "/uploads/" + filename

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Menu created successfully",
		"data": gin.H{
			"id":    menu.ID,
			"name":  menu.Name,
			"image": imageURL,
		},
	})
}

// ==================== UPDATE MENU ====================

// UpdateMenu godoc
// @Summary Update Menu
// @Description Update menu by ID with optional image upload
// @Tags Menus
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Menu ID"
// @Param name formData string true "Menu Name"
// @Param price formData number true "Menu Price"
// @Param description formData string false "Menu Description"
// @Param ingredients formData string true "Ingredients JSON array"
// @Param image formData file false "Menu Image (optional)"
// @Success 200 {object} map[string]interface{}
// @Router /menus/{id} [put]
func UpdateMenu(c *gin.Context) {
	menuID := c.Param("id")

	// Parse form data
	name := c.PostForm("name")
	priceStr := c.PostForm("price")
	description := c.PostForm("description")
	ingredientsStr := c.PostForm("ingredients")

	// Validate required fields
	if name == "" || priceStr == "" || ingredientsStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "name, price, and ingredients are required",
		})
		return
	}

	// Parse price
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid price format",
		})
		return
	}

	// Parse ingredients JSON
	var ingredients []dto.MenuIngredientRequest
	if err := json.Unmarshal([]byte(ingredientsStr), &ingredients); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid ingredients format: " + err.Error(),
		})
		return
	}

	// Validate ingredients
	if len(ingredients) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "At least one ingredient is required",
		})
		return
	}

	// Check if menu exists
	var menu models.Menu
	if err := config.DB.First(&menu, menuID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Menu not found",
		})
		return
	}

	// Store old image filename for deletion if new image is uploaded
	oldImageFilename := menu.Image
	var newFilename string

	// Handle file upload (optional for update)
	file, err := c.FormFile("image")
	if err == nil {
		// Validate file extension
		ext := filepath.Ext(file.Filename)
		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Only image files (jpg, jpeg, png, gif) are allowed",
			})
			return
		}

		uploadDir := "./uploads"
		newFilename = uuid.New().String() + ext // ✅ FIX: Assign ke variable newFilename
		uploadPath := filepath.Join(uploadDir, newFilename)

		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to create upload directory",
			})
			return
		}

		if err := c.SaveUploadedFile(file, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to upload image: " + err.Error(),
			})
			return
		}
	}

	// Start transaction
	tx := config.DB.Begin()

	// Update menu
	menu.Name = name
	menu.Slug = utils.GenerateSlug(name)
	menu.Description = description
	menu.Price = price

	// Update image only if new image uploaded
	if newFilename != "" {
		menu.Image = newFilename
	}

	if err := tx.Save(&menu).Error; err != nil {
		tx.Rollback()
		// Delete new uploaded file if exists
		if newFilename != "" {
			os.Remove("./uploads/" + newFilename)
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Delete old ingredients
	if err := tx.Where("menu_id = ?", menu.ID).Delete(&models.MenuIngredient{}).Error; err != nil {
		tx.Rollback()
		if newFilename != "" {
			os.Remove("./uploads/" + newFilename)
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to delete old ingredients",
		})
		return
	}

	// Insert new ingredients
	for _, item := range ingredients {
		menuIngredient := models.MenuIngredient{
			MenuID:       menu.ID,
			IngredientID: item.IngredientID,
			Quantity:     item.Quantity,
			UnitID:       item.UnitID,
		}

		if err := tx.Create(&menuIngredient).Error; err != nil {
			tx.Rollback()
			if newFilename != "" {
				os.Remove("./uploads/" + newFilename)
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	tx.Commit()

	// Delete old image file if new image was uploaded
	if newFilename != "" && oldImageFilename != "" {
		oldImagePath := filepath.Join("./uploads", oldImageFilename)
		os.Remove(oldImagePath)
	}

	// Generate full image URL
	baseURL := getBaseURL(c)
	imageURL := baseURL + "/uploads/" + menu.Image

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Menu updated successfully",
		"data": gin.H{
			"id":    menu.ID,
			"name":  menu.Name,
			"image": imageURL,
		},
	})
}

// ==================== DELETE MENU ====================

// DeleteMenu godoc
// @Summary Delete Menu
// @Description Delete menu by ID
// @Tags Menus
// @Produce json
// @Param id path int true "Menu ID"
// @Success 200 {object} map[string]interface{}
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

	// Store image filename for deletion
	imageFile := menu.Image

	// Delete ingredients first (foreign key constraint)
	if err := tx.Where("menu_id = ?", menu.ID).
		Delete(&models.MenuIngredient{}).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	// Delete menu
	if err := tx.Delete(&menu).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	tx.Commit()

	// Delete image file after successful database deletion
	if imageFile != "" {
		os.Remove("./uploads/" + imageFile)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Menu deleted successfully",
	})
}

func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host
}
