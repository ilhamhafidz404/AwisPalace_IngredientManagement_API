package seeders

import (
	"AwisPalace_IngredientManagement/models"
	"fmt"

	"gorm.io/gorm"
)

func SeedUnits(db *gorm.DB) error {
	units := []models.Unit{
		{Name: "Kilogram", Symbol: "kg"},
		{Name: "Gram", Symbol: "g"},
		{Name: "Liter", Symbol: "L"},
		{Name: "Meter", Symbol: "m"},
		{Name: "Packages", Symbol: "pcs"},
		{Name: "Persen", Symbol: "%"},
	}

	for _, unit := range units {
		var existing models.Unit
		// Cek apakah data sudah ada berdasarkan nama (agar tidak duplicate)
		if err := db.Where("name = ?", unit.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := db.Create(&unit).Error; err != nil {
				return fmt.Errorf("gagal menambahkan unit %s: %v", unit.Name, err)
			}
			fmt.Printf("✅ Unit %s berhasil ditambahkan\n", unit.Name)
		} else {
			fmt.Printf("⚠️  Unit %s sudah ada, dilewati\n", unit.Name)
		}
	}

	return nil
}
