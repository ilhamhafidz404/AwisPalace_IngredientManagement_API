package dto

import "time"

//
// ===== RESPONSE DTO =====
//

type Menu struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Slug        string           `json:"slug"`
	Image       string           `json:"image"`
	Price       float64          `json:"price"`
	Description string           `json:"description"`
	Ingredients []MenuIngredient `json:"ingredients"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
	DeletedAt   *time.Time       `json:"deleted_at,omitempty"`
}

type MenuIngredient struct {
	ID         uint                     `json:"id"`
	Ingredient MenuIngredientIngredient `json:"ingredient"`
	Quantity   float64                  `json:"quantity"`
	Unit       MenuIngredientUnit       `json:"unit"`
}

type MenuIngredientIngredient struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type MenuIngredientUnit struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

//
// ===== REQUEST DTO =====
//

// untuk CREATE & UPDATE
type MenuCreateRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Image       string                  `json:"image" binding:"required"`
	Price       float64                 `json:"price" binding:"required"`
	Description string                  `json:"description"`
	Ingredients []MenuIngredientRequest `json:"ingredients" binding:"required"`
}

type MenuUpdateRequest struct {
	Name        string                  `json:"name" binding:"required"`
	Image       string                  `json:"image" binding:"required"`
	Price       float64                 `json:"price" binding:"required"`
	Description string                  `json:"description"`
	Ingredients []MenuIngredientRequest `json:"ingredients" binding:"required"`
}

type MenuIngredientRequest struct {
	IngredientID uint    `json:"ingredient_id" binding:"required"`
	Quantity     float64 `json:"quantity" binding:"required"`
	UnitID       uint    `json:"unit_id" binding:"required"`
}
