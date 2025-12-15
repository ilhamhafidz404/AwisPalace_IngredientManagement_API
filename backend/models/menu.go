package models

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name        string  `gorm:"type:varchar(150);not null"`
	Slug        string  `gorm:"type:varchar(150);uniqueIndex"`
	Description string  `gorm:"type:text"`
	Image       string  `gorm:"type:text"`
	Price       float64 `gorm:"type:numeric(12,2)"`

	MenuIngredients  []MenuIngredient
	TransactionItems []TransactionItem
}
