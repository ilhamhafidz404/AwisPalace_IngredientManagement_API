package seeders

import (
	"fmt"

	"gorm.io/gorm"
)
func DatabaseSeeder(db *gorm.DB) {
	fmt.Println("🚀 Menjalankan Database Seeder...")

	// Daftar seeder yang ingin dijalankan
	seeders := []func(*gorm.DB) error{
		SeedUnits,
	}

	for _, seed := range seeders {
		if err := seed(db); err != nil {
			fmt.Println("❌ Gagal menjalankan seeder:", err)
		}
	}

	fmt.Println("✅ Semua seeder berhasil dijalankan!")
}
