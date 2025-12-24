package controllers

import (
	"fmt"
	"net/http"
	"time"

	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/dto"
	"AwisPalace_IngredientManagement/models"

	"github.com/gin-gonic/gin"
)

// GetTransactions godoc
// @Summary Get Transactions
// @Description Get Transactions with optional date filter (default: this week)
// @Tags Transactions
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Router /transactions [get]
func GetTransactions(c *gin.Context) {
	var transactions []models.Transaction

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var startDate, endDate time.Time
	var err error

	if startDateStr == "" || endDateStr == "" {
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}

		startDate = time.Date(
			now.Year(),
			now.Month(),
			now.Day()-weekday+1,
			0, 0, 0, 0,
			now.Location(),
		)

		endDate = startDate.AddDate(0, 0, 7).Add(-time.Nanosecond)
	} else {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid start_date format (YYYY-MM-DD)",
			})
			return
		}

		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid end_date format (YYYY-MM-DD)",
			})
			return
		}

		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	//
	if err := config.DB.
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Preload("TransactionItems").
		Preload("TransactionItems.Menu").
		Preload("TransactionItems.StockReductions").
		Preload("TransactionItems.StockReductions.Ingredient").
		Preload("TransactionItems.StockReductions.Unit").
		Order("created_at DESC").
		Find(&transactions).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	//
	var response []dto.Transaction

	for _, transaction := range transactions {
		transactionDTO := dto.Transaction{
			ID:              transaction.ID,
			TransactionCode: transaction.TransactionCode,
			TransactionDate: transaction.TransactionDate,
			TotalAmount:     transaction.TotalAmount,
			Notes:           transaction.Notes,
			Status:          transaction.Status,
			CreatedAt:       transaction.CreatedAt,
			UpdatedAt:       transaction.UpdatedAt,
		}

		for _, item := range transaction.TransactionItems {
			itemDTO := dto.TransactionItem{
				ID:       item.ID,
				Quantity: item.Quantity,
				Price:    item.Price,
				Menu: dto.TransactionItemMenu{
					ID:    item.Menu.ID,
					Name:  item.Menu.Name,
					Slug:  item.Menu.Slug,
					Image: item.Menu.Image,
				},
			}

			for _, reduction := range item.StockReductions {
				itemDTO.StockReductions = append(itemDTO.StockReductions, dto.StockReduction{
					ID:              reduction.ID,
					QuantityReduced: reduction.QuantityReduced,
					StockBefore:     reduction.StockBefore,
					StockAfter:      reduction.StockAfter,
					Ingredient: dto.StockReductionIngredient{
						ID:   reduction.Ingredient.ID,
						Name: reduction.Ingredient.Name,
						Slug: reduction.Ingredient.Slug,
					},
					Unit: dto.StockReductionUnit{
						ID:   reduction.Unit.ID,
						Name: reduction.Unit.Name,
					},
				})
			}

			transactionDTO.Items = append(transactionDTO.Items, itemDTO)
		}

		response = append(response, transactionDTO)
	}

	if response == nil {
		response = make([]dto.Transaction, 0)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   response,
	})
}

// GetTransaction godoc
// @Summary Get Transaction
// @Description Get Transaction by ID
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Router /transactions/{id} [get]
func GetTransaction(c *gin.Context) {
	id := c.Param("id")
	var transaction models.Transaction

	if err := config.DB.
		Preload("TransactionItems").
		Preload("TransactionItems.Menu").
		Preload("TransactionItems.StockReductions").
		Preload("TransactionItems.StockReductions.Ingredient").
		Preload("TransactionItems.StockReductions.Unit").
		First(&transaction, id).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Transaction not found",
		})
		return
	}

	transactionDTO := dto.Transaction{
		ID:              transaction.ID,
		TransactionCode: transaction.TransactionCode,
		TransactionDate: transaction.TransactionDate,
		TotalAmount:     transaction.TotalAmount,
		Notes:           transaction.Notes,
		Status:          transaction.Status,
		CreatedAt:       transaction.CreatedAt,
		UpdatedAt:       transaction.UpdatedAt,
	}

	for _, item := range transaction.TransactionItems {
		itemDTO := dto.TransactionItem{
			ID:       item.ID,
			Quantity: item.Quantity,
			Price:    item.Price,
			Menu: dto.TransactionItemMenu{
				ID:    item.Menu.ID,
				Name:  item.Menu.Name,
				Slug:  item.Menu.Slug,
				Image: item.Menu.Image,
			},
		}

		for _, reduction := range item.StockReductions {
			itemDTO.StockReductions = append(itemDTO.StockReductions, dto.StockReduction{
				ID:              reduction.ID,
				QuantityReduced: reduction.QuantityReduced,
				StockBefore:     reduction.StockBefore,
				StockAfter:      reduction.StockAfter,
				Ingredient: dto.StockReductionIngredient{
					ID:   reduction.Ingredient.ID,
					Name: reduction.Ingredient.Name,
					Slug: reduction.Ingredient.Slug,
				},
				Unit: dto.StockReductionUnit{
					ID:   reduction.Unit.ID,
					Name: reduction.Unit.Name,
				},
			})
		}

		transactionDTO.Items = append(transactionDTO.Items, itemDTO)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   transactionDTO,
	})
}

// PostTransaction godoc
// @Summary Create Transaction
// @Description Create Transaction
// @Tags Transactions
// @Param transaction body dto.TransactionCreateRequest true "Create transaction"
// @Router /transactions [post]
func PostTransaction(c *gin.Context) {
	var input dto.TransactionCreateRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	tx := config.DB.Begin()

	// Generate transaction code
	transactionCode := fmt.Sprintf("TRX-%s-%d",
		time.Now().Format("20060102"),
		time.Now().Unix()%10000,
	)

	transaction := models.Transaction{
		TransactionCode: transactionCode,
		TransactionDate: time.Now(),
		TotalAmount:     0,
		Notes:           input.Notes,
		Status:          "completed",
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	var totalAmount float64

	for _, item := range input.Items {
		// Get menu data
		var menu models.Menu
		if err := tx.Preload("MenuIngredients").
			Preload("MenuIngredients.Unit").
			First(&menu, item.MenuID).Error; err != nil {

			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": fmt.Sprintf("Menu with ID %d not found", item.MenuID),
			})
			return
		}

		// Create transaction item
		transactionItem := models.TransactionItem{
			TransactionID: transaction.ID,
			MenuID:        menu.ID,
			Quantity:      item.Quantity,
			Price:         menu.Price,
		}

		if err := tx.Create(&transactionItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}

		totalAmount += menu.Price * float64(item.Quantity)

		// Process stock reduction for each ingredient
		for _, menuIngredient := range menu.MenuIngredients {
			// Get current ingredient stock
			var ingredient models.Ingredient
			if err := tx.First(&ingredient, menuIngredient.IngredientID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusNotFound, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("Ingredient with ID %d not found", menuIngredient.IngredientID),
				})
				return
			}

			// Calculate quantity to reduce
			quantityToReduce := menuIngredient.Quantity * float64(item.Quantity)

			// Check if stock is sufficient
			if ingredient.Stock < quantityToReduce {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{
					"status":  "error",
					"message": fmt.Sprintf("Insufficient stock for ingredient: %s. Available: %.2f, Required: %.2f", ingredient.Name, ingredient.Stock, quantityToReduce),
				})
				return
			}

			stockBefore := ingredient.Stock
			stockAfter := stockBefore - quantityToReduce

			// Create stock reduction record
			stockReduction := models.StockReduction{
				TransactionItemID: transactionItem.ID,
				IngredientID:      ingredient.ID,
				QuantityReduced:   quantityToReduce,
				StockBefore:       stockBefore,
				StockAfter:        stockAfter,
				UnitID:            menuIngredient.UnitID,
			}

			if err := tx.Create(&stockReduction).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			// Update ingredient stock
			ingredient.Stock = stockAfter
			if err := tx.Save(&ingredient).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}
		}
	}

	// Update total amount
	transaction.TotalAmount = totalAmount
	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Transaction created successfully",
		"data": gin.H{
			"transaction_code": transactionCode,
			"total_amount":     totalAmount,
		},
	})
}

// DeleteTransaction godoc
// @Summary Delete Transaction
// @Description Delete Transaction by ID (and restore stock)
// @Tags Transactions
// @Param id path int true "Transaction ID"
// @Router /transactions/{id} [delete]
func DeleteTransaction(c *gin.Context) {
	id := c.Param("id")

	tx := config.DB.Begin()

	var transaction models.Transaction
	if err := tx.Preload("TransactionItems").
		Preload("TransactionItems.StockReductions").
		First(&transaction, id).Error; err != nil {

		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Transaction not found",
		})
		return
	}

	// Restore stock for each ingredient
	for _, item := range transaction.TransactionItems {
		for _, reduction := range item.StockReductions {
			var ingredient models.Ingredient
			if err := tx.First(&ingredient, reduction.IngredientID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			// Restore stock
			ingredient.Stock += reduction.QuantityReduced
			if err := tx.Save(&ingredient).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}

			// Delete stock reduction record
			if err := tx.Delete(&reduction).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": err.Error(),
				})
				return
			}
		}

		// Delete transaction item
		if err := tx.Delete(&item).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
	}

	// Delete transaction
	if err := tx.Delete(&transaction).Error; err != nil {
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
		"message": "Transaction deleted successfully and stock restored",
	})
}
