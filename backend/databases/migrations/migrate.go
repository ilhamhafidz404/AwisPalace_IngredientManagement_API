package migrations

import (
	"AwisPalace_IngredientManagement/config"
	"AwisPalace_IngredientManagement/models"
	"fmt"
)

func Migrate() {
	err := config.DB.AutoMigrate(
		&models.User{},
	)

	if err != nil {
		fmt.Println("❌ Migration failed:", err)
	} else {
		fmt.Println("✅ Migration success!")
	}
}
