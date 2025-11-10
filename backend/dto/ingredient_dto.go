package dto

import (
	"time"
)

type Ingredient struct {
	ID       	uint 		`json:"id"`
	Name  		string  	`json:"name"`
	Slug   		string  	`json:"slug"`
	Stock  		float64 	`json:"stock"`
	UnitID 		uint    	`json:"unit_id"`
	Unit   		Unit    	`json:"unit"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdatedAt 	time.Time 	`json:"updated_at"`
	DeletedAt 	time.Time 	`json:"deleted_at"`
}


type IngredientParamRequest struct {
	Name  		string  	`json:"name"`
	Stock  		float64 	`json:"stock"`
	UnitID 		uint    	`json:"unit_id"`
}
