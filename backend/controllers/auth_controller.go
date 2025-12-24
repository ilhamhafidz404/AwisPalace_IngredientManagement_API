package controllers

import (
	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/dto"
	"AwisPalace_IngredientManagement/models"
	"AwisPalace_IngredientManagement/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GoogleAuth godoc
// @Summary Authenticate with Google
// @Description Authenticate user with Google ID token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param auth body dto.GoogleAuthRequest true "Google authentication data"
// @Success 200 {object} dto.AuthResponse "Login successful"
// @Success 201 {object} dto.AuthResponse "User created and logged in"
// @Failure 400 {object} map[string]interface{} "Invalid request"
// @Failure 401 {object} map[string]interface{} "Invalid ID token"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /auth/google [post]
func GoogleAuth(c *gin.Context) {
	var req dto.GoogleAuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request: " + err.Error(),
		})
		return
	}

	// Verify Google ID token (optional but recommended for production)
	// Uncomment this in production
	/*
		payload, err := idtoken.Validate(context.Background(), req.IDToken, "YOUR_GOOGLE_CLIENT_ID")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Invalid ID token: " + err.Error(),
			})
			return
		}

		// Verify email matches
		if payload.Claims["email"] != req.Email {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "Email mismatch",
			})
			return
		}
	*/

	// Check if user exists
	var user models.User
	result := config.DB.Where("email = ?", req.Email).First(&user)

	if result.Error != nil {
		// User doesn't exist, create new user
		user = models.User{
			Email:    req.Email,
			Name:     req.Name,
			PhotoURL: req.PhotoURL,
			GoogleID: req.Email, // Using email as GoogleID for simplicity
		}

		if err := config.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to create user: " + err.Error(),
			})
			return
		}

		// Generate token
		token, err := utils.GenerateToken(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": "Failed to generate token: " + err.Error(),
			})
			return
		}

		// Return response for new user
		c.JSON(http.StatusCreated, dto.AuthResponse{
			Status:  "success",
			Message: "User created and logged in successfully",
			Data: dto.AuthData{
				Token: token,
				User: dto.UserData{
					ID:       user.ID,
					Email:    user.Email,
					Name:     user.Name,
					PhotoURL: user.PhotoURL,
				},
			},
		})
		return
	}

	// User exists, update photo if changed
	if req.PhotoURL != "" && req.PhotoURL != user.PhotoURL {
		user.PhotoURL = req.PhotoURL
		config.DB.Save(&user)
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate token: " + err.Error(),
		})
		return
	}

	// Return response for existing user
	c.JSON(http.StatusOK, dto.AuthResponse{
		Status:  "success",
		Message: "Login successful",
		Data: dto.AuthData{
			Token: token,
			User: dto.UserData{
				ID:       user.ID,
				Email:    user.Email,
				Name:     user.Name,
				PhotoURL: user.PhotoURL,
			},
		},
	})
}

// VerifyToken godoc
// @Summary Verify JWT token
// @Description Verify if the provided JWT token is valid
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Token is valid"
// @Failure 401 {object} map[string]interface{} "Invalid or missing token"
// @Router /auth/verify [get]
func VerifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Authorization header required",
		})
		return
	}

	// Extract token
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Validate token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid token: " + err.Error(),
		})
		return
	}

	// Get user from database
	var user models.User
	if err := config.DB.First(&user, claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token is valid",
		"data": dto.UserData{
			ID:       user.ID,
			Email:    user.Email,
			Name:     user.Name,
			PhotoURL: user.PhotoURL,
		},
	})
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Description Generate a new JWT token from existing valid token
// @Tags Authentication
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "New token generated"
// @Failure 401 {object} map[string]interface{} "Invalid or missing token"
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Authorization header required",
		})
		return
	}

	// Extract token
	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Validate token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "Invalid token: " + err.Error(),
		})
		return
	}

	// Generate new token
	newToken, err := utils.GenerateToken(claims.UserID, claims.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate new token: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Token refreshed successfully",
		"data": gin.H{
			"token": newToken,
		},
	})
}
