package seeders

import (
	"AwisPalace_IngredientManagement/models"
	"AwisPalace_IngredientManagement/utils"
	"fmt"

	"gorm.io/gorm"
)

func MenuSeeder(db *gorm.DB) error {
	menus := []models.Menu{
		{
			Name:        "Nasi Goreng Spesial",
			Slug:        utils.GenerateSlug("Nasi Goreng Spesial"),
			Description: "Nasi goreng khas Awis Palace",
			Image:       "nasi-goreng.jpg",
			Price:       25000,
		},
	}

	for _, menu := range menus {
		var existing models.Menu
		if err := db.Where("slug = ?", menu.Slug).First(&existing).Error; err == nil {
			fmt.Printf("⚠️ Menu %s sudah ada, skip\n", menu.Name)
			continue
		}

		if err := db.Create(&menu).Error; err != nil {
			fmt.Printf("❌ Gagal membuat menu %s: %v\n", menu.Name, err)
			continue
		}

		menuIngredients := []models.MenuIngredient{
			{
				MenuID:       menu.ID,
				IngredientID: 1, // Nasi
				Quantity:     200,
				UnitID:       3, // gram
			},
			{
				MenuID:       menu.ID,
				IngredientID: 2, // Telur
				Quantity:     1,
				UnitID:       1, // pcs
			},
		}

		for _, mi := range menuIngredients {
			db.Create(&mi)
		}

		fmt.Printf("✅ Menu %s berhasil ditambahkan\n", menu.Name)
	}

	return nil
}
