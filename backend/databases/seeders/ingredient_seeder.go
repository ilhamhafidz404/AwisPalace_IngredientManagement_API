package seeders

import (
	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/models"
	"fmt"

	"gorm.io/gorm"
)

func IngredientSeeder(db *gorm.DB) error {
	ingredients := []models.Ingredient{
		{
			Name:   "Garam",
			Slug:   "garam",
			UnitID: 1,
			Stock:  100,
		},
		{
			Name:   "Gula",
			Slug:   "gula",
			UnitID: 2,
			Stock:  80,
		},
		{
			Name:   "Tepung Terigu",
			Slug:   "tepung-terigu",
			UnitID: 3,
			Stock:  50,
		},
	}

	for _, ingredient := range ingredients {
		var existing models.Ingredient
		if err := config.DB.Where("slug = ?", ingredient.Slug).First(&existing).Error; err != nil {
			if err := config.DB.Create(&ingredient).Error; err != nil {
				fmt.Printf("❌ Gagal menambahkan ingredient %s: %v\n", ingredient.Name, err)
			} else {
				fmt.Printf("✅ Berhasil menambahkan ingredient: %s\n", ingredient.Name)
			}
		} else {
			fmt.Printf("⚠️ Ingredient %s sudah ada, skip.\n", ingredient.Name)
		}
	}

	return nil
}
