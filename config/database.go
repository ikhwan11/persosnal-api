package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func ConnectDatabase() {
	// Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env file tidak ditemukan, pakai environment system")
	}

	// Ambil variabel dari .env
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")

	// Format DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, name)

	// Buka koneksi
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Gagal membuka koneksi:", err)
	}

	// Test koneksi
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Gagal koneksi database:", err)
	}

	log.Println("✅ Database connected")
}

