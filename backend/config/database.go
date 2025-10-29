// package config

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"os"

// 	_ "github.com/lib/pq"
// )

// var DB *sql.DB

// func ConnectDB() {
// 	host := os.Getenv("DB_HOST")
// 	port := os.Getenv("POSTGRES_PORT")
// 	user := os.Getenv("POSTGRES_USER")
// 	password := os.Getenv("POSTGRES_PASSWORD")
// 	dbname := os.Getenv("POSTGRES_DB")

// 	connStr := fmt.Sprintf(
// 		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbname,
// 	)

// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Fatal("❌ Failed to connect DB:", err)
// 	}

// 	if err := db.Ping(); err != nil {
// 		log.Fatal("❌ DB unreachable:", err)
// 	}

// 	DB = db
// 	log.Println("✅ Database connected")
// }

package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASS"),
        os.Getenv("DB_NAME"),
        os.Getenv("DB_PORT"),
    )

    var err error
    for i := 0; i < 10; i++ {
        DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            log.Println("✅ Connected to database")
            return
        }

        log.Printf("❌ DB not ready, retrying in 3s... (%d/10)", i+1)
        time.Sleep(3 * time.Second)
    }

    log.Fatalf("❌ Failed to connect to database after retries: %v", err)
}
