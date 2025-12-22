package controllers

import (
	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/models"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

// ExportTransactionDTO represents the data structure for export
type ExportTransactionDTO struct {
	TransactionID   uint
	TransactionCode string
	TransactionDate time.Time
	MenuName        string
	MenuQuantity    int
	MenuPrice       float64
	MenuSubtotal    float64
	TotalAmount     float64
	Status          string
	Notes           string
}

// IngredientUsageDTO represents ingredient usage per transaction
type IngredientUsageDTO struct {
	TransactionID   uint
	TransactionCode string
	TransactionDate time.Time
	MenuName        string
	IngredientName  string
	QuantityReduced float64
	UnitName        string
	StockBefore     float64
	StockAfter      float64
}

// ExportTransactions godoc
// @Summary Export Transactions to Excel
// @Description Export transaction data with ingredient usage to Excel file. Returns Excel file with 3 sheets: Transactions (all transaction details), Ingredient Usage (ingredients used per transaction with stock changes), and Summary (statistics and top selling items). Supports date range filtering, defaults to last 30 days if dates not specified.
// @Tags Transactions
// @Accept json
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param start_date query string false "Start date in YYYY-MM-DD format. Defaults to 30 days ago if not specified." example(2024-01-01)
// @Param end_date query string false "End date in YYYY-MM-DD format. Defaults to today if not specified." example(2024-12-31)
// @Router /export/transactions [get]
func ExportTransactions(c *gin.Context) {
	// Get query parameters
	startDateStr := c.Query("start_date") // Format: 2024-01-01
	endDateStr := c.Query("end_date")     // Format: 2024-12-31

	// Parse dates
	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid start_date format. Use YYYY-MM-DD",
			})
			return
		}
	} else {
		// Default: 30 days ago
		startDate = time.Now().AddDate(0, 0, -30)
	}

	if endDateStr != "" {
		endDate, err = time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": "Invalid end_date format. Use YYYY-MM-DD",
			})
			return
		}
		// Set to end of day
		endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	} else {
		// Default: today end of day
		endDate = time.Now().Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	}

	// Fetch transactions with details
	var transactions []models.Transaction
	if err := config.DB.
		Preload("TransactionItems.Menu").
		Preload("TransactionItems.StockReductions.Ingredient").
		Preload("TransactionItems.StockReductions.Unit").
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Order("transaction_date DESC").
		Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch transactions: " + err.Error(),
		})
		return
	}

	if len(transactions) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "No transactions found in the specified date range",
		})
		return
	}

	// Create Excel file
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	// Create sheets
	f.SetSheetName("Sheet1", "Transactions")
	f.NewSheet("Ingredient Usage")
	f.NewSheet("Summary")

	// Fill Transactions sheet
	fillTransactionsSheet(f, transactions)

	// Fill Ingredient Usage sheet
	fillIngredientUsageSheet(f, transactions)

	// Fill Summary sheet
	fillSummarySheet(f, transactions, startDate, endDate)

	// Generate filename
	filename := fmt.Sprintf("transactions_%s_to_%s.xlsx",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02"))

	// Set headers for download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	// Write to response
	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to generate Excel file: " + err.Error(),
		})
		return
	}
}

// fillTransactionsSheet fills the transactions sheet
func fillTransactionsSheet(f *excelize.File, transactions []models.Transaction) {
	sheet := "Transactions"

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 15)
	f.SetColWidth(sheet, "B", "B", 20)
	f.SetColWidth(sheet, "C", "C", 20)
	f.SetColWidth(sheet, "D", "D", 25)
	f.SetColWidth(sheet, "E", "E", 12)
	f.SetColWidth(sheet, "F", "F", 15)
	f.SetColWidth(sheet, "G", "G", 15)
	f.SetColWidth(sheet, "H", "H", 15)
	f.SetColWidth(sheet, "I", "I", 30)

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"4472C4"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Set headers
	headers := []string{
		"ID", "Transaction Code", "Date", "Menu Name",
		"Quantity", "Price", "Subtotal", "Total Amount", "Notes",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Fill data
	row := 2
	for _, trx := range transactions {
		for _, item := range trx.TransactionItems {
			subtotal := float64(item.Quantity) * item.Price

			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), trx.ID)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), trx.TransactionCode)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), trx.TransactionDate.Format("2006-01-02 15:04:05"))
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Menu.Name)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Quantity)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.Price)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), subtotal)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), trx.TotalAmount)
			f.SetCellValue(sheet, fmt.Sprintf("I%d", row), trx.Notes)
			row++
		}
	}

	// Add total row
	f.SetCellValue(sheet, fmt.Sprintf("F%d", row), "TOTAL:")
	f.SetCellFormula(sheet, fmt.Sprintf("H%d", row), fmt.Sprintf("SUM(H2:H%d)", row-1))

	// Style total row
	totalStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"D9E1F2"}, Pattern: 1},
	})
	f.SetCellStyle(sheet, fmt.Sprintf("F%d", row), fmt.Sprintf("H%d", row), totalStyle)
}

// fillIngredientUsageSheet fills the ingredient usage sheet
func fillIngredientUsageSheet(f *excelize.File, transactions []models.Transaction) {
	sheet := "Ingredient Usage"

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 15)
	f.SetColWidth(sheet, "B", "B", 20)
	f.SetColWidth(sheet, "C", "C", 20)
	f.SetColWidth(sheet, "D", "D", 25)
	f.SetColWidth(sheet, "E", "E", 25)
	f.SetColWidth(sheet, "F", "F", 15)
	f.SetColWidth(sheet, "G", "G", 12)
	f.SetColWidth(sheet, "H", "H", 15)
	f.SetColWidth(sheet, "I", "I", 15)

	// Create header style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"70AD47"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	// Set headers
	headers := []string{
		"Transaction ID", "Transaction Code", "Date", "Menu Name",
		"Ingredient Name", "Qty Reduced", "Unit", "Stock Before", "Stock After",
	}

	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, header)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	// Fill data
	row := 2
	for _, trx := range transactions {
		for _, item := range trx.TransactionItems {
			for _, reduction := range item.StockReductions {
				f.SetCellValue(sheet, fmt.Sprintf("A%d", row), trx.ID)
				f.SetCellValue(sheet, fmt.Sprintf("B%d", row), trx.TransactionCode)
				f.SetCellValue(sheet, fmt.Sprintf("C%d", row), trx.TransactionDate.Format("2006-01-02 15:04:05"))
				f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.Menu.Name)
				f.SetCellValue(sheet, fmt.Sprintf("E%d", row), reduction.Ingredient.Name)
				f.SetCellValue(sheet, fmt.Sprintf("F%d", row), reduction.QuantityReduced)
				f.SetCellValue(sheet, fmt.Sprintf("G%d", row), reduction.Unit.Name)
				f.SetCellValue(sheet, fmt.Sprintf("H%d", row), reduction.StockBefore)
				f.SetCellValue(sheet, fmt.Sprintf("I%d", row), reduction.StockAfter)
				row++
			}
		}
	}
}

// fillSummarySheet fills the summary sheet
func fillSummarySheet(f *excelize.File, transactions []models.Transaction, startDate, endDate time.Time) {
	sheet := "Summary"

	// Set column widths
	f.SetColWidth(sheet, "A", "A", 30)
	f.SetColWidth(sheet, "B", "B", 20)

	// Title style
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 16},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "center"},
	})

	// Label style
	labelStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"E7E6E6"}, Pattern: 1},
	})

	row := 1
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "TRANSACTION REPORT")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), titleStyle)

	row += 2
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Period")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%s to %s",
		startDate.Format("2006-01-02"),
		endDate.Format("2006-01-02")))

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Generated At")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), time.Now().Format("2006-01-02 15:04:05"))

	row += 2
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "STATISTICS")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), titleStyle)

	// Calculate statistics
	totalTransactions := len(transactions)
	totalRevenue := 0.0
	totalItems := 0
	totalMenusSold := 0

	for _, trx := range transactions {
		totalRevenue += trx.TotalAmount
		totalItems += len(trx.TransactionItems)
		for _, item := range trx.TransactionItems {
			totalMenusSold += item.Quantity
		}
	}

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total Transactions")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), totalTransactions)

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total Menu Items")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), totalItems)

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total Menus Sold")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), totalMenusSold)

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total Revenue")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("Rp %.2f", totalRevenue))

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Average Transaction Value")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), labelStyle)
	if totalTransactions > 0 {
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("Rp %.2f", totalRevenue/float64(totalTransactions)))
	}

	// Ingredient usage summary
	row += 2
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "TOP INGREDIENTS USED")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), titleStyle)

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Ingredient")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Total Quantity Reduced")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), labelStyle)

	// Aggregate ingredient usage from StockReductions
	ingredientMap := make(map[string]float64)
	ingredientUnit := make(map[string]string)

	for _, trx := range transactions {
		for _, item := range trx.TransactionItems {
			for _, reduction := range item.StockReductions {
				key := reduction.Ingredient.Name
				ingredientMap[key] += reduction.QuantityReduced
				ingredientUnit[key] = reduction.Unit.Name
			}
		}
	}

	// Sort and display top ingredients
	for name, qty := range ingredientMap {
		row++
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), fmt.Sprintf("%.2f %s", qty, ingredientUnit[name]))
	}

	// Menu popularity
	row += 2
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "TOP SELLING MENUS")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), titleStyle)

	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Menu Name")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), "Total Sold")
	f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("B%d", row), labelStyle)

	// Aggregate menu sales
	menuMap := make(map[string]int)

	for _, trx := range transactions {
		for _, item := range trx.TransactionItems {
			menuMap[item.Menu.Name] += item.Quantity
		}
	}

	// Display menu sales
	for name, qty := range menuMap {
		row++
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), qty)
	}
}
