package models

import "gorm.io/gorm"

type MenuIngredient struct {
	gorm.Model

	MenuID uint
	Menu   Menu

	IngredientID uint
	Ingredient   Ingredient

	Quantity float64 `gorm:"type:numeric(10,2);not null"`
	UnitID   uint
	Unit     Unit
}
