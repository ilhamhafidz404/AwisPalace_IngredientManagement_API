package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load ENV
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  Tidak menemukan file .env, menggunakan environment variable sistem...")
	}

	// Get Env
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	// Validation ENV
	if dbHost == "" || dbPort == "" || dbUser == "" || dbName == "" {
		log.Fatal("❌ Environment variable database belum lengkap. Pastikan DB_HOST, POSTGRES_PORT, POSTGRES_USER, dan POSTGRES_DB sudah diset.")
	}

	// Connection String
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName,
	)

	log.Println("🔗 Connecting to PostgreSQL with:", connStr)

	// Connection DB
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("❌ Gagal membuka koneksi:", err)
	}
	defer db.Close()

	// PING to DB
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Tidak bisa terhubung ke database:", err)
	}
	log.Println("✅ Berhasil terhubung ke PostgreSQL!")

	// Router
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			log.Println("❌ DB connection failed:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"message": "DB connection failed"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Connected to PostgreSQL!"})
	})

	r.Run(":8080")
}
