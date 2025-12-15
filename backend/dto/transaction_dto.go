package dto

import "time"

// Transaction DTOs
type Transaction struct {
	ID              uint              `json:"id"`
	TransactionCode string            `json:"transaction_code"`
	TransactionDate time.Time         `json:"transaction_date"`
	TotalAmount     float64           `json:"total_amount"`
	Notes           string            `json:"notes"`
	Status          string            `json:"status"`
	Items           []TransactionItem `json:"items"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

type TransactionItem struct {
	ID              uint                `json:"id"`
	Menu            TransactionItemMenu `json:"menu"`
	Quantity        int                 `json:"quantity"`
	Price           float64             `json:"price"`
	StockReductions []StockReduction    `json:"stock_reductions"`
}

type TransactionItemMenu struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Slug  string `json:"slug"`
	Image string `json:"image"`
}

type StockReduction struct {
	ID              uint                     `json:"id"`
	Ingredient      StockReductionIngredient `json:"ingredient"`
	QuantityReduced float64                  `json:"quantity_reduced"`
	StockBefore     float64                  `json:"stock_before"`
	StockAfter      float64                  `json:"stock_after"`
	Unit            StockReductionUnit       `json:"unit"`
}

type StockReductionIngredient struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type StockReductionUnit struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Request DTOs
type TransactionCreateRequest struct {
	Items []TransactionItemRequest `json:"items" binding:"required"`
	Notes string                   `json:"notes"`
}

type TransactionItemRequest struct {
	MenuID   uint `json:"menu_id" binding:"required"`
	Quantity int  `json:"quantity" binding:"required,min=1"`
}
