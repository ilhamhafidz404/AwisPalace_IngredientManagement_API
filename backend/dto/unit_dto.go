package dto

import "time"

type Unit struct {
	ID        uint `json:"id"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
