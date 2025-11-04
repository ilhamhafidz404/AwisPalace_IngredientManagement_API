package models

import "gorm.io/gorm"

type Ingredient struct {
	gorm.Model
	Name   string  `gorm:"type:varchar(100);not null"`
	Slug   string  `gorm:"type:varchar(100);uniqueIndex"`
	Stock  float64 `gorm:"type:numeric(10,2)"`
	UnitID uint    `gorm:"type:integer"`
	Unit   Unit    `gorm:"foreignKey:UnitID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
