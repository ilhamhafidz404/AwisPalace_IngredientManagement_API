package models

import "gorm.io/gorm"

type Unit struct {
	gorm.Model
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}
