package models

import "gorm.io/gorm"

type Unit struct {
	gorm.Model
	Name   string `gorm:"type:varchar(50);not null"`
	Symbol string `gorm:"type:varchar(10);uniqueIndex"`
}
