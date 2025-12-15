package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction adalah record penjualan menu
type Transaction struct {
	gorm.Model
	TransactionCode string    `gorm:"type:varchar(50);uniqueIndex;not null"` // Kode transaksi unik
	TransactionDate time.Time `gorm:"not null"`                              // Tanggal transaksi
	TotalAmount     float64   `gorm:"type:numeric(12,2)"`                    // Total harga transaksi
	Notes           string    `gorm:"type:text"`                             // Catatan tambahan (opsional)
	Status          string    `gorm:"type:varchar(20);default:'completed'"`  // Status: completed, cancelled, etc.

	TransactionItems []TransactionItem // Detail item yang terjual
}

// TransactionItem adalah detail menu yang terjual dalam satu transaksi
type TransactionItem struct {
	gorm.Model
	TransactionID uint
	Transaction   Transaction

	MenuID   uint
	Menu     Menu
	Quantity int     `gorm:"not null"`                    // Jumlah menu yang terjual
	Price    float64 `gorm:"type:numeric(12,2);not null"` // Harga menu saat transaksi (untuk historical data)

	StockReductions []StockReduction // Detail pengurangan stok per ingredient
}

// StockReduction adalah record pengurangan stok ingredient akibat transaksi
type StockReduction struct {
	gorm.Model
	TransactionItemID uint
	TransactionItem   TransactionItem

	IngredientID    uint
	Ingredient      Ingredient
	QuantityReduced float64 `gorm:"type:numeric(10,2);not null"` // Jumlah yang dikurangi
	StockBefore     float64 `gorm:"type:numeric(10,2);not null"` // Stok sebelum pengurangan
	StockAfter      float64 `gorm:"type:numeric(10,2);not null"` // Stok setelah pengurangan
	UnitID          uint
	Unit            Unit
}
