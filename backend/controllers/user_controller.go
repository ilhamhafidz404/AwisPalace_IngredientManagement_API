package controllers

import (
	"net/http"

	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/models"

	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary Get Users
// @Description Get Users
// @Tags Users
// @Success 200 {object} map[string]string
// @Router /users [get]
func GetUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}
